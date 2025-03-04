package portainer

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/dotcreep/go-automate-deploy/internal/utils"
)

func (p *Portainer) UpdateStackByName(ctx context.Context, name string) (string, error) {
	if name == "" {
		return "", errors.New("name is required")
	}
	// get list of stack
	lists, err := p.GetStack(ctx)
	if err != nil {
		return "", err
	}
	data := struct {
		Name string
		ID   int
	}{}
	var isExist bool
	for _, v := range lists.Stacks {
		if v.Name == name {
			data.Name = v.Name
			data.ID = v.Id
			isExist = true
			break
		}
	}
	if !isExist {
		return "", errors.New("stack not found")
	}
	// update stack
	custom := &CustomInput{}
	custom.UpdateStack.ID = data.ID
	custom.UpdateStack.Name = data.Name
	custom.UpdateStack.Prune = true
	custom.UpdateStack.PullImage = true

	cfg, err := utils.Open()
	if err != nil {
		return "", err
	}

	dataInput := &AddStackPortainer{
		FromAppTemplate: cfg.Portainer.FromAppTemplate,
		Name:            data.Name,
		// StackFileContent: string(cfg.Portainer.StackFileContent), // TODO: fix this to file
		SwarmID:   cfg.Portainer.SwarmID,
		Type:      cfg.Portainer.Type,
		Method:    cfg.Portainer.Method,
		Prune:     custom.UpdateStack.Prune,
		PullImage: custom.UpdateStack.PullImage,
	}

	jsonData, err := json.Marshal(dataInput)
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("%s/stacks/%d?endpointId=%d", p.BaseURL, custom.UpdateStack.ID, cfg.Portainer.EndpointId)

	resp, err := p.PutPortainer(ctx, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body = io.NopCloser(bytes.NewReader(b))
		return "", fmt.Errorf("unexpected status code: %v\ndata: %s", resp.StatusCode, string(b))
	}
	resp.Body = io.NopCloser(bytes.NewReader(b))

	return "success update stack", nil
}
