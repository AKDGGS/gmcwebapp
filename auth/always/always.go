package always

import (
	"fmt"
	"gmc/auth/user"
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

func (al *Always) AuthOptional(w http.ResponseWriter, r *http.Request) (*user.User, error) {
	return &user.User{Username: al.user}, nil
}

func (al *Always) AuthRequired(w http.ResponseWriter, r *http.Request) (*user.User, error) {
	return &user.User{Username: al.user}, nil
}
