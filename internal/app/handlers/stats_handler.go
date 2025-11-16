package handlers

import (
	"avito-backend-trainee-assignment-autumn-2025/internal/app/handlers/response"
	"avito-backend-trainee-assignment-autumn-2025/internal/domain/service"
	"net/http"
)

type StatsHandler struct {
	s service.StatsServicer
}

func NewStatsHandler(s service.StatsServicer) *StatsHandler {
	return &StatsHandler{
		s: s,
	}
}

func (h StatsHandler) UserStatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.MethodNotAllowed(w, response.MethodIsNotAllowed)
		return
	}
	stats, err := h.s.GetUserStats(r.Context())
	if err != nil {
		response.InternalServerError(w, response.InternalError, "internal error")
		return
	}
	response.Success(w, http.StatusOK, stats)
}

func (h StatsHandler) PrStatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.MethodNotAllowed(w, response.MethodIsNotAllowed)
		return
	}
	stats, err := h.s.GetPrStats(r.Context())
	if err != nil {
		response.InternalServerError(w, response.InternalError, "internal error")
		return
	}
	response.Success(w, http.StatusOK, stats)
}

func (h StatsHandler) OverallStatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.MethodNotAllowed(w, response.MethodIsNotAllowed)
		return
	}
	stats, err := h.s.GetOverallStats(r.Context())
	if err != nil {
		response.InternalServerError(w, response.InternalError, "internal error")
		return
	}
	response.Success(w, http.StatusOK, stats)
}
