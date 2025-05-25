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

// CreateInteractiveJSONDetailView creates a JSON detail view with copy and edit functionality
func CreateInteractiveJSONDetailView(title string, data interface{}, yankTracker *YankTracker, onCopy func(), onEdit func()) *tview.TextView {
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

	// Enable mouse support for text selection
	textView.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		// Allow default mouse handling for text selection
		return action, event
	})

	// Set up key handling
	textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'y':
			if yankTracker.HandleYankKey() {
				// Double-y detected, copy to clipboard
				if onCopy != nil {
					onCopy()
				}
			}
			return nil
		case 'e':
			// Open in nvim
			if onEdit != nil {
				onEdit()
			}
			return nil
		}
		return event
	})

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
		SetText("Press 'Esc' or 'q' to go back, 'Q' to quit, 'yy' to copy JSON, 'e' to edit in nvim").
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

// UpdateModeLineWithPageInfo updates the mode line with profile and page information
func UpdateModeLineWithPageInfo(modeLine *tview.TextView, profileName string, pageInfo string) {
	// Calculate spacing to right-align page info
	leftText := fmt.Sprintf(" Profile: %s | Press 'O' to switch profile ", profileName)

	// Get terminal width (approximate)
	width := 120 // Default width, will be adjusted dynamically

	// Create spacing
	spacingNeeded := width - len(leftText) - len(pageInfo) - 1
	if spacingNeeded < 1 {
		spacingNeeded = 1
	}
	spacing := ""
	for i := 0; i < spacingNeeded; i++ {
		spacing += " "
	}

	modeLineText := leftText + spacing + pageInfo
	modeLine.SetText(modeLineText)
}
