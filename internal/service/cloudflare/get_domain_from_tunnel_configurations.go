package cloudflare

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

func (c *Cloudflare) GetDomainFromTunnelConfiguration(ctx context.Context, domain, tunnelID string) (*ResponseObject, error) {
	if domain == "" {
		return nil, errors.New("domain is empty")
	}
	if tunnelID == "" {
		return nil, errors.New("tunnelID is empty")
	}
	url := fmt.Sprintf("%s/accounts/%s/cfd_tunnel/%s/configurations", c.BaseURL, c.AccountID, tunnelID)
	resp, err := c.GetCloudflare(ctx, url)
	if err != nil {
		return nil, err
	}

	var tunnel ResponseObject
	err = json.NewDecoder(resp.Body).Decode(&tunnel)
	if err != nil {
		return nil, err
	}
	return &tunnel, nil
}
