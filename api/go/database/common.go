package database

import (
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	db *sqlx.DB
)

func init() {
	sql_dsn, ok := os.LookupEnv("SQL_DSN")
	// if ok: execute in docker
	if !ok {
		// execute in local
		file, err := os.Open("database/dsn.secret")
		if err != nil {
			log.Printf("dsn file open warning, %v", err)
			file, err = os.Open("dsn.secret")
			if err != nil {
				log.Printf("dsn file open error, %v", err)
				return
			}
		}
		buf := make([]byte, 1024)
		n, err := file.Read(buf)
		if err != nil {
			log.Printf("dsn file read error, %v", err)
			return
		}
		sql_dsn = string(buf[:n])
	}
	var err error
	db, err = sqlx.Connect("mysql", sql_dsn)
	if err != nil {
		log.Fatal("sql error", err)
	}
}

func Close() {
	db.Close()
}
