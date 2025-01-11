package deploy_api

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/dotcreep/go-automate-deploy/internal/service/cloudflare"
)

func connectCloudflare() (*cloudflare.Cloudflare, error) {
	cf := cloudflare.Cloudflare{}
	secrets := &Secrets{}
	secrets.GetSecret()
	if secrets.Cloudflare.Token == "" {
		return nil, errors.New("cloudflare token is required")
	}
	connect, err := cf.NewCloudflare(secrets.Cloudflare.Token, secrets.Cloudflare.Email)
	if err != nil {
		return nil, err
	}
	return connect, nil
}

func DeployCloudflare(ctx context.Context, data *cloudflare.Subdomains) (string, error) {
	connect, err := connectCloudflare()
	if err != nil {
		return "", err
	}
	if data.Service == "" {
		return "", errors.New("service is required")
	}
	if data.Domain == "" {
		return "", errors.New("domain is required")
	}
	if data.TunnelID == "" {
		return "", errors.New("tunnel id is required")
	}
	if data.Path == "" {
		data.Path = ""
	}
	if !data.LoadBalancer {
		data.LoadBalancer = false
	}
	// Parse domain if subdomain or domain
	parts := strings.Split(data.Domain, ".")
	baseDomain := data.Domain
	var domainTypes string
	if len(parts) > 2 {
		domainTypes = "Subdomain"
		baseDomain = strings.Join(parts[len(parts)-2:], ".")
	} else {
		domainTypes = "Domain"
	}
	// For new Domain
	if domainTypes == "Domain" {
		_, err := connect.GetZone(ctx, baseDomain)
		if err != nil {
			_, errData := connect.Register(ctx, data)
			if errData != nil {
				return "", errData
			}
		}
		zone, err := connect.GetZone(ctx, baseDomain)
		if err != nil {
			return "", err
		}
		if zone.Status == "pending" {
			return "", errors.New("domain is pending")
		}
	}
	_, err = connect.AddDomainToTunnelConfiguration(ctx, data)
	if err != nil {
		return "", err
	}
	_, err = connect.AddDNSRecord(ctx, data)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("success add domain %s", data.Domain), nil
}
