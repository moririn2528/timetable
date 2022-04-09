package database

import (
	"strconv"
	"timetable/errors"
	"timetable/usecase"
)

type Teacher struct {
	Id    int
	Name  string
	Avoid string
}

func (t *Teacher) parse() (*usecase.Teacher, error) {
	s := &usecase.Teacher{
		Id:    t.Id,
		Avoid: make([]int, len(t.Avoid)),
	}
	var err error
	for i := range t.Avoid {
		s.Avoid[i], err = strconv.Atoi(string(t.Avoid[i]))
		if err != nil {
			return nil, errors.ErrorWrap(err)
		}
	}
	return s, nil
}

func (*DatabaseAny) GetTeacher() ([]usecase.Teacher, error) {
	rows, err := db.Query("SELECT t.id, t.name, t.avoid FROM teacher AS t")
	if err != nil {
		return nil, errors.ErrorWrap(err)
	}
	var res []usecase.Teacher
	for rows.Next() {
		var t Teacher
		err := rows.Scan(&t.Id, &t.Name, &t.Avoid)
		if err != nil {
			return nil, errors.ErrorWrap(err)
		}
		ut, err := t.parse()
		if err != nil {
			return nil, errors.ErrorWrap(err)
		}
		res = append(res, *ut)
	}
	return res, nil
}
