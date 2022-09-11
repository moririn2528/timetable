package database

import (
	"log"
	"os"
	"time"

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
	for i := 0; i < 10; i++ {
		db, err = sqlx.Connect("mysql", sql_dsn)
		if err == nil {
			break
		}
		log.Printf("sql connect tried: %v", i+1)
		time.Sleep(time.Second * 2)
	}
	if err != nil {
		log.Fatal("sql connect error", err)
	}
}

func Close() {
	db.Close()
}
