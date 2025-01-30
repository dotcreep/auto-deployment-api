package utils

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type YamlStruct struct {
	Cloudflare cloudflareConfig `yaml:"cloudflare"`
	Portainer  PortainerConfig  `yaml:"portainer"`
	Jenkins    jenkinsConfig    `yaml:"jenkins"`
	Config     commonConfig     `yaml:"config"`
}

type cloudflareConfig struct {
	AccountID string `yaml:"account_id"`
	TunnelID  string `yaml:"tunnel_id"`
	BaseURL   string `yaml:"base_url"`
	Key       string `yaml:"key"`
	Email     string `yaml:"email"`
}

type PortainerConfig struct {
	BaseURL         string `yaml:"base_url"`
	APIKey          string `yaml:"api_key"`
	Username        string `yaml:"username"`
	SwarmID         string `yaml:"swarm_id"`
	Type            string `yaml:"type"`
	Method          string `yaml:"method"`
	FromAppTemplate bool   `yaml:"fromAppTemplate"`
	EndpointId      int    `yaml:"endpointId"`
	MaxClientWeb    int    `yaml:"max_client_web"`
	MaxClientDB     int    `yaml:"max_client_db"`
}

type commonConfig struct {
	ExcludeDomain   []string `yaml:"exclude_domain"`
	ImageAPI        string   `yaml:"image_api"`
	ImageWeb        string   `yaml:"image_web"`
	RegisAdmin      string   `yaml:"regis_admin"`
	RegisMerch      string   `yaml:"regis_merch"`
	TokenRegis      string   `yaml:"token_regis"`
	XToken          string   `yaml:"x_token"`
	XTokenX         string   `yaml:"x_api_key"`
	DomainAPI       string   `yaml:"domain_api"`
	PORT            string   `yaml:"port_api"`
	PathClient      string   `yaml:"path_client"`
	PathEnvironment string   `yaml:"path_environment"`
	ChatAPI         string   `yaml:"chat_api"`
}

type jenkinsConfig struct {
	BaseURL           string `yaml:"base_url"`
	Username          string `yaml:"username"`
	APIKey            string `yaml:"api_key"`
	DomainCredentials string `yaml:"domain_credentials"`
	APIURL            string `yaml:"api_url"`
	PrefixCredentials string `yaml:"prefix_credentials"`
}

type Cfgx struct{}

func (c Cfgx) ShortTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

func (c Cfgx) DefaultTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 30*time.Second)
}

func (c Cfgx) LongTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 60*time.Second)
}

func Open() (*YamlStruct, error) {
	var pathFile string
	if _, err := os.Stat("config.yaml"); err == nil {
		pathFile = "config.yaml"
	} else if _, err := os.Stat("config.yml"); err == nil {
		pathFile = "config.yml"
	} else {
		return nil, errors.New("config yaml file not found")
	}

	data, err := os.ReadFile(pathFile)
	if err != nil {
		return nil, err
	}
	var config YamlStruct
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func GetStatusOK(ctx context.Context, domain string) (int, error) {
	if domain == "" {
		return http.StatusBadRequest, errors.New("domain is required")
	}
	var newDomain string
	if !strings.Contains(domain, "http") {
		newDomain = fmt.Sprintf("https://%s", domain)
	} else {
		newDomain = domain
	}
	req, err := http.NewRequestWithContext(ctx, "GET", newDomain, nil)
	if err != nil {
		return http.StatusBadRequest, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}

func YmlConfig() (*YamlStruct, error) {
	config, err := Open()
	if config == nil {
		return nil, err
	}
	return config, nil
}
