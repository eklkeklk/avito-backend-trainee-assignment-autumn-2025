package handlers

import (
	"avito-backend-trainee-assignment-autumn-2025/internal/app/handlers/response"
	"avito-backend-trainee-assignment-autumn-2025/internal/domain/delivery"
	"avito-backend-trainee-assignment-autumn-2025/internal/domain/service"
	"encoding/json"
	"io"
	"net/http"
)

type UserHandler struct {
	s service.UserServicer
}

func NewUserHandler(s service.UserServicer) *UserHandler {
	return &UserHandler{
		s: s,
	}
}

func (h *UserHandler) SetIsActiveHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		response.MethodNotAllowed(w, response.MethodIsNotAllowed)
		return
	}
	var req delivery.SetActiveStatusRequest
	body, err := io.ReadAll(r.Body)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(r.Body)
	if err != nil {
		response.BadRequest(w, response.BodyReadError, "read body error")
		return
	}
	err = json.Unmarshal(body, &req)
	if err != nil {
		response.BadRequest(w, response.InvalidRequest, "json unmarshal error")
	}
	user, err := h.s.SetActiveStatus(r.Context(), req.UserId, req.IsActive)
	if err != nil {
		if err.Error() == "internal error" {
			response.InternalServerError(w, response.InternalError, "internal error")
			return
		} else if err.Error() == "invalid request" {
			response.BadRequest(w, response.InvalidRequest, "invalid request")
			return
		}
		response.NotFound(w, response.NotFoundError, "resource not found")
		return
	}
	setResponse := delivery.SetActiveStatusResponse{}
	setResponse.User.Id = user.Id
	setResponse.User.Name = user.Name
	setResponse.User.TeamName = user.TeamName
	setResponse.User.IsActive = user.IsActive
	response.Success(w, http.StatusOK, setResponse)
}

func (h *UserHandler) GetReviewHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodGet {
		response.MethodNotAllowed(w, response.MethodIsNotAllowed)
		return
	}
	userId := r.URL.Query().Get("user_id")
	userPullRequests, err := h.s.GetReviewList(r.Context(), userId)
	if err != nil {
		if err.Error() == "internal error" {
			response.InternalServerError(w, response.InternalError, "internal error")
			return
		} else if err.Error() == "invalid request" {
			response.BadRequest(w, response.InvalidRequest, "invalid request")
			return
		}
		response.NotFound(w, response.NotFoundError, "resource not found")
		return
	}
	response.Success(w, http.StatusOK, userPullRequests)
}
