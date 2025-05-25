package service

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
)

// SLBService handles SLB operations
type SLBService struct {
	client *slb.Client
}

// NewSLBService creates a new SLB service
func NewSLBService(client *slb.Client) *SLBService {
	return &SLBService{client: client}
}

// FetchInstances retrieves all SLB instances
func (s *SLBService) FetchInstances() ([]slb.LoadBalancer, error) {
	request := slb.CreateDescribeLoadBalancersRequest()
	request.Scheme = "https"
	response, err := s.client.DescribeLoadBalancers(request)
	if err != nil {
		return nil, fmt.Errorf("describing SLB instances: %w", err)
	}
	return response.LoadBalancers.LoadBalancer, nil
}
