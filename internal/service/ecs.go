package service

import (
	"fmt"

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

// FetchInstances retrieves all ECS instances
func (s *ECSService) FetchInstances() ([]ecs.Instance, error) {
	request := ecs.CreateDescribeInstancesRequest()
	request.Scheme = "https"
	response, err := s.client.DescribeInstances(request)
	if err != nil {
		return nil, fmt.Errorf("describing ECS instances: %w", err)
	}
	return response.Instances.Instance, nil
}
