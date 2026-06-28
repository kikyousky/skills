package domain

import _ "example.com/boundaries/internal/service" // want `domain package must not import "example.com/boundaries/internal/service"\.`

type User struct{}
