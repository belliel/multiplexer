package api

import (
	"context"
	"net/http"
)

func addErrorToContext(r *http.Request, value interface{}) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), "error", value))
}
