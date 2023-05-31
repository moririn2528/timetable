package communicate

import (
	"net/http"

	"timetable/errors"
	"timetable/usecase"
)

func login(w http.ResponseWriter, req *http.Request) error {
	type User struct {
		usecase.User
		Password string `json:"password"`
	}
	var user User
	err := parserPostJson(req, &user)
	if err != nil {
		return errors.ErrorWrap(err)
	}
	token, err := usecase.Login(usecase.User{
		Id:   user.Id,
		Name: user.Name,
	}, user.Password)
	if err != nil {
		return errors.ErrorWrap(err)
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(token))
	if err != nil {
		return errors.ErrorWrap(err)
	}
	return nil
}

func LoginHandle(w http.ResponseWriter, req *http.Request) {
	var err error
	switch req.Method {
	case "POST":
		err = login(w, req)
	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err == nil {
		return
	}
	my_err, ok := err.(*errors.MyError)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error("wrap error")
		return
	}
	w.WriteHeader(my_err.GetCode())
	logger.Error(my_err.Error())
}
