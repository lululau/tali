package ui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// AppControlInterface defines methods the VimSearchHandler needs from the App
type AppControlInterface interface {
	SetActiveSearchHandler(handler *VimSearchHandler)
	SetSearchBarVisibility(visible bool)
	GetAppSearchBar() *tview.InputField // To get the query text for PerformSearch
}

// SearchState holds the current search state
type SearchState struct {
	Query         string
	IsActive      bool
	CurrentIndex  int
	TotalMatches  int
	Matches       []SearchMatch
	CaseSensitive bool
}

// SearchMatch represents a search match location
type SearchMatch struct {
	Row    int
	Column int
	Text   string
}

// NewSearchState creates a new search state
func NewSearchState() *SearchState {
	return &SearchState{
		Query:         "",
		IsActive:      false,
		CurrentIndex:  -1,
		TotalMatches:  0,
		Matches:       []SearchMatch{},
		CaseSensitive: false,
	}
}

// Reset resets the search state
func (s *SearchState) Reset() {
	s.Query = ""
	s.IsActive = false
	s.CurrentIndex = -1
	s.TotalMatches = 0
	s.Matches = []SearchMatch{}
}

// NextMatch moves to the next search match
func (s *SearchState) NextMatch() {
	if s.TotalMatches == 0 {
		return
	}
	s.CurrentIndex = (s.CurrentIndex + 1) % s.TotalMatches
}

// PrevMatch moves to the previous search match
func (s *SearchState) PrevMatch() {
	if s.TotalMatches == 0 {
		return
	}
	s.CurrentIndex = (s.CurrentIndex - 1 + s.TotalMatches) % s.TotalMatches
}

// GetCurrentMatch returns the current match
func (s *SearchState) GetCurrentMatch() *SearchMatch {
	if s.CurrentIndex >= 0 && s.CurrentIndex < len(s.Matches) {
		return &s.Matches[s.CurrentIndex]
	}
	return nil
}

// SearchInTable searches for text in a table and returns matches
func SearchInTable(table *tview.Table, query string, caseSensitive bool) []SearchMatch {
	var matches []SearchMatch
	if query == "" {
		return matches
	}

	searchQuery := query
	if !caseSensitive {
		searchQuery = strings.ToLower(query)
	}

	rowCount := table.GetRowCount()
	colCount := table.GetColumnCount()

	for row := 0; row < rowCount; row++ {
		for col := 0; col < colCount; col++ {
			cell := table.GetCell(row, col)
			if cell != nil {
				cellText := cell.Text
				searchText := cellText
				if !caseSensitive {
					searchText = strings.ToLower(cellText)
				}

				if strings.Contains(searchText, searchQuery) {
					matches = append(matches, SearchMatch{
						Row:    row,
						Column: col,
						Text:   cellText,
					})
				}
			}
		}
	}

	return matches
}

// SearchInTextView searches for text in a TextView and returns matches.
// If pristineText is provided and non-empty, it searches in that instead of textView.GetText().
func SearchInTextView(textView *tview.TextView, query string, caseSensitive bool, pristineText ...string) []SearchMatch {
	var matches []SearchMatch
	if query == "" {
		return matches
	}

	var textToSearch string
	if len(pristineText) > 0 && pristineText[0] != "" {
		textToSearch = pristineText[0]
	} else {
		textToSearch = textView.GetText(false) // GetText(false) gets text without tags
	}

	lines := strings.Split(textToSearch, "\n")

	searchQuery := query
	if !caseSensitive {
		searchQuery = strings.ToLower(query)
	}

	for lineNum, line := range lines {
		searchLine := line
		if !caseSensitive {
			searchLine = strings.ToLower(line)
		}

		index := 0
		for {
			pos := strings.Index(searchLine[index:], searchQuery)
			if pos == -1 {
				break
			}
			actualPos := index + pos
			// Ensure the match text comes from the original line content, not the lowercased one.
			// And its length is based on the original query length.
			matchText := ""
			if actualPos+len(query) <= len(line) {
				matchText = line[actualPos : actualPos+len(query)]
			} else if actualPos < len(line) { // Partial match at end of line (should ideally not happen with Index)
				matchText = line[actualPos:]
			}

			matches = append(matches, SearchMatch{
				Row:    lineNum,
				Column: actualPos,
				Text:   matchText,
			})
			index = actualPos + len(searchQuery) // Advance by length of searchQuery, not original query if case-insensitive
			if index > len(searchLine) {         // Safety break if index somehow overshoots
				break
			}
		}
	}

	return matches
}

