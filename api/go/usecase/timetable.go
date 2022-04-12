package usecase

import (
	"time"
	"timetable/errors"
)

type NormalTimetable struct {
	Id           int    `json:"id"`
	ClassId      int    `json:"class_id"`
	ClassName    string `json:"class_name"`
	DurationId   int    `json:"duration_id"`
	DurationName string `json:"duration_name"`
	FrameId      int    `json:"frame_id"`
	FrameDayWeek int    `json:"frame_day_week"`
	FramePeriod  int    `json:"frame_period"`
	SubjectId    int    `json:"subject_id"`
	SubjectName  string `json:"subject_name"`
	TeacherId    int    `json:"teacher_id"`
	TeacherName  string `json:"teacher_name"`
	PlaceId      int    `json:"place_id"`
}

type Timetable struct {
	NormalTimetable
	Day time.Time `json:"day"`
}

func GetTimetableByClass() {

}

func ChangeTimetable(duration_id int, search_day time.Time, change_id int) ([]Timetable, error) {
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
	var change_unit Timetable
	flag := false
	for _, t := range tt {
		if t.Id == change_id {
			flag = true
			change_unit = t
		}
	}
	if !flag {
		return nil, errors.NewError(400, "change id can't match")
	}
	places, err := Db_any.GetPlace()
	if err != nil {
		return nil, errors.ErrorWrap(err)
	}
	teachers, err := Db_any.GetTeacher()
	if err != nil {
		return nil, errors.ErrorWrap(err)
	}
	holidays, err := Db_any.GetHolidays()
	if err != nil {
		return nil, errors.ErrorWrap(err)
	}

	changes, _, err := Solver.TimetableChange(tt, *class, &change_unit, places, teachers, search_day, holidays)
	if err != nil {
		return nil, errors.ErrorWrap(err)
	}
	return changes, nil
}
