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
	case ui.PageEcsList, ui.PageSecurityGroups, ui.PageDnsDomains, ui.PageSlbList, ui.PageOssBuckets, ui.PageRdsList, ui.PageRedisList, ui.PageRocketMQList:
		a.handleNavigation(ui.PageMainMenu, a.mainMenu)
	case ui.PageEcsDetail:
		a.handleNavigation(ui.PageEcsList, a.ecsInstanceTable)
	case ui.PageSecurityGroupDetail:
		a.handleNavigation(ui.PageSecurityGroups, a.securityGroupTable)
	case ui.PageSecurityGroupRules:
		a.handleNavigation(ui.PageSecurityGroups, a.securityGroupTable)
	case ui.PageSecurityGroupInstances:
		a.handleNavigation(ui.PageSecurityGroups, a.securityGroupTable)
	case ui.PageInstanceSecurityGroups:
		a.handleNavigation(ui.PageEcsList, a.ecsInstanceTable)
	case ui.PageDnsRecords:
		a.handleNavigation(ui.PageDnsDomains, a.dnsDomainsTable)
	case ui.PageSlbDetail:
		a.handleNavigation(ui.PageSlbList, a.slbInstanceTable)
	case ui.PageSlbListeners:
		a.handleNavigation(ui.PageSlbList, a.slbInstanceTable)
	case ui.PageSlbVServerGroups:
		a.handleNavigation(ui.PageSlbList, a.slbInstanceTable)
	case ui.PageSlbVServerGroupBackendServers:
		a.handleNavigation(ui.PageSlbVServerGroups, a.slbVServerGroupsTable)
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
	case ui.PageRocketMQTopics:
		a.handleNavigation(ui.PageRocketMQList, a.rocketmqInstanceTable)
	case ui.PageRocketMQGroups:
		a.handleNavigation(ui.PageRocketMQList, a.rocketmqInstanceTable)
	case "rocketmqDetail":
		a.handleNavigation(ui.PageRocketMQList, a.rocketmqInstanceTable)
	case "rocketmqTopicDetail":
		a.handleNavigation(ui.PageRocketMQTopics, a.rocketmqTopicsTable)
	case "rocketmqGroupDetail":
		a.handleNavigation(ui.PageRocketMQGroups, a.rocketmqGroupsTable)
	}
}

// handleBackKey handles 'q' key navigation
func (a *App) handleBackKey(currentPageName string) {
	switch currentPageName {
	case ui.PageMainMenu:
		return
	case ui.PageEcsList, ui.PageSecurityGroups, ui.PageDnsDomains, ui.PageSlbList, ui.PageOssBuckets, ui.PageRdsList, ui.PageRedisList, ui.PageRocketMQList:
		a.handleNavigation(ui.PageMainMenu, a.mainMenu)
	case ui.PageEcsDetail:
		a.handleNavigation(ui.PageEcsList, a.ecsInstanceTable)
	case ui.PageSecurityGroupDetail:
		a.handleNavigation(ui.PageSecurityGroups, a.securityGroupTable)
	case ui.PageSecurityGroupRules:
		a.handleNavigation(ui.PageSecurityGroups, a.securityGroupTable)
	case ui.PageSecurityGroupInstances:
		a.handleNavigation(ui.PageSecurityGroups, a.securityGroupTable)
	case ui.PageInstanceSecurityGroups:
		a.handleNavigation(ui.PageEcsList, a.ecsInstanceTable)
	case ui.PageDnsRecords:
		a.handleNavigation(ui.PageDnsDomains, a.dnsDomainsTable)
	case ui.PageSlbDetail:
		a.handleNavigation(ui.PageSlbList, a.slbInstanceTable)
	case ui.PageSlbListeners:
		a.handleNavigation(ui.PageSlbList, a.slbInstanceTable)
	case ui.PageSlbVServerGroups:
		a.handleNavigation(ui.PageSlbList, a.slbInstanceTable)
	case ui.PageSlbVServerGroupBackendServers:
		a.handleNavigation(ui.PageSlbVServerGroups, a.slbVServerGroupsTable)
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
	case ui.PageRocketMQTopics:
		a.handleNavigation(ui.PageRocketMQList, a.rocketmqInstanceTable)
	case ui.PageRocketMQGroups:
		a.handleNavigation(ui.PageRocketMQList, a.rocketmqInstanceTable)
	case "rocketmqDetail":
		a.handleNavigation(ui.PageRocketMQList, a.rocketmqInstanceTable)
	case "rocketmqTopicDetail":
		a.handleNavigation(ui.PageRocketMQTopics, a.rocketmqTopicsTable)
	case "rocketmqGroupDetail":
		a.handleNavigation(ui.PageRocketMQGroups, a.rocketmqGroupsTable)
	}
}

