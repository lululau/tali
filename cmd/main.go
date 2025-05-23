package main

import (
	"fmt"
	"os"
	"strings"
	// "github.com/aliyun/aliyun-oss-go-sdk/oss" // STUBBED: Go version compatibility issue

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds" // Added RDS
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Global variables
var (
	app             *tview.Application
	pages           *tview.Pages
	mainMenu        *tview.List
	ecsInstanceTable *tview.Table
	ecsDetailForm   *tview.Form
	dnsDomainsTable *tview.Table
	dnsRecordsTable *tview.Table
	slbInstanceTable *tview.Table
	slbDetailForm   *tview.Form
	ossBucketTable  *tview.Table
	ossObjectTable  *tview.Table
	rdsInstanceTable *tview.Table // For RDS list
	rdsDetailForm   *tview.Form   // Specific for RDS details


	allECSInstances []ecs.Instance
	allDomains      []alidns.Domain
	allSLBInstances []slb.LoadBalancer
	allRDSInstances []rds.DBInstance // To store RDS instances
	// allOssBuckets   []StubBucketInfo // STUBBED
	// currentBucketName string       // STUBBED

	accessKeyID     string
	accessKeySecret string
	regionID        string
	ossEndpoint     string
)

const (
	pageMainMenu      = "mainMenu"
	pageEcsList       = "ecsList"
	pageEcsDetail     = "ecsDetail"
	pageDnsDomains    = "dnsDomains"
	pageDnsRecords    = "dnsRecords"
	pageSlbList       = "slbList"
	pageSlbDetail     = "slbDetail"
	pageOssBuckets    = "ossBuckets"
	pageOssObjects    = "ossObjects"
	pageRdsList       = "rdsList"   // New page for RDS list
	pageRdsDetail     = "rdsDetail" // New page for RDS detail
)

// --- SDK Stubs for OSS ---
type StubBucketInfo struct { Name, Location, CreationDate, StorageClass string }
type StubObjectInfo struct { Key string; Size int64; LastModified, StorageClass string }
func fetchOssBuckets() ([]StubBucketInfo, error) { /* STUBBED */ return []StubBucketInfo{{Name: "stub-bucket-1", Location: "oss-cn-hangzhou"}}, nil }
func fetchOssObjects(bucketName string) ([]StubObjectInfo, error) { /* STUBBED */ return []StubObjectInfo{{Key: "stub-object.txt", Size: 123}}, nil }


// --- ECS SDK Functions ---
func fetchECSInstances() ([]ecs.Instance, error) {
	client, err := ecs.NewClientWithAccessKey(regionID, accessKeyID, accessKeySecret)
	if err != nil { return nil, fmt.Errorf("creating ECS client: %w", err) }
	request := ecs.CreateDescribeInstancesRequest()
	request.Scheme = "https"
	response, err := client.DescribeInstances(request)
	if err != nil { return nil, fmt.Errorf("describing ECS instances: %w", err) }
	return response.Instances.Instance, nil
}

// --- AliDNS SDK Functions ---
func fetchAlidnsDomains() ([]alidns.Domain, error) {
	client, err := alidns.NewClientWithAccessKey(regionID, accessKeyID, accessKeySecret)
	if err != nil { return nil, fmt.Errorf("creating AliDNS client: %w", err) }
	request := alidns.CreateDescribeDomainsRequest()
	request.Scheme = "https"
	response, err := client.DescribeDomains(request)
	if err != nil { return nil, fmt.Errorf("describing AliDNS domains: %w", err) }
	return response.Domains.Domain, nil
}

func fetchAlidnsDomainRecords(domainName string) ([]alidns.Record, error) {
	client, err := alidns.NewClientWithAccessKey(regionID, accessKeyID, accessKeySecret)
	if err != nil { return nil, fmt.Errorf("creating AliDNS client: %w", err) }
	request := alidns.CreateDescribeDomainRecordsRequest()
	request.Scheme = "https"
	request.DomainName = domainName
	response, err := client.DescribeDomainRecords(request)
	if err != nil { return nil, fmt.Errorf("describing AliDNS domain records for %s: %w", domainName, err) }
	return response.DomainRecords.Record, nil
}

// --- SLB SDK Functions ---
func fetchSLBInstances() ([]slb.LoadBalancer, error) {
	client, err := slb.NewClientWithAccessKey(regionID, accessKeyID, accessKeySecret)
	if err != nil { return nil, fmt.Errorf("creating SLB client: %w", err) }
	request := slb.CreateDescribeLoadBalancersRequest()
	request.Scheme = "https"
	response, err := client.DescribeLoadBalancers(request)
	if err != nil { return nil, fmt.Errorf("describing SLB instances: %w", err) }
	return response.LoadBalancers.LoadBalancer, nil
}

// --- RDS SDK Functions ---
func fetchRDSInstances() ([]rds.DBInstance, error) {
	client, err := rds.NewClientWithAccessKey(regionID, accessKeyID, accessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("creating RDS client: %w", err)
	}
	request := rds.CreateDescribeDBInstancesRequest()
	request.Scheme = "https"
	// request.RegionId = regionID // Already part of client config
	response, err := client.DescribeDBInstances(request)
	if err != nil {
		return nil, fmt.Errorf("describing RDS instances: %w", err)
	}
	return response.Items.DBInstance, nil
}


// --- TUI Creation Functions ---

// Main Menu
func createMainMenu() *tview.List {
	list := tview.NewList().
		AddItem("ECS Instances", "View ECS instances", '1', func() { switchToEcsListView() }).
		AddItem("DNS Management", "View AliDNS domains and records", '2', func() { switchToDnsDomainsListView() }).
		AddItem("SLB Instances", "View SLB instances", '3', func() { switchToSlbListView() }).
		AddItem("OSS Management", "Browse OSS buckets (STUBBED)", '4', func() { switchToOssBucketListView() }).
		AddItem("RDS Instances", "View RDS instances", '5', func() { // Added RDS
			switchToRdsListView()
		}).
		AddItem("Quit", "Exit the application", 'q', func() { app.Stop() })
	list.SetBorder(true).SetTitle("Main Menu")
	return list
}

// Generic helper for adding items to a detail form
func addFormItem(form *tview.Form, label, value string) {
	form.AddTextView(label, value, 0, 1, true, false)
}


// ECS List View (no changes)
func createEcsListView(instances []ecs.Instance) *tview.Table {
	table := tview.NewTable().
		SetBorders(true).
		SetSelectable(true, false)
	headers := []string{"Instance ID", "Status", "IP Address", "Name"}
	for c, header := range headers {
		table.SetCell(0, c, tview.NewTableCell(header).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter).SetSelectable(false))
	}
	if len(instances) == 0 {
		table.SetCell(1, 0, tview.NewTableCell("No ECS instances found.").SetSelectable(false).SetExpansion(len(headers)).SetAlign(tview.AlignCenter))
	} else {
		for r, instance := range instances {
			ipAddress := "N/A"
			if len(instance.VpcAttributes.PrivateIpAddress.IpAddress) > 0 {
				ipAddress = instance.VpcAttributes.PrivateIpAddress.IpAddress[0]
			} else if len(instance.InnerIpAddress.IpAddress) > 0 {
				ipAddress = instance.InnerIpAddress.IpAddress[0]
			} else if len(instance.PublicIpAddress.IpAddress) > 0 {
				ipAddress = instance.PublicIpAddress.IpAddress[0]
			}
			table.SetCell(r+1, 0, tview.NewTableCell(instance.InstanceId).SetTextColor(tcell.ColorWhite).SetReference(instance.InstanceId))
			table.SetCell(r+1, 1, tview.NewTableCell(instance.Status).SetTextColor(tcell.ColorWhite))
			table.SetCell(r+1, 2, tview.NewTableCell(ipAddress).SetTextColor(tcell.ColorWhite))
			table.SetCell(r+1, 3, tview.NewTableCell(instance.InstanceName).SetTextColor(tcell.ColorWhite))
		}
	}
	return table
}

