package app

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"aliyun-tui-viewer/internal/client"
	"aliyun-tui-viewer/internal/config"
	"aliyun-tui-viewer/internal/service"
	"aliyun-tui-viewer/internal/ui"
)

// setupGlobalInputCapture sets up global input handling
func (a *App) setupGlobalInputCapture() {
	a.tviewApp.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		currentFocus := a.tviewApp.GetFocus()
		if modal, isModal := currentFocus.(*tview.Modal); isModal && modal.HasFocus() {
			return event
		}
		if form, isForm := currentFocus.(*tview.Form); isForm && form.HasFocus() &&
			(event.Key() == tcell.KeyTab || event.Key() == tcell.KeyBacktab || event.Key() == tcell.KeyEnter) {
			return event
		}

		if event.Key() == tcell.KeyCtrlC {
			a.Stop()
			return nil
		}

		currentPageName, _ := a.pages.GetFrontPage()

		switch event.Key() {
		case tcell.KeyEscape:
			a.handleEscapeKey(currentPageName)
			return nil
		case tcell.KeyRune:
			if event.Rune() == 'Q' { // Only uppercase Q exits the program
				a.Stop()
				return nil
			} else if event.Rune() == 'q' { // lowercase q goes back
				a.handleBackKey(currentPageName)
				return nil
			} else if event.Rune() == 'O' { // Uppercase O opens profile selection
				a.showProfileSelectionDialog()
				return nil
			}
		}
		return event
	})
}

// handleEscapeKey handles escape key navigation
func (a *App) handleEscapeKey(currentPageName string) {
	switch currentPageName {
	case ui.PageEcsList, ui.PageDnsDomains, ui.PageSlbList, ui.PageOssBuckets, ui.PageRdsList:
		a.handleNavigation(ui.PageMainMenu, a.mainMenu)
	case ui.PageEcsDetail:
		a.handleNavigation(ui.PageEcsList, a.ecsInstanceTable)
	case ui.PageDnsRecords:
		a.handleNavigation(ui.PageDnsDomains, a.dnsDomainsTable)
	case ui.PageSlbDetail:
		a.handleNavigation(ui.PageSlbList, a.slbInstanceTable)
	case ui.PageOssObjects:
		a.handleNavigation(ui.PageOssBuckets, a.ossBucketTable)
	case "ossObjectDetail":
		a.handleNavigation(ui.PageOssObjects, a.ossObjectTable)
	case ui.PageRdsDetail:
		a.handleNavigation(ui.PageRdsList, a.rdsInstanceTable)
	}
}

// handleBackKey handles 'q' key navigation
func (a *App) handleBackKey(currentPageName string) {
	switch currentPageName {
	case ui.PageMainMenu:
		// On main menu, q does nothing (only Q exits)
		return
	case ui.PageEcsList, ui.PageDnsDomains, ui.PageSlbList, ui.PageOssBuckets, ui.PageRdsList:
		// On list pages, q goes back to main menu
		a.handleNavigation(ui.PageMainMenu, a.mainMenu)
	case ui.PageEcsDetail:
		a.handleNavigation(ui.PageEcsList, a.ecsInstanceTable)
	case ui.PageDnsRecords:
		a.handleNavigation(ui.PageDnsDomains, a.dnsDomainsTable)
	case ui.PageSlbDetail:
		a.handleNavigation(ui.PageSlbList, a.slbInstanceTable)
	case ui.PageOssObjects:
		a.handleNavigation(ui.PageOssBuckets, a.ossBucketTable)
	case "ossObjectDetail":
		a.handleNavigation(ui.PageOssObjects, a.ossObjectTable)
	case ui.PageRdsDetail:
		a.handleNavigation(ui.PageRdsList, a.rdsInstanceTable)
	}
}

// handleNavigation handles page navigation
func (a *App) handleNavigation(targetPage string, focusItem tview.Primitive) {
	a.pages.SwitchToPage(targetPage)
	if focusItem != nil {
		a.tviewApp.SetFocus(focusItem)
	} else if targetPage == ui.PageMainMenu {
		a.tviewApp.SetFocus(a.mainMenu)
	}
}

