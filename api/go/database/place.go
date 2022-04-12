package database

import (
	"timetable/errors"
	"timetable/usecase"
)

func (*DatabaseAny) GetPlace() ([]usecase.Place, error) {
	rows, err := db.Query("SELECT t.id, t.count FROM place AS t")
	if err != nil {
		return nil, errors.ErrorWrap(err)
	}
	var res []usecase.Place
	for rows.Next() {
		var t usecase.Place
		err := rows.Scan(&t.Id, &t.Count)
		if err != nil {
			return nil, errors.ErrorWrap(err)
		}
		res = append(res, t)
	}
	return res, nil
}
