package service

import "context"

type CreateUserInput struct {
	Email string
}

func Good(ctx context.Context, input CreateUserInput) error {
	_ = ctx
	_ = input
	return nil
}

func BadService(ctx context.Context, req struct{ Email string }) error { // want `service function "BadService" accepts forbidden parameter type struct\{Email string\}`
	_ = ctx
	_ = req
	return nil
}
