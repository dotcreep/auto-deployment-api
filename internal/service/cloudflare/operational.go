package cloudflare

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/dotcreep/go-automate-deploy/internal/utils"
)

// Add domain on tunnel configuration
func (c *Cloudflare) AddDomainToTunnelConfiguration(ctx context.Context, s *Subdomains) (string, error) {
	// 1. Check input
	if s.TunnelID == "" {
		return "", errors.New("tunnelID is empty")
	}
	if s.Service == "" {
		return "", errors.New("service is empty")
	}
	if s.Domain == "" {
		return "", errors.New("domain is empty")
	}
	if s.Path == "" {
		s.Path = ""
	}
	if s.LoadBalancer && !s.LoadBalancer {
		return "", errors.New("loadBalancer is required")
	}

	// 2. Create new Connection
	connect, err := c.NewCloudflare(c.Key, c.Email)
	if err != nil {
		return "", err
	}

	// Get current tunnel configuration
	tunnelConfig, err := connect.GetCurrentTunnelConfiguration(ctx, s.TunnelID)
	if err != nil {
		return "", err
	}
	currentIngress := tunnelConfig.Result

	// Check subdomain is valid
	d, err := url.Parse(s.Domain)
	if err != nil {
		return "", err
	}
	subDomains := strings.Split(d.String(), ".")
	var baseDomain string
	if len(subDomains) < 2 {
		return "", errors.New("invalid domain")
	}

	if len(subDomains) > 2 {
		if len(subDomains[len(subDomains)-1]) == 2 && len(subDomains[len(subDomains)-2]) == 2 {
			baseDomain = subDomains[len(subDomains)-3] + "." + subDomains[len(subDomains)-2] + "." + subDomains[len(subDomains)-1]
		} else {
			baseDomain = subDomains[len(subDomains)-2] + "." + subDomains[len(subDomains)-1]
		}
	} else {
		baseDomain = subDomains[0] + "." + subDomains[1]
	}

	// Get ZoneID from Domain to checking up
	zoneStruct, err := connect.GetZone(ctx, baseDomain)
	if err != nil {
		return "", err
	}

	// Get DNS Record
	dnsRecord, err := c.GetDNSRecord(ctx, zoneStruct.ID)
	if err != nil {
		return "", err
	}

	// Check subdomain is exist in DNSRecord
	for _, v := range dnsRecord {
		if v == s.Domain && !s.LoadBalancer {
			return "", errors.New("subdomain is exist in DNS Record")
		}
	}

	// Check subdomain is exist in Ingress
	for _, ingress := range currentIngress.Config.Ingress {
		if ingress.Hostname == s.Domain && !s.LoadBalancer {
			return "", fmt.Errorf("subdomain is existed in tunnel configuration id: %s", s.TunnelID)
		}
	}

	// Check scheme domain
	var secureScheme bool
	if strings.Contains(s.Service, "https") {
		secureScheme = true
	} else {
		secureScheme = false
	}
	// if d.Scheme == "https" {
	// secureScheme = true
	// } else {
	// secureScheme = false
	// }

	// New Ingress
	newIngress := Ingress{
		Service:  s.Service,
		Hostname: s.Domain,
		Path:     s.Path,
		OriginRequest: OriginRequest{
			NoTLSVerify: secureScheme,
			HTTP2Origin: secureScheme,
		},
	}

	// Append to Ingress before http_status:404
	for i := 1; i < len(currentIngress.Config.Ingress); i++ {
		if currentIngress.Config.Ingress[i].Service == "http_status:404" {
			currentIngress.Config.Ingress = append(
				currentIngress.Config.Ingress[:i],
				append([]Ingress{newIngress},
					currentIngress.Config.Ingress[i:]...)...,
			)
			break
		}
	}

	// Update Tunnel Configurations
	url := fmt.Sprintf("%s/accounts/%s/cfd_tunnel/%s/configurations", connect.BaseURL, connect.AccountID, s.TunnelID)
	result := ResultTunnel{
		Config: currentIngress.Config,
	}
	jsonData, err := json.Marshal(result)
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

	return string(b), nil
}

