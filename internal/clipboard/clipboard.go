package clipboard

import (
	"context"
	"fmt"
	"time"

	"github.com/atotto/clipboard"
)

func GetClipboardContent() (string, error) {
	content, err := clipboard.ReadAll()
	if err != nil {
		return "", err
	}
	return content, nil
}

func SetClipboardContent(content string) error {
	return clipboard.WriteAll(content)
}

var lastContent string

func MonitorClipboard(ctx context.Context, onCopy func(content string)) {
	for {
		select {
		case <-ctx.Done():
			// Stop monitoring when the context is canceled
			fmt.Println("Stopping clipboard monitoring.")
			return
		default:
			// Get the current clipboard content
			content, err := clipboard.ReadAll()
			if err != nil {
				fmt.Println("Error reading clipboard:", err)
				time.Sleep(1 * time.Second)
				continue
			}

			// Check if the content has changed
			if content != lastContent {
				lastContent = content
				onCopy(content) // Trigger the callback with the new content
			}

			time.Sleep(500 * time.Millisecond) // Poll every 500ms
		}
	}
}
