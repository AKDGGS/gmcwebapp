package always

import (
	authu "gmc/auth/util"
)

type Always struct {
	name  string
	allow bool
}

func New(cfg map[string]interface{}) (*Always, error) {
	name, ok := cfg["name"].(string)
	if !ok {
		name = "always"
	}

	allow, ok := cfg["allow"].(bool)
	if !ok {
		allow = false
	}

	return &Always{
		name:  name,
		allow: allow,
	}, nil
}

func (al *Always) Name() string {
	return al.name
}

func (al *Always) Check(u string, p string) (*authu.User, error) {
	if al.allow {
		return &authu.User{Username: u}, nil
	}
	return nil, nil
}
