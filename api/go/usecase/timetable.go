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
	FrameDayWeek int      `json:"frame_day_week"`
	FramePeriod  int      `json:"frame_period"`
	SubjectId    int      `json:"subject_id"`
	SubjectName  string   `json:"subject_name"`
	TeacherIds   []int    `json:"teacher_id"`
	TeacherNames []string `json:"teacher_name"`
	PlaceId      int      `json:"place_id"`
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

func ChangeTimetable(duration_id int, search_day time.Time, change_id int) ([]TimetableMove, error) {
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

	changes, _, err := Solver.TimetableChange(tt, *class, &change_unit, places, teachers, search_day, holidays, 500)
	if err != nil {
		return nil, errors.ErrorWrap(err)
	}
	return changes, nil
}
