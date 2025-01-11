package portainer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func (p *Portainer) GetStack(ctx context.Context) (*PortainerResult, error) {
	url := fmt.Sprintf("%s/stacks", p.BaseURL)
	resp, err := p.GetPortainer(ctx, url)
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

	// resp.Body = io.NopCloser(bytes.NewReader(b))
	// return b, nil

	var stacks []Stacks
	err = json.Unmarshal(b, &stacks)
	if err != nil {
		return nil, err
	}
	return &PortainerResult{Stacks: stacks}, nil
}
