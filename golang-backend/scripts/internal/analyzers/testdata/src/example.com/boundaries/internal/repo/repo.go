package repo

import (
	_ "example.com/boundaries/internal/service"        // want `repo package must not import "example.com/boundaries/internal/service"\.`
	_ "example.com/boundaries/internal/transport/http" // want `repo package must not import "example.com/boundaries/internal/transport/http"\.`
)

type UserRepository struct{}
