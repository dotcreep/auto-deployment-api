package cloudflare_api

import (
	"log"

	"github.com/dotcreep/go-automate-deploy/internal/service/cloudflare"
	"github.com/dotcreep/go-automate-deploy/internal/utils"
)

type RequestInput struct {
	Domain string `json:"domain" example:"sub.example.com"`
}

// This function is used to connect to cloudflare with credential and header
func newCloudflare() *cloudflare.Cloudflare {
	cf := cloudflare.Cloudflare{}
	yamlConfig, err := utils.Open()
	if err != nil {
		log.Println(err)
		return nil
	}
	connect, err := cf.NewCloudflare(yamlConfig.Cloudflare.Key, yamlConfig.Cloudflare.Email)
	if err != nil {
		log.Println(err)
		return nil
	}
	return connect
}
