package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ToshihitoKon/moneytemaa/src/constants"
	_ "github.com/go-sql-driver/mysql"
)

var (
	db *sql.DB
)

func NewDB() {
	var err error
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", constants.DBUser, constants.DBPass, constants.DBName))
	if err != nil {
		log.Fatal("err sql.Open ", err.Error())
	}
	return
}

func GetDB() *sql.DB {
	return db
}
