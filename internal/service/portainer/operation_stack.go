package portainer

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/dotcreep/go-automate-deploy/internal/utils"
)

func (p *Portainer) OperationStack(id int, action string) (*http.Response, error) {
	if id == 0 {
		return nil, errors.New("id is required")
	}
	if action == "" {
		return nil, errors.New("action is required")
	}
	ctx, cancel := utils.Cfgx{}.LongTimeout()
	defer cancel()
	yamlConfig, err := utils.Open()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/stacks/%d/%s?endpointId=%d", p.BaseURL, id, action, yamlConfig.Portainer.EndpointId)
	resp, err := p.PostPortainer(ctx, url, nil)
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