func (c *Cloudflare) Delete(ctx context.Context, s *Subdomains) (string, error) {
	if s.Domain == "" {
		return "", errors.New("domain is required")
	}

	if s.TunnelID == "" {
		return "", errors.New("tunnelID is required")
	}

	// Connect to Cloudflare
	connect, err := c.NewCloudflare(c.Key, c.Email)
	if err != nil {
		return "", err
	}

	// Get Tunnel Configuration
	currentIngress, err := c.GetCurrentTunnelConfiguration(ctx, s.TunnelID)
	if err != nil {
		return "", err
	}

	// Get Ingress
	result := currentIngress.Result
	data := result.Config

	// Remove if exitst
	var isExist bool
	for i, ingress := range data.Ingress {
		if ingress.Hostname == s.Domain {
			data.Ingress = append(data.Ingress[:i], data.Ingress[i+1:]...)
			isExist = true
			break
		}
	}

	if !isExist {
		return "Subdomain is not exist", nil
	}
	// Modify
	url := fmt.Sprintf("%s/accounts/%s/cfd_tunnel/%s/configurations", connect.BaseURL, connect.AccountID, s.TunnelID)
	resultData := ResultTunnel{
		Config: data,
	}
	jsonData, err := json.Marshal(resultData)
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
	resp.Body = io.NopCloser(bytes.NewReader(b))
	return string(b), nil
}

func (c *Cloudflare) Register(ctx context.Context, s *Subdomains) (string, error) {
	if s.Domain == "" {
		return "", errors.New("domain is required")
	}

	if c.AccountID == "" {
		return "", errors.New("accountID is required")
	}

	// Connect to Cloudflare

	connect, err := c.NewCloudflare(c.Key, c.Email)
	if err != nil {
		return "", err
	}

	// Check Domain if already registered
	_, err = c.GetZone(ctx, s.Domain)
	if err == nil {
		return "Domain is already registered", nil
	}

	data := struct {
		Account struct {
			ID string `json:"id"`
		} `json:"account"`
		Name string `json:"name"`
		Type string `json:"type"`
	}{
		Account: struct {
			ID string `json:"id"`
		}{
			ID: c.AccountID,
		},
		Name: s.Domain,
		Type: "full",
	}

	url := fmt.Sprintf("%s/zones", connect.BaseURL)
	jsonByte, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	resp, err := c.PostCloudflare(ctx, url, bytes.NewReader(jsonByte))
	if err != nil {
		return "", err
	}
	var tunnel ResponseObject
	err = json.NewDecoder(resp.Body).Decode(&tunnel)
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
	resp.Body = io.NopCloser(bytes.NewReader(b))
	return string(b), nil
}

func (c *Cloudflare) StatusDomain(ctx context.Context, s *Subdomains) (string, error) {
	if s.Domain == "" {
		return "", errors.New("domain is required")
	}

	// Get ZoneID by Domain name
	zoneStruct, err := c.GetZone(ctx, s.Domain)
	if err != nil {
		return "", err
	}

	domain, err := c.GetZoneDetail(zoneStruct.ID)
	if err != nil {
		return "", err
	}

	return domain.Result.Status, nil
}

// Add DNS record like CNAME with domain example.com for www.example.com
func (c *Cloudflare) AddDNSRecord(ctx context.Context, s *Subdomains) (string, error) {
	if s.Domain == "" {
		return "", errors.New("domain is required")
	}

	baseDomain := utils.GetBaseDomain(s.Domain)

	zone, err := c.GetZone(ctx, baseDomain)
	if err != nil {
		return "", err
	}

	// Connect to Cloudflare
	url := fmt.Sprintf("%s/zones/%s/dns_records", c.BaseURL, zone.ID)

	data := struct {
		Comment string `json:"comment"`
		Content string `json:"content"`
		Name    string `json:"name"`
		Proxied bool   `json:"proxied"`
		TTL     int    `json:"ttl"`
		Type    string `json:"type"`
	}{
		Comment: fmt.Sprintf("DNS for %s", s.Domain),
		Content: fmt.Sprintf("%s.cfargotunnel.com", s.TunnelID),
		Name:    s.Domain,
		Proxied: true,
		TTL:     3600,
		Type:    "CNAME",
	}

	dataByte, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	resp, err := c.PostCloudflare(ctx, url, bytes.NewBuffer(dataByte))
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
	resp.Body = io.NopCloser(bytes.NewReader(b))
	return string(b), nil
}
