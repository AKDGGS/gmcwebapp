package flag

import "math"

// Flags used to control how much additional data is pulled
// in with queries
const (
	FILES = 1 << iota
	INVENTORY_SUMMARY
	MINING_DISTRICTS
	QUADRANGLES
	PRIVATE
	GEOJSON
	ORGANIZATION
)

// Option for everything
const ALL int = math.MaxInt

// Option for everything except private items
const ALL_NOPRIVATE int = ALL &^ PRIVATE

// Option for minimal return
const MINIMAL int = 0