// handleNavigation handles page navigation
func (a *App) handleNavigation(targetPage string, focusItem tview.Primitive) {
	a.pages.SwitchToPage(targetPage)

	// Update mode line with shortcuts for the current page
	ui.UpdateModeLineWithShortcuts(a.modeLine, a.currentProfile, targetPage)

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
				err := ui.OpenInNvimWithSuspend(selectedInstance, a.tviewApp)
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

		// Update mode line with shortcuts for ECS detail page
		ui.UpdateModeLineWithShortcuts(a.modeLine, a.currentProfile, ui.PageEcsDetail)

		a.tviewApp.SetFocus(a.ecsDetailView)
	})

	a.setupTableYankFunctionality(a.ecsInstanceTable, a.allECSInstances)
	a.setupEcsKeyHandlers(a.ecsInstanceTable)
	ecsListFlex := ui.WrapTableInFlex(a.ecsInstanceTable)
	a.pages.AddPage(ui.PageEcsList, ecsListFlex, true, true)

	// Update mode line with shortcuts for ECS list page
	ui.UpdateModeLineWithShortcuts(a.modeLine, a.currentProfile, ui.PageEcsList)

	a.tviewApp.SetFocus(a.ecsInstanceTable)
}

// setupEcsKeyHandlers sets up key handlers for ECS specific actions
func (a *App) setupEcsKeyHandlers(table *tview.Table) {
	originalInputCapture := table.GetInputCapture()

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'g': // g key handler for security groups of this instance
			row, _ := table.GetSelection()
			if row > 0 { // Skip header row
				if cell := table.GetCell(row, 0); cell != nil {
					if instanceId, ok := cell.GetReference().(string); ok {
						a.switchToInstanceSecurityGroupsView(instanceId)
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

// switchToSecurityGroupsListView switches to security groups list view
func (a *App) switchToSecurityGroupsListView() {
	if a.allSecurityGroups == nil {
		securityGroups, err := a.services.ECS.FetchSecurityGroups()
		if err != nil {
			a.showErrorModal(err.Error())
			return
		}
		a.allSecurityGroups = securityGroups
	}
	a.securityGroupTable = ui.CreateSecurityGroupsListView(a.allSecurityGroups)
	ui.SetupTableNavigationWithSearch(a.securityGroupTable, a, func(row, col int) {
		securityGroupId := a.securityGroupTable.GetCell(row, 0).GetReference().(string)
		// 回车键进入安全组规则列表
		a.switchToSecurityGroupRulesView(securityGroupId)
	})

	a.setupTableYankFunctionality(a.securityGroupTable, a.allSecurityGroups)
	a.setupSecurityGroupKeyHandlers(a.securityGroupTable)
	securityGroupListFlex := ui.WrapTableInFlex(a.securityGroupTable)
	a.pages.AddPage(ui.PageSecurityGroups, securityGroupListFlex, true, true)

	// Update mode line with shortcuts for security groups page
	ui.UpdateModeLineWithShortcuts(a.modeLine, a.currentProfile, ui.PageSecurityGroups)

	a.tviewApp.SetFocus(a.securityGroupTable)
}

// setupSecurityGroupKeyHandlers sets up key handlers for security group specific actions
func (a *App) setupSecurityGroupKeyHandlers(table *tview.Table) {
	originalInputCapture := table.GetInputCapture()

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 's': // s key handler for instances using this security group
			row, _ := table.GetSelection()
			if row > 0 { // Skip header row
				if cell := table.GetCell(row, 0); cell != nil {
					if securityGroupId, ok := cell.GetReference().(string); ok {
						a.switchToSecurityGroupInstancesView(securityGroupId)
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

// switchToSecurityGroupRulesView switches to security group rules view
func (a *App) switchToSecurityGroupRulesView(securityGroupId string) {
	rulesResponse, err := a.services.ECS.FetchSecurityGroupRules(securityGroupId)
	if err != nil {
		a.showErrorModal(fmt.Sprintf("Failed to fetch security group rules for %s: %v", securityGroupId, err))
		return
	}

	a.securityGroupRulesTable = ui.CreateSecurityGroupRulesView(rulesResponse)
	ui.SetupTableNavigationWithSearch(a.securityGroupRulesTable, a, nil)

	a.setupTableYankFunctionality(a.securityGroupRulesTable, rulesResponse)
	securityGroupRulesListFlex := ui.WrapTableInFlex(a.securityGroupRulesTable)
	a.pages.AddPage(ui.PageSecurityGroupRules, securityGroupRulesListFlex, true, true)

	// Update mode line with shortcuts for security group rules page
	ui.UpdateModeLineWithShortcuts(a.modeLine, a.currentProfile, ui.PageSecurityGroupRules)

	a.tviewApp.SetFocus(a.securityGroupRulesTable)
}

// switchToSecurityGroupInstancesView switches to instances using a security group
func (a *App) switchToSecurityGroupInstancesView(securityGroupId string) {
	instances, err := a.services.ECS.FetchInstancesBySecurityGroup(securityGroupId)
	if err != nil {
		a.showErrorModal(fmt.Sprintf("Failed to fetch instances for security group %s: %v", securityGroupId, err))
		return
	}

	a.securityGroupInstancesTable = ui.CreateSecurityGroupInstancesView(instances, securityGroupId)
	ui.SetupTableNavigationWithSearch(a.securityGroupInstancesTable, a, func(row, col int) {
		instanceId := a.securityGroupInstancesTable.GetCell(row, 0).GetReference().(string)
		var selectedInstance interface{}
		for _, inst := range instances {
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
				err := ui.OpenInNvimWithSuspend(selectedInstance, a.tviewApp)
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

		// Update mode line with shortcuts for ECS detail page
		ui.UpdateModeLineWithShortcuts(a.modeLine, a.currentProfile, ui.PageEcsDetail)

		a.tviewApp.SetFocus(a.ecsDetailView)
	})

	a.setupTableYankFunctionality(a.securityGroupInstancesTable, instances)
	securityGroupInstancesListFlex := ui.WrapTableInFlex(a.securityGroupInstancesTable)
	a.pages.AddPage(ui.PageSecurityGroupInstances, securityGroupInstancesListFlex, true, true)

	// Update mode line with shortcuts for security group instances page
	ui.UpdateModeLineWithShortcuts(a.modeLine, a.currentProfile, ui.PageSecurityGroupInstances)

	a.tviewApp.SetFocus(a.securityGroupInstancesTable)
}

// switchToInstanceSecurityGroupsView switches to security groups for an instance
func (a *App) switchToInstanceSecurityGroupsView(instanceId string) {
	securityGroups, err := a.services.ECS.FetchSecurityGroupsByInstance(instanceId)
	if err != nil {
		a.showErrorModal(fmt.Sprintf("Failed to fetch security groups for instance %s: %v", instanceId, err))
		return
	}

	a.instanceSecurityGroupsTable = ui.CreateInstanceSecurityGroupsView(securityGroups, instanceId)
	ui.SetupTableNavigationWithSearch(a.instanceSecurityGroupsTable, a, func(row, col int) {
		securityGroupId := a.instanceSecurityGroupsTable.GetCell(row, 0).GetReference().(string)
		// 进入安全组规则视图
		a.switchToSecurityGroupRulesView(securityGroupId)
	})

	a.setupTableYankFunctionality(a.instanceSecurityGroupsTable, securityGroups)
	instanceSecurityGroupsListFlex := ui.WrapTableInFlex(a.instanceSecurityGroupsTable)
	a.pages.AddPage(ui.PageInstanceSecurityGroups, instanceSecurityGroupsListFlex, true, true)

	// Update mode line with shortcuts for instance security groups page
	ui.UpdateModeLineWithShortcuts(a.modeLine, a.currentProfile, ui.PageInstanceSecurityGroups)

	a.tviewApp.SetFocus(a.instanceSecurityGroupsTable)
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

	// Update mode line with shortcuts for DNS domains page
	ui.UpdateModeLineWithShortcuts(a.modeLine, a.currentProfile, ui.PageDnsDomains)

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

	// Update mode line with shortcuts for DNS records page
	ui.UpdateModeLineWithShortcuts(a.modeLine, a.currentProfile, ui.PageDnsRecords)

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
				err := ui.OpenInNvimWithSuspend(selectedSlb, a.tviewApp)
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

		// Update mode line with shortcuts for SLB detail page
		ui.UpdateModeLineWithShortcuts(a.modeLine, a.currentProfile, ui.PageSlbDetail)

		a.tviewApp.SetFocus(a.slbDetailView)
	})

	a.setupTableYankFunctionality(a.slbInstanceTable, a.allSLBInstances)
	a.setupSlbKeyHandlers(a.slbInstanceTable)
	slbListFlex := ui.WrapTableInFlex(a.slbInstanceTable)
	a.pages.AddPage(ui.PageSlbList, slbListFlex, true, true)

	// Update mode line with shortcuts for SLB list page
	ui.UpdateModeLineWithShortcuts(a.modeLine, a.currentProfile, ui.PageSlbList)

	a.tviewApp.SetFocus(a.slbInstanceTable)
}

// setupSlbKeyHandlers sets up key handlers for SLB specific actions
func (a *App) setupSlbKeyHandlers(table *tview.Table) {
	originalInputCapture := table.GetInputCapture()

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'l': // l key handler for listeners of this SLB instance
			row, _ := table.GetSelection()
			if row > 0 { // Skip header row
				if cell := table.GetCell(row, 0); cell != nil {
					if loadBalancerId, ok := cell.GetReference().(string); ok {
						a.switchToSlbListenersView(loadBalancerId)
					}
				}
			}
			return nil
		case 'v': // v key handler for virtual server groups of this SLB instance
			row, _ := table.GetSelection()
			if row > 0 { // Skip header row
				if cell := table.GetCell(row, 0); cell != nil {
					if loadBalancerId, ok := cell.GetReference().(string); ok {
						a.switchToSlbVServerGroupsView(loadBalancerId)
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

// switchToSlbListenersView switches to SLB listeners view
func (a *App) switchToSlbListenersView(loadBalancerId string) {
	detailedListeners, err := a.services.SLB.FetchDetailedListeners(loadBalancerId)
	if err != nil {
		a.showErrorModal(fmt.Sprintf("Failed to fetch listeners for SLB %s: %v", loadBalancerId, err))
		return
	}

	a.slbListenersTable = ui.CreateSlbDetailedListenersView(detailedListeners, loadBalancerId)
	ui.SetupTableNavigationWithSearch(a.slbListenersTable, a, nil)

	a.setupTableYankFunctionality(a.slbListenersTable, detailedListeners)
	slbListenersListFlex := ui.WrapTableInFlex(a.slbListenersTable)
	a.pages.AddPage(ui.PageSlbListeners, slbListenersListFlex, true, true)

	// Update mode line with shortcuts for SLB listeners page
	ui.UpdateModeLineWithShortcuts(a.modeLine, a.currentProfile, ui.PageSlbListeners)

	a.tviewApp.SetFocus(a.slbListenersTable)
}

// switchToSlbVServerGroupsView switches to SLB virtual server groups view
func (a *App) switchToSlbVServerGroupsView(loadBalancerId string) {
	detailedVServerGroups, err := a.services.SLB.FetchDetailedVServerGroups(loadBalancerId)
	if err != nil {
		a.showErrorModal(fmt.Sprintf("Failed to fetch virtual server groups for SLB %s: %v", loadBalancerId, err))
		return
	}

	a.slbVServerGroupsTable = ui.CreateSlbDetailedVServerGroupsView(detailedVServerGroups, loadBalancerId)
	ui.SetupTableNavigationWithSearch(a.slbVServerGroupsTable, a, func(row, col int) {
		vServerGroupId := a.slbVServerGroupsTable.GetCell(row, 0).GetReference().(string)
		a.switchToSlbVServerGroupBackendServersView(vServerGroupId)
	})

	a.setupTableYankFunctionality(a.slbVServerGroupsTable, detailedVServerGroups)
	slbVServerGroupsListFlex := ui.WrapTableInFlex(a.slbVServerGroupsTable)
	a.pages.AddPage(ui.PageSlbVServerGroups, slbVServerGroupsListFlex, true, true)

	// Update mode line with shortcuts for SLB VServer groups page
	ui.UpdateModeLineWithShortcuts(a.modeLine, a.currentProfile, ui.PageSlbVServerGroups)

	a.tviewApp.SetFocus(a.slbVServerGroupsTable)
}

// switchToSlbVServerGroupBackendServersView switches to SLB virtual server group backend servers view
func (a *App) switchToSlbVServerGroupBackendServersView(vServerGroupId string) {
	detailedBackendServers, err := a.services.SLB.FetchDetailedBackendServers(vServerGroupId, a.clients.ECS)
	if err != nil {
		a.showErrorModal(fmt.Sprintf("Failed to fetch backend servers for virtual server group %s: %v", vServerGroupId, err))
		return
	}

	a.slbVServerGroupBackendServersTable = ui.CreateSlbDetailedBackendServersView(detailedBackendServers, vServerGroupId)
	ui.SetupTableNavigationWithSearch(a.slbVServerGroupBackendServersTable, a, nil)

	a.setupTableYankFunctionality(a.slbVServerGroupBackendServersTable, detailedBackendServers)
	slbVServerGroupBackendServersListFlex := ui.WrapTableInFlex(a.slbVServerGroupBackendServersTable)
	a.pages.AddPage(ui.PageSlbVServerGroupBackendServers, slbVServerGroupBackendServersListFlex, true, true)

	// Update mode line with shortcuts for SLB backend servers page
	ui.UpdateModeLineWithShortcuts(a.modeLine, a.currentProfile, ui.PageSlbVServerGroupBackendServers)

	a.tviewApp.SetFocus(a.slbVServerGroupBackendServersTable)
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

	// Update mode line with shortcuts for OSS buckets page
	ui.UpdateModeLineWithShortcuts(a.modeLine, a.currentProfile, ui.PageOssBuckets)

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
	ui.UpdateModeLineWithPageInfoAndShortcuts(a.modeLine, a.currentProfile, ui.PageOssObjects, pageInfo)

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
						err := ui.OpenInNvimWithSuspend(obj, a.tviewApp)
						if err != nil {
							a.showErrorModal(fmt.Sprintf("Failed to edit: %v", err))
						}
					},
				)
				a.ossDetailView = view
				a.pages.AddPage("ossObjectDetail", a.ossDetailView, true, true)

				// Update mode line with shortcuts for OSS object detail page
				ui.UpdateModeLineWithShortcuts(a.modeLine, a.currentProfile, "ossObjectDetail")

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
							case []slb.VServerGroup:
								for _, vsg := range items {
									if vsg.VServerGroupId == ref.(string) {
										rowData = vsg
										break
									}
								}
							case []slb.BackendServerInDescribeVServerGroupAttribute:
								for _, server := range items {
									if server.ServerId == ref.(string) {
										rowData = server
										break
									}
								}
							case *slb.DescribeLoadBalancerAttributeResponse:
								// For listeners response, we'll copy the entire response
								rowData = items
							case []service.ListenerDetail:
								for _, listener := range items {
									if fmt.Sprintf("%d", listener.Port) == ref.(string) {
										rowData = listener
										break
									}
								}
							case []service.VServerGroupDetail:
								for _, vsg := range items {
									if vsg.VServerGroupId == ref.(string) {
										rowData = vsg
										break
									}
								}
							case []service.BackendServerDetail:
								for _, server := range items {
									if server.ServerId == ref.(string) {
										rowData = server
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
		ECS:      service.NewECSService(newClients.ECS),
		DNS:      service.NewDNSService(newClients.DNS),
		SLB:      service.NewSLBService(newClients.SLB),
		RDS:      service.NewRDSService(newClients.RDS),
		OSS:      service.NewOSSServiceWithCredentials(newClients.OSS, cfg.AccessKeyID, cfg.AccessKeySecret, cfg.OssEndpoint),
		Redis:    service.NewRedisService(newClients.Redis),
		RocketMQ: service.NewRocketMQService(newClients.RocketMQ),
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
	a.allSecurityGroups = nil
	a.allDomains = nil
	a.allSLBInstances = nil
	a.allRDSInstances = nil
	a.allRedisInstances = nil
	a.allRocketMQInstances = nil
	a.allOssBuckets = nil
	a.currentBucketName = ""
	a.currentRdsInstanceId = ""
	a.currentRedisInstanceId = ""
	a.currentRocketMQInstanceId = ""

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
				err := ui.OpenInNvimWithSuspend(selectedInstance, a.tviewApp)
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

		// Update mode line with shortcuts for RDS detail page
		ui.UpdateModeLineWithShortcuts(a.modeLine, a.currentProfile, ui.PageRdsDetail)

		a.tviewApp.SetFocus(a.rdsDetailView)
	})

	a.setupTableYankFunctionality(a.rdsInstanceTable, a.allRDSInstances)
	a.setupRdsKeyHandlers(a.rdsInstanceTable)
	rdsListFlex := ui.WrapTableInFlex(a.rdsInstanceTable)
	a.pages.AddPage(ui.PageRdsList, rdsListFlex, true, true)

	// Update mode line with shortcuts for RDS list page
	ui.UpdateModeLineWithShortcuts(a.modeLine, a.currentProfile, ui.PageRdsList)

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
				err := ui.OpenInNvimWithSuspend(selectedDatabase, a.tviewApp)
				if err != nil {
					a.showErrorModal(fmt.Sprintf("Edit failed: %v", err))
				}
			},
		)
		detailViewWithInstructions := ui.CreateDetailViewWithInstructions(detailView)
		a.pages.AddPage("rdsDatabaseDetail", detailViewWithInstructions, true, true)

		// Update mode line with shortcuts for RDS database detail page
		ui.UpdateModeLineWithShortcuts(a.modeLine, a.currentProfile, "rdsDatabaseDetail")

		a.tviewApp.SetFocus(detailView)
	})

	a.setupTableYankFunctionality(a.rdsDatabaseTable, databases)
	rdsDatabaseListFlex := ui.WrapTableInFlex(a.rdsDatabaseTable)
	a.pages.AddPage(ui.PageRdsDatabases, rdsDatabaseListFlex, true, true)

	// Update mode line with shortcuts for RDS databases page
	ui.UpdateModeLineWithShortcuts(a.modeLine, a.currentProfile, ui.PageRdsDatabases)

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
				err := ui.OpenInNvimWithSuspend(selectedAccount, a.tviewApp)
				if err != nil {
					a.showErrorModal(fmt.Sprintf("Edit failed: %v", err))
				}
			},
		)
		detailViewWithInstructions := ui.CreateDetailViewWithInstructions(detailView)
		a.pages.AddPage("rdsAccountDetail", detailViewWithInstructions, true, true)

		// Update mode line with shortcuts for RDS account detail page
		ui.UpdateModeLineWithShortcuts(a.modeLine, a.currentProfile, "rdsAccountDetail")

		a.tviewApp.SetFocus(detailView)
	})

	a.setupTableYankFunctionality(a.rdsAccountTable, accounts)
	rdsAccountListFlex := ui.WrapTableInFlex(a.rdsAccountTable)
	a.pages.AddPage(ui.PageRdsAccounts, rdsAccountListFlex, true, true)

	// Update mode line with shortcuts for RDS accounts page
	ui.UpdateModeLineWithShortcuts(a.modeLine, a.currentProfile, ui.PageRdsAccounts)

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
				err := ui.OpenInNvimWithSuspend(selectedInstance, a.tviewApp)
				if err != nil {
					a.showErrorModal(fmt.Sprintf("Failed to edit: %v", err))
				}
			},
		)
		detailViewWithInstructions := ui.CreateDetailViewWithInstructions(detailView)
		a.pages.AddPage("redisDetail", detailViewWithInstructions, true, true)

		// Update mode line with shortcuts for Redis detail page
		ui.UpdateModeLineWithShortcuts(a.modeLine, a.currentProfile, "redisDetail")

		a.tviewApp.SetFocus(detailView)
	})

	a.setupTableYankFunctionality(a.redisInstanceTable, instances)
	a.setupRedisKeyHandlers(a.redisInstanceTable, searchHandler)

	redisListFlex := ui.WrapTableInFlex(a.redisInstanceTable)
	a.pages.AddPage(ui.PageRedisList, redisListFlex, true, true)

	// Update mode line with shortcuts for Redis list page
	ui.UpdateModeLineWithShortcuts(a.modeLine, a.currentProfile, ui.PageRedisList)

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
				err := ui.OpenInNvimWithSuspend(selectedAccount, a.tviewApp)
				if err != nil {
					a.showErrorModal(fmt.Sprintf("Failed to edit: %v", err))
				}
			},
		)
		detailViewWithInstructions := ui.CreateDetailViewWithInstructions(detailView)
		a.pages.AddPage("redisAccountDetail", detailViewWithInstructions, true, true)

		// Update mode line with shortcuts for Redis account detail page
		ui.UpdateModeLineWithShortcuts(a.modeLine, a.currentProfile, "redisAccountDetail")

		a.tviewApp.SetFocus(detailView)
	})

	a.setupTableYankFunctionality(a.redisAccountTable, accounts)

	redisAccountListFlex := ui.WrapTableInFlex(a.redisAccountTable)
	a.pages.AddPage(ui.PageRedisAccounts, redisAccountListFlex, true, true)

	// Update mode line with shortcuts for Redis accounts page
	ui.UpdateModeLineWithShortcuts(a.modeLine, a.currentProfile, ui.PageRedisAccounts)

	a.tviewApp.SetFocus(a.redisAccountTable)
}

