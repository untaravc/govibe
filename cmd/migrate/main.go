package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"govibe/configs"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pressly/goose/v3"
)

func main() {
	if len(os.Args) < 2 {
		usage(1)
	}

	cmd := os.Args[1]
	dir := "database/migrations"

	if cmd == "create" {
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "missing migration name")
			usage(1)
		}
		name := strings.Join(os.Args[2:], "_")
		if err := createSQLMigration(dir, name); err != nil {
			log.Fatal(err)
		}
		return
	}

	cfg, err := configs.LoadDatabaseConfig()
	if err != nil {
		log.Fatal(err)
	}

	dsn := cfg.MySQLDSN()
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("db ping failed: %v", err)
	}

	goose.SetDialect("mysql")

	switch cmd {
	case "up":
		err = goose.Up(db, dir)
	case "up-to":
		if len(os.Args) < 3 {
			log.Fatal("missing version for up-to")
		}
		version, parseErr := parseVersion(os.Args[2])
		if parseErr != nil {
			log.Fatal(parseErr)
		}
		err = goose.UpTo(db, dir, version)
	case "down":
		err = goose.Down(db, dir)
	case "down-to":
		if len(os.Args) < 3 {
			log.Fatal("missing version for down-to")
		}
		version, parseErr := parseVersion(os.Args[2])
		if parseErr != nil {
			log.Fatal(parseErr)
		}
		err = goose.DownTo(db, dir, version)
	case "status":
		err = goose.Status(db, dir)
	case "version":
		err = goose.Version(db, dir)
	default:
		usage(1)
	}

	if err != nil {
		log.Fatal(err)
	}
}

func usage(exitCode int) {
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintln(os.Stderr, "  go run ./cmd/migrate create <name>         # create SQL migration file")
	fmt.Fprintln(os.Stderr, "  go run ./cmd/migrate up                    # apply all pending migrations")
	fmt.Fprintln(os.Stderr, "  go run ./cmd/migrate up-to <version>       # apply up to version")
	fmt.Fprintln(os.Stderr, "  go run ./cmd/migrate down                  # rollback 1 migration")
	fmt.Fprintln(os.Stderr, "  go run ./cmd/migrate down-to <version>     # rollback down to version")
	fmt.Fprintln(os.Stderr, "  go run ./cmd/migrate status                # show migration status")
	fmt.Fprintln(os.Stderr, "  go run ./cmd/migrate version               # show current version")
	os.Exit(exitCode)
}

func parseVersion(s string) (int64, error) {
	var v int64
	_, err := fmt.Sscanf(strings.TrimSpace(s), "%d", &v)
	if err != nil || v <= 0 {
		return 0, fmt.Errorf("invalid version %q", s)
	}
	return v, nil
}

func createSQLMigration(dir, name string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	version := time.Now().UTC().Format("20060102150405")
	slug := slugify(name)
	filename := fmt.Sprintf("%s_%s.sql", version, slug)
	path := filepath.Join(dir, filename)

	contents := strings.Join([]string{
		"-- +goose Up",
		"-- SQL in section 'Up' is executed when this migration is applied.",
		"",
		"-- +goose Down",
		"-- SQL section 'Down' is executed when this migration is rolled back.",
		"",
	}, "\n")

	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		return err
	}

	fmt.Println(path)
	return nil
}

var nonSlug = regexp.MustCompile(`[^a-z0-9_]+`)

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, " ", "_")
	s = nonSlug.ReplaceAllString(s, "_")
	s = strings.Trim(s, "_")
	if s == "" {
		return "migration"
	}
	return s
}
