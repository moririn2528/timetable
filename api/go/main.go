package main

import (
	"log"
	"net/http"
	"os"

	"timetable/communicate"
	"timetable/database"
	"timetable/solve"
	"timetable/usecase"
)

func init() {
	f, err := os.Create("1.log")
	if err != nil {
		panic(err)
	}
	log.SetOutput(f)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	defer database.Close()

	database.SetFrames()
	usecase.Db_any = &database.DatabaseAny{}
	usecase.Db_class = &database.DatabaseClass{}
	usecase.Db_timetabale = &database.DatabaseTimetable{}
	usecase.Solver = &solve.SolverClass{}

	http.Handle("/", http.FileServer(http.Dir("../../front")))
	http.HandleFunc("/api/timetable/class", communicate.ClassTimetableHandle)
	http.HandleFunc("/api/timetable/change", communicate.ChangeTimetableHandle)
	http.HandleFunc("/api/class", communicate.Class_structure)
	log.Print("start")
	log.Fatal(http.ListenAndServe(":54321", nil))
}
