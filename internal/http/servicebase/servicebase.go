package servicebase

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type ServiceBase struct {
	BaseUrl string
}

type HttpRequestMethod string

const (
	MethodGet     HttpRequestMethod = "GET"
	MethodHead    HttpRequestMethod = "HEAD"
	MethodPost    HttpRequestMethod = "POST"
	MethodPut     HttpRequestMethod = "PUT"
	MethodPatch   HttpRequestMethod = "PATCH"
	MethodDelete  HttpRequestMethod = "DELETE"
	MethodConnect HttpRequestMethod = "CONNECT"
	MethodOptions HttpRequestMethod = "OPTIONS"
	MethodTrace   HttpRequestMethod = "TRACE"
)

func (s *ServiceBase) GET(route string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, s.BaseUrl+route, nil)
	if err != nil {
		return nil, errors.New("unable to get resource")
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	jsonResponse, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %s\n", err)
		return nil, err
	}

	return jsonResponse, nil
}

func (s *ServiceBase) PUT_POST_PATCH(
	method HttpRequestMethod,
	route string,
	body any) ([]byte, error) {

	jsonPayload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		string(method),
		s.BaseUrl+route,
		bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, errors.New("error occurred while making request")
	}
	defer response.Body.Close()

	jsonResponse, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %s\n", err)
		return nil, err
	}

	return jsonResponse, nil
}
