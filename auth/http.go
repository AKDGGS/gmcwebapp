package auth

import (
	"bytes"
	"fmt"
	"net/http"
	"regexp"

	"gmc/assets"
	"gmc/auth/securecookie"
	authu "gmc/auth/util"
)

var redirect_rx *regexp.Regexp = regexp.MustCompile(`^(?:\.|borehole|inventory|outcrop|prospect|shotline|well|wells|qa)\/?\d*$`)

func (auths *Auths) Logout(w http.ResponseWriter, r *http.Request) error {
	cookie, err := securecookie.New(
		"gmc-session", auths.key,
		securecookie.Params{
			MaxAge: auths.maxage, Secure: false,
			Path: "/", SameSite: securecookie.Lax,
		},
	)
	if err != nil {
		return err
	}

	cookie.Delete(w)

	redirect := r.URL.Query().Get("redirect")
	if redirect == "" {
		redirect = "."
	}
	if !redirect_rx.MatchString(redirect) {
		return fmt.Errorf("invalid redirect")
	}

	http.Redirect(w, r, redirect, http.StatusFound)
	return nil
}

func (auths *Auths) CheckRequest(w http.ResponseWriter, r *http.Request) (*authu.User, error) {
	// Try to authenticate the user with a secure cookie
	cookie, err := securecookie.New(
		"gmc-session", auths.key,
		securecookie.Params{
			MaxAge: auths.maxage, Secure: false,
			Path: "/", SameSite: securecookie.Lax,
		},
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

	// If the user can't be authenticate with a secure cookie,
	// try to autenticate the request with a token
	tk := r.Header.Get("GMC-Token")
	user, err := auths.Check("", tk)
	if err != nil {
		return nil, err
	}

	if user != nil {
		return user, nil
	}

	return nil, nil
}

func (auths *Auths) CheckForm(w http.ResponseWriter, r *http.Request) error {
	// Try to authenticate the user with a secure cookie
	cookie, err := securecookie.New(
		"gmc-session", auths.key,
		securecookie.Params{
			MaxAge: auths.maxage, Secure: false,
			Path: "/", SameSite: securecookie.Lax,
		},
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

	redirect := r.FormValue("redirect")
	if redirect == "" {
		redirect = r.URL.Query().Get("redirect")
	}
	if !redirect_rx.MatchString(redirect) {
		return fmt.Errorf("invalid redirect")
	}

	params := map[string]interface{}{}
	if redirect != "" {
		params["redirect"] = redirect
	}
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

			if redirect == "" {
				redirect = "."
			}
			http.Redirect(w, r, redirect, http.StatusFound)
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
