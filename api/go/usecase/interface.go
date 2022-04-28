package usecase

import "time"

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
}

type DatabaseAny interface {
	GetDurationId(date time.Time) (int, error)
	GetTeacher() ([]Teacher, error)
	GetPlace() ([]Place, error)
	GetHolidays() ([]time.Time, error)
}

type SolverClass interface {
	TimetableChange(
		tt_all []Timetable,
		graph ClassGraph,
		change_unit *Timetable,
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

	Frames map[int]Frame
)
