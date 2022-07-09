package database

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func SetHoliday() {
	f, err := os.Open("data/syukujitsu.csv")
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(transform.NewReader(f, japanese.ShiftJIS.NewDecoder()))
	var vals []string
	r.Read() // Header 削除
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		t, err := time.Parse("2006/1/2", record[0])
		if err != nil {
			log.Fatal(err)
		}
		vals = append(vals, "(\""+t.Format("2006-1-2")+"\")")
	}
	_, err = db.Exec("TRUNCATE TABLE holiday")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("INSERT INTO holiday values " + strings.Join(vals, ","))
	if err != nil {
		log.Fatal(err)
	}
}
