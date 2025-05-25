package service

import (
	"fmt"

	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
)

// RedisService handles r-kvstore related operations
type RedisService struct {
	client *r_kvstore.Client
}

// NewRedisService creates a new RedisService
func NewRedisService(client *r_kvstore.Client) *RedisService {
	return &RedisService{client: client}
}

// FetchInstances fetches all Redis instances
func (s *RedisService) FetchInstances() ([]r_kvstore.KVStoreInstance, error) {
	request := r_kvstore.CreateDescribeInstancesRequest()
	request.Scheme = "https"
	// Set PageSize to a large number to fetch all instances in one go,
	// as pagination might be complex to implement quickly for this new feature.
	// Max PageSize is typically 50 or 100. Let's use 100.
	request.PageSize = "100"

	response, err := s.client.DescribeInstances(request)
	if err != nil {
		return nil, fmt.Errorf("fetching redis instances: %w", err)
	}
	return response.Instances.KVStoreInstance, nil
}

// FetchAccounts fetches all accounts for a specific Redis instance
func (s *RedisService) FetchAccounts(instanceID string) ([]r_kvstore.Account, error) {
	request := r_kvstore.CreateDescribeAccountsRequest()
	request.Scheme = "https"
	request.InstanceId = instanceID

	response, err := s.client.DescribeAccounts(request)
	if err != nil {
		return nil, fmt.Errorf("fetching redis accounts for instance %s: %w", instanceID, err)
	}
	return response.Accounts.Account, nil
}

// TODO: Add methods for fetching Redis instance details if needed, similar to ECS/RDS.
// For now, the instance list provides most common details.
