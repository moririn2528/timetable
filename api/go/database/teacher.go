package database

import (
	"strconv"
	"strings"
	"time"
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
		Name:  t.Name,
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

type AddTeacherAvoid struct {
	TeacherId int       `db:"teacher_id"`
	Date      time.Time `db:"date"`
	FrameId   int       `db:"frame_id"`
	Avoid     int       `db:"avoid"`
}

func (base *DatabaseAny) GetTeacherAvoid(id int, date time.Time, end_date time.Time) ([]usecase.TeacherAvoidRes, error) {
	teach, err := base.GetTeacher()
	if err != nil {
		return nil, errors.ErrorWrap(err)
	}
	var res []usecase.TeacherAvoidRes
	date_index := map[time.Time]int{}
	if date.After(end_date) {
		return nil, errors.NewError(400, "input error, start date is after end date")
	}
	for _, t := range teach {
		if t.Id != id {
			continue
		}
		if len(res) != 0 {
			return nil, errors.NewError(500, "sql data pollution, teacher id is not unique")
		}
		for d := date; d.Before(end_date); d = d.AddDate(0, 0, 1) {
			dayweek := int(d.Weekday())
			if dayweek == 0 {
				continue
			}
			date_index[d] = len(res)
			av := usecase.TeacherAvoidRes{
				Day:   d,
				Avoid: make([]int, usecase.PERIOD),
			}
			for i := 0; i < usecase.PERIOD; i++ {
				fid := (dayweek-1)*usecase.PERIOD + i
				if fid < len(t.Avoid) {
					av.Avoid[i] = t.Avoid[fid]
				}
			}
			res = append(res, av)
		}
	}
	if len(res) == 0 {
		return nil, errors.NewError(400, "input error, teacher id not found")
	}

	var adds []AddTeacherAvoid
	err = db.Select(&adds, strings.Join([]string{
		"SELECT teacher_id, date, frame_id, avoid FROM additional_teacher_avoid",
		"WHERE teacher_id = ? AND date >= ? AND date < ?",
	}, " "), id, date, end_date)
	if err != nil {
		return nil, errors.ErrorWrap(err)
	}
	for _, av := range adds {
		p := av.FrameId % usecase.PERIOD
		if !isDate(av.Date) {
			logger.Warning("day in additional_teacher_avoid is not date", av)
		}
		d, ok := date_index[av.Date]
		if !ok {
			logger.Warning("day in additional_teacher_avoid is out of range", av)
			continue
		}
		res[d].Avoid[p] = av.Avoid
	}
	return res, nil
}

func (*DatabaseAny) SetTeacherAvoid(id int, avoids []usecase.ChangingTeacherAvoid) error {
	var ins []AddTeacherAvoid
	for _, av := range avoids {
		ins = append(ins, AddTeacherAvoid{
			TeacherId: id,
			Date:      av.Date,
			FrameId:   av.Period,
			Avoid:     av.Avoid,
		})
	}
	_, err := db.NamedExec(strings.Join([]string{
		"INSERT INTO additional_teacher_avoid (teacher_id, date, frame_id, avoid)",
		"VALUES (:teacher_id, :date, :frame_id, :avoid)",
		"ON DUPLICATE KEY UPDATE avoid = :avoid",
	}, " "), ins)
	if err != nil {
		return errors.ErrorWrap(err)
	}
	return nil
}

func (*DatabaseAny) UpdateTeacher(teacher usecase.Teacher) error {
	avoid := ""
	for _, a := range teacher.Avoid {
		if a < 0 || 9 < a {
			logger.Error("avoid is out of range", teacher)
			return errors.NewError(400, "input error, avoid is out of range")
		}
		avoid += strconv.Itoa(a)
	}
	_, err := db.Exec("UPDATE teacher SET name = ?, avoid = ? WHERE id = ?", teacher.Name, avoid, teacher.Id)
	if err != nil {
		return errors.ErrorWrap(err)
	}
	return nil
}
