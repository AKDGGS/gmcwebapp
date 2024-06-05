package ldap

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"text/template"

	authu "gmc/auth/util"

	"github.com/go-ldap/ldap/v3"
)

type LDAPAuth struct {
	name                       string
	ldap_url                   string
	base_dn                    string
	bind_dn                    string
	bind_password              string
	user_search                string
	bind_as_user               bool
	disable_certificate_verify bool
}

func New(cfg map[string]interface{}) (*LDAPAuth, error) {
	name, ok := cfg["name"].(string)
	if !ok {
		name = "ldap"
	}
	ldap_url, ok := cfg["ldap_url"].(string)
	if !ok {
		return nil, fmt.Errorf("ldap_url is required and must be a string")
	}
	bind_as_user, _ := cfg["bind_as_user"].(bool)
	base_dn, ok := cfg["base_dn"].(string)
	if !ok && !bind_as_user {
		return nil, fmt.Errorf("base_dn is required and must be a string")
	}
	bind_dn, ok := cfg["bind_dn"].(string)
	if !ok && !bind_as_user {
		return nil, fmt.Errorf("bind_dn is required and must be a string")
	}
	bind_password, ok := cfg["bind_password"].(string)
	if !ok && !bind_as_user {
		return nil, fmt.Errorf("when bind_as_user is false, bind_password is required and must be a string")
	}
	user_search, ok := cfg["user_search"].(string)
	if !ok && !bind_as_user {
		return nil, fmt.Errorf("when bind_as_user is false, user_search is required and must be a string")
	}
	disable_certificate_verify, _ := cfg["disable_certificate_verify"].(bool)

	a := &LDAPAuth{
		name:                       name,
		ldap_url:                   ldap_url,
		base_dn:                    base_dn,
		bind_dn:                    bind_dn,
		bind_password:              bind_password,
		user_search:                user_search,
		bind_as_user:               bind_as_user,
		disable_certificate_verify: disable_certificate_verify,
	}
	return a, nil
}

func (a *LDAPAuth) Name() string {
	return a.name
}

func (a *LDAPAuth) Check(u string, p string) (*authu.User, error) {
	conn, err := ldap.DialURL(a.ldap_url,
		ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: a.disable_certificate_verify}))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	var buf bytes.Buffer
	if a.bind_as_user {
		t, err := template.New("bind_dn_tmpl").Parse(a.bind_dn)
		if err != nil {
			return nil, err
		}
		if err := t.Execute(&buf, ldap.EscapeDN(u)); err != nil {
			return nil, err
		}
		if err := conn.Bind(buf.String(), p); err != nil {
			return nil, nil
		}
	} else {
		t, err := template.New("user_search_tmpl").Parse(a.user_search)
		if err != nil {
			return nil, err
		}
		if err := t.Execute(&buf, ldap.EscapeFilter(u)); err != nil {
			return nil, err
		}
		if err := conn.Bind(a.bind_dn, a.bind_password); err != nil {
			return nil, nil
		}
		search_request := ldap.NewSearchRequest(
			a.base_dn,
			ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
			buf.String(),
			[]string{},
			nil,
		)
		search_result, err := conn.Search(search_request)
		if err != nil {
			return nil, err
		}
		if len(search_result.Entries) != 1 {
			return nil, nil
		}
		if err := conn.Bind(search_result.Entries[0].DN, p); err != nil {
			return nil, nil
		}
	}
	return &authu.User{Username: u}, nil
}
