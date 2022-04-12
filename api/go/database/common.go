package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db *sql.DB
)

func init() {
	var err error
	db, err = sql.Open("mysql", "h2JwdKNb:t7DVSdTB@tcp(localhost:3306)/timetable?parseTime=true&loc=Local")
	if err != nil {
		log.Fatal("sql error", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("sql ping error: ", err)
	}
}

func Close() {
	db.Close()
}
