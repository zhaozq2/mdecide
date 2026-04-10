package handler

import (
	"encoding/json"
	"net/http"

	"mdecide/app/round/model"
	"mdecide/common/response"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var roundDB *gorm.DB

func InitRoundDB(database *gorm.DB) {
	roundDB = database
}

type CreateRoundRequest struct {
	TopicID        uint   `json:"topicId"`
	OptionStrategy string `json:"optionStrategy"`
}

func RegisterRoundRoutes(r *mux.Router) {
	h := NewRoundHandler()
	r.HandleFunc("/api/rounds", h.List).Methods("GET")
	r.HandleFunc("/api/rounds", h.Create).Methods("POST")
	r.HandleFunc("/api/rounds/{id}", h.Get).Methods("GET")
	r.HandleFunc("/api/rounds/{id}", h.Update).Methods("PUT")
	r.HandleFunc("/api/rounds/{id}", h.Delete).Methods("DELETE")
	r.HandleFunc("/api/rounds/{id}/complete", h.Complete).Methods("POST")
	r.HandleFunc("/api/rounds/active", h.GetActive).Methods("GET")
}

type RoundHandler struct{}

func NewRoundHandler() *RoundHandler {
	return &RoundHandler{}
}

func (h *RoundHandler) List(w http.ResponseWriter, r *http.Request) {
	var rounds []model.Round
	roundDB.Find(&rounds)
	response.Ok(rounds).WriteTo(w)
}

func (h *RoundHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateRoundRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(http.StatusBadRequest, "invalid request body").WriteTo(w)
		return
	}

	optionStrategy := req.OptionStrategy
	if optionStrategy == "" {
		optionStrategy = "reuse"
	}

	var count int64
	roundDB.Model(&model.Round{}).Where("topic_id = ?", req.TopicID).Count(&count)

	round := model.Round{
		TopicID:        req.TopicID,
		RoundNumber:    int(count) + 1,
		Status:         "pending",
		IsActive:       true,
		OptionStrategy: optionStrategy,
	}

	roundDB.Model(&model.Round{}).Where("topic_id = ?", req.TopicID).Update("is_active", false)
	roundDB.Create(&round)

	if optionStrategy == "copy" {
		var options []model.Option
		roundDB.Where("topic_id = ?", req.TopicID).Find(&options)
		for _, opt := range options {
			newOpt := model.Option{
				Title:      opt.Title,
				TopicID:    req.TopicID,
				RoundID:    round.ID,
				Importance: 0,
				Necessity:  0,
				SortOrder:  opt.SortOrder,
				IsActive:   true,
			}
			roundDB.Create(&newOpt)
		}
	}

	response.Ok(round).WriteTo(w)
}

func (h *RoundHandler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var round model.Round
	if err := roundDB.First(&round, vars["id"]).Error; err != nil {
		response.Error(http.StatusNotFound, "round not found").WriteTo(w)
		return
	}
	response.Ok(round).WriteTo(w)
}

func (h *RoundHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var round model.Round
	if err := roundDB.First(&round, vars["id"]).Error; err != nil {
		response.Error(http.StatusNotFound, "round not found").WriteTo(w)
		return
	}

	if round.Status == "locked" {
		response.Error(http.StatusForbidden, "round is locked, cannot be modified").WriteTo(w)
		return
	}

	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(http.StatusBadRequest, "invalid request body").WriteTo(w)
		return
	}

	roundDB.Model(&round).Updates(req)
	response.Ok(round).WriteTo(w)
}

func (h *RoundHandler) Complete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var round model.Round
	if err := roundDB.First(&round, vars["id"]).Error; err != nil {
		response.Error(http.StatusNotFound, "round not found").WriteTo(w)
		return
	}

	if round.Status == "locked" {
		response.Error(http.StatusForbidden, "round is already locked").WriteTo(w)
		return
	}

	roundDB.Model(&round).Updates(map[string]interface{}{
		"status":            "locked",
		"is_active":         false,
		"importance_status": "completed",
		"necessity_status":  "completed",
	})

	response.Ok(round).WriteTo(w)
}

func (h *RoundHandler) GetActive(w http.ResponseWriter, r *http.Request) {
	var round model.Round
	if err := roundDB.Where("is_active = ?", true).First(&round).Error; err != nil {
		response.Ok(nil).WriteTo(w)
		return
	}
	response.Ok(round).WriteTo(w)
}

func (h *RoundHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var round model.Round
	if err := roundDB.First(&round, vars["id"]).Error; err != nil {
		response.Error(http.StatusNotFound, "round not found").WriteTo(w)
		return
	}

	// 删除轮次关联的选项
	roundDB.Where("round_id = ?", round.ID).Delete(&model.Option{})

	// 删除轮次
	roundDB.Delete(&round)

	response.Ok(map[string]string{"message": "round deleted"}).WriteTo(w)
}
