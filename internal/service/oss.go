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

// FetchBuckets retrieves all OSS buckets using pagination
func (s *OSSService) FetchBuckets() ([]oss.BucketProperties, error) {
	var allBuckets []oss.BucketProperties
	marker := ""
	for {
		options := []oss.Option{
			oss.MaxKeys(100),
			oss.Marker(marker),
		}
		result, err := s.client.ListBuckets(options...)
		if err != nil {
			return nil, fmt.Errorf("listing OSS buckets (marker: %s): %w", marker, err)
		}
		allBuckets = append(allBuckets, result.Buckets...)

		if !result.IsTruncated {
			break
		}
		marker = result.NextMarker
	}
	return allBuckets, nil
}

// FetchObjects retrieves objects from a specific bucket using pagination
func (s *OSSService) FetchObjects(bucketName string) ([]oss.ObjectProperties, error) {
	bucket, err := s.client.Bucket(bucketName)
	if err != nil {
		return nil, fmt.Errorf("getting bucket %s: %w", bucketName, err)
	}

	var allObjects []oss.ObjectProperties
	marker := ""
	for {
		options := []oss.Option{
			oss.MaxKeys(100),
			oss.Marker(marker),
		}
		result, err := bucket.ListObjects(options...)
		if err != nil {
			return nil, fmt.Errorf("listing objects in bucket %s (marker: %s): %w", bucketName, marker, err)
		}

		allObjects = append(allObjects, result.Objects...)

		if !result.IsTruncated {
			break
		}
		marker = result.NextMarker
	}
	return allObjects, nil
}
