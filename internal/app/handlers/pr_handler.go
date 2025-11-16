package handlers

import (
	"avito-backend-trainee-assignment-autumn-2025/internal/app/handlers/response"
	"avito-backend-trainee-assignment-autumn-2025/internal/domain/delivery"
	"avito-backend-trainee-assignment-autumn-2025/internal/domain/service"
	"encoding/json"
	"io"
	"net/http"
)

type PullRequestHandler struct {
	s service.PullRequestServicer
}

func NewPullRequestHandler(s service.PullRequestServicer) *PullRequestHandler {
	return &PullRequestHandler{
		s: s,
	}
}

func (h *PullRequestHandler) CreatePullRequestHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		response.MethodNotAllowed(w, response.MethodIsNotAllowed)
		return
	}
	var req delivery.CreatePullRequestRequest
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
		response.BadRequest(w, response.JsonParseError, "json unmarshal error")
		return
	}
	newPullRequest, err := h.s.CreatePullRequest(r.Context(), req.PullRequestId, req.PullRequestName, req.AuthorId)
	if err != nil {
		if err.Error() == "invalid request" {
			response.BadRequest(w, response.InvalidRequest, "invalid request")
			return
		} else if err.Error() == "internal error" {
			response.InternalServerError(w, response.InternalError, "internal error")
			return
		} else if err.Error() == "resource not found" {
			response.NotFound(w, response.NotFoundError, "resource not found")
			return
		}
		response.Conflict(w, response.PrExists, "PR id already exists")
		return
	}
	prResponse := &delivery.CreatePullRequestResponse{}
	prResponse.Pr.PullRequestId = newPullRequest.Id
	prResponse.Pr.PullRequestName = newPullRequest.Name
	prResponse.Pr.AuthorId = newPullRequest.Author
	prResponse.Pr.Status = newPullRequest.Status
	prResponse.Pr.AssignedReviewers = newPullRequest.Reviewers
	response.Success(w, http.StatusCreated, prResponse)
}

func (h *PullRequestHandler) MergePullRequestHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		response.MethodNotAllowed(w, response.MethodIsNotAllowed)
		return
	}
	var req delivery.MergePullRequestRequest
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
		response.BadRequest(w, response.JsonParseError, "json unmarshal error")
		return
	}
	updatedPullRequest, err := h.s.MergePullRequest(r.Context(), req.PullRequestId)
	if err != nil {
		if err.Error() == "invalid request" {
			response.BadRequest(w, response.InvalidRequest, "invalid request")
			return
		} else if err.Error() == "internal error" {
			response.InternalServerError(w, response.InternalError, "internal error")
			return
		}
		response.NotFound(w, response.NotFoundError, "resource not found")
		return
	}
	prResponse := &delivery.MergePullRequestResponse{}
	prResponse.Pr.PullRequestId = updatedPullRequest.Id
	prResponse.Pr.PullRequestName = updatedPullRequest.Name
	prResponse.Pr.AuthorId = updatedPullRequest.Author
	prResponse.Pr.AssignedReviewers = updatedPullRequest.Reviewers
	prResponse.Pr.Status = updatedPullRequest.Status
	prResponse.Pr.MergedAt = updatedPullRequest.MergedAt
	response.Success(w, http.StatusOK, prResponse)
}

func (h *PullRequestHandler) ReassignPullRequestHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		response.MethodNotAllowed(w, response.MethodIsNotAllowed)
		return
	}
	var req delivery.ReassignPullRequestRequest
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
		response.BadRequest(w, response.JsonParseError, "json unmarshal error")
		return
	}
	updatedPullRequest, newReviewId, err := h.s.ReassignPullRequest(r.Context(), req.PullRequestId, req.OldReviewerId)
	if err != nil {
		if err.Error() == "invalid request" {
			response.BadRequest(w, response.InvalidRequest, "invalid request")
			return
		} else if err.Error() == "internal error" {
			response.InternalServerError(w, response.InternalError, "internal error")
			return
		} else if err.Error() == "no active replacement candidate in team" {
			response.Conflict(w, response.NoCandidate, "no active replacement candidate in team")
			return
		} else if err.Error() == "cannot reassign on merged PR" {
			response.Conflict(w, response.PrMerged, "cannot reassign on merged PR")
			return
		} else if err.Error() == "reviewer is not assigned to this PR" {
			response.Conflict(w, response.NotAssigned, "reviewer is not assigned to this PR")
			return
		}
		response.NotFound(w, response.NotFoundError, "resource not found")
		return
	}
	prResponse := &delivery.ReassignPullRequestResponse{}
	prResponse.Pr.PullRequestId = updatedPullRequest.Id
	prResponse.Pr.PullRequestName = updatedPullRequest.Name
	prResponse.Pr.AuthorId = updatedPullRequest.Author
	prResponse.Pr.AssignedReviewers = updatedPullRequest.Reviewers
	prResponse.Pr.Status = updatedPullRequest.Status
	prResponse.ReplacedBy = newReviewId
	response.Success(w, http.StatusOK, prResponse)
}
