package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/belliel/multiplexer/internal/services"
	"net/http"
)

const MaxUrls = 20

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	error   error
}

type UrlsToProcessRequest struct {
	Urls []string `json:"urls"`
}

type UrlsToProcessResponse struct {
	Urls map[string]interface{} `json:"urls"`
}

func Error(w http.ResponseWriter, r *http.Request) {
	errorResponse := &ErrorResponse{
		Status: http.StatusInternalServerError,
		error:  ErrInternalServerError,
	}

	if e, ok := r.Context().Value("error").(*ErrorResponse); ok {
		errorResponse = e
		if e.Status == 0 {
			e.Status = http.StatusInternalServerError
		}
	}

	w.WriteHeader(errorResponse.Status)
	errorResponse.Message = errorResponse.error.Error()

	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		errorMessage := fmt.Sprintf(
			"[%d] error: %s",
			errorResponse.Status,
			err.Error(),
		)

		w.Write([]byte(errorMessage))
	}
}

func ProcessUrls(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		r = addErrorToContext(r, &ErrorResponse{
			Status: http.StatusBadRequest,
			error:  ErrMethodIsNotPost,
		})
		Error(w, r)
		return
	}

	var urlsToProcessData UrlsToProcessRequest

	if err := json.NewDecoder(r.Body).Decode(&urlsToProcessData); err != nil {
		r = addErrorToContext(r, &ErrorResponse{
			error: err,
		})
		Error(w, r)
		return
	}

	if len(urlsToProcessData.Urls) > MaxUrls {
		r = addErrorToContext(r, &ErrorResponse{
			Status: http.StatusBadRequest,
			error:  ErrUrlsToProcessGreaterThanInt(MaxUrls),
		})
		Error(w, r)
		return
	}

	result, err := services.ProcessUrls(r.Context(), urlsToProcessData.Urls)

	if err != nil {
		errorMessage := ""
		status := http.StatusInternalServerError
		switch {
		case err == context.DeadlineExceeded:
			errorMessage = "request timeout"
			status = http.StatusRequestTimeout
		default:
			errorMessage = err.Error()
		}

		http.Error(w, errorMessage, status)
		return
	}

	if err := json.NewEncoder(w).Encode(&result); err != nil {
		r = addErrorToContext(r, &ErrorResponse{
			error: err,
		})
		Error(w, r)
		return
	}
}
