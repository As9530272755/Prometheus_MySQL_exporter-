// logger/logger.go
package logger

import (
	"fmt"
	"log"
	"os"
)

var Log *log.Logger

func InitLog(filePath string, exitOnFatal bool) {
	logFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Failed to open log file: %s", err)
		Fatal("Failed to open log file: %s", err)
	}
	defer logFile.Close()
	Log = log.New(logFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Error(format string, v ...interface{}) {
	Log.Printf(format, v...)
}

func Fatal(format string, v ...interface{}) {
	Log.Fatalf(format, v...)
}
