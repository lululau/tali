package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Global variables
var (
	app              *tview.Application
	pages            *tview.Pages
	mainMenu         *tview.List
	ecsInstanceTable *tview.Table
	ecsDetailView    *tview.TextView
	dnsDomainsTable  *tview.Table
	dnsRecordsTable  *tview.Table
	slbInstanceTable *tview.Table
	slbDetailView    *tview.TextView
	ossBucketTable   *tview.Table
	ossObjectTable   *tview.Table
	ossDetailView    *tview.TextView
	rdsInstanceTable *tview.Table
	rdsDetailView    *tview.TextView

	allECSInstances   []ecs.Instance
	allDomains        []alidns.DomainInDescribeDomains
	allSLBInstances   []slb.LoadBalancer
	allRDSInstances   []rds.DBInstance
	allOssBuckets     []oss.BucketProperties
	currentBucketName string

	accessKeyID     string
	accessKeySecret string
	regionID        string
	ossEndpoint     string
)

const (
	pageMainMenu   = "mainMenu"
	pageEcsList    = "ecsList"
	pageEcsDetail  = "ecsDetail"
	pageDnsDomains = "dnsDomains"
	pageDnsRecords = "dnsRecords"
	pageSlbList    = "slbList"
	pageSlbDetail  = "slbDetail"
	pageOssBuckets = "ossBuckets"
	pageOssObjects = "ossObjects"
	pageRdsList    = "rdsList"   // New page for RDS list
	pageRdsDetail  = "rdsDetail" // New page for RDS detail
)

// --- OSS SDK Functions ---
func fetchOssBuckets() ([]oss.BucketProperties, error) {
	client, err := oss.New(ossEndpoint, accessKeyID, accessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("creating OSS client: %w", err)
	}

	result, err := client.ListBuckets()
	if err != nil {
		return nil, fmt.Errorf("listing OSS buckets: %w", err)
	}

	return result.Buckets, nil
}

func fetchOssObjects(bucketName string) ([]oss.ObjectProperties, error) {
	client, err := oss.New(ossEndpoint, accessKeyID, accessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("creating OSS client: %w", err)
	}

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return nil, fmt.Errorf("getting bucket %s: %w", bucketName, err)
	}

	result, err := bucket.ListObjects()
	if err != nil {
		return nil, fmt.Errorf("listing objects in bucket %s: %w", bucketName, err)
	}

	return result.Objects, nil
}

// --- ECS SDK Functions ---
func fetchECSInstances() ([]ecs.Instance, error) {
	client, err := ecs.NewClientWithAccessKey(regionID, accessKeyID, accessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("creating ECS client: %w", err)
	}
	request := ecs.CreateDescribeInstancesRequest()
	request.Scheme = "https"
	response, err := client.DescribeInstances(request)
	if err != nil {
		return nil, fmt.Errorf("describing ECS instances: %w", err)
	}
	return response.Instances.Instance, nil
}

// --- AliDNS SDK Functions ---
func fetchAlidnsDomains() ([]alidns.DomainInDescribeDomains, error) {
	client, err := alidns.NewClientWithAccessKey(regionID, accessKeyID, accessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("creating AliDNS client: %w", err)
	}
	request := alidns.CreateDescribeDomainsRequest()
	request.Scheme = "https"
	response, err := client.DescribeDomains(request)
	if err != nil {
		return nil, fmt.Errorf("describing AliDNS domains: %w", err)
	}
	return response.Domains.Domain, nil
}

func fetchAlidnsDomainRecords(domainName string) ([]alidns.Record, error) {
	client, err := alidns.NewClientWithAccessKey(regionID, accessKeyID, accessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("creating AliDNS client: %w", err)
	}
	request := alidns.CreateDescribeDomainRecordsRequest()
	request.Scheme = "https"
	request.DomainName = domainName
	response, err := client.DescribeDomainRecords(request)
	if err != nil {
		return nil, fmt.Errorf("describing AliDNS domain records for %s: %w", domainName, err)
	}
	return response.DomainRecords.Record, nil
}

