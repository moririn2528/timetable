package database

import (
	"fmt"
	"log"
	"os"
	"time"
	"timetable/library/logging"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	db     *sqlx.DB
	logger *logging.Logger = logging.NewLogger()
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
		logger.Errorf("sql connect tried: %v", i+1)
		time.Sleep(time.Second * 2)
	}
	if err != nil {
		log.Fatal("sql connect error", err)
	}
}

func Close() {
	db.Close()
}

func isDate(t time.Time) bool {
	return t.Hour() == 0 && t.Minute() == 0 && t.Second() == 0 && t.Nanosecond() == 0
}
