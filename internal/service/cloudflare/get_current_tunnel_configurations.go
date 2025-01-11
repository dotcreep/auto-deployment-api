package cloudflare

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

func (c *Cloudflare) GetCurrentTunnelConfiguration(ctx context.Context, tunnelID string) (*ResponseObject, error) {
	if tunnelID == "" {
		return nil, errors.New("tunnelID is empty")
	}

	connect, err := c.NewCloudflare(c.Key, c.Email)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/accounts/%s/cfd_tunnel/%s/configurations", connect.BaseURL, connect.AccountID, tunnelID)
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
