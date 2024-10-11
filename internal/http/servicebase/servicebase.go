package servicebase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ServiceBase struct {
	BaseUrl string
}

type httpRequestMethod string

const (
	MethodGet     httpRequestMethod = "GET"
	MethodHead    httpRequestMethod = "HEAD"
	MethodPost    httpRequestMethod = "POST"
	MethodPut     httpRequestMethod = "PUT"
	MethodPatch   httpRequestMethod = "PATCH"
	MethodDelete  httpRequestMethod = "DELETE"
	MethodConnect httpRequestMethod = "CONNECT"
	MethodOptions httpRequestMethod = "OPTIONS"
	MethodTrace   httpRequestMethod = "TRACE"
)

func (s *ServiceBase) GET(route string) ([]byte, error) {
	return []byte{}, nil
}

func (s *ServiceBase) PUT_POST_PATCH(
	method httpRequestMethod,
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
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	jsonResponse, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %s\n", err)
		return nil, err
	}

	return jsonResponse, nil
}
