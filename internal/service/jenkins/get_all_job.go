package jenkins

import (
	"context"
	"encoding/json"
	"errors"
	"io"
)

func (j *Jenkins) GetAllJobItem(ctx context.Context) ([]string, error) {
	url := "/api/json?tree=jobs[name,color,url]"

	resp, err := j.GetJenkins(ctx, &JenkinsData{PathURL: url})
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	jobs, ok := result["jobs"].([]interface{})
	if !ok {
		return nil, errors.New("failed to parse jobs")
	}

	var jobNames []string
	for _, job := range jobs {
		jobMap, ok := job.(map[string]interface{})
		if !ok {
			return nil, errors.New("failed to parse job")
		}
		name, ok := jobMap["name"].(string)
		if !ok {
			return nil, errors.New("failed to parse job name")
		}
		jobNames = append(jobNames, name)
	}

	return jobNames, nil

}