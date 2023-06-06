package file

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	authu "gmc/auth/util"

	"golang.org/x/crypto/bcrypt"
)

var alphabet []byte = []byte{
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r',
	's', 't', 'u', 'v', 'w', 'x', 'y', 'z', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
}

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

	path, _ := cfg["path"].(string)

	pw := &PasswordFile{
		name:  name,
		path:  path,
		users: make(map[string]*authu.User),
	}

	if err := pw.readSource(); err != nil {
		return nil, err
	}

	// If there's no users in the file, generate one
	if len(pw.users) < 1 {
		// Generate a random password
		rand.Seed(time.Now().UnixNano())
		pass := make([]byte, 24)
		for i, _ := range pass {
			pass[i] = alphabet[rand.Intn(len(alphabet))]
		}

		// bcrypt the password
		encpass, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}

		// Set the admin user with the generated password
		pw.users["admin"] = &authu.User{Username: "admin", Password: encpass}

		fmt.Printf("NOTICE: admin user password set to %s\n", pass)
	}

	return pw, nil
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

func (pw *PasswordFile) readSource() error {
	// Ignore empty paths
	if pw.path == "" {
		return nil
	}

	f, err := os.Open(pw.path)
	if err != nil {
		// Don't throw an error if the file doesn't exist,
		// just don't load any users.
		if errors.Is(err, os.ErrNotExist) {
			fmt.Printf("WARNING: password file (%s) not found.\n", pw.path)
			return nil
		}
		return err
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
			return fmt.Errorf("Syntax error in %s on line %d", pw.path, lc)
		}

		user := string(ba[0])

		if _, exists := pw.users[user]; exists {
			return fmt.Errorf(
				"Duplicate user %s in %s on line %d", user, pw.path, lc,
			)
		}

		pw.users[user] = &authu.User{Username: user, Password: ba[1]}
	}
	return nil
}
