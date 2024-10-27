package utils

import (
	"encoding/json"
)

type Response struct {
	Success bool
	Message string
	Data    any
	Error   any
}

func (r *Response) SuccessResponse(data any, message string) ([]byte, error) {
	if message == "" {
		r.Message = "Successfully completed the request"
	}
	r.Success = true
	r.Message = message
	r.Data = data
	r.Error = nil

	response, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (r *Response) ErrorResponse(errorJson any, message string) ([]byte, error) {
	if message == "" {
		r.Message = "Failed to complete the request"
	}
	r.Success = false
	r.Message = message
	r.Data = nil
	r.Error = errorJson
	response, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	return response, nil
}
