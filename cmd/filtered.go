package main

import (
	"bufio"
	"encoding/json"
	"os"
)

type LogEntry struct {
	User     string `json:"user"`
	Action   string `json:"action"`
	Resource string `json:"resource"`
}

func getFilteredLogs(user, action, resource string) ([]LogEntry, error) {
	file, err := os.Open(conf.Log_file)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var filteredLogs []LogEntry
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		var logEntry LogEntry

		if err := json.Unmarshal([]byte(line), &logEntry); err != nil {
			continue
		}

		// Фильтрация
		if user != "" && logEntry.User != user {
			continue
		}

		if action != "" && logEntry.Action != action {
			continue
		}

		if resource != "" && logEntry.Resource != resource {
			continue
		}

		filteredLogs = append(filteredLogs, logEntry)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return filteredLogs, nil
}
