package jenkins

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
)

func (j *Jenkins) AddItem(ctx context.Context, data *JenkinsData) (string, error) {
	if data.Body == nil {
		return "", errors.New("data body is required")
	}
	if data.Username == "" {
		return "", errors.New("data name is required")
	}
	xmlData, err := xml.Marshal(data.JenkinsItem)
	if err != nil {
		return "", err
	}
	body := bytes.NewBuffer(xmlData)
	data.PathURL = fmt.Sprintf("/createItem?name=%s", data.Username)
	data.Body = body

	resp, err := j.PostJenkins(ctx, data)
	if err != nil {
		return "", err
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}

	return "ok", nil
}
