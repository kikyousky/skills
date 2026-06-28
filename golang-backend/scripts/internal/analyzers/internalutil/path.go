package internalutil

import "strings"

func IsTransportPath(path string) bool {
	return strings.Contains(path, "/internal/transport/") || strings.HasSuffix(path, "/internal/transport")
}

func IsServicePath(path string) bool {
	return strings.Contains(path, "/internal/service/") || strings.HasSuffix(path, "/internal/service")
}

func IsDomainPath(path string) bool {
	return strings.Contains(path, "/internal/domain/") || strings.HasSuffix(path, "/internal/domain")
}

func IsRepoPath(path string) bool {
	return strings.Contains(path, "/internal/repo/") || strings.HasSuffix(path, "/internal/repo")
}

func IsConfigPath(path string) bool {
	return strings.Contains(path, "/internal/config/") || strings.HasSuffix(path, "/internal/config")
}

func IsTransportFrameworkImport(path string) bool {
	switch {
	case path == "net/http":
		return true
	case strings.HasPrefix(path, "google.golang.org/grpc"):
		return true
	case strings.Contains(path, "websocket"):
		return true
	default:
		return false
	}
}

func HasForbiddenSuffix(name string) bool {
	for _, suffix := range []string{"Manager", "Helper", "Util", "Data"} {
		if strings.HasSuffix(name, suffix) {
			return true
		}
	}
	return false
}
