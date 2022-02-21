package auth

import (
	"fmt"
	"gmc/auth/always"
	authu "gmc/auth/util"
	"gmc/config"
	"net/http"
)

type Auth interface {
	// Authentication is optional - accept authentication cookies
	// and return user if available
	Optional(http.ResponseWriter, *http.Request) (*authu.User, error)

	// Authentication is required - hard stop and prompt for
	// username/password if authentication cookie is not present
	Required(http.ResponseWriter, *http.Request) (*authu.User, error)
}

func NewAuth(cfg config.AuthConfig) (Auth, error) {
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

type Auths []Auth

func NewAuths(cfgs []config.AuthConfig) (Auths, error) {
	var err error
	auths := make(Auths, len(cfgs))
	for i, v := range cfgs {
		auths[i], err = NewAuth(v)
		if err != nil {
			return nil, err
		}
	}
	return auths, nil
}

func (auths Auths) Required(w http.ResponseWriter, r *http.Request) (*authu.User, error) {
	for _, auth := range auths {
		user, err := auth.Required(w, r)
		if err != nil {
			return nil, err
		}
		if user != nil {
			return user, nil
		}
	}
	return nil, nil
}

func (auths Auths) Optional(w http.ResponseWriter, r *http.Request) (*authu.User, error) {
	for _, auth := range auths {
		user, err := auth.Optional(w, r)
		if err != nil {
			return nil, err
		}
		if user != nil {
			return user, nil
		}
	}
	return nil, nil
}
