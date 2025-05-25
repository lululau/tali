package service

import (
	"fmt"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// OSSService handles OSS operations
type OSSService struct {
	client *oss.Client
}

// NewOSSService creates a new OSS service
func NewOSSService(client *oss.Client) *OSSService {
	return &OSSService{client: client}
}

// FetchBuckets retrieves all OSS buckets
func (s *OSSService) FetchBuckets() ([]oss.BucketProperties, error) {
	result, err := s.client.ListBuckets()
	if err != nil {
		return nil, fmt.Errorf("listing OSS buckets: %w", err)
	}
	return result.Buckets, nil
}

// FetchObjects retrieves objects from a specific bucket
func (s *OSSService) FetchObjects(bucketName string) ([]oss.ObjectProperties, error) {
	bucket, err := s.client.Bucket(bucketName)
	if err != nil {
		return nil, fmt.Errorf("getting bucket %s: %w", bucketName, err)
	}

	result, err := bucket.ListObjects()
	if err != nil {
		return nil, fmt.Errorf("listing objects in bucket %s: %w", bucketName, err)
	}

	return result.Objects, nil
}
