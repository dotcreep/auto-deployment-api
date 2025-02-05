package deploy_api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dotcreep/go-automate-deploy/internal/service/jenkins"
	"github.com/dotcreep/go-automate-deploy/internal/utils"
)

func connectJenkins() (*jenkins.Jenkins, error) {
	j := jenkins.Jenkins{}
	secrets := &Secrets{}
	secrets.GetSecret()
	j.Username = secrets.Jenkins.Username
	j.Token = secrets.Jenkins.Password
	connect, err := j.NewJenkins()
	if err != nil {
		return nil, err
	}
	return connect, nil
}

func DeployJenkins(ctx context.Context, data *jenkins.JenkinsData, domain string) (string, error) {
	connect, err := connectJenkins()
	if err != nil {
		return "", err
	}
	if domain == "" {
		return "", errors.New("domain is required")
	}
	if data.DomainCredentials == "" {
		return "", errors.New("domain credentials is required")
	}
	if data.Username == "" {
		return "", errors.New("name is required")
	}
	if data.MerchantName == "" {
		return "", errors.New("merchant name is required")
	}
	if data.ID == "" {
		return "", errors.New("data id is required")
	}
	// if data.Files == "" {
	// return "", errors.New("data files is required")
	// }
	// Check Domain of client
	readDomain, err := connect.GetCredentialOperation(ctx, data, "ReadDomain")
	if err != nil {
		return "", err
	}
	defer readDomain.Body.Close()

	var domainList jenkins.CredentialsStore
	err = json.NewDecoder(readDomain.Body).Decode(&domainList)
	if err != nil {
		return "", err
	}
	for _, v := range domainList.Domains {
		if v.Name == data.DomainCredentials {
			data.DomainCredentials = v.Name
			break
		}
	}
	// Check Domain of client
	readAll, err := connect.GetCredentialOperation(ctx, data, "ReadAllCredentials")
	if err != nil {
		return "", err
	}
	defer readAll.Body.Close()
	var domainListAll jenkins.DomainWrapper
	err = json.NewDecoder(readAll.Body).Decode(&domainListAll)
	if err != nil {
		return "", err
	}
	for _, v := range domainListAll.Credentials {
		if v.ID == data.ID {
			return "", errors.New("credentials already exists")
		}
	}
	// Deploy Jenkins
	//------------------------------------
	// Add Credentials
	envAndroidPath := "storage/src/android/.env"
	fileContent, err := os.ReadFile(envAndroidPath)
	if err != nil {
		return "", err
	}
	packageName := utils.GeneratePackageName(data.Username, domain)
	// mId := strconv.Itoa(data.MerchantID)
	replaceEnv := strings.ReplaceAll(string(fileContent), "<app_id>", packageName)
	replaceEnv = strings.ReplaceAll(replaceEnv, "<url_api>", data.APIURL)
	replaceEnv = strings.ReplaceAll(replaceEnv, "<url_web>", fmt.Sprintf("https://%s", domain))
	// replaceEnv = strings.ReplaceAll(replaceEnv, "<merchant_id>", mId)
	getNameLength := strings.Split(data.MerchantName, " ")
	var labelName string
	if len(getNameLength) > 2 {
		labelName = fmt.Sprintf("%s %s", getNameLength[0], getNameLength[1])
	} else {
		labelName = data.MerchantName
	}
	replaceEnv = strings.ReplaceAll(replaceEnv, "<label_name>", labelName)
	replaceEnv = strings.ReplaceAll(replaceEnv, "<host_name>", domain)
	replaceEnv = strings.ReplaceAll(replaceEnv, "<app_title>", fmt.Sprintf("\"%s\"", data.MerchantName))
	replaceEnv = strings.ReplaceAll(replaceEnv, "<username>", data.Username)
	replaceEnv = strings.ReplaceAll(replaceEnv, "<packet_merchant>", data.PaketMerchant)
	base64env := base64.StdEncoding.EncodeToString([]byte(replaceEnv))
	xmlCred, err := os.Open("storage/src/jenkins/credentials.xml")
	if err != nil {
		return "", err
	}
	defer xmlCred.Close()
	credData, _ := io.ReadAll(xmlCred)
	err = xml.Unmarshal(credData, &data.JenkinsCredentials.JenFiles)
	if err != nil {
		return "", err
	}
	data.JenkinsCredentials.JenFiles.Id = fmt.Sprintf("client_env_%s", data.Username)
	data.JenkinsCredentials.JenFiles.Filename = ".env"
	data.JenkinsCredentials.JenFiles.SecretBytes = base64env
	data.JenkinsCredentials.JenFiles.Description = fmt.Sprintf("Environment for client %s", data.Username)
	cred, err := connect.AddCredentials(data, "file")
	if err != nil {
		return "", err
	}
	defer cred.Body.Close()
	if cred.StatusCode == http.StatusConflict {
		return "", errors.New("credentials already exists")
	} else if cred.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", cred.StatusCode)
	}
	// Open Template and Apply
	xmlFile, err := os.Open("storage/src/jenkins/config.xml")
	if err != nil {
		return "", err
	}
	defer xmlFile.Close()
	byteValue, _ := io.ReadAll(xmlFile)
	err = xml.Unmarshal(byteValue, &data.JenkinsItem)
	if err != nil {
		return "", err
	}
	data.JenkinsItem.Description = data.Description
	for x, v := range data.JenkinsItem.BuildWrappers {
		for i, j := range v.SecretBinding.FileBindings {
			if j.Variable == "ENV" {
				data.JenkinsItem.BuildWrappers[x].SecretBinding.FileBindings[i].CredentialsId = data.JenkinsCredentials.JenFiles.Id
				break
			}
		}
	}
	for i, v := range data.JenkinsItem.Builders {
		if strings.Contains(v.Command, "client-x") {
			newName := strings.ReplaceAll(
				v.Command, "client-x", fmt.Sprintf("client-%s", data.Username),
			)
			data.JenkinsItem.Builders[i].Command = newName
		}
		if strings.Contains(v.Command, "{{username}}") {
			apkName := strings.ReplaceAll(
				v.Command, "{{username}}", data.Username,
			)
			data.JenkinsItem.Builders[i].Command = apkName
		}
	}
	job, _ := connect.GetJobOperation(data, "DetailStatus")
	if job.StatusCode == http.StatusOK {
		return "job is exist", nil
	} else if job.StatusCode == http.StatusNotFound {
		_, err = connect.AddItem(ctx, data)
		if err != nil {
			return "", err
		}
	}
	defer job.Body.Close()

	build, err := connect.PostJobOperation(data, "Build")
	if err != nil {
		return "", err
	}
	defer build.Body.Close()
	if build.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("unexpected status code: %v", build.StatusCode)
	}
	// Check Job is exist
	result, err := connect.GetJobOperation(data, "AllJob")
	if err != nil {
		return "", err
	}
	defer result.Body.Close()
	if result.StatusCode != http.StatusOK {
		log.Printf("status code: %d with result:\n----------------------------\n%s\n----------------------------\n", result.StatusCode, result.Body)
		return "", fmt.Errorf("unexpected status code: %v", result.StatusCode)
	}
	allJob := jenkins.JenkinsJob{}
	b, _ := io.ReadAll(result.Body)
	err = json.Unmarshal(b, &allJob)
	if err != nil {
		return "", err
	}
	if len(allJob.Jobs) == 0 {
		return "", errors.New("job not found")
	}
	var jobIs string
	for _, v := range allJob.Jobs {
		if v.Name == data.Username {
			jobIs = "found"
			break
		}
	}
	if jobIs != "found" {
		return "", errors.New("job not found")
	}
	timeStart := time.Now()
	for {
		status, err := connect.GetJobOperation(data, "Status")
		if err != nil {
			return "", err
		}
		defer status.Body.Close()
		if status.StatusCode == http.StatusOK {
			break
		}
		time.Sleep(1 * time.Second)
		if time.Since(timeStart).Seconds() > 30 {
			break
		}
	}

	newStatus := jenkins.JenkinsJob{}
	getStatus, err := connect.GetJobOperation(data, "AllJob")
	if err != nil {
		return "", err
	}
	defer getStatus.Body.Close()
	b, err = io.ReadAll(getStatus.Body)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(b, &newStatus)
	if err != nil {
		return "", err
	}
	var jobStatus string
	for _, v := range newStatus.Jobs {
		if v.Name == data.Username {
			jobStatus = v.Color
			break
		}
	}
	var realStatus string
	if jobStatus == "notbuilt" {
		realStatus = "not build"
	} else if jobStatus == "notbuilt_anime" {
		realStatus = "build in proccess"
	} else if jobStatus == "blue" {
		realStatus = "success"
	} else if jobStatus == "red" {
		realStatus = "failed"
	} else {
		realStatus = "unknown"
	}

	return fmt.Sprintf("success deploy jenkins with status %s", realStatus), nil
}
