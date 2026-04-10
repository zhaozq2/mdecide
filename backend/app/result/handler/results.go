package handler

import (
	"encoding/json"
	"net/http"

	"mdecide/app/options/model"
	resultModel "mdecide/app/result/model"
	roundModel "mdecide/app/round/model"
	"mdecide/common/response"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var resultDB *gorm.DB

func InitResultDB(database *gorm.DB) {
	resultDB = database
}

type ResultHandler struct{}

func NewResultHandler() *ResultHandler {
	return &ResultHandler{}
}

func RegisterResultRoutes(r *mux.Router) {
	h := NewResultHandler()
	r.HandleFunc("/api/results", h.List).Methods("GET")
	r.HandleFunc("/api/results", h.Calculate).Methods("POST")
	r.HandleFunc("/api/results/summary", h.GetSummary).Methods("GET")
	r.HandleFunc("/api/results", h.Create).Methods("POST")
	r.HandleFunc("/api/results/{id}", h.Get).Methods("GET")
	r.HandleFunc("/api/results/winner", h.GetWinner).Methods("GET")
}

func (h *ResultHandler) List(w http.ResponseWriter, r *http.Request) {
	topicID := r.URL.Query().Get("topicId")
	roundID := r.URL.Query().Get("roundId")

	var results []resultModel.Result
	query := resultDB.Model(&resultModel.Result{})

	if topicID != "" {
		query = query.Where("topic_id = ?", topicID)
	}
	if roundID != "" {
		query = query.Where("round_id = ?", roundID)
	}

	query.Order("total_score DESC").Find(&results)
	response.Ok(results).WriteTo(w)
}

type CalculateRequest struct {
	RoundID uint `json:"roundId"`
	TopicID uint `json:"topicId"`
}

func (h *ResultHandler) Calculate(w http.ResponseWriter, r *http.Request) {
	var req CalculateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(http.StatusBadRequest, "invalid request body").WriteTo(w)
		return
	}

	var options []model.Option
	if err := resultDB.Where("round_id = ? AND topic_id = ?", req.RoundID, req.TopicID).Find(&options).Error; err != nil {
		response.Error(http.StatusInternalServerError, "failed to query options").WriteTo(w)
		return
	}

	resultDB.Where("round_id = ? AND topic_id = ?", req.RoundID, req.TopicID).Delete(&resultModel.Result{})

	var results []resultModel.Result
	for _, opt := range options {
		totalScore := opt.Importance * opt.Necessity
		res := resultModel.Result{
			OptionID:        opt.ID,
			RoundID:         int(req.RoundID),
			TopicID:         int(req.TopicID),
			ImportanceScore: opt.Importance,
			NecessityScore:  opt.Necessity,
			TotalScore:      totalScore,
		}
		results = append(results, res)
	}
	resultDB.Create(&results)

	matchScore := calculateMatchScore(options)
	resultDB.Model(&roundModel.Round{}).Where("id = ?", req.RoundID).Update("match_score", matchScore)

	response.Ok(map[string]any{
		"results":    results,
		"matchScore": matchScore,
	}).WriteTo(w)
}

func calculateMatchScore(options []model.Option) int {
	if len(options) == 0 {
		return 0
	}
	totalWeight := 0
	necessaryWeight := 0
	for _, opt := range options {
		totalWeight += opt.Importance
		if opt.Necessity == 1 {
			necessaryWeight += opt.Importance
		}
	}
	if totalWeight == 0 {
		return 0
	}
	return (necessaryWeight * 100) / totalWeight
}

func (h *ResultHandler) GetSummary(w http.ResponseWriter, r *http.Request) {
	topicID := r.URL.Query().Get("topicId")
	roundID := r.URL.Query().Get("roundId")

	if topicID == "" {
		response.Error(http.StatusBadRequest, "topicId is required").WriteTo(w)
		return
	}

	var results []resultModel.Result
	query := resultDB.Model(&resultModel.Result{}).Where("topic_id = ?", topicID)

	if roundID != "" {
		query = query.Where("round_id = ?", roundID)
	}

	if err := query.Order("total_score DESC").Find(&results).Error; err != nil {
		response.Error(http.StatusInternalServerError, "failed to query results").WriteTo(w)
		return
	}

	if len(results) == 0 {
		response.Ok(map[string]any{
			"winnerName":   "",
			"winnerScore":  0,
			"description":  "暂无评分数据",
			"totalRounds":  0,
			"totalOptions": 0,
		}).WriteTo(w)
		return
	}

	winner := results[0]
	var options []model.Option
	resultDB.Where("topic_id = ?", topicID).Find(&options)

	type Summary struct {
		WinnerName   string `json:"winnerName"`
		WinnerScore  int    `json:"winnerScore"`
		Description  string `json:"description"`
		TotalRounds  int    `json:"totalRounds"`
		TotalOptions int    `json:"totalOptions"`
	}

	desc := "经过重要性评分和必要性筛选，"
	if winner.TotalScore > 0 {
		var opt model.Option
		if resultDB.First(&opt, winner.OptionID).Error == nil {
			desc += opt.Title + "综合得分最高"
		} else {
			desc += "综合得分最高"
		}
	} else {
		desc += "暂无获胜选项"
	}

	response.Ok(Summary{
		WinnerName:   "",
		WinnerScore:  winner.TotalScore,
		Description:  desc,
		TotalRounds:  0,
		TotalOptions: len(options),
	}).WriteTo(w)
}

func (h *ResultHandler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var result resultModel.Result
	if err := resultDB.First(&result, vars["id"]).Error; err != nil {
		response.Error(http.StatusNotFound, "result not found").WriteTo(w)
		return
	}
	response.Ok(result).WriteTo(w)
}

func (h *ResultHandler) GetWinner(w http.ResponseWriter, r *http.Request) {
	topicID := r.URL.Query().Get("topicId")
	roundID := r.URL.Query().Get("roundId")

	if topicID == "" {
		response.Error(http.StatusBadRequest, "topicId is required").WriteTo(w)
		return
	}

	query := resultDB.Model(&resultModel.Result{}).Where("topic_id = ?", topicID).Order("total_score DESC").Limit(1)

	if roundID != "" {
		query = query.Where("round_id = ?", roundID)
	}

	var result resultModel.Result
	if err := query.First(&result).Error; err != nil {
		response.Ok(nil).WriteTo(w)
		return
	}

	var opt model.Option
	if err := resultDB.First(&opt, result.OptionID).Error; err != nil {
		response.Ok(map[string]any{
			"optionId":   result.OptionID,
			"optionName": "",
			"totalScore": result.TotalScore,
		}).WriteTo(w)
		return
	}

	response.Ok(map[string]any{
		"optionId":   result.OptionID,
		"optionName": opt.Title,
		"totalScore": result.TotalScore,
	}).WriteTo(w)
}

func (h *ResultHandler) Create(w http.ResponseWriter, r *http.Request) {
	var res resultModel.Result
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		response.Error(http.StatusBadRequest, "invalid request body").WriteTo(w)
		return
	}

	resultDB.Create(&res)
	response.Ok(res).WriteTo(w)
}
