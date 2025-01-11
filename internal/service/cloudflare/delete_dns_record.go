package cloudflare

import (
	"context"
	"errors"
	"fmt"
)

func (c *Cloudflare) DeleteDNSRecord(ctx context.Context, zoneId, id string) error {
	if id == "" {
		return errors.New("id is required")
	}
	if zoneId == "" {
		return errors.New("zone id is required")
	}
	connect, err := c.NewCloudflare(c.Key, c.Email)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/zones/%s/dns_records/%s", connect.BaseURL, zoneId, id)
	_, err = c.DeleteCloudflare(ctx, url)
	if err != nil {
		return err
	}
	return nil
}
