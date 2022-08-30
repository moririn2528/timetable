package database

import (
	"timetable/errors"
	"timetable/usecase"
)

func (*DatabaseAny) FindUser(user usecase.User, password string) error {
	var cnt int
	err := db.Get(&cnt, "SELECT COUNT(*) FROM user WHERE id=? AND name=? AND password=?", user.Id, user.Name, password)
	if err != nil {
		return errors.ErrorWrap(err)
	}
	if cnt == 1 {
		return nil
	} else {
		return errors.NewError("database error, user count is not one", cnt)
	}
}

func (*DatabaseAny) InsertUser(user usecase.User, password string) error {
	res, err := db.Exec("INSERT INTO user(id,name,password) VALUE (?,?,?)", user.Id, user.Name, password)
	if err != nil {
		return errors.ErrorWrap(err)
	}
	cnt, err := res.RowsAffected()
	if err != nil {
		return errors.ErrorWrap(err)
	}
	if cnt != 1 {
		return errors.NewError("database error, cannot insert user")
	}
	return nil
}
