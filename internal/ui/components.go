package ui

import (
	"encoding/json"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// CreateJSONDetailView creates a generic JSON detail view
func CreateJSONDetailView(title string, data interface{}) *tview.TextView {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		jsonData = []byte(fmt.Sprintf("Error marshaling JSON: %v", err))
	}

	textView := tview.NewTextView().
		SetText(string(jsonData)).
		SetScrollable(true).
		SetWrap(false).
		SetTextStyle(tcell.StyleDefault.Background(tcell.ColorReset))
	textView.SetBorder(true).SetTitle(title).SetBackgroundColor(tcell.ColorReset)

	return textView
}

// SetupTableWithFixedWidth configures a table with full width
func SetupTableWithFixedWidth(table *tview.Table) *tview.Table {
	table.SetFixed(1, 0) // Fix header row, allow all columns to be flexible
	table.SetSelectable(true, false)
	table.SetBorder(true).SetBackgroundColor(tcell.ColorDefault)
	return table
}

// CreateTableHeaders creates table headers with expansion
func CreateTableHeaders(table *tview.Table, headers []string) {
	for c, header := range headers {
		cell := tview.NewTableCell(header).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter).SetSelectable(false)
		// Set all columns to expand proportionally
		cell.SetExpansion(1)
		table.SetCell(0, c, cell)
	}
}

// WrapTableInFlex wraps table in full-width flex container
func WrapTableInFlex(table *tview.Table) tview.Primitive {
	// Create a flex container that forces the table to use full width
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.AddItem(table, 0, 1, true)
	flex.SetBorder(false).SetBackgroundColor(tcell.ColorDefault)
	return flex
}

// SetupTableNavigation sets up j/k navigation for tables
func SetupTableNavigation(table *tview.Table, onSelect func(row, column int)) {
	table.SetSelectedFunc(func(row, column int) {
		if row > 0 && onSelect != nil {
			onSelect(row, column)
		}
	})
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		currentRow, _ := table.GetSelection()
		rowCount := table.GetRowCount()
		switch event.Rune() {
		case 'j':
			if currentRow < rowCount-1 {
				table.Select(currentRow+1, 0)
			}
			return nil
		case 'k':
			if currentRow > 1 {
				table.Select(currentRow-1, 0)
			} else if rowCount > 1 {
				table.Select(1, 0)
			}
			return nil
		}
		return event
	})
}

// CreateDetailViewWithInstructions creates a detail view with navigation instructions
func CreateDetailViewWithInstructions(detailView *tview.TextView) *tview.Flex {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	// Instructions
	instructions := tview.NewTextView().
		SetText("Press 'Esc' or 'q' to go back, 'Q' to quit").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetBackgroundColor(tcell.ColorReset)
	instructions.SetBorder(false)

	flex.AddItem(instructions, 1, 0, false)
	flex.AddItem(detailView, 0, 1, true)
	flex.SetBackgroundColor(tcell.ColorReset)

	return flex
}

// CreateModeLine creates a mode line component showing current profile and shortcuts
func CreateModeLine(profileName string) *tview.TextView {
	modeLineText := fmt.Sprintf(" Profile: %s | Press 'O' to switch profile ", profileName)

	modeLine := tview.NewTextView()
	modeLine.SetText(modeLineText)
	modeLine.SetTextAlign(tview.AlignLeft)
	modeLine.SetDynamicColors(true)
	modeLine.SetBackgroundColor(tcell.ColorReset)
	modeLine.SetBorder(false)

	return modeLine
}

// UpdateModeLine updates the mode line with new profile information
func UpdateModeLine(modeLine *tview.TextView, profileName string) {
	modeLineText := fmt.Sprintf(" Profile: %s | Press 'O' to switch profile ", profileName)
	modeLine.SetText(modeLineText)
}
