package api

import (
	"errors"
	"strconv"
)

var (
	ErrInternalServerError         = errors.New("internal server error")
	ErrMethodIsNotPost             = errors.New("request method must be \"port\"")
	ErrUrlsToProcessGreaterThanInt = func(count int) error {
		return errors.New("request urls greater than " + strconv.Itoa(count))
	}
)
