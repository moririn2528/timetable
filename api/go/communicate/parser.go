package communicate

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"timetable/errors"
)

// json 入力の POST で、入力されたものを取得する
// value: json を入れるもののアドレス
func parserPostJson(req *http.Request, value interface{}) error {
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
	err = json.Unmarshal(body[:leng], value)
	if err != nil {
		return errors.ErrorWrap(err, http.StatusBadRequest, "json parse error")
	}
	return nil
}
