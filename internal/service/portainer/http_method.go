package portainer

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

func (p *Portainer) GetPortainer(ctx context.Context, url string) (*http.Response, error) {
	if url == "" {
		return nil, errors.New("url is required")
	}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range p.Headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
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

func (p *Portainer) PostPortainer(ctx context.Context, url string, data io.Reader) (*http.Response, error) {
	if url == "" {
		return nil, errors.New("url is required")
	}
	if data == nil {
		return nil, errors.New("data is empty")
	}
	req, err := http.NewRequestWithContext(ctx, "POST", url, data)
	if err != nil {
		return nil, err
	}
	for k, v := range p.Headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
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

func (p *Portainer) PutPortainer(ctx context.Context, url string, data io.Reader) (*http.Response, error) {
	if url == "" {
		return nil, errors.New("url is required")
	}
	if data == nil {
		return nil, errors.New("data is empty")
	}
	req, err := http.NewRequestWithContext(ctx, "PUT", url, data)
	if err != nil {
		return nil, err
	}

	for k, v := range p.Headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
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

func (p *Portainer) DeletePortainer(ctx context.Context, url string) (*http.Response, error) {
	if url == "" {
		return nil, errors.New("url is required")
	}
	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range p.Headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusNoContent {
		log.Println(string(b))
		resp.Body = io.NopCloser(bytes.NewReader(b))
		return nil, fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}
	resp.Body = io.NopCloser(bytes.NewReader(b))
	return resp, nil
}
