package auth

import (
	"bytes"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"gmc/assets"
	"gmc/auth/securecookie"
	authu "gmc/auth/util"
	webu "gmc/web/util"
)

var redirect_rx *regexp.Regexp = regexp.MustCompile(`^(?:\.|borehole|inventory|outcrop|prospect|shotline|well|wells|qa)\/?(\d*|search|search\-help)$`)

func (auths *Auths) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := securecookie.New(
		"gmc-session", auths.key,
		securecookie.Params{
			MaxAge: auths.maxage, Secure: false,
			Path: "/", SameSite: securecookie.Lax,
		},
	)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("secure cookie error: %s", err),
			http.StatusInternalServerError,
		)
		return
	}

	cookie.Delete(w)

	redirect := r.URL.Query().Get("redirect")
	if redirect == "" {
		redirect = "."
	}
	if !redirect_rx.MatchString(redirect) {
		http.Error(
			w,
			fmt.Sprintf("invalid redirect"),
			http.StatusInternalServerError,
		)
		return
	}

	webu.Redirect(w, redirect, http.StatusFound)
}

func (auths *Auths) CheckRequest(w http.ResponseWriter, r *http.Request) (*authu.User, error) {
	var user *authu.User
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

	if user == nil {
		if authorization := r.Header.Get("Authorization"); authorization != "" {
			// If the user can't be authenticate with a secure cookie,
			// try to autenticate the request with a token
			if strings.HasPrefix(authorization, "Token ") {
				tk := strings.TrimPrefix(authorization, "Token ")
				user, err = auths.Check("", tk)
				if err != nil {
					return nil, err
				}
			} else {
				// if there is no secure cookie and no token,
				// check for an authorization header in the request
				username, password, ok := r.BasicAuth()
				if ok {
					user, err = auths.Check(username, password)
					if err != nil {
						return nil, err
					}
				}
			}
		}

		if user != nil {
			return user, nil
		}
	}

	return nil, nil
}

func (auths *Auths) CheckForm(w http.ResponseWriter, r *http.Request) {
	// Try to authenticate the user with a secure cookie
	cookie, err := securecookie.New(
		"gmc-session", auths.key,
		securecookie.Params{
			MaxAge: auths.maxage, Secure: false,
			Path: "/", SameSite: securecookie.Lax,
		},
	)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("secure cookie error: %s", err),
			http.StatusInternalServerError,
		)
		return
	}

	uj, err := cookie.GetValue(nil, r)
	if err == nil {
		user, err := authu.UnmarshalUser(uj)
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("json unmarshal error: %s", err),
				http.StatusInternalServerError,
			)
			return
		}

		if user != nil {
			// Refresh the cookie, extending the timeout
			// and ignore any failures.
			cookie.SetValue(w, uj)

			webu.Redirect(w, ".", http.StatusFound)
			return
		}
	}

	// If the user can't be authenticate with a secure cookie,
	// try to read a POSTed username and password to authenticate with
	err = r.ParseForm()
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("parse error: %s", err),
			http.StatusInternalServerError,
		)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	redirect := r.FormValue("redirect")
	if redirect == "" {
		redirect = r.URL.Query().Get("redirect")
	}
	if !redirect_rx.MatchString(redirect) {
		http.Error(
			w,
			fmt.Sprintf("invalid redirect"),
			http.StatusInternalServerError,
		)
		return
	}

	params := map[string]interface{}{}
	if redirect != "" {
		params["redirect"] = redirect
	}
	if username != "" && password != "" {
		params["username"] = username

		user, err := auths.Check(username, password)
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("authentication error: %s", err),
				http.StatusBadRequest,
			)
			return
		}

		if user != nil {
			uj, err := authu.MarshalUser(user)
			if err != nil {
				http.Error(
					w,
					fmt.Sprintf("json marshal error: %s", err),
					http.StatusInternalServerError,
				)
				return
			}

			err = cookie.SetValue(w, uj)
			if err != nil {
				http.Error(
					w,
					fmt.Sprintf("cookie error: %s", err),
					http.StatusInternalServerError,
				)
				return
			}

			if redirect == "" {
				redirect = "."
			}
			webu.Redirect(w, redirect, http.StatusFound)
			return
		}

		params["error"] = "invalid username or password"
	}

	// If there's no secure cookie, and no POSTed credentials,
	// serve up the login page
	buf := bytes.Buffer{}
	if err := assets.ExecuteTemplate("tmpl/login.html", &buf, params); err != nil {
		http.Error(
			w,
			fmt.Sprintf("parse error: %s", err),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Length", fmt.Sprintf("%d", buf.Len()))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(buf.Bytes())
}
