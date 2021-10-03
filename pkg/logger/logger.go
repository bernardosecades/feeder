package logger

import (
	"log"
	"os"
)

type Logger interface {
	Log(v ...interface{})
}

type fileLogger struct {
	logger *log.Logger
}

// NewFileLogger create new instance of Logger with file handler
func NewFileLogger(fileName string) Logger {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	return &fileLogger{ logger: log.New(file, "INFO: ", log.Ldate|log.Ltime)}
}

// Log print any value into file
func (l * fileLogger) Log(v ...interface{}) {
	l.logger.Print(v...)
}
