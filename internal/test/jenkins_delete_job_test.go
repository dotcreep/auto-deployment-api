package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/dotcreep/go-automate-deploy/internal/service/jenkins"
	"github.com/dotcreep/go-automate-deploy/internal/utils"
)

func connectJenkins() (*jenkins.Jenkins, error) {
	j := jenkins.Jenkins{}
	cfg, err := utils.Open()
	if err != nil {
		return nil, err
	}
	j.Username = cfg.Jenkins.Username
	j.Token = cfg.Jenkins.APIKey
	connect, err := j.NewJenkins()
	if err != nil {
		return nil, err
	}
	return connect, nil
}

func TestJenkinsDeleteJob(t *testing.T) {
	connect, err := connectJenkins()
	if err != nil {
		t.Error(err)
	}
	// NOTE: url required add / in the last
	// Example 1 {{jenkins_url}}/job/aspirasihomman  <- this not working
	// Example 2 {{jenkins_url}}/job/aspirasihomman/ <- this working
	jenkinData := jenkins.JenkinsData{}
	jenkinData.Name = "aspirasihomman"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := connect.DeleteJob(ctx, &jenkinData)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(resp)
}

func TestJenkinsDeleteCredentials(t *testing.T) {
	connect, err := connectJenkins()
	if err != nil {
		t.Error(err)
	}
	jenkinData := jenkins.JenkinsData{}
	jenkinData.Name = "dika12"
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	resp, err := connect.DeleteCredential(ctx, &jenkinData)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(resp)
}
