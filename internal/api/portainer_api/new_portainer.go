package portainer_api

import (
	"log"

	"github.com/dotcreep/go-automate-deploy/internal/service/portainer"
	"github.com/dotcreep/go-automate-deploy/internal/utils"
)

type RequestInput struct {
	Username string `json:"username" example:"exampleusername"`
}

// This function is used to connect to portainer with credential and header
func newPortainer() *portainer.Portainer {
	pt := portainer.Portainer{}
	yamlConf, err := utils.Open()
	if err != nil {
		return nil
	}
	connect, err := pt.NewPortainer(yamlConf.Portainer.APIKey)
	if err != nil {
		log.Println(err)
		return nil
	}
	return connect
}
