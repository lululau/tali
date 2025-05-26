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

// App represents the main application
type App struct {
	tviewApp *tview.Application
	pages    *tview.Pages
	clients  *client.AliyunClients
	services *Services

	// UI components
	mainMenu                           *tview.List
	ecsInstanceTable                   *tview.Table
	ecsDetailView                      *tview.TextView
	securityGroupTable                 *tview.Table
	securityGroupDetailView            *tview.TextView
	securityGroupRulesTable            *tview.Table
	securityGroupInstancesTable        *tview.Table
	instanceSecurityGroupsTable        *tview.Table
	dnsDomainsTable                    *tview.Table
	dnsRecordsTable                    *tview.Table
	slbInstanceTable                   *tview.Table
	slbDetailView                      *tview.TextView
	slbListenersTable                  *tview.Table
	slbVServerGroupsTable              *tview.Table
	slbVServerGroupBackendServersTable *tview.Table
	ossBucketTable                     *tview.Table
	ossObjectTable                     *tview.Table
	ossDetailView                      *tview.TextView
	rdsInstanceTable                   *tview.Table
	rdsDetailView                      *tview.TextView
	rdsDatabaseTable                   *tview.Table
	rdsAccountTable                    *tview.Table
	redisInstanceTable                 *tview.Table
	redisAccountTable                  *tview.Table
	rocketmqInstanceTable              *tview.Table
	rocketmqTopicsTable                *tview.Table
	rocketmqGroupsTable                *tview.Table
	modeLine                           *tview.TextView
	mainLayout                         *tview.Flex // Keep for now, might remove if root structure changes significantly

	// Shared Search UI
	searchBar           *tview.InputField
	searchBarContainer  *tview.Pages
	activeSearchHandler *ui.VimSearchHandler

	// Data cache
	allECSInstances           []ecs.Instance
	allSecurityGroups         []ecs.SecurityGroup
	allDomains                []alidns.DomainInDescribeDomains
	allSLBInstances           []slb.LoadBalancer
	allRDSInstances           []rds.DBInstance
	allRedisInstances         []r_kvstore.KVStoreInstance
	allRocketMQInstances      []service.RocketMQInstance
	allOssBuckets             []oss.BucketProperties
	currentBucketName         string
	currentRdsInstanceId      string
	currentRedisInstanceId    string
	currentRocketMQInstanceId string

	// OSS pagination state
	ossCurrentMarker   string
	ossPreviousMarkers []string // Stack to track previous markers for backward navigation
	ossCurrentPage     int
	ossPageSize        int
	ossHasNextPage     bool

	// Configuration
	currentProfile string

	// Interaction state
	yankTracker       *ui.YankTracker
	currentDetailData interface{} // Store current detail data for copying/editing
}

// Services holds all service instances
type Services struct {
	ECS      *service.ECSService
	DNS      *service.DNSService
	SLB      *service.SLBService
	RDS      *service.RDSService
	OSS      *service.OSSService
	Redis    *service.RedisService
	RocketMQ *service.RocketMQService
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
		ECS:      service.NewECSService(clients.ECS),
		DNS:      service.NewDNSService(clients.DNS),
		SLB:      service.NewSLBService(clients.SLB),
		RDS:      service.NewRDSService(clients.RDS),
		OSS:      service.NewOSSServiceWithCredentials(clients.OSS, cfg.AccessKeyID, cfg.AccessKeySecret, cfg.OssEndpoint),
		Redis:    service.NewRedisService(clients.Redis),
		RocketMQ: service.NewRocketMQService(clients.RocketMQ),
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
		yankTracker:    ui.NewYankTracker(),

		// Search handlers will be initialized when creating views
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
		a.switchToSecurityGroupsListView,
		a.switchToDnsDomainsListView,
		a.switchToSlbListView,
		a.switchToOssBucketListView,
		a.switchToRdsListView,
		a.switchToRedisListView,
		a.switchToRocketMQListView,
		a.Stop,
	)

	// Create shared search bar
	a.searchBar = ui.CreateSearchBar()
	a.searchBar.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if a.activeSearchHandler == nil {
			return event
		}
		switch event.Key() {
		case tcell.KeyEnter:
			a.activeSearchHandler.PerformSearch(a.searchBar.GetText())
			return nil
		case tcell.KeyEscape:
			a.activeSearchHandler.ExitSearchMode()
			return nil
		}
		return event
	})

	// Create search bar container (for visibility control)
	a.searchBarContainer = tview.NewPages()
	emptyBox := tview.NewBox().SetBorder(false) // Placeholder for hidden search bar
	a.searchBarContainer.AddPage("visible", a.searchBar, true, false)
	a.searchBarContainer.AddPage("hidden", emptyBox, true, true) // Initially hidden

	// Create main layout
	a.mainLayout = tview.NewFlex().SetDirection(tview.FlexRow)
	a.mainLayout.AddItem(a.pages, 0, 1, true)               // Main content
	a.mainLayout.AddItem(a.searchBarContainer, 1, 0, false) // Search bar (or empty space)
	a.mainLayout.AddItem(a.modeLine, 1, 0, false)           // Mode line

	// Add main menu to pages
	a.pages.AddPage(ui.PageMainMenu, a.mainMenu, true, true)

	// Set the main layout as root
	a.tviewApp.SetRoot(a.mainLayout, true)

	// Set up global input capture
	a.setupGlobalInputCapture()
}

// SetSearchBarVisibility controls the visibility of the shared search bar
func (a *App) SetSearchBarVisibility(visible bool) {
	if visible {
		a.searchBar.SetText("")
		a.searchBarContainer.SwitchToPage("visible")
		a.tviewApp.SetFocus(a.searchBar)
	} else {
		a.searchBarContainer.SwitchToPage("hidden")
		if a.activeSearchHandler != nil && a.activeSearchHandler.GetMainComponent() != nil {
			a.tviewApp.SetFocus(a.activeSearchHandler.GetMainComponent())
		}
	}
	// It might be necessary to call a.tviewApp.Draw() here if updates are not immediate
}

// SetActiveSearchHandler sets the currently active search handler
func (a *App) SetActiveSearchHandler(handler *ui.VimSearchHandler) {
	a.activeSearchHandler = handler
}

// GetAppSearchBar returns the shared search bar instance
func (a *App) GetAppSearchBar() *tview.InputField {
	return a.searchBar
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
