package database

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db *sql.DB
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
	db, err = sql.Open("mysql", sql_dsn)
	if err != nil {
		log.Fatal("sql error", err)
	}
	for i := 0; i < 20; i++ {
		err = db.Ping()
		if err == nil {
			break
		}
		log.Printf("sql ping, retry %v", i)
		time.Sleep(time.Second * 2)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("sql ping error: ", err)
	}
}

func Close() {
	db.Close()
}