// HighlightTableMatch highlights a match in a table
func HighlightTableMatch(table *tview.Table, match *SearchMatch) {
	if match != nil {
		table.Select(match.Row, match.Column)
		table.ScrollToBeginning()
	}
}

// CreateSearchBar creates a vim-style search bar at the bottom
func CreateSearchBar() *tview.InputField {
	searchBar := tview.NewInputField()
	searchBar.SetLabel("/")
	searchBar.SetFieldBackgroundColor(tcell.ColorDefault)
	searchBar.SetLabelColor(tcell.ColorYellow)
	searchBar.SetFieldTextColor(tcell.ColorWhite)
	searchBar.SetBackgroundColor(tcell.ColorDefault)
	searchBar.SetBorder(false)
	return searchBar
}

// VimSearchHandler handles vim-style search for any component
type VimSearchHandler struct {
	searchState    *SearchState
	mainComponent  tview.Primitive
	searchCallback func(string) // Called when search is performed
	appRef         AppControlInterface
}

// NewVimSearchHandler creates a new vim-style search handler
func NewVimSearchHandler(mainComponent tview.Primitive, appRef AppControlInterface, searchCallback func(string)) *VimSearchHandler {
	return &VimSearchHandler{
		searchState:    NewSearchState(),
		mainComponent:  mainComponent,
		searchCallback: searchCallback,
		appRef:         appRef,
	}
}

// GetMainComponent returns the main component this handler is attached to.
func (h *VimSearchHandler) GetMainComponent() tview.Primitive {
	return h.mainComponent
}

// GetSearchState returns the search state
func (h *VimSearchHandler) GetSearchState() *SearchState {
	return h.searchState
}

// HighlightTextViewMatch highlights a match in a TextView by scrolling to it
func HighlightTextViewMatch(textView *tview.TextView, match *SearchMatch) {
	if match != nil {
		// Scroll to the line containing the match
		textView.ScrollTo(match.Row, 0)
	}
}

// HighlightTextInTextView highlights all matches in a TextView.
// It takes the pristine (unhighlighted) text as a base.
func HighlightTextInTextView(textView *tview.TextView, pristineText string, query string, caseSensitive bool) {
	if query == "" {
		textView.SetText(pristineText) // Restore original if query is cleared
		return
	}

	var highlightedText string
	if !caseSensitive {
		// The third argument to replaceAllCaseInsensitive is ignored by the function itself;
		// it constructs the replacement using originalMatch to preserve casing.
		highlightedText = replaceAllCaseInsensitive(pristineText, query, "")
	} else {
		// For case-sensitive, we need to ensure original casing is preserved in the highlighted segment.
		// strings.ReplaceAll would use the casing from 'query'.
		// So, we use a similar approach to replaceAllCaseInsensitive but without ToLower.
		var result strings.Builder
		lastIndex := 0
		searchText := pristineText
		for {
			index := strings.Index(searchText[lastIndex:], query)
			if index == -1 {
				break
			}
			actualIndex := lastIndex + index
			result.WriteString(pristineText[lastIndex:actualIndex])             // Text before match
			originalMatch := pristineText[actualIndex : actualIndex+len(query)] // The segment from pristineText
			result.WriteString("[yellow::b]" + originalMatch + "[white::-]")    // Highlighted match
			lastIndex = actualIndex + len(query)
		}
		result.WriteString(pristineText[lastIndex:]) // Remaining text
		highlightedText = result.String()
	}

	textView.SetText(highlightedText)
}

