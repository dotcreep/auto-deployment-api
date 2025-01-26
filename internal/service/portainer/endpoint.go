package portainer

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/dotcreep/go-automate-deploy/internal/utils"
)

func (p *Portainer) Endpoint(id int) (*http.Response, error) {
	if id == 0 {
		return nil, errors.New("id is required")
	}
	if p.BaseURL == "" {
		return nil, errors.New("base url is required")
	}
	ctx, cancel := utils.Cfgx{}.LongTimeout()
	defer cancel()

	url := fmt.Sprintf("%s/endpoints/%d", p.BaseURL, id)
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
		resp.Body = io.NopCloser(bytes.NewReader(b))
		return nil, fmt.Errorf("unexpected status code: %v\ndata: %s", resp.StatusCode, string(b))
	}
	resp.Body = io.NopCloser(bytes.NewReader(b))
	return resp, nil
}
