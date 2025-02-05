package jenkins

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/dotcreep/go-automate-deploy/internal/utils"
)

//--------------------------------------

type Project struct {
	XMLName            xml.Name       `xml:"project"`
	Actions            []Action       `xml:"actions"`
	Description        string         `xml:"description"`
	KeepDependencies   bool           `xml:"keepDependencies"`
	Properties         Properties     `xml:"properties"`
	SCM                SCM            `xml:"scm"`
	CanRoam            bool           `xml:"canRoam"`
	Disabled           bool           `xml:"disabled"`
	BlockBuildWhenDown bool           `xml:"blockBuildWhenDownstreamBuilding"`
	BlockBuildWhenUp   bool           `xml:"blockBuildWhenUpstreamBuilding"`
	Triggers           []Trigger      `xml:"triggers>com.cloudbees.jenkins.GitHubPushTrigger"`
	ConcurrentBuild    bool           `xml:"concurrentBuild"`
	Builders           []Builder      `xml:"builders>hudson.tasks.Shell"`
	Publishers         []Publisher    `xml:"publishers"`
	BuildWrappers      []BuildWrapper `xml:"buildWrappers"`
}

type Action struct {
	// Define fields based on the <actions> tag structure if available
}

type Properties struct {
	GithubProjectProperty GithubProjectProperty `xml:"com.coravy.hudson.plugins.github.GithubProjectProperty"`
}

type GithubProjectProperty struct {
	ProjectUrl  string `xml:"projectUrl"`
	DisplayName string `xml:"displayName"`
}

type SCM struct {
	Class             string             `xml:"class,attr"`
	Plugin            string             `xml:"plugin,attr"`
	ConfigVersion     int                `xml:"configVersion"`
	UserRemoteConfigs []UserRemoteConfig `xml:"userRemoteConfigs>hudson.plugins.git.UserRemoteConfig"`
	Branches          []BranchSpec       `xml:"branches>hudson.plugins.git.BranchSpec"`
}

type UserRemoteConfig struct {
	Url           string `xml:"url"`
	CredentialsId string `xml:"credentialsId"`
}

type BranchSpec struct {
	Name string `xml:"name"`
}

type Trigger struct {
	Plugin string `xml:"plugin,attr"`
	Spec   string `xml:"spec"` // Spec bisa berisi string jika ada kontennya
}

type Builder struct {
	Command string `xml:"command"`
}

type Publisher struct {
	// Define fields based on the <publishers> tag structure if available
}

type BuildWrapper struct {
	Cleanup       *PreBuildCleanup    `xml:"hudson.plugins.ws__cleanup.PreBuildCleanup,omitempty"`
	SecretBinding *SecretBuildWrapper `xml:"org.jenkinsci.plugins.credentialsbinding.impl.SecretBuildWrapper,omitempty"`
}

type PreBuildCleanup struct {
	Plugin     string `xml:"plugin"`
	DeleteDirs bool   `xml:"deleteDirs"`
}

type SecretBuildWrapper struct {
	FileBindings   []Binding `xml:"bindings>org.jenkinsci.plugins.credentialsbinding.impl.FileBinding"`
	StringBindings []Binding `xml:"bindings>org.jenkinsci.plugins.credentialsbinding.impl.StringBinding"`
}

type Binding struct {
	CredentialsId string `xml:"credentialsId"`
	Variable      string `xml:"variable"`
}

//--------------------------------------

type JenkinsCredentials struct {
	JenUserPassword JenUserPassword
	JenFiles        JenFiles
	JenStrings      JenStrings
	JenGithub       JenGithub
	JenCertificate  JenCertificate
	JenSSHUsername  JenSSHUsername
}

type JenUserPassword struct {
	XMLName        xml.Name `xml:"com.cloudbees.plugins.credentials.impl.UsernamePasswordCredentialsImpl"`
	Scope          string   `xml:"scope"`
	Id             string   `xml:"id"`
	Description    string   `xml:"description"`
	UsernameSecret string   `xml:"usernameSecret"`
	Username       string   `xml:"username"`
	Password       string   `xml:"password"`
}

