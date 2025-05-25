package app

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
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
	case ui.PageEcsList, ui.PageDnsDomains, ui.PageSlbList, ui.PageOssBuckets, ui.PageRdsList, ui.PageRedisList:
		a.handleNavigation(ui.PageMainMenu, a.mainMenu)
	case ui.PageEcsDetail:
		a.handleNavigation(ui.PageEcsList, a.ecsInstanceTable)
	case ui.PageDnsRecords:
		a.handleNavigation(ui.PageDnsDomains, a.dnsDomainsTable)
	case ui.PageSlbDetail:
		a.handleNavigation(ui.PageSlbList, a.slbInstanceTable)
	case ui.PageOssObjects:
		ui.UpdateModeLine(a.modeLine, a.currentProfile)
		a.handleNavigation(ui.PageOssBuckets, a.ossBucketTable)
	case "ossObjectDetail":
		a.handleNavigation(ui.PageOssObjects, a.ossObjectTable)
	case ui.PageRdsDetail:
		a.handleNavigation(ui.PageRdsList, a.rdsInstanceTable)
	case ui.PageRdsDatabases:
		a.handleNavigation(ui.PageRdsList, a.rdsInstanceTable)
	case ui.PageRdsAccounts:
		a.handleNavigation(ui.PageRdsList, a.rdsInstanceTable)
	case "rdsDatabaseDetail":
		a.handleNavigation(ui.PageRdsDatabases, a.rdsDatabaseTable)
	case "rdsAccountDetail":
		a.handleNavigation(ui.PageRdsAccounts, a.rdsAccountTable)
	case ui.PageRedisAccounts:
		a.handleNavigation(ui.PageRedisList, a.redisInstanceTable)
	case "redisDetail":
		a.handleNavigation(ui.PageRedisList, a.redisInstanceTable)
	case "redisAccountDetail":
		a.handleNavigation(ui.PageRedisAccounts, a.redisAccountTable)
	}
}

