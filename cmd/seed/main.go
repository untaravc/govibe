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
		if err := seeder.SeedRoles(db); err != nil {
			log.Fatal(err)
		}
		if err := seeder.SeedUsers(db); err != nil {
			log.Fatal(err)
		}
		if err := seeder.SeedMenus(db); err != nil {
			log.Fatal(err)
		}
	case "menu":
		if err := seeder.SeedMenus(db); err != nil {
			log.Fatal(err)
		}
	case "regions":
		if err := seeder.SeedRegions(db); err != nil {
			log.Fatal(err)
		}
	case "roles":
		if err := seeder.SeedRoles(db); err != nil {
			log.Fatal(err)
		}
	case "users":
		if err := seeder.SeedUsers(db); err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Fprintln(os.Stderr, "Usage:")
		fmt.Fprintln(os.Stderr, "  go run ./cmd/seed all")
		fmt.Fprintln(os.Stderr, "  go run ./cmd/seed menu")
		fmt.Fprintln(os.Stderr, "  go run ./cmd/seed regions")
		fmt.Fprintln(os.Stderr, "  go run ./cmd/seed roles")
		fmt.Fprintln(os.Stderr, "  go run ./cmd/seed users")
		os.Exit(1)
	}

	fmt.Println("seed completed")
}
