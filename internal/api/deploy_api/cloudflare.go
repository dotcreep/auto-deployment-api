package deploy_api

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/dotcreep/go-automate-deploy/internal/service/cloudflare"
	"github.com/dotcreep/go-automate-deploy/internal/utils"
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
	cfg, err := utils.YmlConfig()
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
	var basedom string
	parts := strings.Split(data.Domain, ".")
	if len(parts) > 2 {
		if parts[len(parts)-2] == "co" || parts[len(parts)-2] == "biz" || parts[len(parts)-2] == "com" {
			basedom = strings.Join(parts[len(parts)-3:], ".")
		} else {
			basedom = strings.Join(parts[len(parts)-2:], ".")
		}
	} else {
		basedom = strings.Join(parts[len(parts)-2:], ".")
	}

	// For new Domain
	var excludedDomain bool = false
	for _, v := range cfg.Config.ExcludeDomain {
		if v == basedom {
			excludedDomain = true
		}
	}
	if !excludedDomain {
		// This is for new domain
		_, err := connect.GetZone(ctx, data.Domain)
		if err != nil {
			data.Domain = basedom
			_, errData := connect.Register(ctx, data)
			if errData != nil {
				return "", errData
			}
		}
		zone, err := connect.GetZone(ctx, basedom)
		if err != nil {
			return "", err
		}
		if zone.Status == "pending" {
			return fmt.Sprintf("success add domain %s", data.Domain), nil
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
