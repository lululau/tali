package ui

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// CreateEcsListView creates ECS instances list view
func CreateEcsListView(instances []ecs.Instance) *tview.Table {
	table := tview.NewTable().
		SetBorders(true).
		SetSelectable(true, false)
	table = SetupTableWithFixedWidth(table)
	headers := []string{"Instance ID", "Status", "Zone", "CPU/RAM", "Private IP", "Public IP", "Name"}
	CreateTableHeaders(table, headers)

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
			} else if instance.EipAddress.IpAddress != "" {
				publicIP = instance.EipAddress.IpAddress
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

// CreateEcsDetailView creates ECS detail view
func CreateEcsDetailView(instance interface{}) *tview.Flex {
	ecsInstance := instance.(ecs.Instance)
	detailView := CreateJSONDetailView(fmt.Sprintf("ECS Details: %s", ecsInstance.InstanceId), instance)
	return CreateDetailViewWithInstructions(detailView)
}

// CreateSecurityGroupsListView creates security groups list view
func CreateSecurityGroupsListView(securityGroups []ecs.SecurityGroup) *tview.Table {
	table := tview.NewTable().
		SetBorders(true).
		SetSelectable(true, false)
	table = SetupTableWithFixedWidth(table)
	headers := []string{"Security Group ID", "Name", "Description", "VPC ID", "Type", "Creation Time"}
	CreateTableHeaders(table, headers)

	if len(securityGroups) == 0 {
		table.SetCell(1, 0, tview.NewTableCell("No security groups found.").SetSelectable(false).SetExpansion(len(headers)).SetAlign(tview.AlignCenter))
	} else {
		for r, sg := range securityGroups {
			table.SetCell(r+1, 0, tview.NewTableCell(sg.SecurityGroupId).SetTextColor(tcell.ColorWhite).SetReference(sg.SecurityGroupId).SetExpansion(1))
			table.SetCell(r+1, 1, tview.NewTableCell(sg.SecurityGroupName).SetTextColor(tcell.ColorWhite).SetExpansion(1))
			table.SetCell(r+1, 2, tview.NewTableCell(sg.Description).SetTextColor(tcell.ColorWhite).SetExpansion(1))
			table.SetCell(r+1, 3, tview.NewTableCell(sg.VpcId).SetTextColor(tcell.ColorWhite).SetExpansion(1))
			table.SetCell(r+1, 4, tview.NewTableCell(sg.SecurityGroupType).SetTextColor(tcell.ColorWhite).SetExpansion(1))
			table.SetCell(r+1, 5, tview.NewTableCell(sg.CreationTime).SetTextColor(tcell.ColorWhite).SetExpansion(1))
		}
	}
	table.SetTitle("Security Groups").SetBorder(true)
	return table
}

// CreateSecurityGroupDetailView creates security group detail view
func CreateSecurityGroupDetailView(securityGroup interface{}) *tview.Flex {
	sg := securityGroup.(ecs.SecurityGroup)
	detailView := CreateJSONDetailView(fmt.Sprintf("Security Group Details: %s", sg.SecurityGroupId), securityGroup)
	return CreateDetailViewWithInstructions(detailView)
}

// CreateSecurityGroupRulesView creates security group rules view
func CreateSecurityGroupRulesView(rulesResponse *ecs.DescribeSecurityGroupAttributeResponse) *tview.Table {
	table := tview.NewTable().
		SetBorders(true).
		SetSelectable(true, false)
	table = SetupTableWithFixedWidth(table)

	// 合并入站和出站规则
	allRules := []interface{}{}

	// 添加入站规则
	for _, rule := range rulesResponse.Permissions.Permission {
		ruleData := map[string]interface{}{
			"Direction":     "Ingress",
			"IpProtocol":    rule.IpProtocol,
			"PortRange":     rule.PortRange,
			"SourceCidrIp":  rule.SourceCidrIp,
			"SourceGroupId": rule.SourceGroupId,
			"Policy":        rule.Policy,
			"Priority":      rule.Priority,
			"NicType":       rule.NicType,
			"Description":   rule.Description,
		}
		allRules = append(allRules, ruleData)
	}

	headers := []string{"Direction", "Protocol", "Port Range", "Source/Dest", "Policy", "Priority", "Description"}
	CreateTableHeaders(table, headers)

	if len(allRules) == 0 {
		table.SetCell(1, 0, tview.NewTableCell("No security group rules found.").SetSelectable(false).SetExpansion(len(headers)).SetAlign(tview.AlignCenter))
	} else {
		for r, ruleInterface := range allRules {
			rule := ruleInterface.(map[string]interface{})

			// 确定源/目标
			sourceDest := ""
			if rule["SourceCidrIp"] != nil && rule["SourceCidrIp"].(string) != "" {
				sourceDest = rule["SourceCidrIp"].(string)
			} else if rule["SourceGroupId"] != nil && rule["SourceGroupId"].(string) != "" {
				sourceDest = rule["SourceGroupId"].(string)
			}

			table.SetCell(r+1, 0, tview.NewTableCell(rule["Direction"].(string)).SetTextColor(tcell.ColorWhite).SetExpansion(1))
			table.SetCell(r+1, 1, tview.NewTableCell(rule["IpProtocol"].(string)).SetTextColor(tcell.ColorWhite).SetExpansion(1))
			table.SetCell(r+1, 2, tview.NewTableCell(rule["PortRange"].(string)).SetTextColor(tcell.ColorWhite).SetExpansion(1))
			table.SetCell(r+1, 3, tview.NewTableCell(sourceDest).SetTextColor(tcell.ColorWhite).SetExpansion(1))
			table.SetCell(r+1, 4, tview.NewTableCell(rule["Policy"].(string)).SetTextColor(tcell.ColorWhite).SetExpansion(1))
			table.SetCell(r+1, 5, tview.NewTableCell(fmt.Sprintf("%v", rule["Priority"])).SetTextColor(tcell.ColorWhite).SetExpansion(1))
			table.SetCell(r+1, 6, tview.NewTableCell(rule["Description"].(string)).SetTextColor(tcell.ColorWhite).SetExpansion(1))
		}
	}
	table.SetTitle(fmt.Sprintf("Security Group Rules: %s", rulesResponse.SecurityGroupId)).SetBorder(true)
	return table
}

// CreateSecurityGroupInstancesView creates view for instances using a security group
func CreateSecurityGroupInstancesView(instances []ecs.Instance, securityGroupId string) *tview.Table {
	table := tview.NewTable().
		SetBorders(true).
		SetSelectable(true, false)
	table = SetupTableWithFixedWidth(table)
	headers := []string{"Instance ID", "Status", "Zone", "CPU/RAM", "Private IP", "Public IP", "Name"}
	CreateTableHeaders(table, headers)

	if len(instances) == 0 {
		table.SetCell(1, 0, tview.NewTableCell("No instances using this security group.").SetSelectable(false).SetExpansion(len(headers)).SetAlign(tview.AlignCenter))
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
			} else if instance.EipAddress.IpAddress != "" {
				publicIP = instance.EipAddress.IpAddress
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
	table.SetTitle(fmt.Sprintf("Instances using Security Group: %s", securityGroupId)).SetBorder(true)
	return table
}

// CreateInstanceSecurityGroupsView creates view for security groups of an instance
func CreateInstanceSecurityGroupsView(securityGroups []ecs.SecurityGroup, instanceId string) *tview.Table {
	table := tview.NewTable().
		SetBorders(true).
		SetSelectable(true, false)
	table = SetupTableWithFixedWidth(table)
	headers := []string{"Security Group ID", "Name", "Description", "VPC ID", "Type", "Creation Time"}
	CreateTableHeaders(table, headers)

	if len(securityGroups) == 0 {
		table.SetCell(1, 0, tview.NewTableCell("No security groups found for this instance.").SetSelectable(false).SetExpansion(len(headers)).SetAlign(tview.AlignCenter))
	} else {
		for r, sg := range securityGroups {
			table.SetCell(r+1, 0, tview.NewTableCell(sg.SecurityGroupId).SetTextColor(tcell.ColorWhite).SetReference(sg.SecurityGroupId).SetExpansion(1))
			table.SetCell(r+1, 1, tview.NewTableCell(sg.SecurityGroupName).SetTextColor(tcell.ColorWhite).SetExpansion(1))
			table.SetCell(r+1, 2, tview.NewTableCell(sg.Description).SetTextColor(tcell.ColorWhite).SetExpansion(1))
			table.SetCell(r+1, 3, tview.NewTableCell(sg.VpcId).SetTextColor(tcell.ColorWhite).SetExpansion(1))
			table.SetCell(r+1, 4, tview.NewTableCell(sg.SecurityGroupType).SetTextColor(tcell.ColorWhite).SetExpansion(1))
			table.SetCell(r+1, 5, tview.NewTableCell(sg.CreationTime).SetTextColor(tcell.ColorWhite).SetExpansion(1))
		}
	}
	table.SetTitle(fmt.Sprintf("Security Groups for Instance: %s", instanceId)).SetBorder(true)
	return table
}

// CreateDnsDomainsListView creates DNS domains list view
func CreateDnsDomainsListView(domains []alidns.DomainInDescribeDomains) *tview.Table {
	table := tview.NewTable().SetBorders(true).SetSelectable(true, false)
	table = SetupTableWithFixedWidth(table)
	headers := []string{"Domain Name", "Record Count", "Version Code"}
	CreateTableHeaders(table, headers)

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

// CreateDnsRecordsListView creates DNS records list view
func CreateDnsRecordsListView(records []alidns.Record, domainName string) *tview.Table {
	table := tview.NewTable().SetBorders(true).SetSelectable(true, false)
	table = SetupTableWithFixedWidth(table)
	headers := []string{"Record ID", "RR", "Type", "Value", "TTL", "Status"}
	CreateTableHeaders(table, headers)

	if len(records) == 0 {
		table.SetCell(1, 0, tview.NewTableCell("No DNS records found.").SetSelectable(false).SetExpansion(len(headers)).SetAlign(tview.AlignCenter))
	} else {
		for r, record := range records {
			table.SetCell(r+1, 0, tview.NewTableCell(record.RecordId).SetTextColor(tcell.ColorWhite).SetReference(record.RecordId).SetExpansion(1))
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

// CreateSlbListView creates SLB instances list view
func CreateSlbListView(slbs []slb.LoadBalancer) *tview.Table {
	table := tview.NewTable().SetBorders(true).SetSelectable(true, false)
	table = SetupTableWithFixedWidth(table)
	headers := []string{"SLB ID", "Name", "IP Address", "Type", "Status"}
	CreateTableHeaders(table, headers)

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

// CreateSlbDetailView creates SLB detail view
func CreateSlbDetailView(lb interface{}) *tview.Flex {
	slbInstance := lb.(slb.LoadBalancer)
	detailView := CreateJSONDetailView(fmt.Sprintf("SLB Details: %s", slbInstance.LoadBalancerId), lb)
	return CreateDetailViewWithInstructions(detailView)
}

// CreateOssBucketListView creates OSS buckets list view
func CreateOssBucketListView(buckets []oss.BucketProperties) *tview.Table {
	table := tview.NewTable().SetBorders(true).SetSelectable(true, false)
	table = SetupTableWithFixedWidth(table)
	headers := []string{"Bucket Name", "Location", "Creation Date", "Storage Class"}
	CreateTableHeaders(table, headers)

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

// CreateOssObjectListView creates OSS objects list view
func CreateOssObjectListView(objects []oss.ObjectProperties, bucketName string) *tview.Table {
	table := tview.NewTable().SetBorders(true).SetSelectable(true, false)
	table = SetupTableWithFixedWidth(table)
	headers := []string{"Object Key", "Size (Bytes)", "Last Modified", "Storage Class", "ETag"}
	CreateTableHeaders(table, headers)

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

// CreateOssObjectPaginatedView creates OSS objects list view with pagination controls
func CreateOssObjectPaginatedView(objects []oss.ObjectProperties, bucketName string, currentPage int, hasNext, hasPrev bool) *tview.Flex {
	// Create the table
	table := tview.NewTable().SetBorders(true).SetSelectable(true, false)
	table = SetupTableWithFixedWidth(table)
	headers := []string{"Object Key", "Size (Bytes)", "Last Modified", "Storage Class", "ETag"}
	CreateTableHeaders(table, headers)

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

	// Create pagination info
	paginationInfo := ""
	if hasPrev {
		paginationInfo += "[ (Prev) "
	}
	paginationInfo += fmt.Sprintf("Page %d", currentPage)
	if hasNext {
		paginationInfo += " ] (Next)"
	}
	paginationInfo += " | Press '[' for previous, ']' for next, '0' for first page"

	// Create pagination status bar
	statusBar := tview.NewTextView().
		SetText(paginationInfo).
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetBackgroundColor(tcell.ColorReset)
	statusBar.SetBorder(false)

	// Create flex container
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.AddItem(table, 0, 1, true)
	flex.AddItem(statusBar, 1, 0, false)
	flex.SetTitle(fmt.Sprintf("Objects in %s", bucketName)).SetBorder(true)
	flex.SetBackgroundColor(tcell.ColorReset)

	return flex
}

// CreateRdsListView creates RDS instances list view
func CreateRdsListView(instances []rds.DBInstance) *tview.Table {
	table := tview.NewTable().
		SetBorders(true).
		SetSelectable(true, false)
	table = SetupTableWithFixedWidth(table)

	headers := []string{"Instance ID", "Engine", "Version", "Class", "Status", "Description"}
	CreateTableHeaders(table, headers)

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

// CreateRdsDetailView creates RDS detail view
func CreateRdsDetailView(instance interface{}) *tview.Flex {
	rdsInstance := instance.(rds.DBInstance)
	detailView := CreateJSONDetailView(fmt.Sprintf("RDS Details: %s", rdsInstance.DBInstanceId), instance)
	return CreateDetailViewWithInstructions(detailView)
}

// CreateRdsDatabasesListView creates RDS databases list view
func CreateRdsDatabasesListView(databases []rds.Database, instanceId string) *tview.Table {
	table := tview.NewTable().
		SetBorders(true).
		SetSelectable(true, false)
	table = SetupTableWithFixedWidth(table)

	headers := []string{"Database Name", "Status", "Character Set", "Bound Accounts", "Description"}
	CreateTableHeaders(table, headers)

	if len(databases) == 0 {
		table.SetCell(1, 0, tview.NewTableCell("No databases found.").SetSelectable(false).SetExpansion(len(headers)).SetAlign(tview.AlignCenter))
	} else {
		for r, db := range databases {
			// Format bound accounts
			boundAccounts := ""
			if len(db.Accounts.AccountPrivilegeInfo) > 0 {
				for i, account := range db.Accounts.AccountPrivilegeInfo {
					if i > 0 {
						boundAccounts += ", "
					}
					boundAccounts += account.Account
				}
			} else {
				boundAccounts = "--"
			}

			table.SetCell(r+1, 0, tview.NewTableCell(db.DBName).SetTextColor(tcell.ColorWhite).SetReference(db.DBName).SetExpansion(1))
			table.SetCell(r+1, 1, tview.NewTableCell(db.DBStatus).SetTextColor(tcell.ColorWhite).SetExpansion(1))
			table.SetCell(r+1, 2, tview.NewTableCell(db.CharacterSetName).SetTextColor(tcell.ColorWhite).SetExpansion(1))
			table.SetCell(r+1, 3, tview.NewTableCell(boundAccounts).SetTextColor(tcell.ColorWhite).SetExpansion(1))
			table.SetCell(r+1, 4, tview.NewTableCell(db.DBDescription).SetTextColor(tcell.ColorWhite).SetExpansion(1))
		}
	}
	table.SetTitle(fmt.Sprintf("Databases for RDS Instance: %s", instanceId)).SetBorder(true)
	return table
}

// CreateRdsAccountsListView creates RDS accounts list view
func CreateRdsAccountsListView(accounts []rds.DBInstanceAccount, instanceId string) *tview.Table {
	table := tview.NewTable().
		SetBorders(true).
		SetSelectable(true, false)
	table = SetupTableWithFixedWidth(table)

	headers := []string{"Account Name", "Type", "Status", "Bound Databases", "Description"}
	CreateTableHeaders(table, headers)

	if len(accounts) == 0 {
		table.SetCell(1, 0, tview.NewTableCell("No accounts found.").SetSelectable(false).SetExpansion(len(headers)).SetAlign(tview.AlignCenter))
	} else {
		for r, account := range accounts {
			// Format bound databases
			boundDatabases := ""
			if len(account.DatabasePrivileges.DatabasePrivilege) > 0 {
				for i, dbPriv := range account.DatabasePrivileges.DatabasePrivilege {
					if i > 0 {
						boundDatabases += ", "
					}
					boundDatabases += fmt.Sprintf("%s(%s)", dbPriv.DBName, dbPriv.AccountPrivilege)
				}
			} else {
				boundDatabases = "--"
			}

			table.SetCell(r+1, 0, tview.NewTableCell(account.AccountName).SetTextColor(tcell.ColorWhite).SetReference(account.AccountName).SetExpansion(1))
			table.SetCell(r+1, 1, tview.NewTableCell(account.AccountType).SetTextColor(tcell.ColorWhite).SetExpansion(1))
			table.SetCell(r+1, 2, tview.NewTableCell(account.AccountStatus).SetTextColor(tcell.ColorWhite).SetExpansion(1))
			table.SetCell(r+1, 3, tview.NewTableCell(boundDatabases).SetTextColor(tcell.ColorWhite).SetExpansion(1))
			table.SetCell(r+1, 4, tview.NewTableCell(account.AccountDescription).SetTextColor(tcell.ColorWhite).SetExpansion(1))
		}
	}
	table.SetTitle(fmt.Sprintf("Accounts for RDS Instance: %s", instanceId)).SetBorder(true)
	return table
}