// switchToEcsListView switches to ECS list view
func (a *App) switchToEcsListView() {
	if a.allECSInstances == nil {
		instances, err := a.services.ECS.FetchInstances()
		if err != nil {
			a.showErrorModal(err.Error())
			return
		}
		a.allECSInstances = instances
	}
	a.ecsInstanceTable = ui.CreateEcsListView(a.allECSInstances)
	ui.SetupTableNavigation(a.ecsInstanceTable, func(row, col int) {
		instanceId := a.ecsInstanceTable.GetCell(row, 0).GetReference().(string)
		var selectedInstance interface{}
		for _, inst := range a.allECSInstances {
			if inst.InstanceId == instanceId {
				selectedInstance = inst
				break
			}
		}
		detailView := ui.CreateEcsDetailView(selectedInstance)
		a.pages.AddPage(ui.PageEcsDetail, detailView, true, true)
		// Extract the detail view from the flex container
		if detailView.GetItemCount() > 1 {
			a.ecsDetailView = detailView.GetItem(1).(*tview.TextView)
		}
		a.tviewApp.SetFocus(a.ecsDetailView)
	})
	ecsListFlex := ui.WrapTableInFlex(a.ecsInstanceTable)
	a.pages.AddPage(ui.PageEcsList, ecsListFlex, true, true)
	a.tviewApp.SetFocus(a.ecsInstanceTable)
}

// switchToDnsDomainsListView switches to DNS domains list view
func (a *App) switchToDnsDomainsListView() {
	if a.allDomains == nil {
		domains, err := a.services.DNS.FetchDomains()
		if err != nil {
			a.showErrorModal(err.Error())
			return
		}
		a.allDomains = domains
	}
	a.dnsDomainsTable = ui.CreateDnsDomainsListView(a.allDomains)
	ui.SetupTableNavigation(a.dnsDomainsTable, func(row, col int) {
		domainName := a.dnsDomainsTable.GetCell(row, 0).GetReference().(string)
		a.switchToDnsRecordsListView(domainName)
	})
	dnsDomainsListFlex := ui.WrapTableInFlex(a.dnsDomainsTable)
	a.pages.AddPage(ui.PageDnsDomains, dnsDomainsListFlex, true, true)
	a.tviewApp.SetFocus(a.dnsDomainsTable)
}

// switchToDnsRecordsListView switches to DNS records list view
func (a *App) switchToDnsRecordsListView(domainName string) {
	records, err := a.services.DNS.FetchDomainRecords(domainName)
	if err != nil {
		a.showErrorModal(err.Error())
		return
	}
	a.dnsRecordsTable = ui.CreateDnsRecordsListView(records, domainName)
	ui.SetupTableNavigation(a.dnsRecordsTable, nil)
	dnsRecordsListFlex := ui.WrapTableInFlex(a.dnsRecordsTable)
	a.pages.AddPage(ui.PageDnsRecords, dnsRecordsListFlex, true, true)
	a.tviewApp.SetFocus(a.dnsRecordsTable)
}

// switchToSlbListView switches to SLB list view
func (a *App) switchToSlbListView() {
	if a.allSLBInstances == nil {
		slbs, err := a.services.SLB.FetchInstances()
		if err != nil {
			a.showErrorModal(err.Error())
			return
		}
		a.allSLBInstances = slbs
	}
	a.slbInstanceTable = ui.CreateSlbListView(a.allSLBInstances)
	ui.SetupTableNavigation(a.slbInstanceTable, func(row, col int) {
		slbId := a.slbInstanceTable.GetCell(row, 0).GetReference().(string)
		var selectedSlb interface{}
		for _, lb := range a.allSLBInstances {
			if lb.LoadBalancerId == slbId {
				selectedSlb = lb
				break
			}
		}
		detailView := ui.CreateSlbDetailView(selectedSlb)
		a.pages.AddPage(ui.PageSlbDetail, detailView, true, true)
		// Extract the detail view from the flex container
		if detailView.GetItemCount() > 1 {
			a.slbDetailView = detailView.GetItem(1).(*tview.TextView)
		}
		a.tviewApp.SetFocus(a.slbDetailView)
	})
	slbListFlex := ui.WrapTableInFlex(a.slbInstanceTable)
	a.pages.AddPage(ui.PageSlbList, slbListFlex, true, true)
	a.tviewApp.SetFocus(a.slbInstanceTable)
}

// switchToOssBucketListView switches to OSS bucket list view
func (a *App) switchToOssBucketListView() {
	if a.allOssBuckets == nil {
		buckets, err := a.services.OSS.FetchBuckets()
		if err != nil {
			a.showErrorModal(err.Error())
			return
		}
		a.allOssBuckets = buckets
	}
	a.ossBucketTable = ui.CreateOssBucketListView(a.allOssBuckets)
	ui.SetupTableNavigation(a.ossBucketTable, func(row, col int) {
		bucketName := a.ossBucketTable.GetCell(row, 0).GetReference().(string)
		a.currentBucketName = bucketName
		a.switchToOssObjectListView(bucketName)
	})
	ossBucketListFlex := ui.WrapTableInFlex(a.ossBucketTable)
	a.pages.AddPage(ui.PageOssBuckets, ossBucketListFlex, true, true)
	a.tviewApp.SetFocus(a.ossBucketTable)
}

