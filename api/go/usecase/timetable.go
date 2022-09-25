package usecase

import (
	"time"
	"timetable/errors"
)

type NormalTimetable struct {
	Id           int      `json:"id"`
	ClassId      int      `json:"class_id"`
	ClassName    string   `json:"class_name"`
	DurationId   int      `json:"duration_id"`
	DurationName string   `json:"duration_name"`
	FrameId      int      `json:"frame_id"`
	SubjectId    int      `json:"subject_id"`
	SubjectName  string   `json:"subject_name"`
	TeacherIds   []int    `json:"teacher_id"`
	TeacherNames []string `json:"teacher_name"`
	PlaceId      int      `json:"place_id"`
}

type DeletedNormalTimetable struct {
	Id       int       `json:"id" db:"id"`
	NormalId int       `json:"normal_id" db:"normal_id"`
	Day      time.Time `json:"day" db:"day"`
}

type Timetable struct {
	NormalTimetable
	Day time.Time `json:"day"`
}

type TimetableMove struct {
	Unit    Timetable `json:"timetable"`
	Day     time.Time `json:"day"`
	FrameId int       `json:"frame_id"`
}

type BanUnit struct {
	Day     time.Time `json:"day"`
	FrameId int       `json:"frame_id"`
}

func GetTimetableByClass(class_ids []int, date time.Time) ([]Timetable, error) {

	graph, err := NewClassGraph()
	if err != nil {
		return nil, errors.ErrorWrap(err)
	}
	class_idxs, err := graph.Id2IndexArray(class_ids)
	if err != nil {
		return nil, errors.ErrorWrap(err)
	}
	class_idxs_tmp := graph.GetAncestors(class_idxs)
	class_ids_all := make([]int, len(class_idxs_tmp))
	for i, idx := range class_idxs_tmp {
		class_ids_all[i] = graph.Nodes[idx].Id
	}
	duration_id, err := Db_any.GetDurationId(date)
	if err != nil {
		return nil, errors.ErrorWrap(err)
	}
	timetable, err := Db_timetabale.GetTimetable(duration_id, class_ids_all, -1, date, date.AddDate(0, 0, 7))
	if err != nil {
		return nil, errors.ErrorWrap(err)
	}
	return timetable, nil
}

func GetTimetableByTeacher(teacher_id int, date time.Time) ([]Timetable, error) {
	duration_id, err := Db_any.GetDurationId(date)
	if err != nil {
		return nil, errors.ErrorWrap(err)
	}
	timetable, err := Db_timetabale.GetTimetable(duration_id, []int{}, teacher_id, date, date.AddDate(0, 0, 7))
	if err != nil {
		return nil, errors.ErrorWrap(err)
	}
	return timetable, nil
}

func ChangeTimetable(duration_id int, search_day time.Time, teacher_id int, ban_units []BanUnit) ([]TimetableMove, error) {
	tt, err := Db_timetabale.GetTimetable(
		duration_id, []int{}, -1, search_day, search_day.AddDate(0, 0, COUNT_DAY),
	)
	if err != nil {
		return nil, errors.ErrorWrap(err)
	}
	class, err := NewClassGraph()
	if err != nil {
		return nil, errors.ErrorWrap(err)
	}
	places, err := Db_any.GetPlace()
	if err != nil {
		return nil, errors.ErrorWrap(err)
	}
	teachers, err := Db_any.GetTeacher()
	if err != nil {
		return nil, errors.ErrorWrap(err)
	}
	var change_teacher Teacher
	flag := false
	for _, t := range teachers {
		if t.Id == teacher_id {
			flag = true
			change_teacher = t
		}
	}
	if !flag {
		return nil, errors.NewError(400, "change teacher id can't match")
	}
	holidays, err := Db_any.GetHolidays()
	if err != nil {
		return nil, errors.ErrorWrap(err)
	}

	changes, _, err := Solver.TimetableChange(tt, *class, change_teacher, ban_units, places, teachers, search_day, holidays, 500)
	if err != nil {
		return nil, errors.ErrorWrap(err)
	}
	return changes, nil
}

func MoveTimetable(move []TimetableMove) error {
	err := Db_timetabale.MoveTimetable(move)
	if err != nil {
		return errors.ErrorWrap(err)
	}
	return nil
}
