package service

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
)

// DNSService handles DNS operations
type DNSService struct {
	client *alidns.Client
}

// NewDNSService creates a new DNS service
func NewDNSService(client *alidns.Client) *DNSService {
	return &DNSService{client: client}
}

// FetchDomains retrieves all DNS domains
func (s *DNSService) FetchDomains() ([]alidns.DomainInDescribeDomains, error) {
	request := alidns.CreateDescribeDomainsRequest()
	request.Scheme = "https"
	response, err := s.client.DescribeDomains(request)
	if err != nil {
		return nil, fmt.Errorf("describing DNS domains: %w", err)
	}
	return response.Domains.Domain, nil
}

// FetchDomainRecords retrieves DNS records for a specific domain
func (s *DNSService) FetchDomainRecords(domainName string) ([]alidns.Record, error) {
	request := alidns.CreateDescribeDomainRecordsRequest()
	request.Scheme = "https"
	request.DomainName = domainName
	response, err := s.client.DescribeDomainRecords(request)
	if err != nil {
		return nil, fmt.Errorf("describing DNS domain records for %s: %w", domainName, err)
	}
	return response.DomainRecords.Record, nil
}
