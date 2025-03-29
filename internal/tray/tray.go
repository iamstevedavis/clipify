package tray

import (
	"os"
	"path/filepath"

	"github.com/getlantern/systray"
)

// handleQuit handles the Quit menu item click event
func handleQuit(mQuit *systray.MenuItem) {
	<-mQuit.ClickedCh
	systray.Quit()
}

func StartTray(viewHistoryCh chan struct{}, assetsDir string) {
	systray.Run(func() {
		if err := OnReady(viewHistoryCh, assetsDir); err != nil {
			panic("Failed to initialize tray: " + err.Error())
		}
	}, nil)
}

func OnReady(viewHistoryCh chan struct{}, assetsDir string) error {
	// Load the icon
	iconFilePath := filepath.Join(assetsDir, "icon.ico")
	icon, err := os.ReadFile(iconFilePath)
	if err != nil {
		panic("Failed to load icon: " + err.Error())
	}
	systray.SetIcon(icon)
	systray.SetTooltip("Clipboard History Manager")

	// Add Quit menu item
	mQuit := systray.AddMenuItem("Quit", "Quit the app")
	go handleQuit(mQuit)

	// Add View History menu item
	mViewHistory := systray.AddMenuItem("Show History", "View clipboard history")
	go func() {
		for {
			<-mViewHistory.ClickedCh
			// Signal the event to open the history window
			viewHistoryCh <- struct{}{}
		}
	}()

	return nil
}
