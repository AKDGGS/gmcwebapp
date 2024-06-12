package errors

import "errors"

// Authentication
var ErrToken = errors.New("token is incorrect")

// Empty field
var ErrBarcodeCannotBeEmpty = errors.New("barcode cannot be empty")
var ErrOldBarcodeCannotBeEmpty = errors.New("old barcode cannot be empty")
var ErrNewBarcodeCannotBeEmpty = errors.New("new barcode cannot be empty")
var ErrDestinationBarcodeEmpty = errors.New("destination barcode cannot be empty")
var ErrSourceBarcodeEmpty = errors.New("source barcode cannot be empty")
var ErrListOfBarcodesEmpty = errors.New("list of barcodes is empty")

// Existance
var ErrBarcodeExists = errors.New("barcode already exists")

// Failure
var ErrInventoryInsertFailed = errors.New("inventory insert failed")
var ErrInventoryQualityInsertFailed = errors.New("inventory quality insert failed")
var ErrAuditInsertFailed = errors.New("audit group insert failed")
var ErrAuditParamsEmpty = errors.New("both the remark and the items list are empty. Audit insert failed.")
var ErrNothingMoved = errors.New("move failed")
var ErrNothingRecoded = errors.New("record failed")

// Upload
var ErrUploadToFileStoreFailed = errors.New("upload to filestore failed")

// Not found
var ErrBarcodeNotFound = errors.New("barcode not found")
var ErrDestinationNotFound = errors.New("the destination barcode not found")
var ErrNotFoundInInventory = errors.New("barcode not found in Inventory")
var ErrAtLeastOneBarcodeNotFound = errors.New("at least one barcode not found")
var ErrSourceNotFound = errors.New("the source barcode not found")

// Type Conflict
var ErrBarcodeNotContainer = errors.New("barcode is not a container")
var ErrMultipleIDs = errors.New("multiple IDs returned")
var ErrDestinationMultipleContainers = errors.New("the destination barcode refers to multiple containers")
var ErrSourceNotValid = errors.New("the source barcode not valid")

// Miscellaneous
var ErrSrcNoInv = errors.New("the source has no inventory")
