package database

import (
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	db *sqlx.DB
)

func Init() {
	env := os.Getenv("EXEC_ENV")
	host := ""

	if env == "local" {
		// execute in local
		host = "localhost"
	} else if env == "docker" {
		// execute in docker
		host = "db"
	} else {
		log.Fatal("EXEC_ENV should be local or docker")
	}
	var err error
	for i := 0; i < 30; i++ {
		db, err = sqlx.Connect("mysql", fmt.Sprintf(
			"%v:%v@tcp(%v:3306)/%v?parseTime=true&loc=Local",
			os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), host, os.Getenv("MYSQL_DATABASE"),
		))
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
