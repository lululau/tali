package ui

import (
	"aliyun-tui-viewer/internal/config"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/rivo/tview"
)

// CopyToClipboard copies the given data as JSON to the system clipboard
func CopyToClipboard(data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data to JSON: %w", err)
	}

	err = clipboard.WriteAll(string(jsonData))
	if err != nil {
		return fmt.Errorf("failed to copy to clipboard: %w", err)
	}

	return nil
}

// OpenInNvim opens the given data as JSON in nvim
func OpenInNvim(data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data to JSON: %w", err)
	}

	// Create a temporary file
	tmpDir := os.TempDir()
	tmpFile := filepath.Join(tmpDir, fmt.Sprintf("tali_detail_%d.json", time.Now().Unix()))

	err = os.WriteFile(tmpFile, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write temporary file: %w", err)
	}

	// We need to suspend the tview application to avoid terminal conflicts
	// This will be handled by the calling application

	// Open in nvim with proper terminal handling
	cmd := exec.Command("nvim", tmpFile)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		// Clean up temp file
		os.Remove(tmpFile)
		return fmt.Errorf("failed to open nvim: %w", err)
	}

	// Clean up temp file after nvim closes
	os.Remove(tmpFile)
	return nil
}

// OpenInNvimWithSuspend opens the given data as JSON in nvim with proper tview suspension
func OpenInNvimWithSuspend(data interface{}, app *tview.Application) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data to JSON: %w", err)
	}

	// Create a temporary file
	tmpDir := os.TempDir()
	tmpFile := filepath.Join(tmpDir, fmt.Sprintf("tali_detail_%d.json", time.Now().Unix()))

	err = os.WriteFile(tmpFile, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write temporary file: %w", err)
	}

	// Suspend the tview application to release terminal control
	app.Suspend(func() {
		// Open in nvim with proper terminal handling
		cmd := exec.Command("nvim", tmpFile)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		// Run nvim and wait for it to complete
		cmd.Run()
	})

	// Clean up temp file after nvim closes and app resumes
	os.Remove(tmpFile)
	return nil
}

// YankTracker tracks consecutive 'y' key presses for double-y functionality
type YankTracker struct {
	lastYankTime time.Time
	yankCount    int
}

// NewYankTracker creates a new yank tracker
func NewYankTracker() *YankTracker {
	return &YankTracker{}
}

// HandleYankKey handles 'y' key press and returns true if it's a double-y
func (yt *YankTracker) HandleYankKey() bool {
	now := time.Now()

	// If more than 500ms since last y, reset counter
	if now.Sub(yt.lastYankTime) > 500*time.Millisecond {
		yt.yankCount = 1
	} else {
		yt.yankCount++
	}

	yt.lastYankTime = now

	// Return true if this is the second y in quick succession
	return yt.yankCount == 2
}

// OpenInEditor opens the given data as JSON in the configured editor
func OpenInEditor(data interface{}, app *tview.Application) error {
	editorCmd, err := config.GetEditor()
	if err != nil {
		return fmt.Errorf("failed to get editor command: %w", err)
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data to JSON: %w", err)
	}

	// Create a temporary file
	tmpDir := os.TempDir()
	tmpFile := filepath.Join(tmpDir, fmt.Sprintf("tali_detail_%d.json", time.Now().Unix()))

	err = os.WriteFile(tmpFile, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write temporary file: %w", err)
	}

	// Suspend the tview application to release terminal control
	app.Suspend(func() {
		// Parse editor command (might have arguments)
		cmdParts := strings.Fields(editorCmd)
		if len(cmdParts) == 0 {
			return
		}

		// Add the temporary file as the last argument
		cmdParts = append(cmdParts, tmpFile)

		// Open in editor with proper terminal handling
		cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		// Run editor and wait for it to complete
		cmd.Run()
	})

	// Clean up temp file after editor closes and app resumes
	os.Remove(tmpFile)
	return nil
}

// OpenInPager opens the given data as JSON in the configured pager
func OpenInPager(data interface{}, app *tview.Application) error {
	pagerCmd, err := config.GetPager()
	if err != nil {
		return fmt.Errorf("failed to get pager command: %w", err)
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data to JSON: %w", err)
	}

	// Create a temporary file
	tmpDir := os.TempDir()
	tmpFile := filepath.Join(tmpDir, fmt.Sprintf("tali_detail_%d.json", time.Now().Unix()))

	err = os.WriteFile(tmpFile, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write temporary file: %w", err)
	}

	// Suspend the tview application to release terminal control
	app.Suspend(func() {
		// Parse pager command (might have arguments)
		cmdParts := strings.Fields(pagerCmd)
		if len(cmdParts) == 0 {
			return
		}

		// Add the temporary file as the last argument
		cmdParts = append(cmdParts, tmpFile)

		// Open in pager with proper terminal handling
		cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		// Run pager and wait for it to complete
		cmd.Run()
	})

	// Clean up temp file after pager closes and app resumes
	os.Remove(tmpFile)
	return nil
}
