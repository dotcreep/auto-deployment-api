package portainer

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/dotcreep/go-automate-deploy/internal/utils"
	"gopkg.in/yaml.v3"
)

type AutoUpdate struct {
	ForcePullImage bool   `json:"forcePullImage"`
	ForceUpdate    bool   `json:"forceUpdate"`
	Interval       string `json:"interval"`
	JobID          string `json:"jobID"`
	Webhook        string `json:"webhook"`
}

type Env struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type AccessPortainer struct {
	AccessLevel int `json:"AccessLevel"`
	UserID      int `json:"UserId"`
}

type TeamAccess struct {
	AccessLevel int `json:"AccessLevel"`
	TeamID      int `json:"TeamId"`
}

type ResourceControl struct {
	AccessLevel        int               `json:"AccessLevel"`
	AdministratorsOnly bool              `json:"AdministratorsOnly"`
	ID                 int               `json:"Id"`
	OwnerID            int               `json:"OwnerId"`
	Public             bool              `json:"Public"`
	ResourceID         string            `json:"ResourceId"`
	SubResourceIDs     []string          `json:"SubResourceIds"`
	System             bool              `json:"System"`
	TeamAccesses       []TeamAccess      `json:"TeamAccesses"`
	Type               int               `json:"Type"`
	UserAccesses       []AccessPortainer `json:"UserAccesses"`
}

type GitAuth struct {
	GitCredentialID int    `json:"gitCredentialID"`
	Password        string `json:"password"`
	Username        string `json:"username"`
}

type GitConfig struct {
	Authentication GitAuth `json:"authentication"`
	ConfigFilePath string  `json:"configFilePath"`
	ConfigHash     string  `json:"configHash"`
	ReferenceName  string  `json:"referenceName"`
	TLSSkipVerify  bool    `json:"tlsskipVerify"`
	URL            string  `json:"url"`
}

type Stacks struct {
	AdditionalFiles []string   `json:"AdditionalFiles"`
	AutoUpdate      AutoUpdate `json:"AutoUpdate"`
	EndpointId      int        `json:"EndpointId"`
	EntryPoint      string     `json:"EntryPoint"`
	Env             []Env      `json:"Env"`
	Id              int        `json:"Id"`
	Name            string     `json:"Name"`
	Option          struct {
		Prune bool `json:"prune"`
	} `json:"Option"`
	ResourceControl ResourceControl `json:"ResourceControl"`
	Status          int             `json:"Status"`
	SwarmId         string          `json:"SwarmId"`
	Type            int             `json:"Type"`
	CreatedBy       string          `json:"createdBy"`
	CreationDate    int64           `json:"creationDate"`
	FromAppTemplate bool            `json:"fromAppTemplate"`
	GitConfig       GitConfig       `json:"gitConfig"`
	IsComposeFormat bool            `json:"isComposeFormat"`
	Namespace       string          `json:"namespace"`
	ProjectPath     string          `json:"projectPath"`
	UpdateDate      int64           `json:"updateDate"`
	UpdatedBy       string          `json:"updatedBy"`
}

type PortainerResult struct {
	Stacks []Stacks `json:"stacks"`
}

type Portainer struct {
	BaseURL string
	Headers map[string]string
}

// ------------------ Status Container ----------------

type ResponseContainer struct {
	Snapshots []Snapshot `json:"snapshots"`
}

type Snapshot struct {
	DockerSnapshotRaw DockerSnapshotRaw `json:"DockerSnapshotRaw"`
}

type DockerSnapshotRaw struct {
	Containers []Container `json:"Containers"`
}

type Container struct {
	ID     string   `json:"Id"`
	Name   []string `json:"Names"`
	Labels Label    `json:"Labels"`
	State  string   `json:"State"`
	Status string   `json:"Status"`
}

type Label struct {
	ServiceID   string `json:"com.docker.swarm.service.id"`
	Namespace   string `json:"com.docker.stack.namespace"`
	NodeID      string `json:"com.docker.swarm.node.id"`
	ServiceName string `json:"com.docker.swarm.service.name"`
	TaskID      string `json:"com.docker.swarm.task.id"`
}

// ---------------- End Status Container --------------

//-------------------- Compose ------------------------

// DockerCompose represents the overall Docker Compose configuration
type DockerCompose struct {
	Networks map[string]Network `yaml:"networks"`
	Services map[string]Service `yaml:"services"`
}

// Network represents a network configuration
type Network struct {
	Driver     string `yaml:"driver"`
	Name       string `yaml:"name"`
	Attachable bool   `yaml:"attachable"`
	External   bool   `yaml:"external"`
}

// Service represents a service configuration
type Service struct {
	Image       string   `yaml:"image"`
	Hostname    string   `yaml:"hostname"`
	Volumes     []string `yaml:"volumes"`
	Networks    []string `yaml:"networks"`
	Environment []string `yaml:"environment"`
	Deploy      Deploy   `yaml:"deploy"`
}

