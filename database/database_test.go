package database

import "testing"

func setup() {
	// Init test db.
	// Init test bucket.
}

func cleanup() {
	// Remove test db.
}

func TestLookup(t *testing.T) {
	// T1. Path exists in db. Correct url is retreived.

	// T2. Path doesn't exist in db. false is returned.

}

func TestRegisterUrl(t *testing.T) {
	// T1. Add path:url to db. Correct path:url pair has been added to db.
}
