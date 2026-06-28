package domain

import _ "example.com/project/internal/service" // want `domain package must not import "example.com/project/internal/service"\.`

type User struct{}
