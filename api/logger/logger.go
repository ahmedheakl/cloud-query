package logger

import (
	"fmt"
	"log"
	"os"
)

var Info *log.Logger
var Error *log.Logger
var Fatal *log.Logger

func Init() {
	logFile, err := os.OpenFile("../logs.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("logger cant open given file")
	}
	Info = log.New(logFile, "[INFO]:\t\t", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(logFile, "[ERROR]:\t", log.Ldate|log.Ltime|log.Lshortfile)
	Fatal = log.New(logFile, "[FATAL]:\t", log.Ldate|log.Ltime|log.Lshortfile)
}
