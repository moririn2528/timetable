package communicate

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"timetable/errors"
	"timetable/usecase"
)

func getClassTimetable(w http.ResponseWriter, req *http.Request) error {
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
	class_ids_str, ok := req.Form["id"]
	if !ok {
		return errors.NewError(http.StatusBadRequest, "no id")
	}
	var class_ids []int
	for _, ids_str := range class_ids_str {
		ids := strings.Split(ids_str, "_")
		for _, id_str := range ids {
			id, err := strconv.Atoi(id_str)
			if err != nil {
				return errors.ErrorWrap(err, http.StatusBadRequest)
			}
			class_ids = append(class_ids, id)
		}
	}

	timetable, err := usecase.GetTimetableByClass(class_ids, date)
	if err != nil {
		return errors.ErrorWrap(err)
	}

	// to json
	err = ResponseJson(w, timetable)
	if err != nil {
		return errors.ErrorWrap(err)
	}
	return nil
}

func ClassTimetableHandle(w http.ResponseWriter, req *http.Request) {
	var err error
	switch req.Method {
	case "GET":
		err = getClassTimetable(w, req)
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
		log.Print("wrap error")
		return
	}
	w.WriteHeader(my_err.GetCode())
	log.Print(my_err.Error())
}

func getChangeTimetable(w http.ResponseWriter, req *http.Request) error {
	var err error
	err = req.ParseForm()
	if err != nil {
		return errors.ErrorWrap(err)
	}
	day_str, ok := req.Form["day"]
	if !ok || len(day_str) != 1 {
		return errors.NewError(http.StatusBadRequest, "input day error")
	}
	date, err := time.ParseInLocation("2006-01-02", day_str[0], time.Local)
	if err != nil {
		return errors.ErrorWrap(err, http.StatusBadRequest, "input day parse error")
	}
	change_id_str, ok := req.Form["change_id"]
	if !ok || len(change_id_str) != 1 {
		return errors.NewError(http.StatusBadRequest, "no change id")
	}
	change_id, err := strconv.Atoi(change_id_str[0])
	if err != nil {
		return errors.ErrorWrap(err, http.StatusBadRequest, "input change id parse error")
	}
	duration_id_str, ok := req.Form["duration_id"]
	if !ok || len(duration_id_str) != 1 {
		return errors.NewError(http.StatusBadRequest, "no duration id")
	}
	duration_id, err := strconv.Atoi(duration_id_str[0])
	if err != nil {
		return errors.ErrorWrap(err, http.StatusBadRequest, "input duration id parse error")
	}
	changes, err := usecase.ChangeTimetable(duration_id, date, change_id)
	if err != nil {
		return errors.ErrorWrap(err)
	}
	err = ResponseJson(w, changes)
	if err != nil {
		return errors.ErrorWrap(err)
	}
	return nil
}

func ChangeTimetableHandle(w http.ResponseWriter, req *http.Request) {
	var err error
	switch req.Method {
	case "GET":
		err = getChangeTimetable(w, req)
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
		log.Print("wrap error")
		return
	}
	w.WriteHeader(my_err.GetCode())
	log.Print(my_err.Error())
}