// --- SLB SDK Functions ---
func fetchSLBInstances() ([]slb.LoadBalancer, error) {
	client, err := slb.NewClientWithAccessKey(regionID, accessKeyID, accessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("creating SLB client: %w", err)
	}
	request := slb.CreateDescribeLoadBalancersRequest()
	request.Scheme = "https"
	response, err := client.DescribeLoadBalancers(request)
	if err != nil {
		return nil, fmt.Errorf("describing SLB instances: %w", err)
	}
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
		AddItem("OSS Management", "Browse OSS buckets and objects", '4', func() { switchToOssBucketListView() }).
		AddItem("RDS Instances", "View RDS instances", '5', func() { // Added RDS
			switchToRdsListView()
		}).
		AddItem("Quit", "Exit the application (Press 'Q')", 'Q', func() { app.Stop() })
	list.SetBorder(true).SetTitle("Main Menu")
	return list
}

// Generic helper for creating JSON detail view
func createJSONDetailView(title string, data interface{}) *tview.TextView {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		jsonData = []byte(fmt.Sprintf("Error marshaling JSON: %v", err))
	}

	textView := tview.NewTextView().
		SetText(string(jsonData)).
		SetScrollable(true).
		SetWrap(false)
	textView.SetBorder(true).SetTitle(title)

	return textView
}

// Generic helper for setting up table with full width
func setupTableWithFixedWidth(table *tview.Table) *tview.Table {
	table.SetFixed(1, 0) // Fix header row, allow all columns to be flexible
	table.SetSelectable(true, false)
	table.SetBorder(true)
	return table
}

// Generic helper for creating table headers with expansion
func createTableHeaders(table *tview.Table, headers []string) {
	for c, header := range headers {
		cell := tview.NewTableCell(header).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter).SetSelectable(false)
		// Set all columns to expand proportionally
		cell.SetExpansion(1)
		table.SetCell(0, c, cell)
	}
}

// Generic helper for wrapping table in full-width flex container
func wrapTableInFlex(table *tview.Table) tview.Primitive {
	// Create a flex container that forces the table to use full width
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.AddItem(table, 0, 1, true)
	flex.SetBorder(false)
	return flex
}

// ECS List View with enhanced information
func createEcsListView(instances []ecs.Instance) *tview.Table {
	table := tview.NewTable().
		SetBorders(true).
		SetSelectable(true, false)
	table = setupTableWithFixedWidth(table)
	headers := []string{"Instance ID", "Status", "Zone", "CPU/RAM", "Private IP", "Public IP", "Name"}
	createTableHeaders(table, headers)
	if len(instances) == 0 {
		table.SetCell(1, 0, tview.NewTableCell("No ECS instances found.").SetSelectable(false).SetExpansion(len(headers)).SetAlign(tview.AlignCenter))
	} else {
		for r, instance := range instances {
			// Private IP
			privateIP := "N/A"
			if len(instance.VpcAttributes.PrivateIpAddress.IpAddress) > 0 {
				privateIP = instance.VpcAttributes.PrivateIpAddress.IpAddress[0]
			} else if len(instance.InnerIpAddress.IpAddress) > 0 {
				privateIP = instance.InnerIpAddress.IpAddress[0]
			}

			// Public IP
			publicIP := "N/A"
			if len(instance.PublicIpAddress.IpAddress) > 0 {
				publicIP = instance.PublicIpAddress.IpAddress[0]
			}

			// CPU/RAM configuration
			cpuRam := fmt.Sprintf("%dC/%dG", instance.Cpu, instance.Memory/1024)

			table.SetCell(r+1, 0, tview.NewTableCell(instance.InstanceId).SetTextColor(tcell.ColorWhite).SetReference(instance.InstanceId).SetExpansion(1))
			table.SetCell(r+1, 1, tview.NewTableCell(instance.Status).SetTextColor(tcell.ColorWhite).SetExpansion(1))
			table.SetCell(r+1, 2, tview.NewTableCell(instance.ZoneId).SetTextColor(tcell.ColorWhite).SetExpansion(1))
			table.SetCell(r+1, 3, tview.NewTableCell(cpuRam).SetTextColor(tcell.ColorWhite).SetExpansion(1))
			table.SetCell(r+1, 4, tview.NewTableCell(privateIP).SetTextColor(tcell.ColorWhite).SetExpansion(1))
			table.SetCell(r+1, 5, tview.NewTableCell(publicIP).SetTextColor(tcell.ColorWhite).SetExpansion(1))
			table.SetCell(r+1, 6, tview.NewTableCell(instance.InstanceName).SetTextColor(tcell.ColorWhite).SetExpansion(1))
		}
	}
	return table
}

