package service

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
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

// FetchInstances retrieves all SLB instances using pagination
func (s *SLBService) FetchInstances() ([]slb.LoadBalancer, error) {
	var allLoadBalancers []slb.LoadBalancer
	pageNumber := int64(1)
	pageSize := int64(100)

	for {
		request := slb.CreateDescribeLoadBalancersRequest()
		request.Scheme = "https"
		request.PageNumber = requests.NewInteger(int(pageNumber))
		request.PageSize = requests.NewInteger(int(pageSize))

		response, err := s.client.DescribeLoadBalancers(request)
		if err != nil {
			return nil, fmt.Errorf("describing SLB instances (page %d): %w", pageNumber, err)
		}

		allLoadBalancers = append(allLoadBalancers, response.LoadBalancers.LoadBalancer...)

		if pageNumber*pageSize >= int64(response.TotalCount) {
			break
		}

		if len(response.LoadBalancers.LoadBalancer) < int(pageSize) {
			break
		}

		pageNumber++
	}
	return allLoadBalancers, nil
}