// ECS Detail View (no changes)
func createEcsDetailView(instance ecs.Instance) *tview.Flex {
	ecsDetailForm = tview.NewForm() 
	ecsDetailForm.SetBorder(true).SetTitle(fmt.Sprintf("ECS Details: %s", instance.InstanceId)).SetTitleAlign(tview.AlignLeft)
	ecsDetailForm.Clear(true)
	addFormItem(ecsDetailForm, "Instance ID:", instance.InstanceId)
	addFormItem(ecsDetailForm, "Instance Name:", instance.InstanceName)
	addFormItem(ecsDetailForm, "Status:", instance.Status)
	addFormItem(ecsDetailForm, "Region ID:", instance.RegionId)
	addFormItem(ecsDetailForm, "Zone ID:", instance.ZoneId)
	// ... (rest of ECS details - simplified for brevity in this combined file)
	ecsDetailForm.AddButton("Back (Esc or q)", func() {
		pages.SwitchToPage(pageEcsList)
		app.SetFocus(ecsInstanceTable)
	})
	return tview.NewFlex().AddItem(ecsDetailForm, 0, 1, true)
}

// DNS Domains List View (no changes)
func createDnsDomainsListView(domains []alidns.Domain) *tview.Table {
	table := tview.NewTable().SetBorders(true).SetSelectable(true, false)
	headers := []string{"Domain Name", "Record Count", "Version Code"}
	for c, header := range headers {
		table.SetCell(0, c, tview.NewTableCell(header).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter).SetSelectable(false))
	}
	if len(domains) == 0 { /* ... */ } else {
		for r, domain := range domains {
			table.SetCell(r+1, 0, tview.NewTableCell(domain.DomainName).SetTextColor(tcell.ColorWhite).SetReference(domain.DomainName))
			table.SetCell(r+1, 1, tview.NewTableCell(fmt.Sprintf("%d", domain.RecordCount)).SetTextColor(tcell.ColorWhite))
			table.SetCell(r+1, 2, tview.NewTableCell(domain.VersionCode).SetTextColor(tcell.ColorWhite))
		}
	}
	return table
}

