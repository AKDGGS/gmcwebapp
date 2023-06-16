package token

import (
	"fmt"
	"os"

	authu "gmc/auth/util"
	"gmc/db"
)

type TokenAuth struct {
	name  string
	path  string
	users map[string]*authu.User
}

func New(cfg map[string]interface{}) (*TokenAuth, error) {
	name, ok := cfg["name"].(string)
	if !ok {
		name = "token"
	}
	path := cfg["path"].(string)

	a := &TokenAuth{
		name:  name,
		path:  path,
		users: make(map[string]*authu.User),
	}
	return a, nil
}

func (a *TokenAuth) Name() string {
	return a.name
}

func (a *TokenAuth) Check(u string, p string) (*authu.User, error) {
	db, err := db.New(a.path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
	t, err := db.CheckToken(p)
	if err == nil {
		return &authu.User{Username: t.Description, Password: []byte(p)}, nil
	}
	return nil, nil
}
