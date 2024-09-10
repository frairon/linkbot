package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/frairon/linkbot/internal/storage"
)

var dbFile = flag.String("file", "", "database file")

func main() {
	flag.Parse()

	if *dbFile == "" {
		log.Fatalf("specify database file!")
	}

	dbHandler, err := storage.NewStorage(fmt.Sprintf("file:%s", *dbFile))
	if err != nil {
		log.Fatalf("Error creating database %v", err)
	}
	dbHandler.Close()
}