// ECS Detail View in JSON format
func createEcsDetailView(instance ecs.Instance) *tview.Flex {
	ecsDetailView = createJSONDetailView(fmt.Sprintf("ECS Details: %s", instance.InstanceId), instance)

	// Add navigation instructions
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	// Instructions
	instructions := tview.NewTextView().
		SetText("Press 'Esc' or 'q' to go back, 'Q' to quit").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)
	instructions.SetBorder(false)

	flex.AddItem(instructions, 1, 0, false)
	flex.AddItem(ecsDetailView, 0, 1, true)

	return flex
}

// DNS Domains List View
func createDnsDomainsListView(domains []alidns.DomainInDescribeDomains) *tview.Table {
	table := tview.NewTable().SetBorders(true).SetSelectable(true, false)
	table = setupTableWithFixedWidth(table)
	headers := []string{"Domain Name", "Record Count", "Version Code"}
	createTableHeaders(table, headers)
	if len(domains) == 0 {
		table.SetCell(1, 0, tview.NewTableCell("No domains found.").SetSelectable(false).SetExpansion(len(headers)).SetAlign(tview.AlignCenter))
	} else {
		for r, domain := range domains {
			table.SetCell(r+1, 0, tview.NewTableCell(domain.DomainName).SetTextColor(tcell.ColorWhite).SetReference(domain.DomainName).SetExpansion(1))
			table.SetCell(r+1, 1, tview.NewTableCell(fmt.Sprintf("%d", domain.RecordCount)).SetTextColor(tcell.ColorWhite).SetExpansion(1))
			table.SetCell(r+1, 2, tview.NewTableCell(domain.VersionCode).SetTextColor(tcell.ColorWhite).SetExpansion(1))
		}
	}
	return table
}

// DNS Records List View
func createDnsRecordsListView(records []alidns.Record, domainName string) *tview.Table {
	table := tview.NewTable().SetBorders(true).SetSelectable(true, false)
	table = setupTableWithFixedWidth(table)
	headers := []string{"Record ID", "RR", "Type", "Value", "TTL", "Status"}
	createTableHeaders(table, headers)
	if len(records) == 0 {
		table.SetCell(1, 0, tview.NewTableCell("No DNS records found.").SetSelectable(false).SetExpansion(len(headers)).SetAlign(tview.AlignCenter))
	} else {
		for r, record := range records {
			table.SetCell(r+1, 0, tview.NewTableCell(record.RecordId).SetTextColor(tcell.ColorWhite).SetExpansion(1))
			table.SetCell(r+1, 1, tview.NewTableCell(record.RR).SetTextColor(tcell.ColorWhite).SetExpansion(1))
			table.SetCell(r+1, 2, tview.NewTableCell(record.Type).SetTextColor(tcell.ColorWhite).SetExpansion(1))
			table.SetCell(r+1, 3, tview.NewTableCell(record.Value).SetTextColor(tcell.ColorWhite).SetExpansion(1))
			table.SetCell(r+1, 4, tview.NewTableCell(fmt.Sprintf("%d", record.TTL)).SetTextColor(tcell.ColorWhite).SetExpansion(1))
			table.SetCell(r+1, 5, tview.NewTableCell(record.Status).SetTextColor(tcell.ColorWhite).SetExpansion(1))
		}
	}
	table.SetTitle(fmt.Sprintf("DNS Records for %s", domainName)).SetBorder(true)
	return table
}

