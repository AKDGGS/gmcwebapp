package auth

import (
	"fmt"
	"gmc/auth/always"
	"gmc/auth/user"
	"gmc/config"
	"net/http"
)

type Auth interface {
	// Authentication is optional - accept authentication cookies
	// and return user if available
	AuthOptional(http.ResponseWriter, *http.Request) (*user.User, error)

	// Authentication is required - hard stop and prompt for
	// username/password if authentication cookie is not present
	AuthRequired(http.ResponseWriter, *http.Request) (*user.User, error)
}

func New(cfg config.AuthConfig) (Auth, error) {
	var auth Auth
	var err error
	switch cfg.Type {
	case "always":
		auth, err = always.New(cfg.Attrs)
		if err != nil {
			return nil, err
		}
	case "":
		return nil, fmt.Errorf("authentication type cannot be empty")
	default:
		return nil, fmt.Errorf("unknown authentication type: %s", cfg.Type)
	}
	return auth, nil
}
