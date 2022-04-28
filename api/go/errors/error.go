package errors

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var root_path string

func init() {
	root_path, _ = os.Getwd()
}

type MyError struct {
	message string
	next    error
	place   string
	code    int
}

func GetCaller(depth int) string {
	pc, file, line, _ := runtime.Caller(depth)
	f := runtime.FuncForPC(pc)
	rel_path, _ := filepath.Rel(root_path, file) // 相対パス
	return fmt.Sprintf("%s:%d %s", filepath.ToSlash(rel_path), line, f.Name())
}

func NewError(msg ...interface{}) *MyError {
	code := -1
	if len(msg) > 0 {
		c, ok := msg[0].(int)
		if ok {
			code = c
		}
	}
	return &MyError{
		message: fmt.Sprintln(msg...),
		next:    nil,
		place:   GetCaller(2),
		code:    code,
	}
}

func ErrorWrap(err error, args ...interface{}) *MyError {
	code := -1
	if len(args) > 0 {
		c, ok := args[0].(int)
		if ok {
			code = c
		}
	}
	return &MyError{
		message: fmt.Sprintln(args...),
		next:    err,
		place:   GetCaller(2),
		code:    code,
	}
}

func (err *MyError) Error() string {
	var res []string
	if len(err.message) > 1 {
		res = append(res, err.message)
	}
	if err.next != nil {
		res = append(res, err.next.Error())
	}

	return fmt.Sprintf(
		"\n%s: %s",
		err.place, strings.Join(res, ", "),
	)
}

func (err *MyError) GetCode() int {
	code := http.StatusInternalServerError
	var e error
	for e = err; e != nil; {
		my_err, ok := e.(*MyError)
		if !ok {
			return code
		}
		if my_err.code > 0 {
			code = my_err.code
		}
		e = my_err.next
	}
	return code
}
