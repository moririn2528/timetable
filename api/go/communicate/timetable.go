package communicate

import (
	"encoding/json"
	"io"
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
	teacher_id_str, ok := req.Form["teacher_id"]
	if !ok || len(teacher_id_str) != 1 {
		return errors.NewError(http.StatusBadRequest, "no teacher_id")
	}
	teacher_id, err := strconv.Atoi(teacher_id_str[0])
	if err != nil {
		return errors.ErrorWrap(err, http.StatusBadRequest, "input teacher_id parse error")
	}
	duration_id_str, ok := req.Form["duration_id"]
	if !ok || len(duration_id_str) != 1 {
		return errors.NewError(http.StatusBadRequest, "no duration id")
	}
	duration_id, err := strconv.Atoi(duration_id_str[0])
	if err != nil {
		return errors.ErrorWrap(err, http.StatusBadRequest, "input duration id parse error")
	}
	ban_units_str, ok := req.Form["ban_units"]
	if !ok {
		return errors.NewError(http.StatusBadRequest, "no ban_units")
	}
	var ban_units []usecase.BanUnit
	for _, bs := range ban_units_str {
		units_list := strings.Split(bs, "A")
		if len(units_list)%2 != 0 {
			return errors.NewError(http.StatusBadRequest, "ban_units error")
		}
		for i := 0; i < len(units_list); i += 2 {
			d, err := time.ParseInLocation("2006-01-02", units_list[i], time.Local)
			if err != nil && d.Weekday() == 0 {
				return errors.ErrorWrap(err, http.StatusBadRequest, "input ban_units parse date error")
			}
			f, err := strconv.Atoi(units_list[i+1])
			if err != nil || f < 0 || usecase.PERIOD <= f {
				return errors.ErrorWrap(err, http.StatusBadRequest, "input ban_units parse frame id error")
			}
			f += int(d.Weekday()-1) * usecase.PERIOD
			ban_units = append(ban_units, usecase.BanUnit{
				Day:     d,
				FrameId: f,
			})
		}
	}
	changes, err := usecase.ChangeTimetable(duration_id, date, teacher_id, ban_units)
	if err != nil {
		return errors.ErrorWrap(err)
	}
	err = ResponseJson(w, changes)
	if err != nil {
		return errors.ErrorWrap(err)
	}
	return nil
}

func postChangeTimetable(w http.ResponseWriter, req *http.Request) error {
	if req.Header.Get("Content-Type") != "application/json" {
		return errors.NewError(http.StatusBadRequest, "content type")
	}
	leng, err := strconv.Atoi(req.Header.Get("Content-Length"))
	if err != nil {
		return errors.ErrorWrap(err, http.StatusBadRequest, "content length")
	}
	body := make([]byte, leng)
	leng, err = req.Body.Read(body)
	if err != nil && err != io.EOF {
		return errors.ErrorWrap(err, http.StatusBadRequest, "read body error")
	}
	var move []usecase.TimetableMove
	err = json.Unmarshal(body[:leng], &move)
	if err != nil {
		return errors.ErrorWrap(err, http.StatusBadRequest, "json parse error")
	}

	err = usecase.MoveTimetable(move)
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
	case "POST":
		err = postChangeTimetable(w, req)
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

func getTeacherTimetable(w http.ResponseWriter, req *http.Request) error {
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

	timetable, err := usecase.GetTimetableByTeacher(teacher_id, date)
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

func TeacherTimetableHandle(w http.ResponseWriter, req *http.Request) {
	var err error
	switch req.Method {
	case "GET":
		err = getTeacherTimetable(w, req)
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
