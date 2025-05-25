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
		SetDynamicColors(true).
		SetTextStyle(tcell.StyleDefault.Background(tcell.ColorReset))
	textView.SetBorder(true).SetTitle(title).SetBackgroundColor(tcell.ColorReset)

	return textView
}

// CreateInteractiveJSONDetailView creates a JSON detail view with copy and edit functionality
func CreateInteractiveJSONDetailView(title string, data interface{}, onCopy func(), onEdit func()) *tview.TextView {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		jsonData = []byte(fmt.Sprintf("Error marshaling JSON: %v", err))
	}

	textView := tview.NewTextView().
		SetText(string(jsonData)).
		SetScrollable(true).
		SetWrap(false).
		SetDynamicColors(true).
		SetTextStyle(tcell.StyleDefault.Background(tcell.ColorReset))
	textView.SetBorder(true).SetTitle(title).SetBackgroundColor(tcell.ColorReset)

	textView.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		return action, event
	})

	return textView
}

// CreateInteractiveJSONDetailViewWithSearch creates a JSON detail view with copy, edit and search functionality
func CreateInteractiveJSONDetailViewWithSearch(title string, data interface{}, appRef AppControlInterface, onCopy func(), onEdit func()) (*tview.TextView, *VimSearchHandler) {
	textView := CreateInteractiveJSONDetailView(title, data, onCopy, onEdit)

	jsonDataBytes, _ := json.MarshalIndent(data, "", "  ")
	pristineText := string(jsonDataBytes)

	var searchHandler *VimSearchHandler
	searchHandler = NewVimSearchHandler(textView, appRef, func(query string) {
		state := searchHandler.GetSearchState()
		HighlightTextInTextView(textView, pristineText, query, state.CaseSensitive)

		matches := SearchInTextView(textView, query, state.CaseSensitive, pristineText)
		state.Matches = matches
		state.TotalMatches = len(matches)
		state.CurrentIndex = -1

		if state.TotalMatches > 0 {
			searchHandler.NextMatch()
			HighlightTextViewMatch(textView, state.GetCurrentMatch())
		} else {
			if query != "" {
				textView.SetText(pristineText)
			}
		}
	})

	textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		yankTracker := NewYankTracker()

		switch event.Rune() {
		case 'y':
			if yankTracker.HandleYankKey() {
				if onCopy != nil {
					onCopy()
				}
			}
			return nil
		case 'e':
			if onEdit != nil {
				onEdit()
			}
			return nil
		case '/':
			searchHandler.EnterSearchMode()
			return nil
		case 'n':
			state := searchHandler.GetSearchState()
			if state.IsActive && state.TotalMatches > 0 {
				searchHandler.NextMatch()
				HighlightTextViewMatch(textView, state.GetCurrentMatch())
			}
			return nil
		case 'N':
			state := searchHandler.GetSearchState()
			if state.IsActive && state.TotalMatches > 0 {
				searchHandler.PrevMatch()
				HighlightTextViewMatch(textView, state.GetCurrentMatch())
			}
			return nil
		}
		return event
	})

	return textView, searchHandler
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

// SetupTableNavigationWithSearch sets up j/k navigation and search for tables
func SetupTableNavigationWithSearch(table *tview.Table, appRef AppControlInterface, onSelect func(row, column int)) *VimSearchHandler {
	table.SetSelectedFunc(func(row, column int) {
		if row > 0 && onSelect != nil {
			onSelect(row, column)
		}
	})

	var searchHandler *VimSearchHandler // Declare here
	searchHandler = NewVimSearchHandler(table, appRef, func(query string) {
		state := searchHandler.GetSearchState()
		matches := SearchInTable(table, query, state.CaseSensitive)
		state.Matches = matches
		state.TotalMatches = len(matches)
		state.CurrentIndex = -1

		if state.TotalMatches > 0 {
			searchHandler.NextMatch()
			currentMatch := state.GetCurrentMatch()
			HighlightTableCells(table, state.Matches, state.CurrentIndex)
			HighlightTableMatch(table, currentMatch)
		} else {
			HighlightTableCells(table, []SearchMatch{}, -1)
		}
	})

	originalInputCapture := table.GetInputCapture()
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'j':
			currentRow, _ := table.GetSelection()
			rowCount := table.GetRowCount()
			if currentRow < rowCount-1 {
				table.Select(currentRow+1, 0)
			}
			return nil
		case 'k':
			currentRow, _ := table.GetSelection()
			rowCount := table.GetRowCount()
			if currentRow > 1 {
				table.Select(currentRow-1, 0)
			} else if rowCount > 1 {
				table.Select(1, 0)
			}
			return nil
		case '/':
			searchHandler.EnterSearchMode()
			return nil
		case 'n':
			state := searchHandler.GetSearchState()
			if state.IsActive && state.TotalMatches > 0 {
				searchHandler.NextMatch()
				currentMatch := state.GetCurrentMatch()
				HighlightTableCells(table, state.Matches, state.CurrentIndex)
				HighlightTableMatch(table, currentMatch)
			}
			return nil
		case 'N':
			state := searchHandler.GetSearchState()
			if state.IsActive && state.TotalMatches > 0 {
				searchHandler.PrevMatch()
				currentMatch := state.GetCurrentMatch()
				HighlightTableCells(table, state.Matches, state.CurrentIndex)
				HighlightTableMatch(table, currentMatch)
			}
			return nil
		}
		if originalInputCapture != nil {
			return originalInputCapture(event)
		}
		return event
	})

	return searchHandler
}

// CreateDetailViewWithInstructions creates a detail view with navigation instructions
func CreateDetailViewWithInstructions(detailView *tview.TextView) *tview.Flex {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	// Instructions
	instructions := tview.NewTextView().
		SetText("Press 'Esc' or 'q' to go back, 'Q' to quit, 'yy' to copy JSON, 'e' to edit in nvim, '/' to search, 'n/p' for next/prev").
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
