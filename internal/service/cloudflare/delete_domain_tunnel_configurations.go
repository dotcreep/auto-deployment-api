package cloudflare

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func (c *Cloudflare) DeleteDomainFromTunnelConfiguration(ctx context.Context, s *Subdomains) (string, error) {
	// 1. Check input
	if s.TunnelID == "" {
		return "", errors.New("tunnelID is empty")
	}
	if s.Domain == "" {
		return "", errors.New("domain is empty")
	}

	// 2. Create new Connection
	connect, err := c.NewCloudflare(c.Key, c.Email)
	if err != nil {
		return "", err
	}

	// 3. Get current tunnel configuration
	tunnelConfig, err := connect.GetCurrentTunnelConfiguration(ctx, s.TunnelID)
	if err != nil {
		return "", err
	}
	newResponse := &ResponseObject{
		Result: ResultTunnel{
			Config: Config{
				Ingress: []Ingress{},
			},
		},
	}
	var isExist bool = false
	for _, ingress := range tunnelConfig.Result.Config.Ingress {
		if ingress.Hostname != s.Domain {
			newResponse.Result.Config.Ingress = append(newResponse.Result.Config.Ingress, ingress)
		} else {
			isExist = true
		}
	}

	// 4. Update tunnel configuration
	url := fmt.Sprintf("%s/accounts/%s/cfd_tunnel/%s/configurations", connect.BaseURL, connect.AccountID, s.TunnelID)
	jsonData, err := json.Marshal(newResponse.Result)
	if err != nil {
		return "", err
	}
	resp, err := c.PutCloudflare(ctx, url, bytes.NewReader(jsonData))
	if err != nil {
		return "", err
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body = io.NopCloser(bytes.NewReader(b))
		return "", fmt.Errorf("unexpected status code: %v\ndata: %s", resp.StatusCode, string(b))
	}
	if !isExist {
		return "", fmt.Errorf("domain %s is already not exist in tunnel %s", s.Domain, s.TunnelID)
	}
	return fmt.Sprintf("successfully delete domain %s from tunnel %s", s.Domain, s.TunnelID), nil
}
