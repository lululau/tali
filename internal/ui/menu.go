package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// CreateMainMenu creates the main menu
func CreateMainMenu(
	onECS func(),
	onSecurityGroups func(),
	onDNS func(),
	onSLB func(),
	onOSS func(),
	onRDS func(),
	onRedis func(),
	onRocketMQ func(),
	onQuit func(),
) *tview.List {
	list := tview.NewList().
		AddItem("ECS Instances", "View ECS instances", 's', onECS).
		AddItem("Security Groups", "View ECS security groups", 'g', onSecurityGroups).
		AddItem("DNS Management", "View AliDNS domains and records", 'd', onDNS).
		AddItem("SLB Instances", "View SLB instances", 'b', onSLB).
		AddItem("OSS Management", "Browse OSS buckets and objects", 'o', onOSS).
		AddItem("RDS Instances", "View RDS instances", 'r', onRDS).
		AddItem("Redis Instances", "View Redis instances", 'i', onRedis).
		AddItem("RocketMQ Instances", "View RocketMQ instances", 'm', onRocketMQ).
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
