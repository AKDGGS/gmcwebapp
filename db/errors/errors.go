package errors

import "errors"

// Authentication
var ErrToken = errors.New("Token is incorrect")

// Empty field
var ErrBarcodeCannotBeEmpty = errors.New("Barcode cannot be empty")
var ErrOldBarcodeCannotBeEmpty = errors.New("Old barcode cannot be empty")
var ErrNewBarcodeCannotBeEmpty = errors.New("New barcode cannot be empty")
var ErrDestinationBarcodeEmpty = errors.New("Destination barcode cannot be empty")
var ErrSourceBarcodeEmpty = errors.New("Source barcode cannot be empty")
var ErrListOfBarcodesEmpty = errors.New("List of barcodes is empty")

// Existance
var ErrBarcodeExists = errors.New("Barcode already exists")

// Failure
var ErrInventoryInsertFailed = errors.New("Inventory insert failed")
var ErrInventoryQualityInsertFailed = errors.New("Inventory quality insert failed")
var ErrAuditInsertFailed = errors.New("Audit group insert failed")
var ErrAuditParamsEmpty = errors.New("Both the remark and the items list are empty. Audit insert failed.")
var ErrNothingMoved = errors.New("Move failed")
var ErrNothingRecoded = errors.New("Record failed")

// Not found
var ErrBarcodeNotFound = errors.New("Barcode not found")
var ErrDestinationNotFound = errors.New("The destination barcode not found")
var ErrNotFoundInInventory = errors.New("Barcode not found in Inventory")
var ErrAtLeastOneBarcodeNotFound = errors.New("At least one barcode not found")
var ErrSourceNotFound = errors.New("The source barcode not found")

// Type Conflict
var ErrBarcodeNotContainer = errors.New("Barcode is not a container")
var ErrMultipleIDs = errors.New("Multiple IDs returned")
var ErrDestinationMultipleContainers = errors.New("The destination barcode refers to multiple containers")
var ErrSourceNotValid = errors.New("The source barcode not valid")
