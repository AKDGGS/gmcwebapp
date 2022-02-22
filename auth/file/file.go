package file

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	authu "gmc/auth/util"
	"golang.org/x/crypto/bcrypt"
	"os"
)

type PasswordFile struct {
	name  string
	path  string
	key   []byte
	users map[string]*authu.User
}

func New(cfg map[string]interface{}) (*PasswordFile, error) {
	name, ok := cfg["name"].(string)
	if !ok {
		name = "file"
	}

	path, ok := cfg["path"].(string)
	if !ok {
		return nil, fmt.Errorf("auth file path must exist and be a string")
	}

	pf := &PasswordFile{
		name:  name,
		path:  path,
		users: make(map[string]*authu.User),
	}

	f, err := os.Open(path)
	if err != nil {
		// Don't throw an error if the file doesn't exist,
		// just don't load any users.
		if errors.Is(err, os.ErrNotExist) {
			return pf, nil
		}

		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for lc := 1; scanner.Scan(); lc++ {
		line := scanner.Bytes()

		// Ignore blank lines
		if len(line) < 1 {
			continue
		}

		ba := bytes.Split(line, []byte(":"))
		if len(ba) < 1 || len(ba) > 2 {
			return nil, fmt.Errorf("Syntax error in %s on line %d", path, lc)
		}

		user := string(ba[0])

		if _, exists := pf.users[user]; exists {
			return nil, fmt.Errorf(
				"Duplicate user %s in %s on line %d", user, path, lc,
			)
		}

		pf.users[user] = &authu.User{Username: user, Password: ba[1]}
	}

	return pf, nil
}

func (pw *PasswordFile) Name() string {
	return pw.name
}

func (pw *PasswordFile) Check(u string, p string) (*authu.User, error) {
	un, exists := pw.users[u]
	if !exists {
		return nil, nil
	}

	err := bcrypt.CompareHashAndPassword(un.Password, []byte(p))
	if err == nil {
		return un, nil
	}
	return nil, nil
}
