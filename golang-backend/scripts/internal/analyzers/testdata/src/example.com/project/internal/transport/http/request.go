package http

import (
	"context"
	"encoding/json"
	"net/http"

	"example.com/project/internal/service"
)

type userCreator interface {
	CreateUser(context.Context, any) error
}

type CreateUserRequest struct {
	Email string `json:"email"`
}

func (r CreateUserRequest) ValidateAndSanitize() (service.CreateUserInput, error) {
	return service.CreateUserInput{Email: r.Email}, nil
}

type BadInput struct { // want `transport input struct "BadInput" must be named with the Request suffix\.`
	Name string
}

type MissingValidationRequest struct { // want `request "MissingValidationRequest" does not implement ValidateAndSanitize\(\) \(service.XInput, error\)`
	Name string
}

type WrongValidationRequest struct {
	Name string
}

func (r WrongValidationRequest) ValidateAndSanitize() (string, error) { // want `ValidateAndSanitize on "WrongValidationRequest" must return a type from internal/service whose name ends with Input\.`
	return r.Name, nil
}

func Handler(w http.ResponseWriter, r *http.Request) error {
	var req CreateUserRequest
	var svc userCreator
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil { // want `decoder "dec" must call DisallowUnknownFields\(\) before Decode\(\)`
		return err
	}
	_ = w
	return svc.CreateUser(context.Background(), req) // want `service call receives request "CreateUserRequest" directly\.`
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
	input, err := req.ValidateAndSanitize()
	if err != nil {
		return err
	}
	return service.CreateUser(context.Background(), input)
}
