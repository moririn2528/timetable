package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"timetable/communicate"
	"timetable/database"
	"timetable/library/logging"
	"timetable/solve"
	"timetable/usecase"

	"github.com/joho/godotenv"
)

var (
	logger *logging.Logger = logging.NewLogger()
)

func init() {
	const location = "Asia/Tokyo"
	dev, ok := os.LookupEnv("EXEC_ENV")
	if !(ok && dev == "docker") {
		f, err := os.Create("1.log")
		if err != nil {
			panic(err)
		}
		log.SetOutput(f)
	}
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	loc, err := time.LoadLocation(location)
	if err != nil {
		loc = time.FixedZone(location, 9*60*60)
	}
	time.Local = loc

	_, err = os.Stat(".env")
	if !os.IsNotExist(err) {
		err = godotenv.Load(".env")
		if err != nil {
			log.Fatalln(err)
		}
	}
	database.Init()
}

// func test() {
// 	_, err := usecase.ChangeTimetable(0, time.Date(2021, 4, 1, 0, 0, 0, 0, time.Local), 3210405)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	logger.Info("ok!")
// 	os.Exit(1)
// }

func main() {
	defer database.Close()

	usecase.Db_any = &database.DatabaseAny{}
	usecase.Db_class = &database.DatabaseClass{}
	usecase.Db_timetabale = &database.DatabaseTimetable{}
	usecase.Solver = &solve.SolverClass{}

	// parse args for internal command
	f := flag.String("mode", "normal", "実行モード")
	flag.Parse()

	if *f == "init-holiday" {
		database.SetHoliday()
		return
	}
	///////////////////////////

	// test()

	http.Handle("/", http.FileServer(http.Dir("../../front")))
	http.HandleFunc("/api/login", communicate.LoginHandle)
	http.HandleFunc("/api/timetable/class", communicate.ClassTimetableHandle)
	http.HandleFunc("/api/timetable/teacher", communicate.TeacherTimetableHandle)
	http.HandleFunc("/api/timetable/change", communicate.ChangeTimetableHandle)
	http.HandleFunc("/api/class", communicate.Class_structure)
	http.HandleFunc("/api/teacher/avoid", communicate.TeacherAvoidHandle)
	http.HandleFunc("/api/teacher", communicate.TeacherHandle)
	logger.Info("start")
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "80"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