// handleBackKey handles 'q' key navigation
func (a *App) handleBackKey(currentPageName string) {
	switch currentPageName {
	case ui.PageMainMenu:
		return
	case ui.PageEcsList, ui.PageDnsDomains, ui.PageSlbList, ui.PageOssBuckets, ui.PageRdsList, ui.PageRedisList:
		a.handleNavigation(ui.PageMainMenu, a.mainMenu)
	case ui.PageEcsDetail:
		a.handleNavigation(ui.PageEcsList, a.ecsInstanceTable)
	case ui.PageDnsRecords:
		a.handleNavigation(ui.PageDnsDomains, a.dnsDomainsTable)
	case ui.PageSlbDetail:
		a.handleNavigation(ui.PageSlbList, a.slbInstanceTable)
	case ui.PageOssObjects:
		ui.UpdateModeLine(a.modeLine, a.currentProfile)
		a.handleNavigation(ui.PageOssBuckets, a.ossBucketTable)
	case "ossObjectDetail":
		a.handleNavigation(ui.PageOssObjects, a.ossObjectTable)
	case ui.PageRdsDetail:
		a.handleNavigation(ui.PageRdsList, a.rdsInstanceTable)
	case ui.PageRdsDatabases:
		a.handleNavigation(ui.PageRdsList, a.rdsInstanceTable)
	case ui.PageRdsAccounts:
		a.handleNavigation(ui.PageRdsList, a.rdsInstanceTable)
	case "rdsDatabaseDetail":
		a.handleNavigation(ui.PageRdsDatabases, a.rdsDatabaseTable)
	case "rdsAccountDetail":
		a.handleNavigation(ui.PageRdsAccounts, a.rdsAccountTable)
	case ui.PageRedisAccounts:
		a.handleNavigation(ui.PageRedisList, a.redisInstanceTable)
	case "redisDetail":
		a.handleNavigation(ui.PageRedisList, a.redisInstanceTable)
	case "redisAccountDetail":
		a.handleNavigation(ui.PageRedisAccounts, a.redisAccountTable)
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
	ui.SetupTableNavigationWithSearch(a.ecsInstanceTable, a, func(row, col int) {
		instanceId := a.ecsInstanceTable.GetCell(row, 0).GetReference().(string)
		var selectedInstance interface{}
		for _, inst := range a.allECSInstances {
			if inst.InstanceId == instanceId {
				selectedInstance = inst
				break
			}
		}
		a.currentDetailData = selectedInstance
		detailView, _ := ui.CreateInteractiveJSONDetailViewWithSearch(
			fmt.Sprintf("ECS Details: %s", instanceId),
			selectedInstance,
			a,
			func() {
				err := ui.CopyToClipboard(selectedInstance)
				if err != nil {
					a.showErrorModal(fmt.Sprintf("Failed to copy: %v", err))
				} else {
					a.showErrorModal("Copied!")
				}
			},
			func() {
				err := ui.OpenInNvim(selectedInstance)
				if err != nil {
					a.showErrorModal(fmt.Sprintf("Failed to edit: %v", err))
				}
			},
		)
		detailViewWithInstructions := ui.CreateDetailViewWithInstructions(detailView)
		a.pages.AddPage(ui.PageEcsDetail, detailViewWithInstructions, true, true)
		if detailViewWithInstructions.GetItemCount() > 1 {
			a.ecsDetailView = detailViewWithInstructions.GetItem(1).(*tview.TextView)
		}
		a.tviewApp.SetFocus(a.ecsDetailView)
	})

	a.setupTableYankFunctionality(a.ecsInstanceTable, a.allECSInstances)
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
	ui.SetupTableNavigationWithSearch(a.dnsDomainsTable, a, func(row, col int) {
		domainName := a.dnsDomainsTable.GetCell(row, 0).GetReference().(string)
		a.switchToDnsRecordsListView(domainName)
	})

	a.setupTableYankFunctionality(a.dnsDomainsTable, a.allDomains)
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
	ui.SetupTableNavigationWithSearch(a.dnsRecordsTable, a, nil)

	a.setupTableYankFunctionality(a.dnsRecordsTable, records)
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
	ui.SetupTableNavigationWithSearch(a.slbInstanceTable, a, func(row, col int) {
		slbId := a.slbInstanceTable.GetCell(row, 0).GetReference().(string)
		var selectedSlb interface{}
		for _, lb := range a.allSLBInstances {
			if lb.LoadBalancerId == slbId {
				selectedSlb = lb
				break
			}
		}
		a.currentDetailData = selectedSlb
		detailView, _ := ui.CreateInteractiveJSONDetailViewWithSearch(
			fmt.Sprintf("SLB Details: %s", slbId),
			selectedSlb,
			a,
			func() {
				err := ui.CopyToClipboard(selectedSlb)
				if err != nil {
					a.showErrorModal(fmt.Sprintf("Failed to copy: %v", err))
				} else {
					a.showErrorModal("Copied!")
				}
			},
			func() {
				err := ui.OpenInNvim(selectedSlb)
				if err != nil {
					a.showErrorModal(fmt.Sprintf("Failed to edit: %v", err))
				}
			},
		)
		detailViewWithInstructions := ui.CreateDetailViewWithInstructions(detailView)
		a.pages.AddPage(ui.PageSlbDetail, detailViewWithInstructions, true, true)
		if detailViewWithInstructions.GetItemCount() > 1 {
			a.slbDetailView = detailViewWithInstructions.GetItem(1).(*tview.TextView)
		}
		a.tviewApp.SetFocus(a.slbDetailView)
	})

	a.setupTableYankFunctionality(a.slbInstanceTable, a.allSLBInstances)
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
	ui.SetupTableNavigationWithSearch(a.ossBucketTable, a, func(row, col int) {
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
	a.ossPageSize = 20
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

	pageInfo := fmt.Sprintf("Page %d", a.ossCurrentPage)
	if a.ossHasNextPage {
		pageInfo += "+"
	}
	ui.UpdateModeLineWithPageInfo(a.modeLine, a.currentProfile, pageInfo)

	ossObjectView := ui.CreateOssObjectPaginatedView(result.Objects, a.currentBucketName, a.ossCurrentPage, a.ossHasNextPage, hasPrevious)

	if ossObjectView.GetItemCount() > 0 {
		a.ossObjectTable = ossObjectView.GetItem(0).(*tview.Table)
	}

	ui.SetupTableNavigationWithSearch(a.ossObjectTable, a, func(row, col int) {
		objectKey := a.ossObjectTable.GetCell(row, 0).GetReference().(string)
		for _, obj := range result.Objects {
			if obj.Key == objectKey {
				a.currentDetailData = obj
				view, _ := ui.CreateInteractiveJSONDetailViewWithSearch(
					fmt.Sprintf("Object Details: %s", objectKey),
					obj,
					a,
					func() {
						err := ui.CopyToClipboard(obj)
						if err != nil {
							a.showErrorModal(fmt.Sprintf("Failed to copy: %v", err))
						} else {
							a.showErrorModal("Copied!")
						}
					},
					func() {
						err := ui.OpenInNvim(obj)
						if err != nil {
							a.showErrorModal(fmt.Sprintf("Failed to edit: %v", err))
						}
					},
				)
				a.ossDetailView = view
				a.pages.AddPage("ossObjectDetail", a.ossDetailView, true, true)
				a.tviewApp.SetFocus(a.ossDetailView)
				break
			}
		}
	})

	a.setupTableYankFunctionality(a.ossObjectTable, result.Objects)
	a.setupOssPaginationNavigation(ossObjectView, result)
	a.pages.AddPage(ui.PageOssObjects, ossObjectView, true, true)
	a.tviewApp.SetFocus(a.ossObjectTable)
}

// setupTableYankFunctionality adds yank (copy) functionality to tables
func (a *App) setupTableYankFunctionality(table *tview.Table, data interface{}) {
	originalInputCapture := table.GetInputCapture()

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'y' {
			if a.yankTracker.HandleYankKey() {
				// Double-y detected, copy current row
				row, _ := table.GetSelection()
				if row > 0 { // Skip header row
					var rowData interface{}

					// Get the reference from the first cell to identify the item
					if cell := table.GetCell(row, 0); cell != nil {
						if ref := cell.GetReference(); ref != nil {
							// Find the corresponding data item based on type
							switch items := data.(type) {
							case []oss.ObjectProperties:
								for _, obj := range items {
									if obj.Key == ref.(string) {
										rowData = obj
										break
									}
								}
							case []ecs.Instance:
								for _, inst := range items {
									if inst.InstanceId == ref.(string) {
										rowData = inst
										break
									}
								}
							case []slb.LoadBalancer:
								for _, lb := range items {
									if lb.LoadBalancerId == ref.(string) {
										rowData = lb
										break
									}
								}
							case []rds.DBInstance:
								for _, db := range items {
									if db.DBInstanceId == ref.(string) {
										rowData = db
										break
									}
								}
							case []alidns.DomainInDescribeDomains:
								for _, domain := range items {
									if domain.DomainName == ref.(string) {
										rowData = domain
										break
									}
								}
							case []alidns.Record:
								for _, record := range items {
									if record.RecordId == ref.(string) {
										rowData = record
										break
									}
								}
							case []rds.Database:
								for _, db := range items {
									if db.DBName == ref.(string) {
										rowData = db
										break
									}
								}
							case []rds.DBInstanceAccount:
								for _, account := range items {
									if account.AccountName == ref.(string) {
										rowData = account
										break
									}
								}
							case []r_kvstore.KVStoreInstance:
								for _, inst := range items {
									if inst.InstanceId == ref.(string) {
										rowData = inst
										break
									}
								}
							case []r_kvstore.Account:
								for _, account := range items {
									if account.AccountName == ref.(string) {
										rowData = account
										break
									}
								}
							}
						}
					}

					if rowData != nil {
						err := ui.CopyToClipboard(rowData)
						if err != nil {
							a.showErrorModal(fmt.Sprintf("Failed to copy to clipboard: %v", err))
						} else {
							a.showErrorModal("Row data copied to clipboard!")
						}
					}
				}
			}
			return nil
		}

		// Call original input capture if it exists
		if originalInputCapture != nil {
			return originalInputCapture(event)
		}
		return event
	})
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
		ECS:   service.NewECSService(newClients.ECS),
		DNS:   service.NewDNSService(newClients.DNS),
		SLB:   service.NewSLBService(newClients.SLB),
		RDS:   service.NewRDSService(newClients.RDS),
		OSS:   service.NewOSSServiceWithCredentials(newClients.OSS, cfg.AccessKeyID, cfg.AccessKeySecret, cfg.OssEndpoint),
		Redis: service.NewRedisService(newClients.Redis),
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
	a.allRedisInstances = nil
	a.allOssBuckets = nil
	a.currentBucketName = ""
	a.currentRdsInstanceId = ""
	a.currentRedisInstanceId = ""

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
	ui.SetupTableNavigationWithSearch(a.rdsInstanceTable, a, func(row, col int) {
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
		a.currentDetailData = selectedInstance
		detailView, _ := ui.CreateInteractiveJSONDetailViewWithSearch(
			fmt.Sprintf("RDS Details: %s", instanceId),
			selectedInstance,
			a,
			func() {
				err := ui.CopyToClipboard(selectedInstance)
				if err != nil {
					a.showErrorModal(fmt.Sprintf("Copy failed: %v", err))
				} else {
					a.showErrorModal("Copied!")
				}
			},
			func() {
				err := ui.OpenInNvim(selectedInstance)
				if err != nil {
					a.showErrorModal(fmt.Sprintf("Edit failed: %v", err))
				}
			},
		)
		detailViewWithInstructions := ui.CreateDetailViewWithInstructions(detailView)
		a.pages.AddPage(ui.PageRdsDetail, detailViewWithInstructions, true, true)
		if detailViewWithInstructions.GetItemCount() > 1 {
			a.rdsDetailView = detailViewWithInstructions.GetItem(1).(*tview.TextView)
		}
		a.tviewApp.SetFocus(a.rdsDetailView)
	})

	a.setupTableYankFunctionality(a.rdsInstanceTable, a.allRDSInstances)
	a.setupRdsKeyHandlers(a.rdsInstanceTable)
	rdsListFlex := ui.WrapTableInFlex(a.rdsInstanceTable)
	a.pages.AddPage(ui.PageRdsList, rdsListFlex, true, true)
	a.tviewApp.SetFocus(a.rdsInstanceTable)
}

// setupRdsKeyHandlers sets up key handlers for RDS specific actions
func (a *App) setupRdsKeyHandlers(table *tview.Table) {
	originalInputCapture := table.GetInputCapture()

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'D': // D key handler for databases
			row, _ := table.GetSelection()
			if row > 0 { // Skip header row
				if cell := table.GetCell(row, 0); cell != nil {
					if instanceId, ok := cell.GetReference().(string); ok {
						a.switchToRdsDatabasesView(instanceId)
					}
				}
			}
			return nil
		case 'A': // A key handler for accounts
			row, _ := table.GetSelection()
			if row > 0 { // Skip header row
				if cell := table.GetCell(row, 0); cell != nil {
					if instanceId, ok := cell.GetReference().(string); ok {
						a.switchToRdsAccountsView(instanceId)
					}
				}
			}
			return nil
		}

		// Call original input capture if it exists
		if originalInputCapture != nil {
			return originalInputCapture(event)
		}
		return event
	})
}