// DNS Records List View (no changes)
func createDnsRecordsListView(records []alidns.Record, domainName string) *tview.Table {
	table := tview.NewTable().SetBorders(true).SetSelectable(true, false) 
	headers := []string{"Record ID", "RR", "Type", "Value", "TTL", "Status"}
	for c, header := range headers { /* ... */ }
	if len(records) == 0 { /* ... */ } else {
		for r, record := range records { /* ... */ }
	}
	table.SetTitle(fmt.Sprintf("DNS Records for %s", domainName)).SetBorder(true)
	return table
}

// SLB List View (no changes)
func createSlbListView(slbs []slb.LoadBalancer) *tview.Table {
	table := tview.NewTable().SetBorders(true).SetSelectable(true, false)
	headers := []string{"SLB ID", "Name", "IP Address", "Type", "Status"}
	for c, header := range headers { /* ... */ }
	if len(slbs) == 0 { /* ... */ } else {
		for r, lb := range slbs {
			table.SetCell(r+1, 0, tview.NewTableCell(lb.LoadBalancerId).SetTextColor(tcell.ColorWhite).SetReference(lb.LoadBalancerId))
			// ... other cells
		}
	}
	return table
}

// SLB Detail View (no changes)
func createSlbDetailView(lb slb.LoadBalancer) *tview.Flex {
	slbDetailForm = tview.NewForm() 
	slbDetailForm.SetBorder(true).SetTitle(fmt.Sprintf("SLB Details: %s", lb.LoadBalancerId)).SetTitleAlign(tview.AlignLeft)
	slbDetailForm.Clear(true)
	addFormItem(slbDetailForm, "SLB ID:", lb.LoadBalancerId)
	// ... (rest of SLB details)
	slbDetailForm.AddButton("Back (Esc or q)", func() {
		pages.SwitchToPage(pageSlbList)
		app.SetFocus(slbInstanceTable)
	})
	return tview.NewFlex().AddItem(slbDetailForm, 0, 1, true)
}

// OSS Bucket List View (STUBBED - no changes)
func createOssBucketListView(buckets []StubBucketInfo) *tview.Table {
	table := tview.NewTable().SetBorders(true).SetSelectable(true, false)
	headers := []string{"Bucket Name", "Location", "Creation Date", "Storage Class"}
	for c, header := range headers { /* ... */ }
	if len(buckets) == 0 { /* ... */ } else {
		for r, bucket := range buckets {
			table.SetCell(r+1, 0, tview.NewTableCell(bucket.Name).SetTextColor(tcell.ColorWhite).SetReference(bucket.Name))
			// ... other cells
		}
	}
	table.SetTitle("OSS Buckets (SDK STUBBED)").SetBorder(true)
	return table
}

// OSS Object List View (STUBBED - no changes)
func createOssObjectListView(objects []StubObjectInfo, bucketName string) *tview.Table {
	table := tview.NewTable().SetBorders(true).SetSelectable(true, false)
	headers := []string{"Object Key", "Size (Bytes)", "Last Modified", "Storage Class"}
	for c, header := range headers { /* ... */ }
	if len(objects) == 0 { /* ... */ } else {
		for r, object := range objects { /* ... */ }
	}
	table.SetTitle(fmt.Sprintf("Objects in %s (SDK STUBBED)", bucketName)).SetBorder(true)
	return table
}