// SLB List View
func createSlbListView(slbs []slb.LoadBalancer) *tview.Table {
	table := tview.NewTable().SetBorders(true).SetSelectable(true, false)
	table = setupTableWithFixedWidth(table)
	headers := []string{"SLB ID", "Name", "IP Address", "Type", "Status"}
	createTableHeaders(table, headers)
	if len(slbs) == 0 {
		table.SetCell(1, 0, tview.NewTableCell("No SLB instances found.").SetSelectable(false).SetExpansion(len(headers)).SetAlign(tview.AlignCenter))
	} else {
		for r, lb := range slbs {
			table.SetCell(r+1, 0, tview.NewTableCell(lb.LoadBalancerId).SetTextColor(tcell.ColorWhite).SetReference(lb.LoadBalancerId).SetExpansion(1))
			table.SetCell(r+1, 1, tview.NewTableCell(lb.LoadBalancerName).SetTextColor(tcell.ColorWhite).SetExpansion(1))
			table.SetCell(r+1, 2, tview.NewTableCell(lb.Address).SetTextColor(tcell.ColorWhite).SetExpansion(1))
			table.SetCell(r+1, 3, tview.NewTableCell(lb.LoadBalancerSpec).SetTextColor(tcell.ColorWhite).SetExpansion(1))
			table.SetCell(r+1, 4, tview.NewTableCell(lb.LoadBalancerStatus).SetTextColor(tcell.ColorWhite).SetExpansion(1))
		}
	}
	return table
}

// SLB Detail View in JSON format
func createSlbDetailView(lb slb.LoadBalancer) *tview.Flex {
	slbDetailView = createJSONDetailView(fmt.Sprintf("SLB Details: %s", lb.LoadBalancerId), lb)

	// Add navigation instructions
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	// Instructions
	instructions := tview.NewTextView().
		SetText("Press 'Esc' or 'q' to go back, 'Q' to quit").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)
	instructions.SetBorder(false)

	flex.AddItem(instructions, 1, 0, false)
	flex.AddItem(slbDetailView, 0, 1, true)

	return flex
}

// OSS Bucket List View
func createOssBucketListView(buckets []oss.BucketProperties) *tview.Table {
	table := tview.NewTable().SetBorders(true).SetSelectable(true, false)
	table = setupTableWithFixedWidth(table)
	headers := []string{"Bucket Name", "Location", "Creation Date", "Storage Class"}
	createTableHeaders(table, headers)
	if len(buckets) == 0 {
		table.SetCell(1, 0, tview.NewTableCell("No OSS buckets found.").SetSelectable(false).SetExpansion(len(headers)).SetAlign(tview.AlignCenter))
	} else {
		for r, bucket := range buckets {
			table.SetCell(r+1, 0, tview.NewTableCell(bucket.Name).SetTextColor(tcell.ColorWhite).SetReference(bucket.Name).SetExpansion(1))
			table.SetCell(r+1, 1, tview.NewTableCell(bucket.Location).SetTextColor(tcell.ColorWhite).SetExpansion(1))
			table.SetCell(r+1, 2, tview.NewTableCell(bucket.CreationDate.Format("2006-01-02 15:04:05")).SetTextColor(tcell.ColorWhite).SetExpansion(1))
			table.SetCell(r+1, 3, tview.NewTableCell(bucket.StorageClass).SetTextColor(tcell.ColorWhite).SetExpansion(1))
		}
	}
	table.SetTitle("OSS Buckets").SetBorder(true)
	return table
}

