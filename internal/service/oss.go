package service

import (
	"fmt"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// OSSService handles OSS operations
type OSSService struct {
	client          *oss.Client
	accessKeyID     string
	accessKeySecret string
	defaultEndpoint string
}

// NewOSSService creates a new OSS service
func NewOSSService(client *oss.Client) *OSSService {
	return &OSSService{
		client: client,
	}
}

// NewOSSServiceWithCredentials creates a new OSS service with credentials for cross-region access
func NewOSSServiceWithCredentials(client *oss.Client, accessKeyID, accessKeySecret, defaultEndpoint string) *OSSService {
	return &OSSService{
		client:          client,
		accessKeyID:     accessKeyID,
		accessKeySecret: accessKeySecret,
		defaultEndpoint: defaultEndpoint,
	}
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

// getClientForBucket creates an OSS client for the specific bucket's region
func (s *OSSService) getClientForBucket(bucketName string) (*oss.Client, error) {
	// First try with the default client
	bucket, err := s.client.Bucket(bucketName)
	if err != nil {
		return nil, fmt.Errorf("getting bucket %s with default client: %w", bucketName, err)
	}

	// Try to list objects to check if we can access the bucket
	_, err = bucket.ListObjects(oss.MaxKeys(1))
	if err != nil {
		// If we get an access denied error, try to determine the correct endpoint
		if strings.Contains(err.Error(), "AccessDenied") && strings.Contains(err.Error(), "endpoint") {
			// Extract region from bucket location if available
			// For buckets in different regions, we need to create a new client
			if s.accessKeyID != "" && s.accessKeySecret != "" {
				// Try common endpoints based on bucket name patterns or error messages
				endpoints := s.guessEndpointsFromError(err.Error(), bucketName)
				for _, endpoint := range endpoints {
					newClient, clientErr := oss.New(endpoint, s.accessKeyID, s.accessKeySecret)
					if clientErr != nil {
						continue
					}

					// Test if this endpoint works
					testBucket, bucketErr := newClient.Bucket(bucketName)
					if bucketErr != nil {
						continue
					}

					_, listErr := testBucket.ListObjects(oss.MaxKeys(1))
					if listErr == nil {
						return newClient, nil
					}
				}
			}
		}
		return nil, fmt.Errorf("accessing bucket %s: %w", bucketName, err)
	}

	return s.client, nil
}

// guessEndpointsFromError tries to extract the correct endpoint from error messages
func (s *OSSService) guessEndpointsFromError(errorMsg, bucketName string) []string {
	var endpoints []string

	// Extract endpoint from error message if present
	if strings.Contains(errorMsg, "oss-cn-beijing") {
		endpoints = append(endpoints, "oss-cn-beijing.aliyuncs.com")
	}
	if strings.Contains(errorMsg, "oss-cn-shanghai") {
		endpoints = append(endpoints, "oss-cn-shanghai.aliyuncs.com")
	}
	if strings.Contains(errorMsg, "oss-cn-hangzhou") {
		endpoints = append(endpoints, "oss-cn-hangzhou.aliyuncs.com")
	}
	if strings.Contains(errorMsg, "oss-cn-shenzhen") {
		endpoints = append(endpoints, "oss-cn-shenzhen.aliyuncs.com")
	}

	// If no specific endpoint found in error, try common ones
	if len(endpoints) == 0 {
		endpoints = []string{
			"oss-cn-beijing.aliyuncs.com",
			"oss-cn-shanghai.aliyuncs.com",
			"oss-cn-hangzhou.aliyuncs.com",
			"oss-cn-shenzhen.aliyuncs.com",
			"oss-us-west-1.aliyuncs.com",
			"oss-ap-southeast-1.aliyuncs.com",
		}
	}

	return endpoints
}

// FetchObjects retrieves objects from a specific bucket using pagination
func (s *OSSService) FetchObjects(bucketName string) ([]oss.ObjectProperties, error) {
	// Get the appropriate client for this bucket
	client, err := s.getClientForBucket(bucketName)
	if err != nil {
		return nil, err
	}

	bucket, err := client.Bucket(bucketName)
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
