package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"specialdates-backend/internal/models"
	"specialdates-backend/internal/service"

	"github.com/gorilla/mux"
)

type Handler struct {
	Svc service.DateService
}

func NewHandler(s service.DateService) *Handler {
	return &Handler{Svc: s}
}

func (h *Handler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/dates", h.ListDates).Methods("GET")
	r.HandleFunc("/dates/{id}", h.GetDate).Methods("GET")
	r.HandleFunc("/dates", h.CreateDate).Methods("POST")
	r.HandleFunc("/dates/{id}", h.UpdateDate).Methods("PUT")
	r.HandleFunc("/dates/{id}", h.DeleteDate).Methods("DELETE")
}

// request payload (incoming)
type datePayload struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	EventAt     string `json:"event_at"` // ISO string
	DateType    string `json:"date_type"`
	Recurring   bool   `json:"recurring"`
}

func (h *Handler) ListDates(w http.ResponseWriter, r *http.Request) {
	dates, err := h.Svc.ListDates()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, dates)
}

func (h *Handler) GetDate(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idStr)
	d, err := h.Svc.GetDate(id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	writeJSON(w, d)
}

func (h *Handler) CreateDate(w http.ResponseWriter, r *http.Request) {
	var p datePayload
	body, _ := io.ReadAll(r.Body)
	if err := json.Unmarshal(body, &p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// parse time from client (expects RFC3339)
	t, err := time.Parse(time.RFC3339, p.EventAt)
	if err != nil {
		http.Error(w, "invalid event_at format (expect RFC3339)", http.StatusBadRequest)
		return
	}
	d := &models.DateEvent{
		Title:       p.Title,
		Description: p.Description,
		EventAt:     t,
		DateType:    p.DateType,
		Recurring:   p.Recurring,
	}
	id, err := h.Svc.CreateDate(d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	d.ID = id
	writeJSON(w, d)
}

func (h *Handler) UpdateDate(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idStr)
	var p datePayload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	t, err := time.Parse(time.RFC3339, p.EventAt)
	if err != nil {
		http.Error(w, "invalid event_at format (expect RFC3339)", http.StatusBadRequest)
		return
	}
	d := &models.DateEvent{
		ID:          id,
		Title:       p.Title,
		Description: p.Description,
		EventAt:     t,
		DateType:    p.DateType,
		Recurring:   p.Recurring,
	}
	if err := h.Svc.UpdateDate(d); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteDate(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idStr)
	if err := h.Svc.DeleteDate(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// helper to write json response
func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.Encode(v)
}
