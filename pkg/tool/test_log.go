package tool

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
)

type TestLog struct {
	Name    string      `json:"name" yaml:"name"`
	Entries []*LogEntry `json:"log" yaml:"log"`
}

type LogEntry struct {
	Entry string `json:"Entry" yaml:"Entry"`
}

var testlog TestLog

func init() {
	testlog = TestLog{Name: "Chart Verifier Log"}
}

func LogWarning(message string) {
	warning_log_entry := LogEntry{Entry: fmt.Sprintf("[WARNING] %s", message)}
	testlog.Entries = append(testlog.Entries, &warning_log_entry)
}

func LogInfo(message string) {
	info_log_entry := LogEntry{Entry: fmt.Sprintf("[INFO] %s", message)}
	testlog.Entries = append(testlog.Entries, &info_log_entry)
}

func LogError(message string) {
	error_log_entry := LogEntry{Entry: fmt.Sprintf("[ERROR} %s", message)}
	testlog.Entries = append(testlog.Entries, &error_log_entry)
}
func GetLogsOutput(log_format string) (string, error) {

	if len(testlog.Entries) > 0 {
		if log_format == "json" {
			b, err := json.Marshal(&testlog)
			if err != nil {
				return "", err
			}
			return string(b), nil
		} else {
			b, err := yaml.Marshal(&testlog)
			if err != nil {
				return "", err
			}
			return string(b), nil
		}
	}
	return "", nil

}
