package always

import (
	"fmt"
	authu "gmc/auth/util"
	"net/http"
)

type Always struct {
	user string
}

func New(cfg map[string]interface{}) (*Always, error) {
	user, ok := cfg["user"].(string)
	if !ok {
		return nil, fmt.Errorf("auth always user must exist and be a string")
	}

	return &Always{user: user}, nil
}

func (al *Always) Optional(w http.ResponseWriter, r *http.Request) (*authu.User, error) {
	return &authu.User{Username: al.user}, nil
}

func (al *Always) Required(w http.ResponseWriter, r *http.Request) (*authu.User, error) {
	return &authu.User{Username: al.user}, nil
}
