package cloudflare

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

func (c *Cloudflare) GetDNSRecord(ctx context.Context, zoneID string) ([]string, error) {
	if zoneID == "" {
		return nil, errors.New("zone id is required")
	}
	connect, err := c.NewCloudflare(c.Key, c.Email)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/zones/%s/dns_records", connect.BaseURL, zoneID)
	resp, err := c.GetCloudflare(ctx, url)
	if err != nil {
		return nil, err
	}
	var records map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&records)
	if err != nil {
		return nil, err
	}

	record := []string{}
	for _, v := range records["result"].([]interface{}) {
		record = append(record, v.(map[string]interface{})["name"].(string))
	}
	return record, nil
}