// OSS Object List View
func createOssObjectListView(objects []oss.ObjectProperties, bucketName string) *tview.Table {
	table := tview.NewTable().SetBorders(true).SetSelectable(true, false)
	table = setupTableWithFixedWidth(table)
	headers := []string{"Object Key", "Size (Bytes)", "Last Modified", "Storage Class", "ETag"}
	createTableHeaders(table, headers)
	if len(objects) == 0 {
		table.SetCell(1, 0, tview.NewTableCell("No objects found.").SetSelectable(false).SetExpansion(len(headers)).SetAlign(tview.AlignCenter))
	} else {
		for r, object := range objects {
			table.SetCell(r+1, 0, tview.NewTableCell(object.Key).SetTextColor(tcell.ColorWhite).SetReference(object.Key).SetExpansion(1))
			table.SetCell(r+1, 1, tview.NewTableCell(fmt.Sprintf("%d", object.Size)).SetTextColor(tcell.ColorWhite).SetExpansion(1))
			table.SetCell(r+1, 2, tview.NewTableCell(object.LastModified.Format("2006-01-02 15:04:05")).SetTextColor(tcell.ColorWhite).SetExpansion(1))
			table.SetCell(r+1, 3, tview.NewTableCell(object.StorageClass).SetTextColor(tcell.ColorWhite).SetExpansion(1))
			table.SetCell(r+1, 4, tview.NewTableCell(object.ETag).SetTextColor(tcell.ColorWhite).SetExpansion(1))
		}
	}
	table.SetTitle(fmt.Sprintf("Objects in %s", bucketName)).SetBorder(true)
	return table
}

// RDS List View
func createRdsListView(instances []rds.DBInstance) *tview.Table {
	table := tview.NewTable().
		SetBorders(true).
		SetSelectable(true, false)
	table = setupTableWithFixedWidth(table)

	headers := []string{"Instance ID", "Engine", "Version", "Class", "Status", "Description"}
	createTableHeaders(table, headers)

	if len(instances) == 0 {
		table.SetCell(1, 0, tview.NewTableCell("No RDS instances found.").SetSelectable(false).SetExpansion(len(headers)).SetAlign(tview.AlignCenter))
	} else {
		for r, inst := range instances {
			table.SetCell(r+1, 0, tview.NewTableCell(inst.DBInstanceId).SetTextColor(tcell.ColorWhite).SetReference(inst.DBInstanceId).SetExpansion(1))
			table.SetCell(r+1, 1, tview.NewTableCell(inst.Engine).SetTextColor(tcell.ColorWhite).SetExpansion(1))
			table.SetCell(r+1, 2, tview.NewTableCell(inst.EngineVersion).SetTextColor(tcell.ColorWhite).SetExpansion(1))
			table.SetCell(r+1, 3, tview.NewTableCell(inst.DBInstanceClass).SetTextColor(tcell.ColorWhite).SetExpansion(1))
			table.SetCell(r+1, 4, tview.NewTableCell(inst.DBInstanceStatus).SetTextColor(tcell.ColorWhite).SetExpansion(1))
			table.SetCell(r+1, 5, tview.NewTableCell(inst.DBInstanceDescription).SetTextColor(tcell.ColorWhite).SetMaxWidth(40).SetExpansion(1))
		}
	}
	return table
}

// RDS Detail View in JSON format
func createRdsDetailView(instance rds.DBInstance) *tview.Flex {
	rdsDetailView = createJSONDetailView(fmt.Sprintf("RDS Details: %s", instance.DBInstanceId), instance)

	// Add navigation instructions
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	// Instructions
	instructions := tview.NewTextView().
		SetText("Press 'Esc' or 'q' to go back, 'Q' to quit").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)
	instructions.SetBorder(false)

	flex.AddItem(instructions, 1, 0, false)
	flex.AddItem(rdsDetailView, 0, 1, true)

	return flex
}

// --- Page Switching Functions ---
func switchToEcsListView() {
	if allECSInstances == nil { /* fetch */
		instances, err := fetchECSInstances()
		if err != nil {
			showErrorModal(err.Error())
			return
		}
		allECSInstances = instances
	}
	ecsInstanceTable = createEcsListView(allECSInstances)
	setupTableNavigation(ecsInstanceTable, func(row, col int) {
		instanceId := ecsInstanceTable.GetCell(row, 0).GetReference().(string)
		var selectedInstance ecs.Instance
		for _, inst := range allECSInstances {
			if inst.InstanceId == instanceId {
				selectedInstance = inst
				break
			}
		}
		pages.AddPage(pageEcsDetail, createEcsDetailView(selectedInstance), true, true)
		app.SetFocus(ecsDetailView)
	})
	ecsListFlex := wrapTableInFlex(ecsInstanceTable)
	pages.AddPage(pageEcsList, ecsListFlex, true, true)
	app.SetFocus(ecsInstanceTable)
}

