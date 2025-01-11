package cloudflare

import (
	"context"
	"errors"
)

func (c *Cloudflare) GetNameserver(ctx context.Context, s *Subdomains) ([]string, error) {
	if s.Domain == "" {
		return nil, errors.New("domain is required")
	}

	zone, err := c.GetZone(ctx, s.Domain)
	if err != nil {
		return nil, err
	}
	domain, err := c.GetZoneDetail(zone.ID)
	if err != nil {
		return nil, err
	}
	return domain.Result.Nameserver, nil
}
