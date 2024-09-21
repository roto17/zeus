package logs

import (
	"fmt"
	"os"
	"time"
)

func AddLog(logType string, username string, message string) error {
	// Prepare the log entry
	log := fmt.Sprintf("[%s] %s: User: %s, Message: %s\n", time.Now().Format(time.RFC3339), logType, username, message)

	// Open or create the log file, appending the new log
	file, err := os.OpenFile("./lib/logs/logfile.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("could not open or create log file: %v", err)
	}
	defer file.Close()

	// Write the log entry to the file
	if _, err := file.WriteString(log); err != nil {
		return fmt.Errorf("could not write to log file: %v", err)
	}

	return nil
}