// switchToRocketMQListView switches to RocketMQ list view
func (a *App) switchToRocketMQListView() {
	if a.allRocketMQInstances == nil {
		instances, err := a.services.RocketMQ.FetchInstances()
		if err != nil {
			a.showErrorModal(fmt.Sprintf("Failed to fetch RocketMQ instances: %v", err))
			return
		}
		a.allRocketMQInstances = instances
	}

	a.rocketmqInstanceTable = ui.CreateRocketMQListView(a.allRocketMQInstances)
	searchHandler := ui.SetupTableNavigationWithSearch(a.rocketmqInstanceTable, a, func(row, col int) {
		instanceId := a.rocketmqInstanceTable.GetCell(row, 0).GetReference().(string)
		var selectedInstance interface{}
		for _, inst := range a.allRocketMQInstances {
			if inst.InstanceId == instanceId {
				selectedInstance = inst
				break
			}
		}
		a.currentDetailData = selectedInstance
		detailView, _ := ui.CreateInteractiveJSONDetailViewWithSearch(
			fmt.Sprintf("RocketMQ Details: %s", instanceId),
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
				err := ui.OpenInNvimWithSuspend(selectedInstance, a.tviewApp)
				if err != nil {
					a.showErrorModal(fmt.Sprintf("Failed to edit: %v", err))
				}
			},
		)
		detailViewWithInstructions := ui.CreateDetailViewWithInstructions(detailView)
		a.pages.AddPage("rocketmqDetail", detailViewWithInstructions, true, true)

		// Update mode line with shortcuts for RocketMQ detail page
		ui.UpdateModeLineWithShortcuts(a.modeLine, a.currentProfile, "rocketmqDetail")

		a.tviewApp.SetFocus(detailView)
	})

	a.setupTableYankFunctionality(a.rocketmqInstanceTable, a.allRocketMQInstances)
	a.setupRocketMQKeyHandlers(a.rocketmqInstanceTable, searchHandler)

	rocketmqListFlex := ui.WrapTableInFlex(a.rocketmqInstanceTable)
	a.pages.AddPage(ui.PageRocketMQList, rocketmqListFlex, true, true)

	// Update mode line with shortcuts for RocketMQ list page
	ui.UpdateModeLineWithShortcuts(a.modeLine, a.currentProfile, ui.PageRocketMQList)

	a.tviewApp.SetFocus(a.rocketmqInstanceTable)
}

