package cloudflare

import (
	"context"
	"errors"

	"github.com/dotcreep/go-automate-deploy/internal/utils"
)

func (c *Cloudflare) GetDomainIsUsed(ctx context.Context, domain, tunnelID string) (string, error) {
	if domain == "" {
		return "", errors.New("domain is required")
	}
	if tunnelID == "" {
		return "", errors.New("tunnelID is required")
	}
	// 1. Search from Tunnel Configuration
	res, err := c.GetDomainFromTunnelConfiguration(ctx, domain, tunnelID)
	if err != nil {
		return "", err
	}
	var domainFound bool
	domainFound = false
	for _, v := range res.Result.Config.Ingress {
		if v.Hostname == domain {
			domainFound = true
		}
	}

	// 2. Search from Zone
	baseDomain := utils.GetBaseDomain(domain)
	zone, err := c.GetZone(ctx, baseDomain)
	if err != nil {
		return "", err
	}

	// 3. Search in DNS Records
	dnsrecords, err := c.GetDNSRecord(ctx, zone.ID)
	if err != nil {
		return "", err
	}

	for _, v := range dnsrecords {
		if v == domain {
			domainFound = true
		}
	}
	if domainFound {
		return "unavailable", errors.New("domain is already used")
	}

	return "available", nil
}
