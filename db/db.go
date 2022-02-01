package db

import (
	"fmt"
	"math"
	"net/url"
	"strings"
)

// Flags used to control how much additional data is pulled
// in with queries
const (
	FILES = 1 << iota
	INVENTORY_SUMMARY
	MINING_DISTRICTS
	QUADRANGLES
	PRIVATE
	GEOJSON
)

// Option for everything
const ALL int = math.MaxInt

// Option for everything except private items
const ALL_NOPRIVATE int = math.MaxInt &^ PRIVATE

// Option for minimal return
const MINIMAL int = 0

type DB interface {
	Init() error
	Verify() error
	Drop() error
	GetProspect(int, int) (map[string]interface{}, error)
	GetBorehole(int, int) (map[string]interface{}, error)
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
