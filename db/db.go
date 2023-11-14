package db

import (
	"fmt"
	"net/url"
	"strings"

	authu "gmc/auth/util"
	"gmc/db/model"
	"gmc/db/pg"
)

type DB interface {
	// Fetches the complete details for a Prospect
	GetProspect(id int, flags int) (*model.Prospect, error)

	// Fetches the complete details for a Borehole
	GetBorehole(id int, flags int) (*model.Borehole, error)

	// Fetches the complete details for an Outcrop
	GetOutcrop(id int, flags int) (*model.Outcrop, error)

	// Fetches the complete details for a Well
	GetWell(id int, flags int) (*model.Well, error)

	// Fetches the complete details for a Shotline
	GetShotline(id int, flags int) (*model.Shotline, error)

	// Fetches the complete details for an Inventory
	GetInventory(id int, flags int) (*model.Inventory, error)

	// Fetches the complete details for an Inventory by barcode
	GetInventoryByBarcode(barcode string, flags int) ([]*model.Inventory, error)

	// Fetches the complete details for a Summary by barcode
	GetSummaryByBarcode(barcode string, flags int) (*model.Summary, error)

	// Fetches stash for a specific inventory id
	GetStash(id int) (map[string]interface{}, error)

	// Fetches wells point list
	GetWellPoints() ([]map[string]interface{}, error)

	// Fetches file details
	GetFile(id int) (*model.File, error)

	// Put file
	PutFile(*model.File, func() error) error

	// Updates container_id (dest) for inventory in container_list
	MoveByBarcode(dest string, container_list []string, user *authu.User) error

	//Insert a new container
	AddContainer(barcode string, name string, remark string) error

	//Insert a new inventory item
	AddInventory(barcode string, remark string, container_id *int32, issues []string, username string) error

	// List available tokens
	ListTokens() ([]*model.Token, error)
	// Creates a new token
	CreateToken(token *model.Token) error
	// Removes a token
	DeleteToken(id int) error
	// Check token validity
	CheckToken(t string) (*model.Token, error)

	// Lists available keywords
	ListKeywords() ([]string, error)
	// Adds any number of keywords
	AddKeywords(keywords ...string) error
	// Deletes any number of keywords
	DeleteKeywords(keywords ...string) error

	// Lists available issues
	ListIssues() ([]string, error)
	// Add any number of issues
	AddIssues(issues ...string) error
	// Deletes any number of issues
	DeleteIssues(issues ...string) error

	// Lists quality assurance reports
	ListQAReports() ([]map[string]string, error)
	// Runs just the count of a specific QA report
	CountQAReport(id int) (int, error)
	// Runs a specific QA report and returns the results
	RunQAReport(id int) (*model.Table, error)

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
