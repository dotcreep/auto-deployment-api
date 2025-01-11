package jenkins

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type Status struct {
	Color  string `json:"color"`
	Status string `json:"status"`
}

func (j *Jenkins) GetStatusOfItem(ctx context.Context, username string) (*Status, error) {
	if username == "" {
		return nil, errors.New("username is required")
	}
	url := fmt.Sprintf("/job/%s/api/json?depth=1", username)
	resp, err := j.GetJenkins(ctx, &JenkinsData{PathURL: url})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound && resp.StatusCode != http.StatusOK {
		return nil, errors.New("item not found")
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	color, ok := result["color"].(string)
	if !ok {
		return nil, errors.New("failed to parse color")
	}
	response := &Status{}
	if color == "red" {
		response.Color = "red"
		response.Status = "Failed"
		return response, nil
	} else if color == "blue" {
		response.Color = "green"
		response.Status = "Ready"
		return response, nil
	} else if color == "notbuilt" {
		response.Color = "blue"
		response.Status = "Not Built"
		return response, nil
	} else if strings.Contains(color, "anime") {
		response.Color = "blue"
		response.Status = "Process"
		return response, nil
	} else if color == "aborted" {
		response.Color = "gray"
		response.Status = "Aborted"
		return response, nil
	} else {
		response.Color = "gray"
		response.Status = "Unknown"
		return response, nil
	}
}