// RDS List View
func createRdsListView(instances []rds.DBInstance) *tview.Table {
	table := tview.NewTable().
		SetBorders(true).
		SetSelectable(true, false)

	headers := []string{"Instance ID", "Engine", "Version", "Class", "Status", "Description"}
	for c, header := range headers {
		table.SetCell(0, c, tview.NewTableCell(header).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter).SetSelectable(false))
	}

	if len(instances) == 0 {
		table.SetCell(1, 0, tview.NewTableCell("No RDS instances found.").SetSelectable(false).SetExpansion(len(headers)).SetAlign(tview.AlignCenter))
	} else {
		for r, inst := range instances {
			table.SetCell(r+1, 0, tview.NewTableCell(inst.DBInstanceId).SetTextColor(tcell.ColorWhite).SetReference(inst.DBInstanceId))
			table.SetCell(r+1, 1, tview.NewTableCell(inst.Engine).SetTextColor(tcell.ColorWhite))
			table.SetCell(r+1, 2, tview.NewTableCell(inst.EngineVersion).SetTextColor(tcell.ColorWhite))
			table.SetCell(r+1, 3, tview.NewTableCell(inst.DBInstanceClass).SetTextColor(tcell.ColorWhite))
			table.SetCell(r+1, 4, tview.NewTableCell(inst.DBInstanceStatus).SetTextColor(tcell.ColorWhite))
			table.SetCell(r+1, 5, tview.NewTableCell(inst.DBInstanceDescription).SetTextColor(tcell.ColorWhite).SetMaxWidth(40))
		}
	}
	return table
}

// RDS Detail View
func createRdsDetailView(instance rds.DBInstance) *tview.Flex {
	rdsDetailForm = tview.NewForm()
	rdsDetailForm.SetBorder(true).SetTitle(fmt.Sprintf("RDS Details: %s", instance.DBInstanceId)).SetTitleAlign(tview.AlignLeft)
	rdsDetailForm.Clear(true)

	addFormItem(rdsDetailForm, "Instance ID:", instance.DBInstanceId)
	addFormItem(rdsDetailForm, "Description:", instance.DBInstanceDescription)
	addFormItem(rdsDetailForm, "Status:", instance.DBInstanceStatus)
	addFormItem(rdsDetailForm, "Engine:", instance.Engine)
	addFormItem(rdsDetailForm, "Engine Version:", instance.EngineVersion)
	addFormItem(rdsDetailForm, "Instance Class:", instance.DBInstanceClass)
	addFormItem(rdsDetailForm, "Storage Type:", instance.DBInstanceStorageType)
	addFormItem(rdsDetailForm, "Allocated Storage (GB):", fmt.Sprintf("%d", instance.DBInstanceStorage))
	addFormItem(rdsDetailForm, "Connection String:", instance.ConnectionString)
	addFormItem(rdsDetailForm, "Port:", instance.Port)
	addFormItem(rdsDetailForm, "Network Type:", instance.InstanceNetworkType) // VPC or Classic
	addFormItem(rdsDetailForm, "VPC ID:", instance.VpcId)
	// addFormItem(rdsDetailForm, "VSwitch ID:", instance.VSwitchId) // Not directly in DBInstance, might need DescribeDBInstanceAttribute
	addFormItem(rdsDetailForm, "Region ID:", instance.RegionId)
	addFormItem(rdsDetailForm, "Zone ID:", instance.ZoneId)
	addFormItem(rdsDetailForm, "Creation Time:", instance.CreateTime)
	addFormItem(rdsDetailForm, "Expire Time:", instance.ExpireTime)
	addFormItem(rdsDetailForm, "Lock Mode:", instance.LockMode)
	addFormItem(rdsDetailForm, "Pay Type:", instance.PayType)


	rdsDetailForm.AddButton("Back (Esc or q)", func() {
		pages.SwitchToPage(pageRdsList)
		app.SetFocus(rdsInstanceTable)
	})
	return tview.NewFlex().AddItem(rdsDetailForm, 0, 1, true)
}