type JenFiles struct {
	XMLName     xml.Name `xml:"org.jenkinsci.plugins.plaincredentials.impl.FileCredentialsImpl"`
	Scope       string   `xml:"scope"`
	Id          string   `xml:"id"`
	Filename    string   `xml:"fileName"`
	SecretBytes string   `xml:"secretBytes"`
	Description string   `xml:"description"`
}

type JenStrings struct {
	XMLName     xml.Name `xml:"org.jenkinsci.plugins.plaincredentials.impl.StringCredentialsImpl"`
	Scope       string   `xml:"scope"`
	Id          string   `xml:"id"`
	Secret      string   `xml:"secret"`
	Description string   `xml:"description"`
}

type JenGithub struct {
	XMLName     xml.Name `xml:"com.cloudbees.plugins.credentials.impl.BaseStandardCredentials"`
	Scope       string   `xml:"scope"`
	Id          string   `xml:"id"`
	Description string   `xml:"description"`
	Token       string   `xml:"token"`
}

type JenCertificate struct {
	XMLName     xml.Name `xml:"com.cloudbees.plugins.credentials.impl.CertificateCredentialsImpl"`
	Scope       string   `xml:"scope"`
	Id          string   `xml:"id"`
	Description string   `xml:"description"`
	Certificate JenCerts `xml:"certificate"`
}

type JenCerts struct {
	Password              string `xml:"password"`
	UploadedKeystoreBytes string `xml:"uploadedKeystoreBytes"`
	Description           string `xml:"description"`
}

type JenSSHUsername struct {
	XMLName          xml.Name            `xml:"com.cloudbees.jenkins.plugins.sshcredentials.impl.BasicSSHUserPrivateKey"`
	Scope            string              `xml:"scope"`
	Id               string              `xml:"id"`
	Description      string              `xml:"description"`
	Username         string              `xml:"username"`
	PrivateKeySource JenPrivateKeySource `xml:"privateKeySource"`
}

type JenPrivateKeySource struct {
	Class       string `xml:"class,attr"`
	PrivateKey  string `xml:"privateKey"`
	Description string `xml:"description"`
}

type FreeStyleBuild struct {
	Class             string        `json:"_class"`
	Actions           []interface{} `json:"actions"`
	Artifacts         []interface{} `json:"artifacts"`
	Building          bool          `json:"building"`
	Description       interface{}   `json:"description"`
	DisplayName       string        `json:"displayName"`
	Duration          int           `json:"duration"`
	EstimatedDuration int           `json:"estimatedDuration"`
	Executor          interface{}   `json:"executor"`
	FullDisplayName   string        `json:"fullDisplayName"`
	ID                string        `json:"id"`
	InProgress        bool          `json:"inProgress"`
	KeepLog           bool          `json:"keepLog"`
	Number            int           `json:"number"`
	QueueID           int           `json:"queueId"`
	Result            string        `json:"result"`
	Timestamp         int64         `json:"timestamp"`
	URL               string        `json:"url"`
	BuiltOn           string        `json:"builtOn"`
	ChangeSet         ChangeSet     `json:"changeSet"`
	Culpits           []interface{} `json:"culprits"`
}

type ChangeSet struct {
	Class string        `json:"_class"`
	Items []interface{} `json:"items"`
	Kind  string        `json:"kind"`
}

type Jenkins struct {
	BaseURL    string
	Username   string
	Token      string
	HeaderJson map[string]string
	HeaderForm map[string]string
	HeaderXML  map[string]string
}

type JenkinsJob struct {
	Class string `json:"_class"`
	Jobs  []Job  `json:"jobs"`
}

type Job struct {
	Class string `json:"_class"`
	Name  string `json:"name"`
	URL   string `json:"url"`
	Color string `json:"color"`
}

// Controller Interface

