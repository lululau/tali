package service

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
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

// FetchDomains retrieves all DNS domains using pagination
func (s *DNSService) FetchDomains() ([]alidns.DomainInDescribeDomains, error) {
	var allDomains []alidns.DomainInDescribeDomains
	pageNumber := int64(1) // SDK uses int64 for PageNumber in response, so keep consistent
	pageSize := int64(100)

	for {
		request := alidns.CreateDescribeDomainsRequest()
		request.Scheme = "https"
		request.PageNumber = requests.NewInteger(int(pageNumber)) // CreateDescribeDomainsRequest uses requests.Integer
		request.PageSize = requests.NewInteger(int(pageSize))

		response, err := s.client.DescribeDomains(request)
		if err != nil {
			return nil, fmt.Errorf("describing DNS domains (page %d): %w", pageNumber, err)
		}

		allDomains = append(allDomains, response.Domains.Domain...)

		// TotalCount is int64, PageNumber in response is int64, PageSize in response is int64
		if pageNumber*pageSize >= response.TotalCount {
			break
		}

		// Another way to check, if current page has less than pageSize items (safer if TotalCount is not always accurate or if page size is not respected)
		if len(response.Domains.Domain) < int(pageSize) {
			break
		}

		pageNumber++
	}
	return allDomains, nil
}

// FetchDomainRecords retrieves DNS records for a specific domain using pagination
func (s *DNSService) FetchDomainRecords(domainName string) ([]alidns.Record, error) {
	var allRecords []alidns.Record
	pageNumber := int64(1)
	pageSize := int64(100)

	for {
		request := alidns.CreateDescribeDomainRecordsRequest()
		request.Scheme = "https"
		request.DomainName = domainName
		request.PageNumber = requests.NewInteger(int(pageNumber))
		request.PageSize = requests.NewInteger(int(pageSize))

		response, err := s.client.DescribeDomainRecords(request)
		if err != nil {
			return nil, fmt.Errorf("describing DNS domain records for %s (page %d): %w", domainName, pageNumber, err)
		}

		allRecords = append(allRecords, response.DomainRecords.Record...)

		if pageNumber*pageSize >= response.TotalCount {
			break
		}

		if len(response.DomainRecords.Record) < int(pageSize) {
			break
		}

		pageNumber++
	}
	return allRecords, nil
}
