package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ShowErrorModal creates and shows an error modal
func ShowErrorModal(pages *tview.Pages, app *tview.Application, message string, onDone func()) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK"}).
		SetBackgroundColor(tcell.ColorDefault).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			pages.RemovePage("errorModal")
			if onDone != nil {
				onDone()
			}
		})
	pages.AddPage("errorModal", modal, false, true)
	app.SetFocus(modal)
}

// ShowProfileSelectionDialog creates and shows a profile selection dialog
func ShowProfileSelectionDialog(pages *tview.Pages, app *tview.Application, profiles []string, currentProfile string, onSelect func(string), onCancel func()) {
	list := tview.NewList()

	// Add profiles to the list
	for _, profile := range profiles {
		profileName := profile // Capture for closure
		displayText := profileName
		if profileName == currentProfile {
			displayText = profileName + " (current)"
		}
		list.AddItem(displayText, "", 0, func() {
			pages.RemovePage("profileDialog")
			if onSelect != nil {
				onSelect(profileName)
			}
		})
	}

	list.SetBorder(true).
		SetTitle("Select Profile").
		SetBackgroundColor(tcell.ColorDefault)

	// Set up j/k navigation and cancel keys
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

		// Let other events (including ESC and 'q') bubble up to the flex container
		return event
	})

	// Create a flex container to center the list
	flex := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(list, 0, 2, true).
			AddItem(nil, 0, 1, false), 0, 2, true).
		AddItem(nil, 0, 1, false)

	pages.AddPage("profileDialog", flex, true, true)

	// Set up global input capture for this dialog
	originalInputCapture := app.GetInputCapture()
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Only handle events when the profile dialog is active
		if pages.HasPage("profileDialog") {
			// Handle Escape key
			if event.Key() == tcell.KeyEscape {
				pages.RemovePage("profileDialog")
				app.SetInputCapture(originalInputCapture) // Restore original capture
				if onCancel != nil {
					onCancel()
				}
				return nil
			}

			// Handle 'q' key
			if event.Rune() == 'q' {
				pages.RemovePage("profileDialog")
				app.SetInputCapture(originalInputCapture) // Restore original capture
				if onCancel != nil {
					onCancel()
				}
				return nil
			}
		}

		// Pass through to original handler or default behavior
		if originalInputCapture != nil {
			return originalInputCapture(event)
		}
		return event
	})

	app.SetFocus(list)
}
