package database

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/boltdb/bolt"
)

const (
	testDataBase = "test.db"
	testBucket   = "UrlFromPath"
)

func setup() {
	// Create database
	db, err := bolt.Open(testDataBase, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatalf("could not create test db: %s, db: %v", err, db)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		// Create test bucket
		bucket, err := tx.CreateBucketIfNotExists([]byte(testBucket))
		if err != nil {
			return err
		}
		// Create test record
		err = bucket.Put([]byte("path-exists"), []byte("http://example.com"))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func cleanup() {
	os.Remove(testDataBase)
}

func TestLookup(t *testing.T) {
	setup()
	defer cleanup()

	test_db := ConnectWith(testDataBase)

	// T1. Path exists in db. Correct url is retreived.
	url, ok := test_db.Lookup("path-exists")
	if url != "http://example.com" || !ok {
		t.Errorf("expected record found. Got: url: %s, ok: %v", url, ok)
	}

	// T2. Path doesn't exist in db. false is returned.
	url, ok = test_db.Lookup("path-not-exists")
	if url != "" || ok {
		t.Errorf("expected record not found. Got: url: %s, ok: %v", url, ok)
	}
}

func TestRegisterUrl(t *testing.T) {
	setup()
	defer cleanup()

	test_db := ConnectWith(testDataBase)

	// T1. Add path:url to db. Correct path:url pair has been added to db.
	err := test_db.RegisterUrl("test-path", "http://example.com")
	if err != nil {
		t.Errorf("coundn't register url to db, err: %s", err)
	}

	url, ok := test_db.Lookup("test-path")
	if url != "http://example.com" || !ok {
		t.Errorf(`record ("test-path", "http://example.com") was no created correctly`)
	}

	// T2. Add path:url to db, but url doesn't start with http.
	err = test_db.RegisterUrl("test-path", "example.com")
	if err != nil {
		t.Errorf("coundn't register url to db, err: %s", err)
	}

	url, ok = test_db.Lookup("test-path")
	if url != "https://example.com" || !ok {
		t.Errorf("path example.com was not conferted to https:// version")
	}
}
