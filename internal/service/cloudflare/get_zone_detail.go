package cloudflare

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dotcreep/go-automate-deploy/internal/utils"
)

func (c *Cloudflare) GetZoneDetail(zoneID string) (*ResponseObject, error) {
	ctx, cancel := utils.Cfgx{}.DefaultTimeout()
	defer cancel()

	if zoneID == "" {
		return nil, errors.New("zone id is empty")
	}
	connect, err := c.NewCloudflare(c.Key, c.Email)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/zones/%s", connect.BaseURL, zoneID)
	resp, err := c.GetCloudflare(ctx, url)
	if err != nil {
		return nil, err
	}
	var zone ResponseObject
	err = json.NewDecoder(resp.Body).Decode(&zone)
	if err != nil {
		return nil, err
	}
	return &zone, nil
}