// switchToOssObjectListView switches to OSS object list view
func (a *App) switchToOssObjectListView(bucketName string) {
	// Initialize pagination state
	a.currentBucketName = bucketName
	a.ossCurrentMarker = ""
	a.ossPreviousMarkers = []string{}
	a.ossCurrentPage = 1
	a.ossPageSize = 20 // Set page size to 20 objects per page
	a.ossHasNextPage = false

	// Load first page
	a.loadOssObjectPage()
}

// loadOssObjectPage loads the current page of OSS objects
func (a *App) loadOssObjectPage() {
	result, err := a.services.OSS.FetchObjects(a.currentBucketName, a.ossCurrentMarker, a.ossPageSize)
	if err != nil {
		a.showErrorModal(err.Error())
		return
	}

	a.ossHasNextPage = result.IsTruncated
	hasPrevious := len(a.ossPreviousMarkers) > 0

	// Create paginated view
	ossObjectView := ui.CreateOssObjectPaginatedView(result.Objects, a.currentBucketName, a.ossCurrentPage, a.ossHasNextPage, hasPrevious)

	// Extract the table from the flex container for navigation setup
	if ossObjectView.GetItemCount() > 0 {
		a.ossObjectTable = ossObjectView.GetItem(0).(*tview.Table)
	}

	// Setup table navigation with object selection
	ui.SetupTableNavigation(a.ossObjectTable, func(row, col int) {
		objectKey := a.ossObjectTable.GetCell(row, 0).GetReference().(string)
		// Find the object details
		for _, obj := range result.Objects {
			if obj.Key == objectKey {
				a.ossDetailView = ui.CreateJSONDetailView(fmt.Sprintf("Object Details: %s", objectKey), obj)
				a.pages.AddPage("ossObjectDetail", a.ossDetailView, true, true)
				a.tviewApp.SetFocus(a.ossDetailView)
				break
			}
		}
	})

	// Setup pagination navigation
	a.setupOssPaginationNavigation(ossObjectView, result)

	a.pages.AddPage(ui.PageOssObjects, ossObjectView, true, true)
	a.tviewApp.SetFocus(a.ossObjectTable)
}

// setupOssPaginationNavigation sets up pagination key bindings
func (a *App) setupOssPaginationNavigation(view *tview.Flex, result *service.ObjectListResult) {
	view.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case ']': // Next page
			if a.ossHasNextPage {
				a.goToNextOssPage(result.NextMarker)
			}
			return nil
		case '[': // Previous page
			if len(a.ossPreviousMarkers) > 0 {
				a.goToPrevOssPage()
			}
			return nil
		case '0': // First page
			a.goToFirstOssPage()
			return nil
		}
		return event
	})
}

// goToNextOssPage navigates to the next page
func (a *App) goToNextOssPage(nextMarker string) {
	if nextMarker == "" {
		return
	}

	// Save current marker to previous markers stack
	a.ossPreviousMarkers = append(a.ossPreviousMarkers, a.ossCurrentMarker)
	a.ossCurrentMarker = nextMarker
	a.ossCurrentPage++

	a.loadOssObjectPage()
}

// goToPrevOssPage navigates to the previous page
func (a *App) goToPrevOssPage() {
	if len(a.ossPreviousMarkers) == 0 {
		return
	}

	// Pop the last marker from the stack
	lastIndex := len(a.ossPreviousMarkers) - 1
	a.ossCurrentMarker = a.ossPreviousMarkers[lastIndex]
	a.ossPreviousMarkers = a.ossPreviousMarkers[:lastIndex]
	a.ossCurrentPage--

	a.loadOssObjectPage()
}

// goToFirstOssPage navigates to the first page
func (a *App) goToFirstOssPage() {
	a.ossCurrentMarker = ""
	a.ossPreviousMarkers = []string{}
	a.ossCurrentPage = 1

	a.loadOssObjectPage()
}

// showProfileSelectionDialog shows the profile selection dialog
func (a *App) showProfileSelectionDialog() {
	profiles, err := config.ListAllProfiles()
	if err != nil {
		a.showErrorModal(fmt.Sprintf("Failed to load profiles: %v", err))
		return
	}

	ui.ShowProfileSelectionDialog(a.pages, a.tviewApp, profiles, a.currentProfile,
		func(selectedProfile string) {
			// Profile selected callback
			a.switchToProfile(selectedProfile)
		},
		func() {
			// Cancel callback - restore focus to current page
			_, prim := a.pages.GetFrontPage()
			if prim != nil {
				a.tviewApp.SetFocus(prim)
			} else {
				a.tviewApp.SetFocus(a.mainMenu)
			}
		})
}

