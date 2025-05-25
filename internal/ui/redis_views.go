package ui

import (
	"fmt"

	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// CreateRedisListView creates Redis instances list view
func CreateRedisListView(instances []r_kvstore.KVStoreInstance) *tview.Table {
	table := tview.NewTable().
		SetBorders(true).
		SetSelectable(true, false)
	table = SetupTableWithFixedWidth(table)

	headers := []string{"Instance ID", "Instance Name", "Type", "Version", "Status", "Region", "Capacity", "Connection Domain"}
	CreateTableHeaders(table, headers)

	if len(instances) == 0 {
		table.SetCell(1, 0, tview.NewTableCell("No Redis instances found.").SetSelectable(false).SetExpansion(len(headers)).SetAlign(tview.AlignCenter))
	} else {
		for r, inst := range instances {
			table.SetCell(r+1, 0, tview.NewTableCell(inst.InstanceId).SetTextColor(tcell.ColorWhite).SetReference(inst.InstanceId).SetExpansion(1))
			table.SetCell(r+1, 1, tview.NewTableCell(inst.InstanceName).SetTextColor(tcell.ColorWhite).SetExpansion(1))
			table.SetCell(r+1, 2, tview.NewTableCell(inst.InstanceType).SetTextColor(tcell.ColorWhite).SetExpansion(1))
			table.SetCell(r+1, 3, tview.NewTableCell(inst.EngineVersion).SetTextColor(tcell.ColorWhite).SetExpansion(1))
			table.SetCell(r+1, 4, tview.NewTableCell(inst.InstanceStatus).SetTextColor(tcell.ColorWhite).SetExpansion(1))
			table.SetCell(r+1, 5, tview.NewTableCell(inst.RegionId).SetTextColor(tcell.ColorWhite).SetExpansion(1))
			table.SetCell(r+1, 6, tview.NewTableCell(fmt.Sprintf("%d MB", inst.Capacity)).SetTextColor(tcell.ColorWhite).SetExpansion(1))
			table.SetCell(r+1, 7, tview.NewTableCell(inst.ConnectionDomain).SetTextColor(tcell.ColorWhite).SetExpansion(1))
		}
	}
	table.SetTitle("Redis Instances").SetBorder(true)
	return table
}

// CreateRedisAccountsListView creates Redis accounts list view
func CreateRedisAccountsListView(accounts []r_kvstore.Account, instanceId string) *tview.Table {
	table := tview.NewTable().
		SetBorders(true).
		SetSelectable(true, false)
	table = SetupTableWithFixedWidth(table)

	headers := []string{"Account Name", "Status", "Type"}
	CreateTableHeaders(table, headers)

	if len(accounts) == 0 {
		table.SetCell(1, 0, tview.NewTableCell("No Redis accounts found.").SetSelectable(false).SetExpansion(len(headers)).SetAlign(tview.AlignCenter))
	} else {
		for r, acc := range accounts {
			table.SetCell(r+1, 0, tview.NewTableCell(acc.AccountName).SetTextColor(tcell.ColorWhite).SetReference(acc.AccountName).SetExpansion(1))
			table.SetCell(r+1, 1, tview.NewTableCell(acc.AccountStatus).SetTextColor(tcell.ColorWhite).SetExpansion(1))
			table.SetCell(r+1, 2, tview.NewTableCell(acc.AccountType).SetTextColor(tcell.ColorWhite).SetExpansion(1))
		}
	}
	table.SetTitle(fmt.Sprintf("Accounts for Redis Instance: %s", instanceId)).SetBorder(true)
	return table
}
