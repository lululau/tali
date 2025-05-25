package app

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
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

// App represents the main application
type App struct {
	tviewApp *tview.Application
	pages    *tview.Pages
	clients  *client.AliyunClients
	services *Services

	// UI components
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
	modeLine         *tview.TextView

	// Data cache
	allECSInstances   []ecs.Instance
	allDomains        []alidns.DomainInDescribeDomains
	allSLBInstances   []slb.LoadBalancer
	allRDSInstances   []rds.DBInstance
	allOssBuckets     []oss.BucketProperties
	currentBucketName string

	// OSS pagination state
	ossCurrentMarker   string
	ossPreviousMarkers []string // Stack to track previous markers for backward navigation
	ossCurrentPage     int
	ossPageSize        int
	ossHasNextPage     bool

	// Configuration
	currentProfile string
}

// Services holds all service instances
type Services struct {
	ECS *service.ECSService
	DNS *service.DNSService
	SLB *service.SLBService
	RDS *service.RDSService
	OSS *service.OSSService
}

// New creates a new application instance
func New() (*App, error) {
	// Load configuration
	cfg, err := config.LoadAliyunConfig()
	if err != nil {
		return nil, fmt.Errorf("loading config: %w", err)
	}

	// Get current profile name
	currentProfile, err := config.GetCurrentProfileName()
	if err != nil {
		return nil, fmt.Errorf("getting current profile: %w", err)
	}

	// Create clients
	clientConfig := &client.Config{
		AccessKeyID:     cfg.AccessKeyID,
		AccessKeySecret: cfg.AccessKeySecret,
		RegionID:        cfg.RegionID,
		OssEndpoint:     cfg.OssEndpoint,
	}

	clients, err := client.NewAliyunClients(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("creating clients: %w", err)
	}

	// Create services
	services := &Services{
		ECS: service.NewECSService(clients.ECS),
		DNS: service.NewDNSService(clients.DNS),
		SLB: service.NewSLBService(clients.SLB),
		RDS: service.NewRDSService(clients.RDS),
		OSS: service.NewOSSServiceWithCredentials(clients.OSS, cfg.AccessKeyID, cfg.AccessKeySecret, cfg.OssEndpoint),
	}

	// Create tview app and pages
	tviewApp := tview.NewApplication()
	pages := tview.NewPages()
	pages.SetBackgroundColor(tcell.ColorReset)

	app := &App{
		tviewApp:       tviewApp,
		pages:          pages,
		clients:        clients,
		services:       services,
		currentProfile: currentProfile,
	}

	// Initialize UI
	app.initializeUI()

	return app, nil
}

// Run starts the application
func (a *App) Run() error {
	return a.tviewApp.EnableMouse(true).Run()
}

// Stop stops the application
func (a *App) Stop() {
	a.tviewApp.Stop()
}

// initializeUI initializes the user interface
func (a *App) initializeUI() {
	// Create mode line
	a.modeLine = ui.CreateModeLine(a.currentProfile)

	// Create main menu
	a.mainMenu = ui.CreateMainMenu(
		a.switchToEcsListView,
		a.switchToDnsDomainsListView,
		a.switchToSlbListView,
		a.switchToOssBucketListView,
		a.switchToRdsListView,
		a.Stop,
	)

	// Create main layout with mode line at bottom
	mainLayout := tview.NewFlex().SetDirection(tview.FlexRow)
	mainLayout.AddItem(a.pages, 0, 1, true)
	mainLayout.AddItem(a.modeLine, 1, 0, false)

	// Add main menu to pages
	a.pages.AddPage(ui.PageMainMenu, a.mainMenu, true, true)

	// Set the main layout as root instead of pages directly
	a.tviewApp.SetRoot(mainLayout, true)

	// Set up global input capture
	a.setupGlobalInputCapture()
}

// showErrorModal shows an error modal
func (a *App) showErrorModal(message string) {
	ui.ShowErrorModal(a.pages, a.tviewApp, message, func() {
		// Try to focus the current page's main element or fallback to main menu
		currentPageName, prim := a.pages.GetFrontPage()
		if prim != nil && currentPageName != "errorModal" {
			a.tviewApp.SetFocus(prim)
		} else {
			a.tviewApp.SetFocus(a.mainMenu)
		}
	})
}
