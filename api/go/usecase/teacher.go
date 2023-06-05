package usecase

import (
	"time"
	"timetable/errors"
)

type Teacher struct {
	Id    int
	Name  string
	Avoid []int
}

type TeacherAvoidRes struct {
	Day   time.Time `json:"day"`
	Avoid []int     `json:"avoid"` //index: period
}

func GetTeacherAvoid(id int, date time.Time, end_date time.Time) ([]TeacherAvoidRes, error) {
	res, err := Db_any.GetTeacherAvoid(id, date, end_date)
	if err != nil {
		return nil, errors.ErrorWrap(err)
	}
	return res, nil
}

type ChangingTeacherAvoid struct {
	Date   time.Time `json:"date"`
	Period int       `json:"period"`
	Avoid  int       `json:"avoid"`
}

// id: teacher id
func SetTeacherWeeklyAvoid(id int, avoids []ChangingTeacherAvoid) error {
	teach, err := Db_any.GetTeacher()
	if err != nil {
		return errors.ErrorWrap(err)
	}
	for _, t := range teach {
		if t.Id != id {
			continue
		}
		for _, av := range avoids {
			dayweek := int(av.Date.Weekday())
			if dayweek == 0 {
				return errors.NewError(400, "input error")
			}
			fid := (dayweek-1)*PERIOD + av.Period
			if len(t.Avoid) <= fid {
				return errors.NewError(400, "input error")
			}
			t.Avoid[fid] = av.Avoid
		}
		err = Db_any.UpdateTeacher(t)
		if err != nil {
			return errors.ErrorWrap(err)
		}
		return nil
	}
	return errors.NewError(400, "input error")
}

func SetTeacherAvoid(id int, avoids []ChangingTeacherAvoid) error {
	err := Db_any.SetTeacherAvoid(id, avoids)
	if err != nil {
		return errors.ErrorWrap(err)
	}
	return nil
}
