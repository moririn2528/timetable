package communicate

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"timetable/errors"
	"timetable/usecase"
)

func getTeacherAvoid(w http.ResponseWriter, req *http.Request) error {
	var err error
	err = req.ParseForm()
	if err != nil {
		return errors.ErrorWrap(err)
	}
	day_str, ok := req.Form["day"]
	if !ok {
		return errors.NewError(http.StatusBadRequest, "input day error")
	}
	date, err := time.ParseInLocation("2006-01-02", day_str[0], time.Local)
	if err != nil {
		return errors.ErrorWrap(err, http.StatusBadRequest, "input day parse error")
	}
	teacher_id_str, ok := req.Form["id"]
	if !ok {
		return errors.NewError(http.StatusBadRequest, "no id")
	}
	teacher_id, err := strconv.Atoi(teacher_id_str[0])
	if err != nil {
		return errors.ErrorWrap(err, http.StatusBadRequest)
	}

	avoid, err := usecase.GetTeacherAvoid(teacher_id, date, date.AddDate(0, 0, 7))
	if err != nil {
		return errors.ErrorWrap(err)
	}

	// to json
	err = ResponseJson(w, avoid)
	if err != nil {
		return errors.ErrorWrap(err)
	}
	return nil
}

type TeacherAvoidForm struct {
	Avoids    []usecase.ChangingTeacherAvoid `json:"avoids"`
	TeacherId int                            `json:"teacher_id"`
	Weekly    bool                           `json:"weekly"`
}

func setTeacherAvoid(w http.ResponseWriter, req *http.Request) error {
	var input TeacherAvoidForm
	err := json.NewDecoder(req.Body).Decode(&input)
	if err != nil {
		return errors.ErrorWrap(err)
	}
	if input.Weekly {
		err = usecase.SetTeacherWeeklyAvoid(input.TeacherId, input.Avoids)
	} else {
		err = usecase.SetTeacherAvoid(input.TeacherId, input.Avoids)
	}
	if err != nil {
		return errors.ErrorWrap(err)
	}
	return nil
}

func getTeacher(w http.ResponseWriter, req *http.Request) error {
	type Teacher struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}
	var err error
	err = req.ParseForm()
	if err != nil {
		return errors.ErrorWrap(err)
	}
	teacher_use, err := usecase.Db_any.GetTeacher()
	var teachers []Teacher
	for _, t := range teacher_use {
		teachers = append(teachers, Teacher{
			Id:   t.Id,
			Name: t.Name,
		})
	}

	if err != nil {
		return errors.ErrorWrap(err)
	}

	// to json
	err = ResponseJson(w, teachers)
	if err != nil {
		return errors.ErrorWrap(err)
	}
	return nil
}
func TeacherAvoidHandle(w http.ResponseWriter, req *http.Request) {
	var err error
	switch req.Method {
	case "GET":
		err = getTeacherAvoid(w, req)
	case "POST":
		err = setTeacherAvoid(w, req)
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

func TeacherHandle(w http.ResponseWriter, req *http.Request) {
	var err error
	switch req.Method {
	case "GET":
		err = getTeacher(w, req)
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
