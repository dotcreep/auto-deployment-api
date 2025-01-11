package deploy_api

import (
	"fmt"

	"github.com/dotcreep/go-automate-deploy/internal/service/portainer"
)

func connectPortainer() (*portainer.Portainer, error) {
	p := portainer.Portainer{}
	secrets := &Secrets{}
	secrets.GetSecret()
	connect, err := p.NewPortainer(secrets.Portainer.Password)
	if err != nil {
		return nil, err
	}
	return connect, nil
}

func StringChecker(field, str string) error {
	if str == "" {
		return fmt.Errorf("%s is required", field)
	}
	return nil
}

func DeployPortainer(custom *portainer.CustomInput) (string, error) {
	// Check input
	if err := StringChecker("web image", custom.WebImageContainer); err != nil {
		return "", err
	}
	if err := StringChecker("api image", custom.APIImageContainer); err != nil {
		return "", err
	}
	if err := StringChecker("api url", custom.APIURL); err != nil {
		return "", err
	}
	if err := StringChecker("db host", custom.DBHost); err != nil {
		return "", err
	}
	if err := StringChecker("db port", custom.DBPort); err != nil {
		return "", err
	}
	if err := StringChecker("db root user", custom.DBRootUser); err != nil {
		return "", err
	}
	if err := StringChecker("db root pass", custom.DBRootPass); err != nil {
		return "", err
	}
	if err := StringChecker("db web name", custom.DBWebName); err != nil {
		return "", err
	}
	if err := StringChecker("db web user", custom.DBWebUser); err != nil {
		return "", err
	}
	if err := StringChecker("db web pass", custom.DBWebPass); err != nil {
		return "", err
	}
	if err := StringChecker("db api name", custom.DBAPIName); err != nil {
		return "", err
	}
	if err := StringChecker("db api user", custom.DBAPIUser); err != nil {
		return "", err
	}
	if err := StringChecker("db api pass", custom.DBAPIPass); err != nil {
		return "", err
	}
	if err := StringChecker("docker compose", custom.DockerPath.Source); err != nil {
		return "", err
	}

	// Deploy Portainer
	connect, err := connectPortainer()
	if err != nil {
		return "", err
	}

	// Database
	// custom.DockerPath.Source = "storage/src/docker/database.yml"
	// custom.DockerPath.Dist = "storage/dist"
	// _, err = connect.AddStack(custom.Name, custom)
	// if err != nil {
	// return "", err
	// }
	// Web
	custom.DockerPath.Source = "storage/src/docker/web.yml"
	custom.DockerPath.Dist = "storage/dist"
	_, err = connect.AddStack(custom.Name, custom)
	if err != nil {
		return "", err
	}
	return "success deploy portainer", nil
}
