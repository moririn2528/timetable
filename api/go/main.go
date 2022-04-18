package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"timetable/communicate"
	"timetable/database"
	"timetable/solve"
	"timetable/usecase"
)

func init() {
	const location = "Asia/Tokyo"
	f, err := os.Create("1.log")
	if err != nil {
		panic(err)
	}
	log.SetOutput(f)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	loc, err := time.LoadLocation(location)
	if err != nil {
		loc = time.FixedZone(location, 9*60*60)
	}
	time.Local = loc
}

func test() {
	_, err := usecase.ChangeTimetable(0, time.Date(2021, 4, 1, 0, 0, 0, 0, time.Local), 3210405)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("ok!")
	os.Exit(1)
}

func main() {
	defer database.Close()

	database.SetFrames()
	usecase.Db_any = &database.DatabaseAny{}
	usecase.Db_class = &database.DatabaseClass{}
	usecase.Db_timetabale = &database.DatabaseTimetable{}
	usecase.Solver = &solve.SolverClass{}

	// test()

	http.Handle("/", http.FileServer(http.Dir("../../front")))
	http.HandleFunc("/api/timetable/class", communicate.ClassTimetableHandle)
	http.HandleFunc("/api/timetable/change", communicate.ChangeTimetableHandle)
	http.HandleFunc("/api/class", communicate.Class_structure)
	log.Print("start")
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "54321"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
