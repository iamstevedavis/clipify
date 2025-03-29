package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/getlantern/systray"
	"github.com/iamstevedavis/clipify/internal/clipboard"
	"github.com/iamstevedavis/clipify/internal/history"
	"github.com/iamstevedavis/clipify/internal/tray"
	"github.com/iamstevedavis/clipify/internal/ui"
)

func getPaths() (dataDir string, assetsDir string, err error) {
	// Get the current working directory
	baseDir, err := os.Getwd()
	if err != nil {
		return "", "", err
	}

	// Construct the paths
	dataDir = filepath.Join(baseDir, "runtime")
	assetsDir = filepath.Join(baseDir, "assets\\icons")

	// Ensure the directories exist
	if err := os.MkdirAll(dataDir, os.ModePerm); err != nil {
		return "", "", err
	}
	if err := os.MkdirAll(assetsDir, os.ModePerm); err != nil {
		return "", "", err
	}

	return dataDir, assetsDir, nil
}

func main() {
	// Get paths for data and assets directories
	dataDir, assetsDir, err := getPaths()
	if err != nil {
		panic("Failed to get paths: " + err.Error())
	}

	// Create a context for managing goroutines
	ctx, cancel := context.WithCancel(context.Background())

	// Channel to signal "View History" events
	viewHistoryCh := make(chan struct{})

	// Start the tray with the event channel and assets directory
	go tray.StartTray(viewHistoryCh, assetsDir)

	// Start monitoring the clipboard
	go clipboard.MonitorClipboard(ctx, func(content string) {
		fmt.Println("New clipboard content:", content)

		// Save the new content to history
		historyItem := history.HistoryItem{
			Content: content,
			Time:    time.Now().Format(time.RFC3339),
		}

		var historyList []history.HistoryItem

		historyList, err = history.LoadHistory(dataDir)
		if err != nil {
			fmt.Println("Error loading history:", err)
			return
		}

		historyList = append(historyList, historyItem)
		err = history.SaveHistory(historyList, dataDir)
		if err != nil {
			fmt.Println("Error saving history:", err)
		}
	})

	// Listen for "View History" events
	go func() {
		for range viewHistoryCh {
			// Load the clipboard history
			historyList, err := history.LoadHistory(dataDir)
			if err != nil {
				fmt.Println("Error loading history:", err)
				continue
			}

			// Open the window and display the history
			// Convert historyList to a slice of ui.HistoryItem
			uiHistoryList := make([]ui.HistoryItem, len(historyList))
			for i, item := range historyList {
				uiHistoryList[i] = ui.HistoryItem{
					Content: item.Content,
					Time:    item.Time,
				}
			}
			ui.NewWindow(uiHistoryList)
		}
	}()

	// Handle application termination
	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, os.Interrupt, syscall.SIGTERM)

	// Wait for termination signal
	<-quitCh

	// Cancel the context to stop all goroutines
	cancel()

	// Clean up the tray icon
	systray.Quit()

	fmt.Println("Application exited cleanly.")
}
