package cloudflare

import (
	"encoding/json"
	"fmt"

	"github.com/dotcreep/go-automate-deploy/internal/utils"
)

func (c *Cloudflare) GetTunnelConnection() (*ResponseArray, error) {
	ctx, cancel := utils.Cfgx{}.ShortTimeout()
	defer cancel()
	connect, err := c.NewCloudflare(c.Key, c.Email)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/accounts/%s/cfd_tunnel", connect.BaseURL, connect.AccountID)
	resp, err := c.GetCloudflare(ctx, url)
	if err != nil {
		return nil, err
	}
	var tunnel ResponseArray
	err = json.NewDecoder(resp.Body).Decode(&tunnel)
	if err != nil {
		return nil, err
	}
	return &tunnel, nil
}
