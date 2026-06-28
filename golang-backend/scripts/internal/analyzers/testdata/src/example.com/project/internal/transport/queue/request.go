package queue

import "example.com/project/internal/service"

type PublishJobRequest struct {
	Topic string
}

func (r PublishJobRequest) ValidateAndSanitize() (service.CreateUserInput, error) {
	return service.CreateUserInput{Email: r.Topic}, nil
}
