package token

import (
	authu "gmc/auth/util"
	"gmc/db"
)

type DatabaseTokenAuth struct {
	name string
	db   db.DB
}

func New(cfg map[string]interface{}, db db.DB) (*DatabaseTokenAuth, error) {
	name, ok := cfg["name"].(string)
	if !ok {
		name = "token"
	}

	a := &DatabaseTokenAuth{
		name: name,
		db:   db,
	}
	return a, nil
}

func (a *DatabaseTokenAuth) Name() string {
	return a.name
}

func (a *DatabaseTokenAuth) Check(u string, p string) (*authu.User, error) {
	t, err := a.db.CheckToken(p)
	if err == nil {
		return &authu.User{Username: t.Description, Password: nil}, nil
	}
	return nil, nil
}
