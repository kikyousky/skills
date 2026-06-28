package service

import "context"

type CreateUserInput struct {
	Email string
}

type Clock interface {
	Now() int64
}

type BadManager struct{} // want `type "BadManager" uses forbidden suffix\.`

func CreateUser(ctx context.Context, input CreateUserInput) error {
	_ = ctx
	_ = input
	return nil
}

func BadService(ctx context.Context, req struct{ Email string }) error { // want `service function "BadService" accepts forbidden parameter type struct\{ Email string \}`
	_ = ctx
	_ = req
	return nil
}

func TooManyParams(ctx context.Context, a, b, c, d, e int) int { // want `service function "TooManyParams" accepts forbidden parameter type int\.` // want `function "TooManyParams" has 5 parameters excluding context.Context\.`
	return a + b + c + d + e
}
