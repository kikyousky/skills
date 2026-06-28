package http

import (
	"context"

	"example.com/transportflow/internal/service"
)

type CreateUserRequest struct {
	Email string
}

type UserService interface {
	CreateUser(context.Context, any, service.CreateUserInput) error
}

func Handler(svc UserService, req CreateUserRequest) error {
	return svc.CreateUser(context.Background(), req, service.CreateUserInput{}) // want `service call receives request "CreateUserRequest" directly\.`
}
