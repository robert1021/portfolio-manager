package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func connectDb() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dbFilename), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}
