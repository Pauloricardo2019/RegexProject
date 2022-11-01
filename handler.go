package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sync"
)

var gormDb *gorm.DB
var mutexDB sync.Mutex

func GetGormDB() (*gorm.DB, error) {
	mutexDB.Lock()
	defer mutexDB.Unlock()

	if gormDb != nil {
		return gormDb, nil
	}

	dsn := "host=localhost user=user password=password dbname=regex port=5439 sslmode=disable"

	newDb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{FullSaveAssociations: true})
	if err != nil {
		return nil, err
	}

	err = newDb.AutoMigrate(&Table{})
	if err != nil {
		panic(err)
	}
	gormDb = newDb

	return gormDb, nil
}