// switchToRdsDatabasesView switches to RDS databases view
func (a *App) switchToRdsDatabasesView(instanceId string) {
	databases, err := a.services.RDS.FetchDatabases(instanceId)
	if err != nil {
		a.showErrorModal(fmt.Sprintf("Failed to fetch databases for instance %s: %v", instanceId, err))
		return
	}

	a.currentRdsInstanceId = instanceId
	a.rdsDatabaseTable = ui.CreateRdsDatabasesListView(databases, instanceId)

	ui.SetupTableNavigationWithSearch(a.rdsDatabaseTable, a, func(row, col int) {
		dbName := a.rdsDatabaseTable.GetCell(row, 0).GetReference().(string)
		var selectedDatabase interface{}
		for _, db := range databases {
			if db.DBName == dbName {
				selectedDatabase = db
				break
			}
		}
		a.currentDetailData = selectedDatabase
		detailView, _ := ui.CreateInteractiveJSONDetailViewWithSearch(
			fmt.Sprintf("Database Details: %s", dbName),
			selectedDatabase,
			a,
			func() {
				err := ui.CopyToClipboard(selectedDatabase)
				if err != nil {
					a.showErrorModal(fmt.Sprintf("Copy failed: %v", err))
				} else {
					a.showErrorModal("Copied!")
				}
			},
			func() {
				err := ui.OpenInNvim(selectedDatabase)
				if err != nil {
					a.showErrorModal(fmt.Sprintf("Edit failed: %v", err))
				}
			},
		)
		detailViewWithInstructions := ui.CreateDetailViewWithInstructions(detailView)
		a.pages.AddPage("rdsDatabaseDetail", detailViewWithInstructions, true, true)
		a.tviewApp.SetFocus(detailView)
	})

	a.setupTableYankFunctionality(a.rdsDatabaseTable, databases)
	rdsDatabaseListFlex := ui.WrapTableInFlex(a.rdsDatabaseTable)
	a.pages.AddPage(ui.PageRdsDatabases, rdsDatabaseListFlex, true, true)
	a.tviewApp.SetFocus(a.rdsDatabaseTable)
}

