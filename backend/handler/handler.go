package handler

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	//"log"
	"net/http"
)

type JSON map[string]interface{}

type Handler struct {
}

func NewHandler() Handler {
	return Handler{}
}

func (h Handler) Check(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(JSON{"status": "ok"})
}
