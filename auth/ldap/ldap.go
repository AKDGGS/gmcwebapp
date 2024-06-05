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
	name          string
	ldap_url      string
	base_dn       string
	bind_dn       string
	bind_password string
	user_search   string
	bind_as_user  bool
	skip_verify   bool
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
	bind_as_user, ok := cfg["bind_as_user"].(bool)
	if !ok {
		bind_as_user = false
	}
	base_dn, ok := cfg["base_dn"].(string)
	if !ok && !bind_as_user {
		return nil, fmt.Errorf("base_dn is required and must be a string")
	}
	bind_dn, ok := cfg["bind_dn"].(string)
	if !ok && !bind_as_user {
		return nil, fmt.Errorf("bind_dn must be a string")
	}
	bind_password, ok := cfg["bind_password"].(string)
	if !ok && !bind_as_user {
		return nil, fmt.Errorf("bind_password must be a string")
	}
	user_search, ok := cfg["user_search"].(string)
	if !ok && !bind_as_user {
		return nil, fmt.Errorf("user_search must be a string")
	}
	skip_verify, _ := cfg["skip_verify"].(bool)

	a := &LDAPAuth{
		name:          name,
		ldap_url:      ldap_url,
		base_dn:       base_dn,
		bind_dn:       bind_dn,
		bind_password: bind_password,
		user_search:   user_search,
		bind_as_user:  bind_as_user,
		skip_verify:   skip_verify,
	}
	return a, nil
}

func (a *LDAPAuth) Name() string {
	return a.name
}

func (a *LDAPAuth) Check(u string, p string) (*authu.User, error) {
	conn, err := ldap.DialURL(a.ldap_url,
		ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: a.skip_verify}))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	t, err := template.New("bind_dn_tmpl").Parse(a.bind_dn)
	if err != nil {
		return nil, err
	}
	if a.bind_as_user {
		var bind_dn_buf bytes.Buffer
		if err := t.Execute(&bind_dn_buf, ldap.EscapeDN(u)); err != nil {
			return nil, err
		}
		if err := conn.Bind(bind_dn_buf.String(), p); err != nil {
			return nil, nil
		}
	} else {
		if err := conn.Bind(a.bind_dn, a.bind_password); err != nil {
			return nil, nil
		}
		t_filter, err := template.New("user_search_tmpl").Parse(a.user_search)
		if err != nil {
			return nil, err
		}
		var filter_buf bytes.Buffer
		if err := t_filter.Execute(&filter_buf, ldap.EscapeFilter(u)); err != nil {
			return nil, err
		}
		search_request := ldap.NewSearchRequest(
			a.base_dn,
			ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
			filter_buf.String(),
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
