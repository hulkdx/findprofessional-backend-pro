package logger

import (
	"log"
	"os"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils/config"
)

var logDebug = log.New(os.Stdout, "[DEBUG] ", log.LstdFlags|log.Lmsgprefix)
var logError = log.New(os.Stdout, "[ERROR] ", log.LstdFlags|log.Lmsgprefix)

func Debug(v ...any) {
	if config.IsDebug() {
		logDebug.Println(v...)
	}

}

func DebugF(format string, v ...any) {
	if config.IsDebug() {
		logDebug.Printf(format, v...)
	}
}

func Error(str string, err error) {
	logError.Printf("%s\t%s\n", str, err)
}
