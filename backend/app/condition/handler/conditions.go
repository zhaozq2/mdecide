package handler

import (
	"encoding/json"
	"net/http"

	"mdecide/app/condition/model"
	optionModel "mdecide/app/options/model"
	"mdecide/common/response"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var condDB *gorm.DB

func InitCondDB(database *gorm.DB) {
	condDB = database
}

type ConditionHandler struct{}

func NewConditionHandler() *ConditionHandler {
	return &ConditionHandler{}
}

func RegisterConditionRoutes(r *mux.Router) {
	h := NewConditionHandler()
	r.HandleFunc("/api/conditions", h.List).Methods("GET")
	r.HandleFunc("/api/conditions", h.Create).Methods("POST")
	r.HandleFunc("/api/conditions/batch", h.BatchCreate).Methods("POST")
	r.HandleFunc("/api/conditions/importance", h.SetImportance).Methods("POST")
	r.HandleFunc("/api/conditions/necessity", h.SetNecessity).Methods("POST")
	r.HandleFunc("/api/conditions/{id}", h.Update).Methods("PUT")
	r.HandleFunc("/api/conditions/{id}", h.Delete).Methods("DELETE")
}

func (h *ConditionHandler) List(w http.ResponseWriter, r *http.Request) {
	topicID := r.URL.Query().Get("topicId")
	roundID := r.URL.Query().Get("roundId")

	var conditions []model.Condition
	query := condDB.Model(&model.Condition{})

	if topicID != "" {
		query = query.Where("topic_id = ?", topicID)
	}
	if roundID != "" {
		query = query.Where("round_id = ?", roundID)
	}

	query.Find(&conditions)
	response.Ok(conditions).WriteTo(w)
}

func (h *ConditionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var cond model.Condition
	if err := json.NewDecoder(r.Body).Decode(&cond); err != nil {
		response.Error(http.StatusBadRequest, "invalid request body").WriteTo(w)
		return
	}

	condDB.Create(&cond)
	response.Ok(cond).WriteTo(w)
}

func (h *ConditionHandler) BatchCreate(w http.ResponseWriter, r *http.Request) {
	var conditions []model.Condition
	if err := json.NewDecoder(r.Body).Decode(&conditions); err != nil {
		response.Error(http.StatusBadRequest, "invalid request body").WriteTo(w)
		return
	}

	condDB.Create(&conditions)
	response.Ok(conditions).WriteTo(w)
}

func (h *ConditionHandler) SetImportance(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		TopicID    uint `json:"topicId"`
		RoundID    uint `json:"roundId"`
		OptionID   uint `json:"optionId"`
		Importance int  `json:"importance"`
	}

	var req Req
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(http.StatusBadRequest, "invalid request body").WriteTo(w)
		return
	}

	// 更新或创建 importance 条件
	var cond model.Condition
	query := condDB.Where("option_id = ? AND phase = ?", req.OptionID, "importance")

	if req.RoundID > 0 {
		query = query.Where("round_id = ?", req.RoundID)
	}

	if err := query.First(&cond).Error; err != nil {
		cond = model.Condition{
			Name:     "重要性评分",
			Phase:    "importance",
			Score:    req.Importance,
			OptionID: req.OptionID,
			RoundID:  req.RoundID,
			TopicID:  req.TopicID,
		}
		condDB.Create(&cond)
	} else {
		condDB.Model(&cond).Update("score", req.Importance)
	}

	condDB.Model(&optionModel.Option{}).Where("id = ?", req.OptionID).Update("importance", req.Importance)

	response.Ok(cond).WriteTo(w)
}

func (h *ConditionHandler) SetNecessity(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		TopicID   uint `json:"topicId"`
		RoundID   uint `json:"roundId"`
		OptionID  uint `json:"optionId"`
		Necessity int  `json:"necessity"`
	}

	var req Req
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(http.StatusBadRequest, "invalid request body").WriteTo(w)
		return
	}

	// 更新或创建 necessity 条件
	var cond model.Condition
	query := condDB.Where("option_id = ? AND phase = ?", req.OptionID, "necessity")

	if req.RoundID > 0 {
		query = query.Where("round_id = ?", req.RoundID)
	}

	if err := query.First(&cond).Error; err != nil {
		cond = model.Condition{
			Name:     "必要性评分",
			Phase:    "necessity",
			Score:    req.Necessity,
			OptionID: req.OptionID,
			RoundID:  req.RoundID,
			TopicID:  req.TopicID,
		}
		condDB.Create(&cond)
	} else {
		condDB.Model(&cond).Update("score", req.Necessity)
	}

	condDB.Model(&optionModel.Option{}).Where("id = ?", req.OptionID).Update("necessity", req.Necessity)

	response.Ok(cond).WriteTo(w)
}

func (h *ConditionHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var cond model.Condition
	if err := condDB.First(&cond, vars["id"]).Error; err != nil {
		response.Error(http.StatusNotFound, "condition not found").WriteTo(w)
		return
	}

	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(http.StatusBadRequest, "invalid request body").WriteTo(w)
		return
	}

	condDB.Model(&cond).Updates(req)
	response.Ok(cond).WriteTo(w)
}

func (h *ConditionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var cond model.Condition
	if err := condDB.First(&cond, vars["id"]).Error; err != nil {
		response.Error(http.StatusNotFound, "condition not found").WriteTo(w)
		return
	}

	condDB.Delete(&cond)
	response.Ok(map[string]string{"id": vars["id"], "status": "deleted"}).WriteTo(w)
}
