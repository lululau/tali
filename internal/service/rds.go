package service

import (
	"fmt"

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

// FetchInstances retrieves all RDS instances
func (s *RDSService) FetchInstances() ([]rds.DBInstance, error) {
	request := rds.CreateDescribeDBInstancesRequest()
	request.Scheme = "https"
	response, err := s.client.DescribeDBInstances(request)
	if err != nil {
		return nil, fmt.Errorf("describing RDS instances: %w", err)
	}
	return response.Items.DBInstance, nil
}