func (j *Jenkins) JenkinsConfig() (*utils.YamlStruct, error) {
	config, err := utils.Open()
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (j *Jenkins) NewJenkins() (*Jenkins, error) {
	config, err := j.JenkinsConfig()
	if err != nil {
		return nil, errors.New("failed to read config.yml")
	}
	return &Jenkins{
		Username: config.Jenkins.Username,
		Token:    config.Jenkins.APIKey,
		BaseURL:  config.Jenkins.BaseURL,
		HeaderJson: map[string]string{
			"Content-Type": "application/json",
			"User-Agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/237.84.2.178 Safari/537.36",
			"Accept":       "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
		},
		HeaderForm: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
			"User-Agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/237.84.2.178 Safari/537.36",
			"Accept":       "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
		},
		HeaderXML: map[string]string{
			"Content-Type": "application/xml",
			"User-Agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/237.84.2.178 Safari/537.36",
			"Accept":       "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
		},
	}, nil
}

type JenkinsData struct {
	PathURL            string
	Body               io.Reader
	DomainCredentials  string
	APIURL             string
	Files              string
	Username           string
	MerchantName       string
	MerchantID         int
	ID                 string
	Description        string
	CredentialsStore   CredentialsStore
	DomainWrapper      DomainWrapper
	JenkinsItem        Project
	JenkinsCredentials JenkinsCredentials
	PaketMerchant      string
}

type CredentialsStore struct {
	Class   string            `json:"_class"`
	Domains map[string]Domain `json:"domains"`
}

type Domain struct {
	Class string `json:"_class"`
	Name  string `json:"name"`
}

// Domain Credentials Result
type DomainWrapper struct {
	Class       string       `json:"_class"`
	Credentials []Credential `json:"credentials"`
}

type Credential struct {
	Description string      `json:"description"`
	DisplayName string      `json:"displayName"`
	Fingerprint interface{} `json:"fingerprint"`
	FullName    string      `json:"fullName"`
	ID          string      `json:"id"`
	TypeName    string      `json:"typeName"`
}

func createAuthHeader(username, password string) string {
	auth := fmt.Sprintf("%s:%s", username, password)
	encodeAuth := base64.StdEncoding.EncodeToString([]byte(auth))
	return fmt.Sprintf("Basic %s", encodeAuth)
}

// Operational Interface
// ** Operation Credetial

func (j *Jenkins) AddDomainCredentials(data *JenkinsData) (*http.Response, error) {
	if data.Body == nil {
		return nil, errors.New("data body is required")
	}
	ctx, cancel := utils.Cfgx{}.DefaultTimeout()
	defer cancel()
	data = &JenkinsData{
		PathURL: "/scriptText",
		Body:    data.Body,
	}

	url := fmt.Sprintf("/%s", data.PathURL)

	req, err := http.NewRequestWithContext(ctx, "POST", url, data.Body)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	b, _ := io.ReadAll(req.Body)
	if resp.StatusCode != http.StatusOK {
		req.Body = io.NopCloser(bytes.NewReader(b))
		return nil, fmt.Errorf("unexpected status code: %v\ndata: %s", resp.StatusCode, string(b))
	}

	return resp, nil
}