// switchToRdsAccountsView switches to RDS accounts view
func (a *App) switchToRdsAccountsView(instanceId string) {
	accounts, err := a.services.RDS.FetchAccounts(instanceId)
	if err != nil {
		a.showErrorModal(fmt.Sprintf("Failed to fetch accounts for instance %s: %v", instanceId, err))
		return
	}

	a.currentRdsInstanceId = instanceId
	a.rdsAccountTable = ui.CreateRdsAccountsListView(accounts, instanceId)

	ui.SetupTableNavigationWithSearch(a.rdsAccountTable, a, func(row, col int) {
		accountName := a.rdsAccountTable.GetCell(row, 0).GetReference().(string)
		var selectedAccount interface{}
		for _, account := range accounts {
			if account.AccountName == accountName {
				selectedAccount = account
				break
			}
		}
		a.currentDetailData = selectedAccount
		detailView, _ := ui.CreateInteractiveJSONDetailViewWithSearch(
			fmt.Sprintf("Account Details: %s", accountName),
			selectedAccount,
			a,
			func() {
				err := ui.CopyToClipboard(selectedAccount)
				if err != nil {
					a.showErrorModal(fmt.Sprintf("Copy failed: %v", err))
				} else {
					a.showErrorModal("Copied!")
				}
			},
			func() {
				err := ui.OpenInNvim(selectedAccount)
				if err != nil {
					a.showErrorModal(fmt.Sprintf("Edit failed: %v", err))
				}
			},
		)
		detailViewWithInstructions := ui.CreateDetailViewWithInstructions(detailView)
		a.pages.AddPage("rdsAccountDetail", detailViewWithInstructions, true, true)
		a.tviewApp.SetFocus(detailView)
	})

	a.setupTableYankFunctionality(a.rdsAccountTable, accounts)
	rdsAccountListFlex := ui.WrapTableInFlex(a.rdsAccountTable)
	a.pages.AddPage(ui.PageRdsAccounts, rdsAccountListFlex, true, true)
	a.tviewApp.SetFocus(a.rdsAccountTable)
}

