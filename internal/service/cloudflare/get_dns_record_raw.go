package cloudflare

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

func (c *Cloudflare) GetRawDNSRecord(ctx context.Context, zoneID string) (*ResponseDNSRrecords, error) {
	if zoneID == "" {
		return nil, errors.New("zone id is required")
	}
	connect, err := c.NewCloudflare(c.Key, c.Email)
	if err != nil {
		return nil, err
	}
	// Get DNS Record Lists
	url := fmt.Sprintf("%s/zones/%s/dns_records", connect.BaseURL, zoneID)
	resp, err := c.GetCloudflare(ctx, url)
	if err != nil {
		return nil, err
	}
	var records ResponseDNSRrecords
	err = json.NewDecoder(resp.Body).Decode(&records)
	if err != nil {
		return nil, err
	}

	return &records, nil
}
