package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type LogEntry struct {
	Timestamp      string `json:"timestamp"`
	ServiceName    string `json:"service_name"`
	StatusCode     int    `json:"status_code"`
	ResponseTimeMs int    `json:"response_time_ms"`
	UserID         string `json:"user_id"`
	TransactionID  string `json:"transaction_id"`
	AdditionalInfo string `json:"additional_info"`
}

func main() {
	err := parseLogFile("sample.log", "generated.log")
	if err != nil {
		fmt.Println("Error parsing log file:", err)
	}
}

func parseLogFile(inputFile, outputFile string) error {
	file, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	output, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer output.Close()

	writer := bufio.NewWriter(output)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		logEntry, err := parseLogEntry(line)
		if err != nil {
			fmt.Println("Error parsing log entry:", err)
			continue
		}

		jsonData, err := json.Marshal(logEntry)
		if err != nil {
			fmt.Println("Error converting log entry to JSON:", err)
			continue
		}
		_, err = writer.WriteString(string(jsonData) + "\n")
		if err != nil {
			fmt.Println("Error writing log entry to file:", err)
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	writer.Flush()
	return nil
}

func parseLogEntry(line string) (LogEntry, error) {
	re := regexp.MustCompile(`^(.*?) (.*?) (.*?) (\d{3}) (\d+)ms (.*?) (.*?) (.*)$`)
	matches := re.FindStringSubmatch(line)
	if len(matches) != 9 {
		return LogEntry{}, fmt.Errorf("invalid log entry format")
	}

	statusCode, err := strconv.Atoi(matches[4])
	if err != nil {
		return LogEntry{}, fmt.Errorf("invalid status code")
	}

	responseTimeMs, err := strconv.Atoi(matches[5])
	if err != nil {
		return LogEntry{}, fmt.Errorf("invalid response time")
	}

	return LogEntry{
		Timestamp:      matches[1] + " " + matches[2],
		ServiceName:    matches[3],
		StatusCode:     statusCode,
		ResponseTimeMs: responseTimeMs,
		UserID:         matches[6],
		TransactionID:  matches[7],
		AdditionalInfo: matches[8],
	}, nil
}