// switchToRedisListView switches to Redis list view
func (a *App) switchToRedisListView() {
	instances, err := a.services.Redis.FetchInstances()
	if err != nil {
		a.showErrorModal(fmt.Sprintf("Failed to fetch Redis instances: %v", err))
		return
	}
	a.allRedisInstances = instances

	a.redisInstanceTable = ui.CreateRedisListView(instances)
	searchHandler := ui.SetupTableNavigationWithSearch(a.redisInstanceTable, a, func(row, col int) {
		instanceId := a.redisInstanceTable.GetCell(row, 0).GetReference().(string)
		var selectedInstance interface{}
		for _, inst := range instances {
			if inst.InstanceId == instanceId {
				selectedInstance = inst
				break
			}
		}
		a.currentDetailData = selectedInstance
		detailView, _ := ui.CreateInteractiveJSONDetailViewWithSearch(
			fmt.Sprintf("Redis Details: %s", instanceId),
			selectedInstance,
			a,
			func() {
				err := ui.CopyToClipboard(selectedInstance)
				if err != nil {
					a.showErrorModal(fmt.Sprintf("Failed to copy: %v", err))
				} else {
					a.showErrorModal("Copied!")
				}
			},
			func() {
				err := ui.OpenInNvim(selectedInstance)
				if err != nil {
					a.showErrorModal(fmt.Sprintf("Failed to edit: %v", err))
				}
			},
		)
		detailViewWithInstructions := ui.CreateDetailViewWithInstructions(detailView)
		a.pages.AddPage("redisDetail", detailViewWithInstructions, true, true)
		a.tviewApp.SetFocus(detailView)
	})

	a.setupTableYankFunctionality(a.redisInstanceTable, instances)
	a.setupRedisKeyHandlers(a.redisInstanceTable, searchHandler)

	redisListFlex := ui.WrapTableInFlex(a.redisInstanceTable)
	a.pages.AddPage(ui.PageRedisList, redisListFlex, true, true)
	a.tviewApp.SetFocus(a.redisInstanceTable)
}

