package cloudflare

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

func (c *Cloudflare) GetZone(ctx context.Context, domain string) (*domainZone, error) {
	connect, err := c.NewCloudflare(c.Key, c.Email)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/zones", connect.BaseURL)
	resp, err := c.GetCloudflare(ctx, url)
	if err != nil {
		return nil, err
	}
	var zone map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&zone)
	if err != nil {
		return nil, err
	}

	var domainZones domainZone
	for _, v := range zone["result"].([]interface{}) {
		if v.(map[string]interface{})["name"].(string) == domain {
			domainZones = domainZone{
				ID:          v.(map[string]interface{})["id"].(string),
				Name:        v.(map[string]interface{})["name"].(string),
				NameServers: v.(map[string]interface{})["name_servers"].([]interface{}),
				Status:      v.(map[string]interface{})["status"].(string),
				Plan: struct {
					Name  string  `json:"name,omitempty"`
					Price float64 `json:"price,omitempty"`
				}{
					Name:  v.(map[string]interface{})["plan"].(map[string]interface{})["name"].(string),
					Price: v.(map[string]interface{})["plan"].(map[string]interface{})["price"].(float64),
				},
			}
			return &domainZones, nil
		}
	}
	return nil, errors.New("domain not found")
}
