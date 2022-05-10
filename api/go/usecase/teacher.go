package usecase

import (
	"time"
	"timetable/errors"
)

type Teacher struct {
	Id    int
	Avoid []int
}

type TeacherAvoidRes struct {
	Day   time.Time `json:"day"`
	Avoid []int     `json:"avoid"` //index: period
}

func GetTeacherAvoid(id int, date time.Time, end_date time.Time) ([]TeacherAvoidRes, error) {
	teach, err := Db_any.GetTeacher()
	if err != nil {
		return nil, errors.ErrorWrap(err)
	}
	for _, t := range teach {
		if t.Id != id {
			continue
		}
		var res []TeacherAvoidRes
		for d := date; d.Before(end_date); d = d.AddDate(0, 0, 1) {
			dayweek := int(d.Weekday())
			if dayweek == 0 {
				continue
			}
			av := TeacherAvoidRes{
				Day:   d,
				Avoid: make([]int, PERIOD),
			}
			for i := 0; i < PERIOD; i++ {
				fid := (dayweek-1)*PERIOD + i
				if fid < len(t.Avoid) {
					av.Avoid[i] = t.Avoid[fid]
				}
			}
			res = append(res, av)
		}
		return res, nil
	}
	return nil, errors.NewError(400, "input error")
}
