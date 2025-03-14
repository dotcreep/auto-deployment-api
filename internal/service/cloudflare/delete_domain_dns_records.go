package cloudflare

import (
	"context"
	"errors"

	"github.com/dotcreep/go-automate-deploy/internal/utils"
)

func (c *Cloudflare) DeleteDomainDNSRecords(ctx context.Context, s *Subdomains) (string, error) {
	// 1. Check Input
	if s.Domain == "" {
		return "", errors.New("domain is required")
	}

	// 2. Cerate new connection
	connect, err := c.NewCloudflare(c.Key, c.Email)
	if err != nil {
		return "", err
	}

	// 3. Get Zone ID
	baseDom := utils.GetBaseDomain(s.Domain)
	if baseDom == "" {
		return "", errors.New("domain is invalid")
	}
	zone, err := connect.GetZone(ctx, baseDom)
	if err != nil {
		return "", err
	}

	// 4. Get DNS Records list
	DNSRecords, err := connect.GetRawDNSRecord(ctx, zone.ID)
	if err != nil {
		return "", err
	}

	// 5. Delete DNS Records
	// 5.1. Get ID DNS Record from Domain
	var dnsId string
	for _, v := range DNSRecords.Result {
		if v.Name == s.Domain {
			dnsId = v.ID
			break
		}
	}

	// 5.2. Delete DNS Record
	err = connect.DeleteDNSRecord(ctx, zone.ID, dnsId)
	if err != nil {
		return "", err
	}
	return "success", nil
}
