package portainer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/dotcreep/go-automate-deploy/internal/utils"
	"gopkg.in/yaml.v3"
)

func (p *Portainer) UpdateStack(input *CustomInput) (*http.Response, error) {
	if input.UpdateStack.ID == 0 {
		return nil, errors.New("id is required")
	}
	if input.UpdateStack.Name == "" {
		return nil, errors.New("name is required")
	}
	if input.DockerPath.Source == "" {
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
	yaml := dataInternal.Portainer
	// Preaparing Data
	dataInput := &AddStackPortainer{
		FromAppTemplate:  yaml.FromAppTemplate,
		Name:             input.UpdateStack.Name,
		StackFileContent: string(yamlData),
		SwarmID:          yaml.SwarmID,
		Type:             yaml.Type,
		Method:           yaml.Method,
		Prune:            input.UpdateStack.Prune,
		PullImage:        input.UpdateStack.PullImage,
	}

	jsonData, err := json.Marshal(dataInput)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/stacks/%d?endpointId=%d", p.BaseURL, input.UpdateStack.ID, yaml.EndpointId)

	resp, err := p.PutPortainer(ctx, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body = io.NopCloser(bytes.NewReader(b))
		return nil, fmt.Errorf("unexpected status code: %v\ndata: %s", resp.StatusCode, string(b))
	}
	resp.Body = io.NopCloser(bytes.NewReader(b))
	return resp, nil
}
