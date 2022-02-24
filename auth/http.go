package auth

import (
	"bytes"
	"fmt"
	"gmc/assets"
	"gmc/auth/securecookie"
	authu "gmc/auth/util"
	"net/http"
)

func (auths *Auths) Logout(w http.ResponseWriter, r *http.Request) error {
	cookie, err := securecookie.New(
		"session", auths.key,
		securecookie.Params{MaxAge: auths.maxage, Secure: false},
	)
	if err != nil {
		return err
	}

	cookie.Delete(w)
	http.Redirect(w, r, ".", http.StatusFound)
	return nil
}

func (auths *Auths) CheckRequest(w http.ResponseWriter, r *http.Request) (*authu.User, error) {
	// Try to authenticate the user with a secure cookie
	cookie, err := securecookie.New(
		"session", auths.key,
		securecookie.Params{MaxAge: auths.maxage, Secure: false},
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
			// Refresh the cookie, extending the timeout
			// and ignore any failures.
			cookie.SetValue(w, uj)
			return user, nil
		}
	}

	return nil, nil
}

func (auths *Auths) CheckForm(w http.ResponseWriter, r *http.Request) error {
	// Try to authenticate the user with a secure cookie
	cookie, err := securecookie.New(
		"session", auths.key,
		securecookie.Params{MaxAge: auths.maxage, Secure: false},
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
			// Refresh the cookie, extending the timeout
			// and ignore any failures.
			cookie.SetValue(w, uj)

			http.Redirect(w, r, ".", http.StatusFound)
			return nil
		}
	}

	// If the user can't be authenticate with a secure cookie,
	// try to read a POSTed username and password to authenticate with
	err = r.ParseForm()
	if err != nil {
		return err
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	params := map[string]interface{}{}

	if username != "" && password != "" {
		params["username"] = username

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

		params["error"] = "Invalid username or password."
	}

	// If there's no secure cookie, and no POSTed credentials,
	// serve up the login page
	buf := bytes.Buffer{}
	if err := assets.ExecuteTemplate("tmpl/login.html", &buf, params); err != nil {
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
