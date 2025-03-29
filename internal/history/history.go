package history

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type HistoryItem struct {
	Content string `json:"content"`
	Time    string `json:"time"`
}

func SaveHistory(history []HistoryItem, dataDir string) error {
	// Save the history to the JSON file
	dataFilePath := filepath.Join(dataDir, "history.json")
	data, err := json.Marshal(history)
	if err != nil {
		return err
	}
	return os.WriteFile(dataFilePath, data, 0644)
}

func LoadHistory(dataDir string) ([]HistoryItem, error) {
	// Read the history from the JSON file
	dataFilePath := filepath.Join(dataDir, "history.json")

	// Check if the file exists
	if _, err := os.Stat(dataFilePath); os.IsNotExist(err) {
		// Create an empty history.json file
		emptyHistory := []HistoryItem{}
		data, err := json.Marshal(emptyHistory)
		if err != nil {
			return nil, err
		}
		if err := os.WriteFile(dataFilePath, data, 0644); err != nil {
			return nil, err
		}
		// Return an empty history
		return emptyHistory, nil
	}

	// Read the file
	data, err := os.ReadFile(dataFilePath)
	if err != nil {
		return nil, err
	}

	// Handle empty or invalid JSON
	if len(data) == 0 || string(data) == "{}" {
		return []HistoryItem{}, nil // Return an empty slice
	}

	var history []HistoryItem
	err = json.Unmarshal(data, &history)
	if err != nil {
		return nil, err
	}

	// Reverse the history slice to show the most recent first
	for i, j := 0, len(history)-1; i < j; i, j = i+1, j-1 {
		history[i], history[j] = history[j], history[i]
	}

	return history, nil
}
