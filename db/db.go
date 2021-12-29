package db

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

type DB interface {
	GetProspect(int) (map[string]interface{}, error)
	GetFile(int, bool) (int, string, time.Time, error)
	Shutdown()
}

func New(su string) (DB, error) {
	u, err := url.Parse(su)
	if err != nil {
		return nil, err
	}

	if u.Scheme == "" {
		return nil, fmt.Errorf("URL must include a scheme")
	}

	var db DB
	switch strings.ToLower(u.Scheme) {
	case "postgres", "postgresql":
		db, err = newPostgres(u)
		if err != nil {
			return nil, err
		}
	}
	return db, nil
}
