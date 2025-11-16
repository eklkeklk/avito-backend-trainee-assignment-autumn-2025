package handlers

import (
	"avito-backend-trainee-assignment-autumn-2025/internal/app/handlers/response"
	"avito-backend-trainee-assignment-autumn-2025/internal/domain/delivery"
	"avito-backend-trainee-assignment-autumn-2025/internal/domain/models/team"
	"avito-backend-trainee-assignment-autumn-2025/internal/domain/service"
	"encoding/json"
	"io"
	"net/http"
)

type TeamHandler struct {
	s service.TeamServicer
}

func NewTeamHandler(s service.TeamServicer) *TeamHandler {
	return &TeamHandler{
		s: s,
	}
}

func (h *TeamHandler) CreateTeamHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		response.MethodNotAllowed(w, response.MethodIsNotAllowed)
		return
	}
	var req team.Team
	body, err := io.ReadAll(r.Body)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(r.Body)
	if err != nil {
		response.BadRequest(w, response.BodyReadError, "body read error")
		return
	}
	if err := json.Unmarshal(body, &req); err != nil {
		response.BadRequest(w, response.JsonParseError, "invalid json")
		return
	}
	newTeam, err := h.s.CreateTeam(r.Context(), req.Name, req.Members)
	if err != nil {
		if err.Error() == "internal error" {
			response.InternalServerError(w, response.InternalError, err.Error())
			return
		} else if err.Error() == "invalid request" {
			response.BadRequest(w, response.InvalidRequest, "invalid request")
			return
		}
		response.BadRequest(w, response.TeamExists, "team_name already exists")
		return
	}
	createResponse := delivery.CreateTeamResponse{}
	createResponse.Team.Name = newTeam.Name
	createResponse.Team.Members = newTeam.Members
	response.Success(w, http.StatusCreated, createResponse)
}

func (h *TeamHandler) GetTeamHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodGet {
		response.MethodNotAllowed(w, response.MethodIsNotAllowed)
		return
	}
	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		response.NotFound(w, response.NotFoundError, "resource not found")
		return
	}
	newTeam, err := h.s.GetTeam(r.Context(), teamName)
	if err != nil {
		if err.Error() == "internal error" {
			response.InternalServerError(w, response.InternalError, err.Error())
			return
		} else if err.Error() == "invalid request" {
			response.BadRequest(w, response.InvalidRequest, "invalid request")
			return
		}
		response.NotFound(w, response.NotFoundError, "resource not found")
		return
	}
	response.Success(w, http.StatusOK, newTeam)
}
