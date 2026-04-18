package main

import (
	"fmt"
	"log"
	"os"

	"govibe/configs"
	"govibe/database/seeder"
)

func main() {
	cmd := "all"
	if len(os.Args) >= 2 {
		cmd = os.Args[1]
	}

	db, err := configs.OpenGormMySQL()
	if err != nil {
		log.Fatal(err)
	}

	switch cmd {
	case "all":
		if err := seeder.SeedMenus(db); err != nil {
			log.Fatal(err)
		}
	case "menu":
		if err := seeder.SeedMenus(db); err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Fprintln(os.Stderr, "Usage:")
		fmt.Fprintln(os.Stderr, "  go run ./cmd/seed all")
		fmt.Fprintln(os.Stderr, "  go run ./cmd/seed menu")
		os.Exit(1)
	}

	fmt.Println("seed completed")
}