// switchToProfile switches to the selected profile and reinitializes the application
func (a *App) switchToProfile(profileName string) {
	if profileName == a.currentProfile {
		// Same profile, just restore focus
		_, prim := a.pages.GetFrontPage()
		if prim != nil {
			a.tviewApp.SetFocus(prim)
		} else {
			a.tviewApp.SetFocus(a.mainMenu)
		}
		return
	}

	// Store original profile for rollback
	originalProfile := a.currentProfile

	// Switch profile in config
	err := config.SwitchProfile(profileName)
	if err != nil {
		a.showErrorModal(fmt.Sprintf("Failed to switch profile: %v", err))
		return
	}

	// Load new configuration
	cfg, err := config.LoadAliyunConfig()
	if err != nil {
		// Rollback profile change
		config.SwitchProfile(originalProfile)
		a.showErrorModal(fmt.Sprintf("Failed to load new configuration: %v", err))
		return
	}

	// Create new clients with the new configuration
	clientConfig := &client.Config{
		AccessKeyID:     cfg.AccessKeyID,
		AccessKeySecret: cfg.AccessKeySecret,
		RegionID:        cfg.RegionID,
		OssEndpoint:     cfg.OssEndpoint,
	}

	newClients, err := client.NewAliyunClients(clientConfig)
	if err != nil {
		// Rollback profile change
		config.SwitchProfile(originalProfile)
		a.showErrorModal(fmt.Sprintf("Failed to create new clients: %v", err))
		return
	}

	// Create new services with the new clients
	newServices := &Services{
		ECS: service.NewECSService(newClients.ECS),
		DNS: service.NewDNSService(newClients.DNS),
		SLB: service.NewSLBService(newClients.SLB),
		RDS: service.NewRDSService(newClients.RDS),
		OSS: service.NewOSSServiceWithCredentials(newClients.OSS, cfg.AccessKeyID, cfg.AccessKeySecret, cfg.OssEndpoint),
	}

	// Update application state
	a.clients = newClients
	a.services = newServices
	a.currentProfile = profileName

	// Update mode line
	ui.UpdateModeLine(a.modeLine, a.currentProfile)

	// Clear cached data to force reload with new profile
	a.clearCachedData()

	// Navigate back to main menu for better user experience
	a.pages.SwitchToPage(ui.PageMainMenu)
	a.tviewApp.SetFocus(a.mainMenu)

	// Show success message
	a.showErrorModal(fmt.Sprintf("Successfully switched to profile: %s\nNew credentials are now active.", profileName))
}

// clearCachedData clears all cached data to force reload with new profile
func (a *App) clearCachedData() {
	a.allECSInstances = nil
	a.allDomains = nil
	a.allSLBInstances = nil
	a.allRDSInstances = nil
	a.allOssBuckets = nil
	a.currentBucketName = ""

	// Reset OSS pagination state
	a.ossCurrentMarker = ""
	a.ossPreviousMarkers = []string{}
	a.ossCurrentPage = 0
	a.ossPageSize = 0
	a.ossHasNextPage = false
}

// switchToRdsListView switches to RDS list view
func (a *App) switchToRdsListView() {
	if a.allRDSInstances == nil {
		instances, err := a.services.RDS.FetchInstances()
		if err != nil {
			a.showErrorModal(fmt.Sprintf("Failed to fetch RDS instances: %v", err))
			return
		}
		a.allRDSInstances = instances
	}
	a.rdsInstanceTable = ui.CreateRdsListView(a.allRDSInstances)
	ui.SetupTableNavigation(a.rdsInstanceTable, func(row, col int) {
		cell := a.rdsInstanceTable.GetCell(row, 0)
		instanceId, ok := cell.GetReference().(string)
		if !ok {
			return
		}
		var selectedInstance interface{}
		for _, inst := range a.allRDSInstances {
			if inst.DBInstanceId == instanceId {
				selectedInstance = inst
				break
			}
		}
		detailViewContent := ui.CreateRdsDetailView(selectedInstance)
		a.pages.AddPage(ui.PageRdsDetail, detailViewContent, true, true)
		// Extract the detail view from the flex container
		if detailViewContent.GetItemCount() > 1 {
			a.rdsDetailView = detailViewContent.GetItem(1).(*tview.TextView)
		}
		a.tviewApp.SetFocus(a.rdsDetailView)
	})
	rdsListFlex := ui.WrapTableInFlex(a.rdsInstanceTable)
	a.pages.AddPage(ui.PageRdsList, rdsListFlex, true, true)
	a.tviewApp.SetFocus(a.rdsInstanceTable)
}
