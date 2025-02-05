package cloudflare

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/dotcreep/go-automate-deploy/internal/utils"
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
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var zone resultDomainZone
	if err := json.Unmarshal(bodyBytes, &zone); err != nil {
		return nil, err
	}
	baseDomain := utils.GetBaseDomain(domain)
	if baseDomain == "" {
		return nil, errors.New("domain is invalid")
	}
	var domainZones domainZone
	for _, v := range zone.Result {
		if v.Name == baseDomain {
			domainZones = domainZone{
				ID:          v.ID,
				Name:        v.Name,
				NameServers: v.NameServers,
				Status:      v.Status,
				Plan: struct {
					Name  string  `json:"name"`
					Price float64 `json:"price"`
				}{
					Name:  v.Plan.Name,
					Price: v.Plan.Price,
				},
			}
			return &domainZones, nil
		}
	}
	return nil, errors.New("domain is not registered")
}
