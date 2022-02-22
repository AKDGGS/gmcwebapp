package auth

import (
	"bytes"
	"fmt"
	"gmc/assets"
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

func (auths *Auths) CheckForm(w http.ResponseWriter, r *http.Request) error {
	// Try to authenticate the user with a secure cookie
	cookie, err := securecookie.New(
		"session", auths.key,
		securecookie.Params{MaxAge: 86400, Secure: false},
	)
	if err != nil {
		return err
	}

	uj, err := cookie.GetValue(nil, r)
	if err == nil {
		user, err := authu.UnmarshalUser(uj)
		if err != nil {
			return err
		}

		if user != nil {
			http.Redirect(w, r, ".", http.StatusFound)
			return nil
		}
	}

	err = r.ParseForm()
	if err != nil {
		return err
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if username != "" && password != "" {
		user, err := auths.Check(username, password)
		if err != nil {
			return err
		}

		if user != nil {
			uj, err := authu.MarshalUser(user)
			if err != nil {
				return err
			}

			err = cookie.SetValue(w, uj)
			if err != nil {
				return err
			}

			http.Redirect(w, r, ".", http.StatusFound)
			return nil
		}
	}

	buf := bytes.Buffer{}
	if err := assets.ExecuteTemplate("tmpl/login.html", &buf, nil); err != nil {
		http.Error(
			w, fmt.Sprintf("parse error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return nil
	}

	w.Header().Set("Content-Length", fmt.Sprintf("%d", buf.Len()))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(buf.Bytes())
	return nil
}
