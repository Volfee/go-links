// Package database handles communication with BoltDB.
package database

import (
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/boltdb/bolt"
)

var dbOptions *bolt.Options = &bolt.Options{
	Timeout: 1 * time.Second,
}

const (
	dbName     = "urls.db"
	bucketName = "UrlFromPath"
)

func init() {
	// Make sure that the database exists.
	db, err := bolt.Open(dbName, 0600, dbOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// Make sure that UrlFromPath bucket exists.
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

// Lookup checks if `path` key exists in database and returns related `url` value if found. If not
// found second return value `ok` is set to false.
func Lookup(path string) (url string, ok bool) { // TODO: return
	// Open connection to DB.
	db, err := bolt.Open(dbName, 0600, dbOptions)
	if err != nil {
		log.Fatalf("Cannot connect to DB: %s", err)
		return "", false
	}
	defer db.Close()

	// Retreive url from DB.
	var url_bytes []byte
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		url_bytes = b.Get([]byte(path))
		return nil
	})
	if err != nil {
		return "", false
	}
	if url_bytes != nil {
		return string(url_bytes), true
	} else {
		return "", false
	}

}

// RegisterUrl saves provided key:value (path:url) pair to database.
func RegisterUrl(path, url string) error {
	// Open connection to DB.
	db, err := bolt.Open(dbName, 0600, dbOptions)
	if err != nil {
		return err
	}
	defer db.Close()

	// Check is valid URL and starts with http(s).
	url, err = validateAndFixUrl(url)
	if err != nil {
		return err
	}

	// Save path and url.
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		err := b.Put([]byte(path), []byte(url))
		return err
	})
	if err != nil {
		return err
	}

	log.Printf("Database: %s, %s registered.", path, url)
	return nil
}

// validateAndFixUrl validates if URL is valid and starts with http:// or https://. If not then appends
// and returns fixed link.
func validateAndFixUrl(url_string string) (string, error) {
	if _, err := url.ParseRequestURI(url_string); err != nil {
		return "", err
	}
	if !(strings.HasPrefix(url_string, "http://") || strings.HasPrefix(url_string, "https://")) {
		url_string = "https://" + url_string
	}

	return url_string, nil
}
