package db

import (
	"fmt"

	"gmc/config"
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

	// Fetches stash for a specific inventory id
	GetInventoryStash(id int) (interface{}, error)

	// Get a flattened version of all inventory appropriate
	// for search indexing, calling function parameter for
	// each individual item
	GetFlatInventory(func(*model.FlatInventory) error) error

	// Fetches the complete details for a Summary by barcode
	GetSummaryByBarcode(barcode string, flags int) (*model.Summary, error)

	// Fetches wells point list
	GetWellPoints() (interface{}, error)

	// Fetches file details
	GetFile(id int) (*model.File, error)

	// Put file
	PutFile(*model.File, func() error) error

	// Delete a file from the file and linking tables
	DeleteFile(file *model.File, rm_links bool) error

	// Updates container_id (dest) for inventory in barcodes_to_move
	MoveInventoryAndContainers(dest string, barcodes_to_move []string, username string) error

	// Adds Audit to db
	AddAudit(remark string, container_list []string) error

	// Updates container_id (dest) for inventory in source
	MoveInventoryAndContainersContents(src string, dest string) error

	// Insert a new container
	AddContainer(barcode string, name string, remark string) error

	// Insert a new inventory item
	AddInventory(barcode string, remark string, container_id *int32, issues []string, username string) error

	// Insert issues
	AddInventoryQuality(barcode string, remark string, issues []string, username string) error

	// Update barcode
	RecodeInventoryAndContainer(old_barcode string, new_barcode string) error

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
	// WARNING: This is destructive and intended only for use in
	// development.
	SchemaDrop() error

	// Shutdown this database connection
	Shutdown()
}

func New(cfg config.DatabaseConfig) (DB, error) {
	var db DB
	var err error

	switch cfg.Type {
	case "pg", "postgres", "postgresql":
		db, err = pg.New(cfg.Attrs)
		if err != nil {
			return nil, err
		}
	case "":
		return nil, fmt.Errorf("database type cannot be empty")
	default:
		return nil, fmt.Errorf("unknown database type: %s", cfg.Type)
	}
	return db, nil
}