// --- Page Switching Functions ---
func switchToEcsListView() {
	if allECSInstances == nil { /* fetch */ instances, err := fetchECSInstances(); if err != nil { showErrorModal(err.Error()); return }; allECSInstances = instances }
	ecsInstanceTable = createEcsListView(allECSInstances)
	setupTableNavigation(ecsInstanceTable, func(row, col int) { 
		instanceId := ecsInstanceTable.GetCell(row, 0).GetReference().(string)
		var selectedInstance ecs.Instance; for _, inst := range allECSInstances { if inst.InstanceId == instanceId { selectedInstance = inst; break } }
		pages.AddPage(pageEcsDetail, createEcsDetailView(selectedInstance), true, true) 
		app.SetFocus(ecsDetailForm.GetButton(ecsDetailForm.GetButtonCount() - 1))
	})
	pages.AddPage(pageEcsList, ecsInstanceTable, true, true); app.SetFocus(ecsInstanceTable)
}

func switchToDnsDomainsListView() {
	if allDomains == nil { /* fetch */ domains, err := fetchAlidnsDomains(); if err != nil { showErrorModal(err.Error()); return }; allDomains = domains }
	dnsDomainsTable = createDnsDomainsListView(allDomains)
	setupTableNavigation(dnsDomainsTable, func(row, col int) { 
		domainName := dnsDomainsTable.GetCell(row, 0).GetReference().(string)
		switchToDnsRecordsListView(domainName)
	})
	pages.AddPage(pageDnsDomains, dnsDomainsTable, true, true); app.SetFocus(dnsDomainsTable)
}

func switchToDnsRecordsListView(domainName string) {
	records, err := fetchAlidnsDomainRecords(domainName); if err != nil { showErrorModal(err.Error()); return }
	dnsRecordsTable = createDnsRecordsListView(records, domainName)
	setupTableNavigation(dnsRecordsTable, nil)
	pages.AddPage(pageDnsRecords, dnsRecordsTable, true, true); app.SetFocus(dnsRecordsTable)
}

func switchToSlbListView() {
	if allSLBInstances == nil { /* fetch */ slbs, err := fetchSLBInstances(); if err != nil { showErrorModal(err.Error()); return }; allSLBInstances = slbs }
	slbInstanceTable = createSlbListView(allSLBInstances)
	setupTableNavigation(slbInstanceTable, func(row, col int) { 
		slbId := slbInstanceTable.GetCell(row, 0).GetReference().(string)
		var selectedSlb slb.LoadBalancer; for _, lb := range allSLBInstances { if lb.LoadBalancerId == slbId { selectedSlb = lb; break } }
		pages.AddPage(pageSlbDetail, createSlbDetailView(selectedSlb), true, true) 
		app.SetFocus(slbDetailForm.GetButton(slbDetailForm.GetButtonCount() - 1))
	})
	pages.AddPage(pageSlbList, slbInstanceTable, true, true); app.SetFocus(slbInstanceTable)
}

func switchToOssBucketListView() {
	buckets, err := fetchOssBuckets(); if err != nil { showErrorModal(err.Error()) } // STUBBED
	ossBucketTable = createOssBucketListView(buckets)
	setupTableNavigation(ossBucketTable, func(row, col int) { 
		bucketName := ossBucketTable.GetCell(row, 0).GetReference().(string)
		switchToOssObjectListView(bucketName)
	})
	pages.AddPage(pageOssBuckets, ossBucketTable, true, true); app.SetFocus(ossBucketTable)
}

func switchToOssObjectListView(bucketName string) {
	objects, err := fetchOssObjects(bucketName); if err != nil { showErrorModal(err.Error()) } // STUBBED
	ossObjectTable = createOssObjectListView(objects, bucketName)
	setupTableNavigation(ossObjectTable, nil)
	pages.AddPage(pageOssObjects, ossObjectTable, true, true); app.SetFocus(ossObjectTable)
}

func switchToRdsListView() {
	if allRDSInstances == nil {
		instances, err := fetchRDSInstances()
		if err != nil {
			showErrorModal(fmt.Sprintf("Failed to fetch RDS instances: %v", err))
			return
		}
		allRDSInstances = instances
	}
	rdsInstanceTable = createRdsListView(allRDSInstances)
	setupTableNavigation(rdsInstanceTable, func(row, col int) { // RDS item selected
		cell := rdsInstanceTable.GetCell(row, 0)
		instanceId, ok := cell.GetReference().(string)
		if !ok { return }
		var selectedInstance rds.DBInstance
		for _, inst := range allRDSInstances {
			if inst.DBInstanceId == instanceId {
				selectedInstance = inst
				break
			}
		}
		detailViewContent := createRdsDetailView(selectedInstance)
		pages.AddPage(pageRdsDetail, detailViewContent, true, true) // Add and show
		app.SetFocus(rdsDetailForm.GetButton(rdsDetailForm.GetButtonCount() - 1))
	})
	pages.AddPage(pageRdsList, rdsInstanceTable, true, true) // Add and show
	app.SetFocus(rdsInstanceTable)
}


