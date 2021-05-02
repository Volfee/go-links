package main

import (
	"golinks/webserver"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	webserver.Run(port)
}
