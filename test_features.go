package main

import (
	"fmt"

	"aliyun-tui-viewer/internal/ui"
)

func TestYankTracker() {
	tracker := ui.NewYankTracker()

	// First y should not trigger double-y
	if tracker.HandleYankKey() {
		fmt.Println("ERROR: First y should not trigger double-y")
		return
	}

	// Second y should trigger double-y
	if !tracker.HandleYankKey() {
		fmt.Println("ERROR: Second y should trigger double-y")
		return
	}

	fmt.Println("SUCCESS: YankTracker works correctly")
}

func TestCopyToClipboard() {
	testData := map[string]interface{}{
		"test":   "data",
		"number": 123,
	}

	err := ui.CopyToClipboard(testData)
	if err != nil {
		fmt.Printf("ERROR: CopyToClipboard failed: %v\n", err)
		return
	}

	fmt.Println("SUCCESS: CopyToClipboard works correctly")
}

func main() {
	TestYankTracker()
	TestCopyToClipboard()
}
