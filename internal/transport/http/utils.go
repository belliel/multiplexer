package http

import "net/http"

type StatusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *StatusRecorder) WriteHeader(status int) {
	s.ResponseWriter.WriteHeader(status)
	s.status = status
}
