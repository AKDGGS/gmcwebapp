package auth

import (
	"fmt"

	"gmc/auth/file"
	"gmc/auth/token"
	authu "gmc/auth/util"
	"gmc/config"
	"gmc/db"
)

type Auth interface {
	Name() string

	// Uses this auth to check a username/password pair.
	// Authentication failures (such as bad usernames or passwords)
	// result in nil, nil.
	Check(string, string) (*authu.User, error)
}

func NewAuth(cfg config.AuthConfig, db db.DB) (Auth, error) {
	var auth Auth
	var err error
	switch cfg.Type {
	case "file":
		auth, err = file.New(cfg.Attrs)
		if err != nil {
			return nil, err
		}
	case "token":
		auth, err = token.New(cfg.Attrs, db)
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

type Auths struct {
	maxage int
	key    []byte
	auths  []Auth
}

func NewAuths(key []byte, maxage int, cfgs []config.AuthConfig, db db.DB) (*Auths, error) {
	auths := &Auths{
		key: key, maxage: maxage,
		auths: make([]Auth, len(cfgs)),
	}

	var err error
	for i, v := range cfgs {
		auths.auths[i], err = NewAuth(v, db)
		if err != nil {
			return nil, err
		}
	}

	return auths, nil
}

func (auths *Auths) Check(u string, p string) (*authu.User, error) {
	for _, auth := range auths.auths {
		user, err := auth.Check(u, p)
		if err != nil {
			return nil, err
		}
		if user != nil {
			return user, nil
		}
	}
	return nil, nil
}
