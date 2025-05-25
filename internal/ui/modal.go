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
