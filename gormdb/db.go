package gormdb

import (
	"github.com/adneg/golangmodules/logtrace"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

var (
	conf Configuration
)

func Init(configfile string) {

	conf = loadconfig(configfile)
}
func Start() (db *gorm.DB) {
	db, err := gorm.Open(conf.DriverDB, conf.Databasefile)
	db.DB().SetMaxIdleConns(1)
	db.DB().SetMaxOpenConns(1)

	if err != nil {
		logtrace.Error.Fatalln(err.Error())
	}

	db.Exec("PRAGMA foreign_keys = ON")
	db.LogMode(conf.DatabaseDebug)
	return
}
