package response

import (
	"encoding/json"
	"net/http"
)

type Code string

const (
	TeamExists         Code = "TEAM_EXISTS"
	PrExists           Code = "PR_EXISTS"
	PrMerged           Code = "PR_MERGED"
	MethodIsNotAllowed Code = "METHOD_NOT_ALLOWED"
	BodyReadError      Code = "BODY_READ_ERROR"
	JsonParseError     Code = "JSON_PARSE_ERROR"
	InternalError      Code = "INTERNAL_ERROR"
	NotFoundError      Code = "NOT_FOUND"
	InvalidRequest     Code = "INVALID_REQUEST"
	NoCandidate        Code = "NO_CANDIDATE"
	NotAssigned        Code = "NOT_ASSIGNED"
)

type ErrorResponse struct {
	Error struct {
		Code    Code   `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func JsonResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			return
		}
	}
}

func Error(w http.ResponseWriter, statusCode int, code Code, message string) {
	err := &ErrorResponse{}
	err.Error.Code = code
	err.Error.Message = message
	JsonResponse(w, statusCode, err)
}

func Success(w http.ResponseWriter, statusCode int, data interface{}) {
	JsonResponse(w, statusCode, data)
}

func MethodNotAllowed(w http.ResponseWriter, code Code) {
	Error(w, http.StatusMethodNotAllowed, code, "method not allowed")
}

func BadRequest(w http.ResponseWriter, code Code, message string) {
	Error(w, http.StatusBadRequest, code, message)
}

func InternalServerError(w http.ResponseWriter, code Code, message string) {
	Error(w, http.StatusInternalServerError, code, message)
}

func NotFound(w http.ResponseWriter, code Code, message string) {
	Error(w, http.StatusNotFound, code, message)
}

func Conflict(w http.ResponseWriter, code Code, message string) {
	Error(w, http.StatusConflict, code, message)
}