// setupRocketMQKeyHandlers sets up 'T' and 'G' keys for RocketMQ instance list
func (a *App) setupRocketMQKeyHandlers(table *tview.Table, searchHandler *ui.VimSearchHandler) {
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
		case 'T':
			row, _ := table.GetSelection()
			if row > 0 { // Skip header
				cell := table.GetCell(row, 0)
				if instanceId, ok := cell.GetReference().(string); ok {
					a.switchToRocketMQTopicsView(instanceId)
				}
			}
			return nil
		case 'G':
			row, _ := table.GetSelection()
			if row > 0 { // Skip header
				cell := table.GetCell(row, 0)
				if instanceId, ok := cell.GetReference().(string); ok {
					a.switchToRocketMQGroupsView(instanceId)
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

// switchToRocketMQTopicsView switches to RocketMQ topics view for a given instance
func (a *App) switchToRocketMQTopicsView(instanceId string) {
	topics, err := a.services.RocketMQ.FetchTopics(instanceId)
	if err != nil {
		a.showErrorModal(fmt.Sprintf("Failed to fetch topics for RocketMQ instance %s: %v", instanceId, err))
		return
	}

	a.currentRocketMQInstanceId = instanceId
	a.rocketmqTopicsTable = ui.CreateRocketMQTopicsListView(topics, instanceId)

	ui.SetupTableNavigationWithSearch(a.rocketmqTopicsTable, a, func(row, col int) {
		topicName := a.rocketmqTopicsTable.GetCell(row, 0).GetReference().(string)
		var selectedTopic interface{}
		for _, topic := range topics {
			if topic.Topic == topicName {
				selectedTopic = topic
				break
			}
		}
		a.currentDetailData = selectedTopic
		detailView, _ := ui.CreateInteractiveJSONDetailViewWithSearch(
			fmt.Sprintf("Topic Details: %s", topicName),
			selectedTopic,
			a,
			func() {
				err := ui.CopyToClipboard(selectedTopic)
				if err != nil {
					a.showErrorModal(fmt.Sprintf("Failed to copy: %v", err))
				} else {
					a.showErrorModal("Copied!")
				}
			},
			func() {
				err := ui.OpenInNvimWithSuspend(selectedTopic, a.tviewApp)
				if err != nil {
					a.showErrorModal(fmt.Sprintf("Failed to edit: %v", err))
				}
			},
		)
		detailViewWithInstructions := ui.CreateDetailViewWithInstructions(detailView)
		a.pages.AddPage("rocketmqTopicDetail", detailViewWithInstructions, true, true)

		// Update mode line with shortcuts for RocketMQ topic detail page
		ui.UpdateModeLineWithShortcuts(a.modeLine, a.currentProfile, "rocketmqTopicDetail")

		a.tviewApp.SetFocus(detailView)
	})

	a.setupTableYankFunctionality(a.rocketmqTopicsTable, topics)

	rocketmqTopicsListFlex := ui.WrapTableInFlex(a.rocketmqTopicsTable)
	a.pages.AddPage(ui.PageRocketMQTopics, rocketmqTopicsListFlex, true, true)

	// Update mode line with shortcuts for RocketMQ topics page
	ui.UpdateModeLineWithShortcuts(a.modeLine, a.currentProfile, ui.PageRocketMQTopics)

	a.tviewApp.SetFocus(a.rocketmqTopicsTable)
}

// switchToRocketMQGroupsView switches to RocketMQ groups view for a given instance
func (a *App) switchToRocketMQGroupsView(instanceId string) {
	groups, err := a.services.RocketMQ.FetchGroups(instanceId)
	if err != nil {
		a.showErrorModal(fmt.Sprintf("Failed to fetch groups for RocketMQ instance %s: %v", instanceId, err))
		return
	}

	a.currentRocketMQInstanceId = instanceId
	a.rocketmqGroupsTable = ui.CreateRocketMQGroupsListView(groups, instanceId)

	ui.SetupTableNavigationWithSearch(a.rocketmqGroupsTable, a, func(row, col int) {
		groupId := a.rocketmqGroupsTable.GetCell(row, 0).GetReference().(string)
		var selectedGroup interface{}
		for _, group := range groups {
			if group.GroupId == groupId {
				selectedGroup = group
				break
			}
		}
		a.currentDetailData = selectedGroup
		detailView, _ := ui.CreateInteractiveJSONDetailViewWithSearch(
			fmt.Sprintf("Group Details: %s", groupId),
			selectedGroup,
			a,
			func() {
				err := ui.CopyToClipboard(selectedGroup)
				if err != nil {
					a.showErrorModal(fmt.Sprintf("Failed to copy: %v", err))
				} else {
					a.showErrorModal("Copied!")
				}
			},
			func() {
				err := ui.OpenInNvimWithSuspend(selectedGroup, a.tviewApp)
				if err != nil {
					a.showErrorModal(fmt.Sprintf("Failed to edit: %v", err))
				}
			},
		)
		detailViewWithInstructions := ui.CreateDetailViewWithInstructions(detailView)
		a.pages.AddPage("rocketmqGroupDetail", detailViewWithInstructions, true, true)

		// Update mode line with shortcuts for RocketMQ group detail page
		ui.UpdateModeLineWithShortcuts(a.modeLine, a.currentProfile, "rocketmqGroupDetail")

		a.tviewApp.SetFocus(detailView)
	})

	a.setupTableYankFunctionality(a.rocketmqGroupsTable, groups)

	rocketmqGroupsListFlex := ui.WrapTableInFlex(a.rocketmqGroupsTable)
	a.pages.AddPage(ui.PageRocketMQGroups, rocketmqGroupsListFlex, true, true)

	// Update mode line with shortcuts for RocketMQ groups page
	ui.UpdateModeLineWithShortcuts(a.modeLine, a.currentProfile, ui.PageRocketMQGroups)

	a.tviewApp.SetFocus(a.rocketmqGroupsTable)
}
