package portainer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/dotcreep/go-automate-deploy/internal/utils"
	"gopkg.in/yaml.v3"
)

func (p *Portainer) AddStack(name string, input *CustomInput) (*http.Response, error) {
	path := input.DockerPath
	if name == "" {
		return nil, errors.New("name is required")
	}
	if path.Source == "" {
		return nil, errors.New("yaml path is required")
	}
	ctx, cancel := utils.Cfgx{}.LongTimeout()
	defer cancel()

	compose, err := CustomInputDockerCompose(input)
	if err != nil {
		return nil, err
	}

	yamlData, err := yaml.Marshal(compose)
	if err != nil {
		return nil, err
	}
	// Initial Data from config
	dataInternal, err := p.PortainerConfig()
	if err != nil {
		return nil, err
	}
	yamldi := dataInternal.Portainer

	// Preaparing Data
	dataInput := &AddStackPortainer{
		FromAppTemplate:  yamldi.FromAppTemplate,
		Name:             name,
		StackFileContent: utils.YamlIndent(yamlData, 2),
		SwarmID:          yamldi.SwarmID,
		Type:             yamldi.Type,
		Method:           yamldi.Method,
		Environment:      []AddStackPortainerEnv{},
	}

	jsonData, err := json.Marshal(dataInput)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/stacks/create/swarm/string?endpointId=%d", p.BaseURL, yamldi.EndpointId)

	resp, err := p.PostPortainer(ctx, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Println(string(b))
		resp.Body = io.NopCloser(bytes.NewReader(b))
		return nil, fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}
	resp.Body = io.NopCloser(bytes.NewReader(b))
	return resp, nil
}
