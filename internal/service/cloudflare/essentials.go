package cloudflare

import (
	"context"
	"errors"
	"fmt"

	"github.com/dotcreep/go-automate-deploy/internal/utils"
)

type ResponseArray struct {
	Success    bool
	Errors     []Errorscf     `json:"errors,omitempty"`
	Messages   []Messages     `json:"messages,omitempty"`
	Result     []ResultTunnel `json:"result,omitempty"`
	ResultInfo ResultInfo     `json:"result_info,omitempty"`
}

type ResponseObject struct {
	Success  bool
	Errors   []Errorscf   `json:"errors,omitempty"`
	Messages []Messages   `json:"messages,omitempty"`
	Result   ResultTunnel `json:"result,omitempty"`
}

type Errorscf struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type Messages struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type ResultTunnel struct {
	AccountID          string       `json:"account_id,omitempty"`
	Config             Config       `json:"config,omitempty"`
	Connections        []Connection `json:"connections,omitempty"`
	ConnsActiveAt      string       `json:"conns_active_at,omitempty"`
	ConnsInactiveAt    string       `json:"conns_inactive_at,omitempty"`
	CreateAt           string       `json:"created_at,omitempty"`
	DeleteAt           string       `json:"deleted_at,omitempty"`
	ID                 string       `json:"id,omitempty"`
	Metadata           interface{}  `json:"metadata,omitempty"`
	Name               string       `json:"name,omitempty"`
	Nameserver         []string     `json:"name_servers,omitempty"`
	OriginalDNSHost    string       `json:"original_dnshost,omitempty"`
	OriginalNameserver []string     `json:"original_name_servers,omitempty"`
	RemoteConfig       bool         `json:"remote_config,omitempty"`
	Source             string       `json:"source,omitempty"`
	Status             string       `json:"status,omitempty"`
	Types              string       `json:"type,omitempty"`
	TunnelID           string       `json:"tunnel_id,omitempty"`
	Version            int          `json:"version,omitempty"`
}

type Config struct {
	Ingress       []Ingress     `json:"ingress,omitempty"`
	OriginRequest OriginRequest `json:"originRequest,omitempty"`
	WrapRouting   WrapRouting   `json:"wrap-routing,omitempty"`
}

type Ingress struct {
	Hostname      string        `json:"hostname,omitempty"`
	OriginRequest OriginRequest `json:"originRequest,omitempty"`
	Path          string        `json:"path,omitempty"`
	Service       string        `json:"service,omitempty"`
}

type OriginRequest struct {
	Access                 AccessCloudflare `json:"access,omitempty"`
	CAPool                 string           `json:"caPool,omitempty"`
	ConnectTimeout         int              `json:"connectTimeout,omitempty"`
	DisableChunkedEncoding bool             `json:"disableChunkedEncoding,omitempty"`
	HTTP2Origin            bool             `json:"http2Origin,omitempty"`
	HTTPHostHeader         string           `json:"httpHostHeader,omitempty"`
	KeepAliveConnections   int              `json:"keepAliveConnections,omitempty"`
	KeepAliveTimeout       int              `json:"keepAliveTimeout,omitempty"`
	NoHappyEyeballs        bool             `json:"noHappyEyeballs,omitempty"`
	NoTLSVerify            bool             `json:"noTLSVerify,omitempty"`
	OriginServerName       string           `json:"originServerName,omitempty"`
	ProxyType              string           `json:"proxyType,omitempty"`
	TCPKeepAlive           int              `json:"tcpKeepAlive,omitempty"`
	TLSTimeout             int              `json:"tlsTimeout,omitempty"`
}

type AccessCloudflare struct {
	AuthTag  []string `json:"authTag,omitempty"`
	Required bool     `json:"required,omitempty"`
	TeamName string   `json:"teamName,omitempty"`
}

type WrapRouting struct {
	Enabled bool `json:"enabled,omitempty"`
}

type Connection struct {
	ClientID           string `json:"client_id,omitempty"`
	ClientVersion      string `json:"client_version,omitempty"`
	ColoName           string `json:"colo_name,omitempty"`
	ID                 string `json:"id,omitempty"`
	IsPendingReconnect bool   `json:"is_pending_reconnect,omitempty"`
	OpenedAt           string `json:"opened_at,omitempty"`
	OriginIP           string `json:"origin_ip,omitempty"`
	UUID               string `json:"uuid,omitempty"`
}

