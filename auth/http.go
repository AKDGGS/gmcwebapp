package auth

import (
	"gmc/auth/securecookie"
	authu "gmc/auth/util"
	"net/http"
)

func (auths *Auths) CheckRequest(r *http.Request) (*authu.User, error) {
	// Try to authenticate the user with a secure cookie
	cookie, err := securecookie.New(
		"session", auths.key,
		securecookie.Params{MaxAge: 86400, Secure: false},
	)
	if err != nil {
		return nil, err
	}

	uj, err := cookie.GetValue(nil, r)
	if err == nil {
		user, err := authu.UnmarshalUser(uj)
		if err != nil {
			return nil, err
		}

		if user != nil {
			return user, nil
		}
	}

	return nil, nil
}
