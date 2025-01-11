package jenkins

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dotcreep/go-automate-deploy/internal/utils"
)

func (j *Jenkins) GetCredentialsIsNotExists(ctx context.Context, username string) (string, error) {
	if username == "" {
		return "", errors.New("username is required")
	}

	yamlConfig, err := utils.Open()
	if err != nil {
		return "", err
	}
	data := &JenkinsData{}
	data.DomainCredentials = yamlConfig.Jenkins.DomainCredentials
	res, err := j.GetCredentialOperation(ctx, data, "ReadAllCredentials")
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var result DomainWrapper
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return "", err
	}
	var endResults bool
	endResults = false
	for _, v := range result.Credentials {
		if v.ID == fmt.Sprintf("%s%s", yamlConfig.Jenkins.PrefixCredentials, username) {
			endResults = true
			break
		}
	}
	if endResults {
		return "", errors.New("credentials is exists")
	}
	return "credentials is ready for create", nil
}
