package service

type MonitoringAPI interface {
	MonitoringService
}

type MonitoringService interface {
	CloudflareMonitoring() (string, error)
	PortainerMonitoring() (string, error)
	JenkinsMonitoring() (string, error)
}

type Monitoring struct {
	Cloudflare string
	Portainer  string
	Jenkins    string
}

func NewMonitoringService() *Monitoring {
	return &Monitoring{
		Cloudflare: "",
		Portainer:  "",
		Jenkins:    "",
	}
}