type ResultInfo struct {
	Count      int `json:"count,omitempty"`
	Page       int `json:"page,omitempty"`
	PerPage    int `json:"per_page,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
}

type domainZone struct {
	ID          string      `json:"id,omitempty"`
	Name        string      `json:"name,omitempty"`
	NameServers interface{} `json:"name_servers,omitempty"`
	Status      string      `json:"status,omitempty"`
	Plan        struct {
		Name  string  `json:"name,omitempty"`
		Price float64 `json:"price,omitempty"`
	} `json:"plan,omitempty"`
}

// Response DNS Records
type ResultDNSRecord struct {
	ID         string            `json:"id"`
	ZoneID     string            `json:"zone_id"`
	ZoneName   string            `json:"zone_name"`
	Name       string            `json:"name"`
	Type       string            `json:"type"`
	Content    string            `json:"content"`
	Proxiable  bool              `json:"proxiable"`
	Proxied    bool              `json:"proxied"`
	TTL        int               `json:"ttl"`
	Settings   SettingsDNSRecord `json:"settings"`
	Meta       Meta              `json:"meta"`
	Comment    *string           `json:"comment"`
	Tags       []string          `json:"tags"`
	CreatedOn  string            `json:"created_on"`
	ModifiedOn string            `json:"modified_on"`
}

type SettingsDNSRecord struct{}

// Meta struct to represent additional metadata
type Meta struct {
	AutoAdded           bool `json:"auto_added"`
	ManagedByApps       bool `json:"managed_by_apps"`
	ManagedByArgoTunnel bool `json:"managed_by_argo_tunnel"`
}

// Response struct to represent the overall JSON structure
type ResponseDNSRrecords struct {
	Result []ResultDNSRecord `json:"result"`
}

// EOF DNS RECORDS

type CloudflareAPI interface {
	// Always follow this interface for access cloudflare first
	CloudflareInitial
	CloudflareMethod
	CloudflareOperator
	CloudflareController
}

type CloudflareInitial interface {
	NewCloudflare(key, email string) *Cloudflare
	CloudflareConfig() (*utils.YamlStruct, error)
}

type CloudflareOperator interface {
	GetZone(domain string) (*domainZone, error)
	GetZoneByID(zoneID string) (*domainZone, error)
	GetDNSRecord(zoneID string) ([]string, error)
	GetTunnelConnection() (*ResponseArray, error)
	GetTunnelConfiguration(tunnelID string) (*ResponseObject, error)
	GetZoneDetail(zoneID string) (*ResponseObject, error)
	GetDomainFromTunnelConfiguration(ctx context.Context, domain, tunnelID string) (*ResponseObject, error)
}

type CloudflareController interface {
	Add(s *Subdomains) (string, error)
	Delete(s *Subdomains) (string, error)
	Register(s *Subdomains) (string, error)
	StatusDomain(s *Subdomains) (string, error)
}

type ConfigData struct{}
type Cloudflare struct {
	Key       string
	Email     string
	TunnelID  string
	BaseURL   string
	AccountID string
	Headers   map[string]string
}
type App struct{}

type Subdomains struct {
	Service      string
	Domain       string
	TunnelID     string
	Path         string
	LoadBalancer bool
}

// Initial access

func (c *Cloudflare) NewCloudflare(key, email string) (*Cloudflare, error) {
	if key == "" {
		return nil, errors.New("key is required")
	}
	if email == "" {
		return nil, errors.New("email is required")
	}
	config, err := utils.YmlConfig()
	if err != nil {
		return nil, errors.New("failed to read config.yml")
	}
	return &Cloudflare{
		Key:       key,
		Email:     email,
		TunnelID:  config.Cloudflare.TunnelID,
		AccountID: config.Cloudflare.AccountID,
		BaseURL:   config.Cloudflare.BaseURL,
		Headers: map[string]string{
			"X-Auth-Email":  email,
			"Authorization": fmt.Sprintf("Bearer %s", key),
			"Content-Type":  "application/json",
			"User-Agent":    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36",
		},
	}, nil
}
