package user

import (
	"encoding/json"
)

type User struct {
	Username string `json:"user"`
	Password []byte `json:"-"`
}

func UnmarshalUser(u []byte) (*User, error) {
	user := &User{}
	err := json.Unmarshal(u, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func MarshalUser(u *User) ([]byte, error) {
	return json.Marshal(u)
}
