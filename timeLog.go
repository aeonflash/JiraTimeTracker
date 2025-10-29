package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type TimeLogEntry struct {
	JiraID    string    `json:"jiraId"`
	Summary   string    `json:"summary"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	Duration  string    `json:"duration"`
	Comment   string    `json:"comment"`
	LoggedAt  time.Time `json:"loggedAt"`
}

func saveTimeLogEntry(entry TimeLogEntry) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	
	logFile := filepath.Join(homeDir, ".jira_time_log.json")
	
	// Read existing entries
	var entries []TimeLogEntry
	if data, err := os.ReadFile(logFile); err == nil {
		json.Unmarshal(data, &entries)
	}
	
	// Add new entry
	entries = append(entries, entry)
	
	// Write back to file
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(logFile, data, 0644)
}

func getTodaysTimeLog() ([]TimeLogEntry, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	
	logFile := filepath.Join(homeDir, ".jira_time_log.json")
	
	var entries []TimeLogEntry
	if data, err := os.ReadFile(logFile); err == nil {
		json.Unmarshal(data, &entries)
	}
	
	// Filter for today's entries
	today := time.Now().Format("2006-01-02")
	var todaysEntries []TimeLogEntry
	
	for _, entry := range entries {
		if entry.LoggedAt.Format("2006-01-02") == today {
			todaysEntries = append(todaysEntries, entry)
		}
	}
	
	return todaysEntries, nil
}

func formatDurationForJira(duration time.Duration) string {
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	
	if hours > 0 && minutes > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	} else if hours > 0 {
		return fmt.Sprintf("%dh", hours)
	} else if minutes > 0 {
		return fmt.Sprintf("%dm", minutes)
	} else {
		return "1m" // Minimum 1 minute
	}
}