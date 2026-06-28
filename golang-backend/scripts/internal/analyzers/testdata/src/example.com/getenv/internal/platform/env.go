package platform

import "os"

func ReadEnv() string {
	return os.Getenv("APP_ENV") // want `raw os.Getenv usage is forbidden outside internal/config\.`
}