// setupRedisKeyHandlers sets up 'A' key for Redis instance list
func (a *App) setupRedisKeyHandlers(table *tview.Table, searchHandler *ui.VimSearchHandler) {
	originalInputCapture := table.GetInputCapture()
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Check if the app-level search bar is active for the current handler
		isCurrentHandlerSearchActive := false
		if a.activeSearchHandler == searchHandler {
			frontPage, _ := a.searchBarContainer.GetFrontPage()
			if frontPage == "visible" { // "visible" is the page name for the search bar
				isCurrentHandlerSearchActive = true
			}
		}

		if isCurrentHandlerSearchActive {
			switch event.Rune() {
			case 'n', 'N', '/':
				// Delegate to the app's shared search bar input capture
				return a.searchBar.GetInputCapture()(event)
			}
		}

		switch event.Rune() {
		case 'A':
			row, _ := table.GetSelection()
			if row > 0 { // Skip header
				cell := table.GetCell(row, 0)
				if instanceId, ok := cell.GetReference().(string); ok {
					a.switchToRedisAccountsView(instanceId)
				}
			}
			return nil
		}
		if originalInputCapture != nil {
			return originalInputCapture(event)
		}
		return event
	})
}

// switchToRedisAccountsView switches to Redis accounts view for a given instance
func (a *App) switchToRedisAccountsView(instanceId string) {
	accounts, err := a.services.Redis.FetchAccounts(instanceId)
	if err != nil {
		a.showErrorModal(fmt.Sprintf("Failed to fetch accounts for Redis instance %s: %v", instanceId, err))
		return
	}

	a.currentRedisInstanceId = instanceId
	a.redisAccountTable = ui.CreateRedisAccountsListView(accounts, instanceId)

	ui.SetupTableNavigationWithSearch(a.redisAccountTable, a, func(row, col int) {
		accountName := a.redisAccountTable.GetCell(row, 0).GetReference().(string)
		var selectedAccount interface{}
		for _, account := range accounts {
			if account.AccountName == accountName {
				selectedAccount = account
				break
			}
		}
		a.currentDetailData = selectedAccount
		detailView, _ := ui.CreateInteractiveJSONDetailViewWithSearch(
			fmt.Sprintf("Redis Account Details: %s", accountName),
			selectedAccount,
			a,
			func() {
				err := ui.CopyToClipboard(selectedAccount)
				if err != nil {
					a.showErrorModal(fmt.Sprintf("Failed to copy: %v", err))
				} else {
					a.showErrorModal("Copied!")
				}
			},
			func() {
				err := ui.OpenInNvim(selectedAccount)
				if err != nil {
					a.showErrorModal(fmt.Sprintf("Failed to edit: %v", err))
				}
			},
		)
		detailViewWithInstructions := ui.CreateDetailViewWithInstructions(detailView)
		a.pages.AddPage("redisAccountDetail", detailViewWithInstructions, true, true)
		a.tviewApp.SetFocus(detailView)
	})

	a.setupTableYankFunctionality(a.redisAccountTable, accounts)

	redisAccountListFlex := ui.WrapTableInFlex(a.redisAccountTable)
	a.pages.AddPage(ui.PageRedisAccounts, redisAccountListFlex, true, true)
	a.tviewApp.SetFocus(a.redisAccountTable)
}
