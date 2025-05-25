package client

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// AliyunClients holds all Aliyun service clients
type AliyunClients struct {
	ECS    *ecs.Client
	DNS    *alidns.Client
	SLB    *slb.Client
	RDS    *rds.Client
	OSS    *oss.Client
	config *Config
}

// Config represents the client configuration
type Config struct {
	AccessKeyID     string
	AccessKeySecret string
	RegionID        string
	OssEndpoint     string
}

// NewAliyunClients creates and initializes all Aliyun service clients
func NewAliyunClients(cfg *Config) (*AliyunClients, error) {
	clients := &AliyunClients{config: cfg}

	// Initialize ECS client
	ecsClient, err := ecs.NewClientWithAccessKey(cfg.RegionID, cfg.AccessKeyID, cfg.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("creating ECS client: %w", err)
	}
	clients.ECS = ecsClient

	// Initialize DNS client
	dnsClient, err := alidns.NewClientWithAccessKey(cfg.RegionID, cfg.AccessKeyID, cfg.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("creating DNS client: %w", err)
	}
	clients.DNS = dnsClient

	// Initialize SLB client
	slbClient, err := slb.NewClientWithAccessKey(cfg.RegionID, cfg.AccessKeyID, cfg.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("creating SLB client: %w", err)
	}
	clients.SLB = slbClient

	// Initialize RDS client
	rdsClient, err := rds.NewClientWithAccessKey(cfg.RegionID, cfg.AccessKeyID, cfg.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("creating RDS client: %w", err)
	}
	clients.RDS = rdsClient

	// Initialize OSS client
	ossClient, err := oss.New(cfg.OssEndpoint, cfg.AccessKeyID, cfg.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("creating OSS client: %w", err)
	}
	clients.OSS = ossClient

	return clients, nil
}

// GetConfig returns the client configuration
func (c *AliyunClients) GetConfig() *Config {
	return c.config
}
