package http

import "example.com/requests/internal/service"

type CreateUserRequest struct {
	Email string
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
