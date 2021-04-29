package database

import "log"

func Lookup(key string) (url string, ok bool) {
	if key == "google" {
		return "https://www.google.com", true
	} else {
		return "", false
	}

}

func RegisterUrl(key, url string) error {
	log.Printf("Database: %s, %s registered.", key, url)
	return nil
}