// replaceAllCaseInsensitive replaces all occurrences of substr in s with replacement, case insensitive
func replaceAllCaseInsensitive(s, substr, replacement string) string {
	if substr == "" {
		return s
	}

	lowerS := strings.ToLower(s)
	lowerSubstr := strings.ToLower(substr)

	var result strings.Builder
	lastIndex := 0

	for {
		index := strings.Index(lowerS[lastIndex:], lowerSubstr)
		if index == -1 {
			break
		}

		actualIndex := lastIndex + index

		// Add text before match
		result.WriteString(s[lastIndex:actualIndex])

		// Add highlighted match (preserve original case)
		originalMatch := s[actualIndex : actualIndex+len(substr)]
		result.WriteString("[yellow::b]" + originalMatch + "[white::-]")

		lastIndex = actualIndex + len(substr)
	}

	// Add remaining text
	result.WriteString(s[lastIndex:])

	return result.String()
}

// ClearHighlightInTextView removes highlighting from a TextView
func ClearHighlightInTextView(textView *tview.TextView, originalText string) {
	textView.SetText(originalText)
}

// HighlightTableCells highlights matching cells in a table
func HighlightTableCells(table *tview.Table, matches []SearchMatch, currentIndex int) {
	// Reset all cell colors first
	rowCount := table.GetRowCount()
	colCount := table.GetColumnCount()

	for row := 0; row < rowCount; row++ {
		for col := 0; col < colCount; col++ {
			cell := table.GetCell(row, col)
			if cell != nil {
				if row == 0 {
					// Header row
					cell.SetTextColor(tcell.ColorYellow)
				} else {
					// Data row
					cell.SetTextColor(tcell.ColorWhite)
				}
				cell.SetBackgroundColor(tcell.ColorDefault)
			}
		}
	}

	// Highlight matches
	for i, match := range matches {
		cell := table.GetCell(match.Row, match.Column)
		if cell != nil {
			if i == currentIndex {
				// Current match - bright highlight
				cell.SetBackgroundColor(tcell.ColorYellow)
				cell.SetTextColor(tcell.ColorBlack)
			} else {
				// Other matches - dim highlight
				cell.SetBackgroundColor(tcell.ColorDarkCyan)
				cell.SetTextColor(tcell.ColorWhite)
			}
		}
	}
}

// EnterSearchMode is called by the view when '/' is pressed.
func (h *VimSearchHandler) EnterSearchMode() {
	if h.appRef != nil {
		h.appRef.SetActiveSearchHandler(h)
		h.appRef.SetSearchBarVisibility(true)
	}
}

// ExitSearchMode is called by the App's search bar (on Esc) or by PerformSearch.
func (h *VimSearchHandler) ExitSearchMode() {
	if h.appRef != nil {
		h.appRef.SetSearchBarVisibility(false)
	}
}

// PerformSearch is called by the App's search bar (on Enter).
func (h *VimSearchHandler) PerformSearch(query string) { // Query now comes from app's search bar directly
	if query != "" && h.searchCallback != nil {
		h.searchState.Query = query
		h.searchCallback(query)
		h.searchState.IsActive = true
		h.searchState.CurrentIndex = -1
		h.NextMatch()
	}
	h.ExitSearchMode()
}

// NextMatch moves to the next search match
func (h *VimSearchHandler) NextMatch() {
	if h.searchState.TotalMatches == 0 {
		return
	}
	h.searchState.CurrentIndex = (h.searchState.CurrentIndex + 1) % h.searchState.TotalMatches
	// The view itself should handle highlighting the current match based on SearchState
}

// PrevMatch moves to the previous search match
func (h *VimSearchHandler) PrevMatch() {
	if h.searchState.TotalMatches == 0 {
		return
	}
	h.searchState.CurrentIndex = (h.searchState.CurrentIndex - 1 + h.searchState.TotalMatches) % h.searchState.TotalMatches
	// The view itself should handle highlighting the current match based on SearchState
}
