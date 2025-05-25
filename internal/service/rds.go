package service

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
)

// RDSService handles RDS operations
type RDSService struct {
	client *rds.Client
}

// NewRDSService creates a new RDS service
func NewRDSService(client *rds.Client) *RDSService {
	return &RDSService{client: client}
}

// FetchInstances retrieves all RDS instances using pagination
func (s *RDSService) FetchInstances() ([]rds.DBInstance, error) {
	var allInstances []rds.DBInstance
	pageNumber := 1
	pageSize := 100 // 使用最大页面大小以减少请求次数

	for {
		request := rds.CreateDescribeDBInstancesRequest()
		request.Scheme = "https"
		request.PageNumber = requests.NewInteger(pageNumber)
		request.PageSize = requests.NewInteger(pageSize)

		response, err := s.client.DescribeDBInstances(request)
		if err != nil {
			return nil, fmt.Errorf("describing RDS instances (page %d): %w", pageNumber, err)
		}

		// 添加当前页的实例到总列表
		allInstances = append(allInstances, response.Items.DBInstance...)

		// 检查是否还有更多页面
		// 如果当前页的实例数量小于页面大小，说明这是最后一页
		if len(response.Items.DBInstance) < pageSize {
			break
		}

		// 也可以通过TotalRecordCount来判断是否获取完所有数据
		if len(allInstances) >= response.TotalRecordCount {
			break
		}

		pageNumber++
	}

	return allInstances, nil
}
