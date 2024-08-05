package db

import (
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Store struct {
	Db *gorm.DB
}

func NewStore(dbName string) (Store, error) {
	Db, err := getConnection(dbName)
	if err != nil {
		return Store{}, err
	}

	return Store{
		Db,
	}, nil
}

func getConnection(dbName string) (*gorm.DB, error) {
	var (
		err error
		db  *gorm.DB
	)

	if db != nil {
		return db, nil
	}

	// Init SQLite3 database
	db, err = gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		// log.Fatalf("🔥 failed to connect to the database: %s", err.Error())
		return nil, fmt.Errorf("🔥 failed to connect to the database: %s", err)
	}

	log.Println("🚀 Connected Successfully to the Database")

	return db, nil
}
