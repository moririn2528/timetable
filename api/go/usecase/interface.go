package usecase

import (
	"time"
	"timetable/library/logging"
)

type Frame struct {
	Id      int
	DayWeek int
	Period  int
}

type DatabaseClass interface {
	GetClassGraph() (*ClassGraph, error)
	InsertClassroom(classes []ClassNode) error
	InsertClassEdge(edges []ClassEdge) error
}

type DatabaseTimetable interface {
	GetNomalTimetable(duration_id int, class_ids []int, teacher_id int) ([]NormalTimetable, error)
	GetTimetable(
		duration_id int, class_ids []int, teacher_id int,
		start_day time.Time, end_day time.Time,
	) ([]Timetable, error)
	MoveTimetable(move []TimetableMove) error
}

type DatabaseAny interface {
	GetDurationId(date time.Time) (int, error)
	GetTeacher() ([]Teacher, error)
	GetPlace() ([]Place, error)
	GetHolidays() ([]time.Time, error)
	FindUser(user User, password string) error
	InsertUser(user User, password string) error
	GetTeacherAvoid(id int, date time.Time, end_date time.Time) ([]TeacherAvoidRes, error)
	SetTeacherAvoid(id int, avoids []ChangingTeacherAvoid) error
	UpdateTeacher(teacher Teacher) error
}

type SolverClass interface {
	TimetableChange(
		tt_all []Timetable,
		graph ClassGraph,
		change_teacher Teacher,
		ban_units []BanUnit,
		places []Place,
		teachers []Teacher,
		start_day time.Time,
		holidays []time.Time,
		teacher_relax int,
	) ([]TimetableMove, int, error)
}

var (
	Db_class      DatabaseClass
	Db_timetabale DatabaseTimetable
	Db_any        DatabaseAny
	Solver        SolverClass
	logger        *logging.Logger = logging.NewLogger()
)
