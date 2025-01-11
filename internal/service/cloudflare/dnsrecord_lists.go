package cloudflare

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dotcreep/go-automate-deploy/internal/utils"
)

func (c *Cloudflare) RecordList(zoneId string) (*ResponseArray, error) {
	ctx, cancel := utils.Cfgx{}.DefaultTimeout()
	defer cancel()

	if zoneId == "" {
		return nil, errors.New("zoneID is empty")
	}

	url := fmt.Sprintf("%s/zones/%s/dns_records", c.BaseURL, zoneId)

	resp, err := c.GetCloudflare(ctx, url)
	if err != nil {
		return nil, err
	}

	var zone ResponseArray

	err = json.NewDecoder(resp.Body).Decode(&zone)
	if err != nil {
		return nil, err
	}
	return &zone, nil
}
