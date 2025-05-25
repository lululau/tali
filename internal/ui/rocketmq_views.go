package ui

import (
	"time"

	"aliyun-tui-viewer/internal/service"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// CreateRocketMQListView creates a table view for RocketMQ instances
func CreateRocketMQListView(instances []service.RocketMQInstance) *tview.Table {
	table := tview.NewTable().
		SetBorders(true).
		SetSelectable(true, false)
	table = SetupTableWithFixedWidth(table)

	// Set headers - removed Region and Remark since they're not available
	headers := []string{"Instance ID", "Instance Name", "Type", "Status", "Create Time"}
	CreateTableHeaders(table, headers)

	// Add data rows
	if len(instances) == 0 {
		table.SetCell(1, 0, tview.NewTableCell("No RocketMQ instances found.").SetSelectable(false).SetExpansion(len(headers)).SetAlign(tview.AlignCenter))
	} else {
		for i, instance := range instances {
			row := i + 1

			// Instance ID
			table.SetCell(row, 0, tview.NewTableCell(instance.InstanceId).
				SetTextColor(tcell.ColorWhite).
				SetReference(instance.InstanceId).
				SetExpansion(1))

			// Instance Name
			table.SetCell(row, 1, tview.NewTableCell(instance.InstanceName).
				SetTextColor(tcell.ColorWhite).
				SetExpansion(1))

			// Instance Type
			instanceType := "Unknown"
			switch instance.InstanceType {
			case 1:
				instanceType = "Standard"
			case 2:
				instanceType = "Platinum"
			}
			table.SetCell(row, 2, tview.NewTableCell(instanceType).
				SetTextColor(tcell.ColorWhite).
				SetExpansion(1))

			// Instance Status
			status := "Unknown"
			switch instance.InstanceStatus {
			case 0:
				status = "Deploying"
			case 2:
				status = "Arrears"
			case 5:
				status = "Running"
			case 7:
				status = "Upgrading"
			}
			table.SetCell(row, 3, tview.NewTableCell(status).
				SetTextColor(tcell.ColorWhite).
				SetExpansion(1))

			// Create Time
			createTime := ""
			if instance.CreateTime > 0 {
				createTime = time.Unix(instance.CreateTime/1000, 0).Format("2006-01-02 15:04:05")
			}
			table.SetCell(row, 4, tview.NewTableCell(createTime).
				SetTextColor(tcell.ColorWhite).
				SetExpansion(1))
		}
	}

	return table
}

// CreateRocketMQTopicsListView creates a table view for RocketMQ topics
func CreateRocketMQTopicsListView(topics []service.RocketMQTopic, instanceId string) *tview.Table {
	table := tview.NewTable().
		SetBorders(true).
		SetSelectable(true, false)
	table = SetupTableWithFixedWidth(table)

	// Set headers - removed fields that are not available
	headers := []string{"Topic", "Message Type", "Create Time", "Remark"}
	CreateTableHeaders(table, headers)

	// Add data rows
	if len(topics) == 0 {
		table.SetCell(1, 0, tview.NewTableCell("No topics found.").SetSelectable(false).SetExpansion(len(headers)).SetAlign(tview.AlignCenter))
	} else {
		for i, topic := range topics {
			row := i + 1

			// Topic
			table.SetCell(row, 0, tview.NewTableCell(topic.Topic).
				SetTextColor(tcell.ColorWhite).
				SetReference(topic.Topic).
				SetExpansion(1))

			// Message Type
			messageType := "Unknown"
			switch topic.MessageType {
			case 0:
				messageType = "Normal"
			case 1:
				messageType = "Partition Ordered"
			case 2:
				messageType = "Global Ordered"
			case 4:
				messageType = "Transaction"
			case 5:
				messageType = "Scheduled/Delayed"
			}
			table.SetCell(row, 1, tview.NewTableCell(messageType).
				SetTextColor(tcell.ColorWhite).
				SetExpansion(1))

			// Create Time
			createTime := ""
			if topic.CreateTime > 0 {
				createTime = time.Unix(topic.CreateTime/1000, 0).Format("2006-01-02 15:04:05")
			}
			table.SetCell(row, 2, tview.NewTableCell(createTime).
				SetTextColor(tcell.ColorWhite).
				SetExpansion(1))

			// Remark
			table.SetCell(row, 3, tview.NewTableCell(topic.Remark).
				SetTextColor(tcell.ColorWhite).
				SetExpansion(1))
		}
	}

	table.SetTitle("RocketMQ Topics").SetBorder(true)
	return table
}

// CreateRocketMQGroupsListView creates a table view for RocketMQ consumer groups
func CreateRocketMQGroupsListView(groups []service.RocketMQGroup, instanceId string) *tview.Table {
	table := tview.NewTable().
		SetBorders(true).
		SetSelectable(true, false)
	table = SetupTableWithFixedWidth(table)

	// Set headers
	headers := []string{"Group ID", "Group Type", "Create Time", "Update Time", "Remark"}
	CreateTableHeaders(table, headers)

	// Add data rows
	if len(groups) == 0 {
		table.SetCell(1, 0, tview.NewTableCell("No consumer groups found.").SetSelectable(false).SetExpansion(len(headers)).SetAlign(tview.AlignCenter))
	} else {
		for i, group := range groups {
			row := i + 1

			// Group ID
			table.SetCell(row, 0, tview.NewTableCell(group.GroupId).
				SetTextColor(tcell.ColorWhite).
				SetReference(group.GroupId).
				SetExpansion(1))

			// Group Type
			table.SetCell(row, 1, tview.NewTableCell(group.GroupType).
				SetTextColor(tcell.ColorWhite).
				SetExpansion(1))

			// Create Time
			createTime := ""
			if group.CreateTime > 0 {
				createTime = time.Unix(group.CreateTime/1000, 0).Format("2006-01-02 15:04:05")
			}
			table.SetCell(row, 2, tview.NewTableCell(createTime).
				SetTextColor(tcell.ColorWhite).
				SetExpansion(1))

			// Update Time
			updateTime := ""
			if group.UpdateTime > 0 {
				updateTime = time.Unix(group.UpdateTime/1000, 0).Format("2006-01-02 15:04:05")
			}
			table.SetCell(row, 3, tview.NewTableCell(updateTime).
				SetTextColor(tcell.ColorWhite).
				SetExpansion(1))

			// Remark
			table.SetCell(row, 4, tview.NewTableCell(group.Remark).
				SetTextColor(tcell.ColorWhite).
				SetExpansion(1))
		}
	}

	table.SetTitle("RocketMQ Consumer Groups").SetBorder(true)
	return table
}
