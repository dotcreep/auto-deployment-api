package jenkins_api

import (
	"log"

	"github.com/dotcreep/go-automate-deploy/internal/service/jenkins"
)

type RequestInput struct {
	Username string `json:"username" example:"exampleusername"`
}

// This function is used to connect to jenkins with credential and header
func newJenkins() *jenkins.Jenkins {
	jen := jenkins.Jenkins{}
	connect, err := jen.NewJenkins()
	if err != nil {
		log.Println(err)
		return nil
	}
	return connect
}
