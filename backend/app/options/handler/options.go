package handler

import (
	"net/http"

	"mdecide/app/options/model"
	"mdecide/common/response"

	"github.com/gorilla/mux"
)

type OptionsHandler struct{}

func NewOptionsHandler() *OptionsHandler {
	return &OptionsHandler{}
}

func RegisterOptionsRoutes(r *mux.Router) {
	h := NewOptionsHandler()
	r.HandleFunc("/api/options", h.List).Methods("GET")
	r.HandleFunc("/api/options", h.Create).Methods("POST")
	r.HandleFunc("/api/options/{id}", h.Get).Methods("GET")
	r.HandleFunc("/api/options/{id}", h.Update).Methods("PUT")
	r.HandleFunc("/api/options/{id}", h.Delete).Methods("DELETE")
	r.HandleFunc("/api/options/batch", h.BatchCreate).Methods("POST")
}

func (h *OptionsHandler) List(w http.ResponseWriter, r *http.Request) {
	var options []model.Option
	response.Ok(options).WriteTo(w)
}

func (h *OptionsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var opt model.Option
	response.Ok(opt).WriteTo(w)
}

func (h *OptionsHandler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	response.Ok(map[string]string{"id": id}).WriteTo(w)
}

func (h *OptionsHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	response.Ok(map[string]string{"id": id}).WriteTo(w)
}

func (h *OptionsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	response.Ok(map[string]string{"id": id, "status": "deleted"}).WriteTo(w)
}

func (h *OptionsHandler) BatchCreate(w http.ResponseWriter, r *http.Request) {
	var opts []model.Option
	response.Ok(opts).WriteTo(w)
}
