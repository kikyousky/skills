package http

import (
	"encoding/json"
	"net/http"
)

type CreateUserRequest struct {
	Email string `json:"email"`
}

func Handler(r *http.Request) error {
	var req CreateUserRequest
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil { // want `decoder "dec" must call DisallowUnknownFields\(\) before Decode\(\)`
		return err
	}
	return nil
}

func HandlerInline(r *http.Request) error {
	var req CreateUserRequest
	return json.NewDecoder(r.Body).Decode(&req) // want `inline json.NewDecoder\(\.\.\.\)\.Decode\(\.\.\.\) is forbidden`
}

func HandlerUnmarshal(body []byte) error {
	var req CreateUserRequest
	return json.Unmarshal(body, &req) // want `json.Unmarshal is forbidden for HTTP boundary decoding\.`
}

func GoodHandler(r *http.Request) error {
	var req CreateUserRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		return err
	}
	return nil
}
