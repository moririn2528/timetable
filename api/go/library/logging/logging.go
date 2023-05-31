package logging

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

type Logger struct {
	debug *log.Logger
	info  *log.Logger
	std   *log.Logger
}

var (
	infoFile  *os.File
	debugFile *os.File
)

func NewLogger() *Logger {
	var err error
	if infoFile == nil {
		infoFile, err = os.OpenFile("data/log/info.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			panic(err)
		}
	}
	if debugFile == nil {
		debugFile, err = os.OpenFile("data/log/debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			panic(err)
		}
	}
	return &Logger{
		debug: log.New(debugFile, "", log.Ldate|log.Ltime),
		info:  log.New(infoFile, "", log.Ldate|log.Ltime),
		std:   log.New(os.Stdout, "", log.Ldate|log.Ltime),
	}
}

func getCallerInfo() string {
	_, file, line, _ := runtime.Caller(2)
	_, file_name := filepath.Split(file)
	return file_name + ":" + strconv.Itoa(line)
}

func (l *Logger) Debug(v ...interface{}) {
	l.debug.Println(append([]interface{}{"[ DEBUG ]", getCallerInfo()}, v...)...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.debug.Printf("[ DEBUG ] "+getCallerInfo()+" "+format, v...)
}

func (l *Logger) Info(v ...interface{}) {
	l.debug.Println(append([]interface{}{"[ INFO ]", getCallerInfo()}, v...)...)
	l.info.Println(append([]interface{}{"[ INFO ]", getCallerInfo()}, v...)...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.debug.Printf("[ INFO ] "+getCallerInfo()+" "+format, v...)
	l.info.Printf("[ INFO ] "+getCallerInfo()+" "+format, v...)
}

func (l *Logger) Warning(v ...interface{}) {
	l.debug.Println(append([]interface{}{"[ WARNING ]", getCallerInfo()}, v...)...)
	l.info.Println(append([]interface{}{"[ WARNING ]", getCallerInfo()}, v...)...)
	l.std.Println(append([]interface{}{"[ \x1b[0;33WARNING\x1b[0m ]", getCallerInfo()}, v...)...)
}

func (l *Logger) Warningf(format string, v ...interface{}) {
	l.debug.Printf("[ WARNING ] "+getCallerInfo()+" "+format, v...)
	l.info.Printf("[ WARNING ] "+getCallerInfo()+" "+format, v...)
	l.std.Printf("[ \x1b[0;33WARNING\x1b[0m ] "+getCallerInfo()+" "+format, v...)
}

func (l *Logger) Error(v ...interface{}) {
	l.debug.Println(append([]interface{}{"[ ERROR ]", getCallerInfo()}, v...)...)
	l.info.Println(append([]interface{}{"[ ERROR ]", getCallerInfo()}, v...)...)
	l.std.Println(append([]interface{}{"[ \x1b[0;31mERROR\x1b[0m ]", getCallerInfo()}, v...)...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.debug.Printf("[ ERROR ] "+getCallerInfo()+" "+format, v...)
	l.info.Printf("[ ERROR ] "+getCallerInfo()+" "+format, v...)
	l.std.Printf("[ \x1b[0;31mERROR\x1b[0m ] "+getCallerInfo()+" "+format, v...)
}
