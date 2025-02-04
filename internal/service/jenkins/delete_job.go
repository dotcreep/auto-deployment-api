package jenkins

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

func (j *Jenkins) DeleteJob(ctx context.Context, data *JenkinsData) (string, error) {
	if data.Name == "" {
		return "", errors.New("data name is required")
	}
	data.PathURL = fmt.Sprintf("/job/%s/", data.Name)
	resp, err := j.DeleteJenkins(ctx, data)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusNoContent {
		return "", fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}
	return "ok", nil
}
