package service

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
)

// ECSService handles ECS operations
type ECSService struct {
	client *ecs.Client
}

// NewECSService creates a new ECS service
func NewECSService(client *ecs.Client) *ECSService {
	return &ECSService{client: client}
}

// FetchInstances retrieves all ECS instances using pagination
func (s *ECSService) FetchInstances() ([]ecs.Instance, error) {
	var allInstances []ecs.Instance
	pageNumber := 1
	pageSize := 100 // 使用最大页面大小以减少请求次数

	for {
		request := ecs.CreateDescribeInstancesRequest()
		request.Scheme = "https"
		request.PageNumber = requests.NewInteger(pageNumber)
		request.PageSize = requests.NewInteger(pageSize)

		response, err := s.client.DescribeInstances(request)
		if err != nil {
			return nil, fmt.Errorf("describing ECS instances (page %d): %w", pageNumber, err)
		}

		// 添加当前页的实例到总列表
		allInstances = append(allInstances, response.Instances.Instance...)

		// 检查是否还有更多页面
		// 如果当前页的实例数量小于页面大小，说明这是最后一页
		if len(response.Instances.Instance) < pageSize {
			break
		}

		// 也可以通过TotalCount来判断是否获取完所有数据
		if len(allInstances) >= response.TotalCount {
			break
		}

		pageNumber++
	}

	return allInstances, nil
}
