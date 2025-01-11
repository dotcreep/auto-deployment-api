package cloudflare

import (
	"encoding/json"
	"fmt"

	"github.com/dotcreep/go-automate-deploy/internal/utils"
)

func (c *Cloudflare) GetZoneByID(zoneID string) (*domainZone, error) {
	ctx, cancel := utils.Cfgx{}.ShortTimeout()
	defer cancel()
	connect, err := c.NewCloudflare(c.Key, c.Email)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/zones/%s", connect.BaseURL, zoneID)
	resp, err := c.GetCloudflare(ctx, url)
	if err != nil {
		return nil, err
	}
	var zone domainZone
	err = json.NewDecoder(resp.Body).Decode(&zone)
	if err != nil {
		return nil, err
	}
	return &zone, nil
}
