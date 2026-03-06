package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("Authorization header does not exist")
	}

	const prefix = "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		return "", fmt.Errorf("invalid authorization header format")
	}

	token := strings.TrimPrefix(authHeader, prefix)
	token = strings.TrimSpace(token)

	if token == "" {
		return "", fmt.Errorf("no token found in authorization header")
	}

	return token, nil
}

