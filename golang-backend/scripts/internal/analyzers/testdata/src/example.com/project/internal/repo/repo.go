package repo

import (
	_ "example.com/project/internal/service"        // want `repo package must not import "example.com/project/internal/service"\.`
	_ "example.com/project/internal/transport/http" // want `repo package must not import "example.com/project/internal/transport/http"\.`
)

type UserRepository struct{}
