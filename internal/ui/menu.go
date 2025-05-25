package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// CreateMainMenu creates the main menu
func CreateMainMenu(
	onECS func(),
	onDNS func(),
	onSLB func(),
	onOSS func(),
	onRDS func(),
	onQuit func(),
) *tview.List {
	list := tview.NewList().
		AddItem("ECS Instances", "View ECS instances", '1', onECS).
		AddItem("DNS Management", "View AliDNS domains and records", '2', onDNS).
		AddItem("SLB Instances", "View SLB instances", '3', onSLB).
		AddItem("OSS Management", "Browse OSS buckets and objects", '4', onOSS).
		AddItem("RDS Instances", "View RDS instances", '5', onRDS).
		AddItem("Quit", "Exit the application (Press 'Q')", 'Q', onQuit)

	list.SetBorder(true).SetTitle("Main Menu").SetBackgroundColor(tcell.ColorReset)

	// Set text colors with transparent background
	list.SetMainTextStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorReset))
	list.SetSecondaryTextStyle(tcell.StyleDefault.Foreground(tcell.ColorGray).Background(tcell.ColorReset))
	list.SetShortcutStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorReset))
	list.SetSelectedTextColor(tcell.ColorYellow)
	list.SetSelectedBackgroundColor(tcell.ColorReset)

	// Disable full line highlighting to make unselected items transparent
	list.SetHighlightFullLine(false)

	// Force transparent background for all text
	list.SetBackgroundColor(tcell.ColorReset)

	// Add j/k navigation support for main menu
	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'j':
			// Move down
			currentItem := list.GetCurrentItem()
			itemCount := list.GetItemCount()
			if currentItem < itemCount-1 {
				list.SetCurrentItem(currentItem + 1)
			}
			return nil
		case 'k':
			// Move up
			currentItem := list.GetCurrentItem()
			if currentItem > 0 {
				list.SetCurrentItem(currentItem - 1)
			}
			return nil
		}
		return event
	})

	return list
}
