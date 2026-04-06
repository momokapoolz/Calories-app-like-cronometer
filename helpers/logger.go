package helpers

import (
	"fmt"
	"os"
	"time"
)

// log file path
const logFilePath = "./log/error.txt"

func LogError(err error) {
	if err == nil {
		return
	}

	writeLog("ERROR", "", err)
}

// internal file written function
func writeLog(level string, context string, err error) {
	// create log folder if not exists
	_ = os.MkdirAll("log", os.ModePerm)

	file, fileErr := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if fileErr != nil {
		fmt.Println("Cannot open log file:", fileErr)
		return
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Cannot close log file:", err)
		}
	}(file)

	timestamp := time.Now().Format("2006-01-02 15:04:05")

	var logMessage string

	if context != "" {
		logMessage = fmt.Sprintf("[%s] [%s] [%s] %s\n",
			timestamp,
			level,
			context,
			err.Error(),
		)
	} else {
		logMessage = fmt.Sprintf("[%s] [%s] %s\n",
			timestamp,
			level,
			err.Error(),
		)
	}

	_, _ = file.WriteString(logMessage)
}