func switchToDnsDomainsListView() {
	if allDomains == nil { /* fetch */
		domains, err := fetchAlidnsDomains()
		if err != nil {
			showErrorModal(err.Error())
			return
		}
		allDomains = domains
	}
	dnsDomainsTable = createDnsDomainsListView(allDomains)
	setupTableNavigation(dnsDomainsTable, func(row, col int) {
		domainName := dnsDomainsTable.GetCell(row, 0).GetReference().(string)
		switchToDnsRecordsListView(domainName)
	})
	dnsDomainsListFlex := wrapTableInFlex(dnsDomainsTable)
	pages.AddPage(pageDnsDomains, dnsDomainsListFlex, true, true)
	app.SetFocus(dnsDomainsTable)
}

func switchToDnsRecordsListView(domainName string) {
	records, err := fetchAlidnsDomainRecords(domainName)
	if err != nil {
		showErrorModal(err.Error())
		return
	}
	dnsRecordsTable = createDnsRecordsListView(records, domainName)
	setupTableNavigation(dnsRecordsTable, nil)
	dnsRecordsListFlex := wrapTableInFlex(dnsRecordsTable)
	pages.AddPage(pageDnsRecords, dnsRecordsListFlex, true, true)
	app.SetFocus(dnsRecordsTable)
}

func switchToSlbListView() {
	if allSLBInstances == nil { /* fetch */
		slbs, err := fetchSLBInstances()
		if err != nil {
			showErrorModal(err.Error())
			return
		}
		allSLBInstances = slbs
	}
	slbInstanceTable = createSlbListView(allSLBInstances)
	setupTableNavigation(slbInstanceTable, func(row, col int) {
		slbId := slbInstanceTable.GetCell(row, 0).GetReference().(string)
		var selectedSlb slb.LoadBalancer
		for _, lb := range allSLBInstances {
			if lb.LoadBalancerId == slbId {
				selectedSlb = lb
				break
			}
		}
		pages.AddPage(pageSlbDetail, createSlbDetailView(selectedSlb), true, true)
		app.SetFocus(slbDetailView)
	})
	slbListFlex := wrapTableInFlex(slbInstanceTable)
	pages.AddPage(pageSlbList, slbListFlex, true, true)
	app.SetFocus(slbInstanceTable)
}

func switchToOssBucketListView() {
	if allOssBuckets == nil {
		buckets, err := fetchOssBuckets()
		if err != nil {
			showErrorModal(err.Error())
			return
		}
		allOssBuckets = buckets
	}
	ossBucketTable = createOssBucketListView(allOssBuckets)
	setupTableNavigation(ossBucketTable, func(row, col int) {
		bucketName := ossBucketTable.GetCell(row, 0).GetReference().(string)
		currentBucketName = bucketName
		switchToOssObjectListView(bucketName)
	})
	ossBucketListFlex := wrapTableInFlex(ossBucketTable)
	pages.AddPage(pageOssBuckets, ossBucketListFlex, true, true)
	app.SetFocus(ossBucketTable)
}

func switchToOssObjectListView(bucketName string) {
	objects, err := fetchOssObjects(bucketName)
	if err != nil {
		showErrorModal(err.Error())
		return
	}
	ossObjectTable = createOssObjectListView(objects, bucketName)
	setupTableNavigation(ossObjectTable, func(row, col int) {
		objectKey := ossObjectTable.GetCell(row, 0).GetReference().(string)
		// Find the object details
		for _, obj := range objects {
			if obj.Key == objectKey {
				ossDetailView = createJSONDetailView(fmt.Sprintf("Object Details: %s", objectKey), obj)
				pages.AddPage("ossObjectDetail", ossDetailView, true, true)
				app.SetFocus(ossDetailView)
				break
			}
		}
	})
	ossObjectListFlex := wrapTableInFlex(ossObjectTable)
	pages.AddPage(pageOssObjects, ossObjectListFlex, true, true)
	app.SetFocus(ossObjectTable)
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
		if !ok {
			return
		}
		var selectedInstance rds.DBInstance
		for _, inst := range allRDSInstances {
			if inst.DBInstanceId == instanceId {
				selectedInstance = inst
				break
			}
		}
		detailViewContent := createRdsDetailView(selectedInstance)
		pages.AddPage(pageRdsDetail, detailViewContent, true, true) // Add and show
		app.SetFocus(rdsDetailView)
	})
	rdsListFlex := wrapTableInFlex(rdsInstanceTable)
	pages.AddPage(pageRdsList, rdsListFlex, true, true) // Add and show
	app.SetFocus(rdsInstanceTable)
}