func (j *Jenkins) AddCredentials(data *JenkinsData, typeInput string) (*http.Response, error) {
	if data.ID == "" {
		return nil, errors.New("data id is required")
	}
	if data.Username == "" {
		return nil, errors.New("data name is required")
	}
	// if data.Files == "" {
	// return nil, errors.New("data files is required")
	// }
	if typeInput == "" {
		return nil, errors.New("data type is required")
	}
	if data.DomainCredentials == "" {
		return nil, errors.New("data domain credentials is required")
	}
	ctx, cancel := utils.Cfgx{}.DefaultTimeout()
	defer cancel()
	var customBody interface{}
	switch typeInput {
	case "file":
		if data.JenkinsCredentials.JenFiles.Id == "" ||
			data.JenkinsCredentials.JenFiles.Description == "" ||
			data.JenkinsCredentials.JenFiles.Filename == "" ||
			data.JenkinsCredentials.JenFiles.SecretBytes == "" {
			return nil, errors.New("data jenkins files is required")
		}
		customBody = data.JenkinsCredentials.JenFiles
	default:
		return nil, errors.New("type is not in listed")
	}
	xmlData, err := xml.Marshal(customBody)
	if err != nil {
		return nil, err
	}
	data.PathURL = fmt.Sprintf("/credentials/store/system/domain/%s/createCredentials", data.DomainCredentials)
	data.Body = bytes.NewBuffer(xmlData)

	resp, err := j.PostJenkins(ctx, data)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (j *Jenkins) GetJobOperation(data *JenkinsData, action string) (*http.Response, error) {
	var url string
	if action == "" {
		return nil, errors.New("action is required")
	}
	ctx, cancel := utils.Cfgx{}.DefaultTimeout()
	defer cancel()
	switch action {
	case "Status":
		if data.Username == "" {
			return nil, errors.New("data name is required")
		}
		url = fmt.Sprintf("/job/%s/lastBuild/api/json?depth=1", data.Username)
	case "DetailStatus":
		if data.Username == "" {
			return nil, errors.New("data name is required")
		}
		url = fmt.Sprintf("/job/%s/api/json?depth=1", data.Username)
	case "AllJob":
		url = "/api/json?tree=jobs[name,color,url]"
	case "JobConfig":
		if data.Username == "" {
			return nil, errors.New("data name is required")
		}
		url = fmt.Sprintf("/job/%s/config.xml", data.Username)
	default:
		return nil, errors.New("action is not in listed")
	}

	resp, err := j.GetJenkins(ctx, &JenkinsData{PathURL: url})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (j *Jenkins) PostJobOperation(data *JenkinsData, action string) (*http.Response, error) {
	ctx, cancel := utils.Cfgx{}.DefaultTimeout()
	defer cancel()
	var url string
	if action == "" {
		return nil, errors.New("action is required")
	}
	switch action {
	case "Update":
		if data.Username == "" {
			return nil, errors.New("data name is required")
		}
		url = fmt.Sprintf("/job/%s/config.xml", data.Username)
	case "Delete":
		if data.Username == "" {
			return nil, errors.New("data name is required")
		}
		url = fmt.Sprintf("/job/%s/doDelete", data.Username)
	case "Build":
		if data.Username == "" {
			return nil, errors.New("data name is required")
		}
		if data.Body == nil {
			return nil, errors.New("data body is required")
		}
		url = fmt.Sprintf("/job/%s/build", data.Username)
	default:
		return nil, errors.New("action is not in listed")
	}

	JenkinsData := &JenkinsData{
		PathURL: url,
		Body:    data.Body,
	}

	resp, err := j.PostJenkins(ctx, JenkinsData)
	if err != nil {
		return nil, err
	}

	// b, _ := io.ReadAll(resp.Body)
	// if resp.StatusCode != http.StatusOK {
	// 	resp.Body = io.NopCloser(bytes.NewReader(b))
	// 	return "", fmt.Errorf("unexpected status code: %v\ndata: %s", resp.StatusCode, string(b))
	// }
	return resp, nil
}

func (j *Jenkins) GetCredentialOperation(ctx context.Context, data *JenkinsData, action string) (*http.Response, error) {
	if action == "" {
		return nil, errors.New("action is required")
	}
	if j.BaseURL == "" {
		return nil, errors.New("base url is required")
	}
	switch action {
	case "ReadDomain":
		data.PathURL = "/credentials/store/system/api/json"
	case "ReadAllCredentials":
		data.PathURL = fmt.Sprintf("/credentials/store/system/domain/%s/api/json?depth=1", data.DomainCredentials)
	case "ReadDetailCredentials":
		if data.Username == "" {
			return nil, errors.New("data name is required")
		}
		data.PathURL = fmt.Sprintf("/credentials/store/system/domain/%s/credentials/%s/api/json", data.DomainCredentials, data.Username)
	default:
		return nil, errors.New("action is not in listed")
	}
	resp, err := j.GetJenkins(ctx, data)
	if err != nil {
		return nil, err
	}

	// Mengembalikan body sebagai string
	return resp, nil
}
