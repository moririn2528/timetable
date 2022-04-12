package database

import (
	"sort"
	"time"
	"timetable/errors"
)

func getHolidays() ([]time.Time, error) {
	var holidays []time.Time
	rows, err := db.Query("SELECT day FROM holiday")
	if err != nil {
		return nil, errors.ErrorWrap(err)
	}
	for rows.Next() {
		var d time.Time
		err := rows.Scan(&d)
		if err != nil {
			return nil, errors.ErrorWrap(err)
		}
		holidays = append(holidays, d)
	}
	sort.Slice(holidays, func(i, j int) bool { return holidays[i].Before(holidays[j]) })
	return holidays, nil
}
func (*DatabaseAny) GetHolidays() ([]time.Time, error) {
	return getHolidays()
}
