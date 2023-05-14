package logger

import (
	"log"
	"os"
)

var logDebug = log.New(os.Stdout, "[DEBUG] ", log.LstdFlags|log.Lmsgprefix)
var logError = log.New(os.Stdout, "[ERROR] ", log.LstdFlags|log.Lmsgprefix)

func Debug(v ...any) {
	logDebug.Println(v)
}

func Error(str string, err error) {
	logError.Printf("%s\t%s\n", str, err)
}
