package db

import (
	"fmt"
	"gmc/db/model"
	"gmc/db/pg"
	"net/url"
	"strings"
)

type DB interface {
	// Fetches the complete details for a Prospect
	GetProspect(id int, flags int) (map[string]interface{}, error)

	// Fetches the complete details for a Borehole
	GetBorehole(id int, flags int) (map[string]interface{}, error)

	// Fetches the complete details for an Outcrop
	GetOutcrop(id int, flags int) (map[string]interface{}, error)

	// Fetches the complete details for a Well
	GetWell(id int, flags int) (map[string]interface{}, error)

	// Fetches the complete details for a Shotline
	GetShotline(id int, flags int) (map[string]interface{}, error)

	// Fetches the complete details for an Inventory
	GetInventory(id int, flags int) (map[string]interface{}, error)

	// Fetches stash for a specific inventory id
	GetStash(id int) (map[string]interface{}, error)

	// Fetches wells point list for a specific inventory id
	GetWellPoints() ([]map[string]interface{}, error)

	// List available tokens
	ListTokens() ([]*model.Token, error)

	// Creates a new token
	CreateToken(token *model.Token) error

	// Removes a token
	DeleteToken(id int) error

	// Lists quality assurance reports
	ListQAReports() ([]map[string]string, error)

	// Runs just the count of a specific QA report
	CountQAReport(id int) (int, error)

	// Verify the database connection is working.
	// (usually by performing a simple query)
	Verify() error

	// Initializes schema for a new installation. Throws an error
	// if the schema already exists, or the initialization fails.
	SchemaInit() error

	// Removes schema from configured database.
	// WARNING: This is destructive and intended only for use in development.
	SchemaDrop() error

	// Shutdown this database connection
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
	case "pg", "postgres", "postgresql":
		db, err = pg.New(u)
		if err != nil {
			return nil, err
		}
	}
	return db, nil
}
