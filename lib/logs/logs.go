package logs

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

func AddLog(logType string, username string, message string) error {

	// Get the caller file and line number
	_, srcFile, line, ok := runtime.Caller(1)
	if ok {
		// Prepare the log entry
		log := fmt.Sprintf("[%s] %s: User: %s, Message: %s, file: %s => line: %d\n",
			time.Now().Format(time.RFC3339), logType, username, message, srcFile, line)

		// Open or create the log file, appending the new log
		logFile, err := os.OpenFile("./lib/logs/logfile.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("could not open or create log file: %v", err)
		}
		defer logFile.Close()

		// Write the log entry to the file
		if _, err := logFile.WriteString(log); err != nil {
			return fmt.Errorf("could not write to log file: %v", err)
		}
	}

	return nil
}