// Helper for table navigation (j/k)
func setupTableNavigation(table *tview.Table, onSelect func(row, column int)) {
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

// Error Modal
func showErrorModal(message string) {
	modal := tview.NewModal().SetText(message).AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			pages.RemovePage("errorModal")
			// Try to focus the current page's main element or fallback to main menu
			currentPageName, prim := pages.GetFrontPage()
			if prim != nil && currentPageName != "errorModal" {
				app.SetFocus(prim)
			} else {
				app.SetFocus(mainMenu)
			}
		})
	pages.AddPage("errorModal", modal, false, true)
	app.SetFocus(modal)
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
		if modal, isModal := currentFocus.(*tview.Modal); isModal && modal.HasFocus() {
			return event
		}
		if form, isForm := currentFocus.(*tview.Form); isForm && form.HasFocus() &&
			(event.Key() == tcell.KeyTab || event.Key() == tcell.KeyBacktab || event.Key() == tcell.KeyEnter) {
			return event
		}

		if event.Key() == tcell.KeyCtrlC {
			app.Stop()
			return nil
		}

		currentPageName, _ := pages.GetFrontPage()

		handleNavigation := func(targetPage string, focusItem tview.Primitive) {
			pages.SwitchToPage(targetPage)
			if focusItem != nil {
				app.SetFocus(focusItem)
			} else if targetPage == pageMainMenu {
				app.SetFocus(mainMenu)
			}
		}

		switch event.Key() {
		case tcell.KeyEscape:
			switch currentPageName {
			case pageEcsList, pageDnsDomains, pageSlbList, pageOssBuckets, pageRdsList: // Added RdsList
				handleNavigation(pageMainMenu, mainMenu)
			case pageEcsDetail:
				handleNavigation(pageEcsList, ecsInstanceTable)
			case pageDnsRecords:
				handleNavigation(pageDnsDomains, dnsDomainsTable)
			case pageSlbDetail:
				handleNavigation(pageSlbList, slbInstanceTable)
			case pageOssObjects:
				handleNavigation(pageOssBuckets, ossBucketTable)
			case "ossObjectDetail":
				handleNavigation(pageOssObjects, ossObjectTable)
			case pageRdsDetail:
				handleNavigation(pageRdsList, rdsInstanceTable) // Added RdsDetail
			}
			return nil
		case tcell.KeyRune:
			if event.Rune() == 'Q' { // Only uppercase Q exits the program
				app.Stop()
				return nil
			} else if event.Rune() == 'q' { // lowercase q goes back
				switch currentPageName {
				case pageMainMenu, pageEcsList, pageDnsDomains, pageSlbList, pageOssBuckets, pageRdsList:
					// On main pages, q does nothing (only Q exits)
					return nil
				case pageEcsDetail:
					handleNavigation(pageEcsList, ecsInstanceTable)
				case pageDnsRecords:
					handleNavigation(pageDnsDomains, dnsDomainsTable)
				case pageSlbDetail:
					handleNavigation(pageSlbList, slbInstanceTable)
				case pageOssObjects:
					handleNavigation(pageOssBuckets, ossBucketTable)
				case "ossObjectDetail":
					handleNavigation(pageOssObjects, ossObjectTable)
				case pageRdsDetail:
					handleNavigation(pageRdsList, rdsInstanceTable)
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
