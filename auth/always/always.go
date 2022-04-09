package always

import (
	"fmt"

	authu "gmc/auth/util"
)

type Always struct {
	name string
	user string
}

func New(cfg map[string]interface{}) (*Always, error) {
	name, ok := cfg["name"].(string)
	if !ok {
		name = "always"
	}

	user, ok := cfg["user"].(string)
	if !ok {
		return nil, fmt.Errorf("auth always user must exist and be a string")
	}

	return &Always{name: name, user: user}, nil
}

func (al *Always) Name() string {
	return al.name
}

func (al *Always) Check(user string, pass string) (*authu.User, error) {
	return &authu.User{Username: al.user}, nil
}