// Helper for table navigation (j/k)
func setupTableNavigation(table *tview.Table, onSelect func(row, column int)) {
	table.SetSelectedFunc(func(row, column int) { if row > 0 && onSelect != nil { onSelect(row, column) } })
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		currentRow, _ := table.GetSelection(); rowCount := table.GetRowCount()
		switch event.Rune() {
		case 'j': if currentRow < rowCount-1 { table.Select(currentRow+1, 0) }; return nil
		case 'k': if currentRow > 1 { table.Select(currentRow-1, 0) } else if rowCount > 1 { table.Select(1,0)}; return nil
		}
		return event
	})
}

// Error Modal
func showErrorModal(message string) {
	modal := tview.NewModal().SetText(message).AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			pages.RemovePage("errorModal")
			// Try to focus the current page's main element or fallback to main menu
			currentPageName, prim := pages.GetFrontPage()
			if prim != nil && currentPageName != "errorModal" { app.SetFocus(prim) } else { app.SetFocus(mainMenu) }
		})
	pages.AddPage("errorModal", modal, false, true); app.SetFocus(modal)
}


func main() {
	accessKeyID = os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID")
	accessKeySecret = os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET")
	regionID = os.Getenv("ALIBABA_CLOUD_REGION_ID")
	ossEndpoint = os.Getenv("ALIBABA_CLOUD_OSS_ENDPOINT") 

	if accessKeyID == "" || accessKeySecret == "" || regionID == "" {
		fmt.Println("Error: ALIBABA_CLOUD_ACCESS_KEY_ID, ALIBABA_CLOUD_ACCESS_KEY_SECRET, and ALIBABA_CLOUD_REGION_ID environment variables must be set.")
		os.Exit(1)
	}

	app = tview.NewApplication()
	pages = tview.NewPages()

	mainMenu = createMainMenu()
	pages.AddPage(pageMainMenu, mainMenu, true, true)

	// Global input capture
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		currentFocus := app.GetFocus()
		if modal, isModal := currentFocus.(*tview.Modal); isModal && modal.HasFocus() { return event }
		if form, isForm := currentFocus.(*tview.Form); isForm && form.HasFocus() && 
		   (event.Key() == tcell.KeyTab || event.Key() == tcell.KeyBacktab || event.Key() == tcell.KeyEnter) {
			return event 
		}

		if event.Key() == tcell.KeyCtrlC { app.Stop(); return nil }
		
		currentPageName, _ := pages.GetFrontPage()

		handleNavigation := func(targetPage string, focusItem tview.Primitive) {
			pages.SwitchToPage(targetPage)
			if focusItem != nil { app.SetFocus(focusItem) } else if targetPage == pageMainMenu { app.SetFocus(mainMenu) }
		}
		
		switch event.Key() {
		case tcell.KeyEscape:
			switch currentPageName {
			case pageEcsList, pageDnsDomains, pageSlbList, pageOssBuckets, pageRdsList: // Added RdsList
				handleNavigation(pageMainMenu, mainMenu)
			case pageEcsDetail: handleNavigation(pageEcsList, ecsInstanceTable)
			case pageDnsRecords: handleNavigation(pageDnsDomains, dnsDomainsTable)
			case pageSlbDetail: handleNavigation(pageSlbList, slbInstanceTable)
			case pageOssObjects: handleNavigation(pageOssBuckets, ossBucketTable)
			case pageRdsDetail: handleNavigation(pageRdsList, rdsInstanceTable) // Added RdsDetail
			}
			return nil 
		case tcell.KeyRune:
			if event.Rune() == 'q' {
				switch currentPageName {
				case pageMainMenu, pageEcsList, pageDnsDomains, pageSlbList, pageOssBuckets, pageRdsList: // Added RdsList
					app.Stop() 
				case pageEcsDetail: handleNavigation(pageEcsList, ecsInstanceTable)
				case pageDnsRecords: handleNavigation(pageDnsDomains, dnsDomainsTable)
				case pageSlbDetail: handleNavigation(pageSlbList, slbInstanceTable)
				case pageOssObjects: handleNavigation(pageOssBuckets, ossBucketTable)
				case pageRdsDetail: handleNavigation(pageRdsList, rdsInstanceTable) // Added RdsDetail
				}
				return nil 
			}
		}
		return event 
	})

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		fmt.Printf("Error running TUI application: %v\n", err)
		os.Exit(1)
	}
}
