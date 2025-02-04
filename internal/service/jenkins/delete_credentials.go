package jenkins

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/dotcreep/go-automate-deploy/internal/utils"
)

func (j *Jenkins) DeleteCredential(ctx context.Context, data *JenkinsData) (string, error) {
	if data.Name == "" {
		return "", errors.New("data name is required")
	}
	yamlConfig, err := utils.Open()
	if yamlConfig == nil {
		return "", err
	}
	prefixUsername := fmt.Sprintf("%s%s", yamlConfig.Jenkins.PrefixCredentials, data.Name)
	data.PathURL = fmt.Sprintf("/credentials/store/system/domain/%s/credentials/%s/config.xml", yamlConfig.Jenkins.DomainCredentials, prefixUsername)
	resp, err := j.DeleteJenkins(ctx, data)
	if err != nil {
		return "", err
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}
	return "ok", nil
}