// Deploy represents the deployment configuration
type Deploy struct {
	Mode          string        `yaml:"mode"`
	Replicas      int           `yaml:"replicas"`
	UpdateConfig  UpdateConfig  `yaml:"update_config"`
	RestartPolicy RestartPolicy `yaml:"restart_policy"`
}

// UpdateConfig represents the update configuration
type UpdateConfig struct {
	Parallelism int    `yaml:"parallelism"`
	Delay       string `yaml:"delay"`
}

// RestartPolicy represents the restart policy
type RestartPolicy struct {
	Condition string `yaml:"condition"`
}

// -------------------------------------------------------------

type PortainerAPI interface {
	InitialPortainer
	PortainerController
	utils.PortainerConfig
}

type InitialPortainer interface {
	PortainerConfig() (*utils.YamlStruct, error)
	NewPortainer() *Portainer
}

type PortainerController interface {
	GetStack() (*PortainerResult, error)
	AddStack(path, name string, input *CustomInput) (*http.Response, error)
	UpdateStack(id int, path, name string, input *CustomInput) (*http.Response, error)
	OperationStack(id int, action string) (*http.Response, error)
}

type PortainerMethod interface {
	GetPortainer(ctx context.Context, url string) (*http.Response, error)
	PostPortainer(ctx context.Context, url string, data io.Reader) (*http.Response, error)
	PutPortainer(ctx context.Context, url string, data io.Reader) (*http.Response, error)
	DeletePortainer(ctx context.Context, url string) (*http.Response, error)
}

// Controller Docker

type AddStackPortainer struct {
	Environment      []AddStackPortainerEnv `json:"env"`
	FromAppTemplate  bool                   `json:"fromAppTemplate"`
	Name             string                 `json:"name"`
	StackFileContent string                 `json:"stackFileContent"`
	Prune            bool                   `json:"prune"`
	PullImage        bool                   `json:"pullImage"`
	SwarmID          string                 `json:"swarmId"`
	Type             string                 `jsin:"type"`
	Method           string                 `json:"method"`
}

type AddStackPortainerEnv struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type CustomInput struct {
	Name              string
	WebImageContainer string
	APIImageContainer string
	DBRootUser        string
	DBRootPass        string
	DBHost            string
	DBPort            string
	DBWebName         string
	DBWebUser         string
	DBWebPass         string
	DBAPIName         string
	DBAPIUser         string
	DBAPIPass         string
	APIURL            string
	DockerPath        DockerPath
	UpdateStack       UpdateStack
}

type UpdateStack struct {
	ID        int
	Path      string
	Name      string
	Prune     bool
	PullImage bool
}

type DockerPath struct {
	Source string
	Dist   string
}

// Initial

func CustomInputDockerCompose(input *CustomInput) (*DockerCompose, error) {
	var data *DockerCompose
	path := input.DockerPath
	if path.Source == "" {
		return nil, errors.New("yaml path source is required")
	}
	if path.Dist == "" {
		return nil, errors.New("yaml path dist is required")
	}
	file, err := os.ReadFile(path.Source)
	if err != nil {
		return nil, err
	}
	content := string(file)
	replace := map[string]string{
		"<name>":                input.Name,
		"<web_image_container>": input.WebImageContainer,
		"<api_image_container>": input.APIImageContainer,
		"<db_root_user>":        input.DBRootUser,
		"<db_root_pass>":        input.DBRootPass,
		"<db_host>":             input.DBHost,
		"<db_port>":             input.DBPort,
		"<db_web_name>":         input.DBWebName,
		"<db_web_user>":         input.DBWebUser,
		"<db_web_pass>":         input.DBWebPass,
		"<db_api_name>":         input.DBAPIName,
		"<db_api_user>":         input.DBAPIUser,
		"<db_api_pass>":         input.DBAPIPass,
		"<api_url>":             input.APIURL,
	}

	for k, v := range replace {
		content = strings.ReplaceAll(content, k, v)
	}

	err = yaml.Unmarshal([]byte(content), &data)
	if err != nil {
		return nil, err
	}
	directory := fmt.Sprintf("%s/%s/docker", path.Dist, input.Name)
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		err := os.MkdirAll(directory, 0755)
		if err != nil {
			return nil, fmt.Errorf("error creating directory: %w", err)
		}
	}
	outputFile, err := os.Create(filepath.Join(directory, "compose.yml"))
	if err != nil {
		return nil, err
	}
	defer outputFile.Close()

	err = yaml.NewEncoder(outputFile).Encode(data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (p *Portainer) PortainerConfig() (*utils.YamlStruct, error) {
	config, err := utils.Open()
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (p *Portainer) NewPortainer(key string) (*Portainer, error) {
	if key == "" {
		return nil, errors.New("key is required")
	}
	config, err := p.PortainerConfig()
	if err != nil {
		return nil, errors.New("failed to read config.yml")
	}
	return &Portainer{
		BaseURL: config.Portainer.BaseURL,
		Headers: map[string]string{
			"Content-Type": "application/json",
			"User-Agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36",
			"X-API-KEY":    key,
		},
	}, nil
}
