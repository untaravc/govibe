package configs

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func OpenMySQL() (*sql.DB, error) {
	cfg, err := LoadDatabaseConfig()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("mysql", cfg.MySQLDSN())
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("db ping failed: %w", err)
	}

	return db, nil
}

func OpenGormMySQL() (*gorm.DB, error) {
	cfg, err := LoadDatabaseConfig()
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(mysql.Open(cfg.MySQLDSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if err := sqlDB.Ping(); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("db ping failed: %w", err)
	}

	return db, nil
}
