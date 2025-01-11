package jenkins

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

type JenkinsMethod interface {
	JenkinsControllers
}

type JenkinsControllers interface {
	GetJenkins() (*http.Response, error)
	PostJenkins(ctx context.Context, url string, data io.Reader) (*http.Response, error)
	PutJenkins(ctx context.Context, url string, data io.Reader) (*http.Response, error)
	DeleteJenkins(ctx context.Context, url string) (*http.Response, error)
}

func (j *Jenkins) GetJenkins(ctx context.Context, data *JenkinsData) (*http.Response, error) {
	if data.PathURL == "" {
		return nil, errors.New("data path is required")
	}
	url := fmt.Sprintf("%s%s", j.BaseURL, data.PathURL)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	if j.Username == "" || j.Token == "" {
		return nil, errors.New("jenkins username or token is required")
	}

	authHeader := createAuthHeader(j.Username, j.Token)
	req.Header.Set("Authorization", authHeader)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (j *Jenkins) PostJenkins(ctx context.Context, data *JenkinsData) (*http.Response, error) {
	if data.PathURL == "" {
		return nil, errors.New("data path is required")
	}
	if data.Body == nil {
		return nil, errors.New("data body is required")
	}
	if j.Username == "" || j.Token == "" {
		return nil, errors.New("jenkins username or token is required")
	}
	url := fmt.Sprintf("%s%s", j.BaseURL, data.PathURL)

	req, err := http.NewRequestWithContext(ctx, "POST", url, data.Body)
	if err != nil {
		return nil, err
	}
	authHeader := createAuthHeader(j.Username, j.Token)
	req.Header.Set("Content-Type", "application/xml")
	req.Header.Set("Authorization", authHeader)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (j *Jenkins) DeleteJenkins(ctx context.Context, data *JenkinsData) (*http.Response, error) {
	if data.PathURL == "" {
		return nil, errors.New("data path is required")
	}

	url := fmt.Sprintf("%s%s", j.BaseURL, data.PathURL)

	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	if j.Username == "" || j.Token == "" {
		return nil, errors.New("jenkins username or token is required")
	}
	authHeader := createAuthHeader(j.Username, j.Token)
	req.Header.Set("Authorization", authHeader)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		log.Println(string(b))
		resp.Body = io.NopCloser(bytes.NewReader(b))
		return nil, fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}
	return resp, nil
}
