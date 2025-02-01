package deploy_api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/dotcreep/go-automate-deploy/internal/cli"
	"github.com/dotcreep/go-automate-deploy/internal/generator"
	"github.com/dotcreep/go-automate-deploy/internal/service"
	"github.com/dotcreep/go-automate-deploy/internal/service/cloudflare"
	"github.com/dotcreep/go-automate-deploy/internal/service/jenkins"
	"github.com/dotcreep/go-automate-deploy/internal/service/portainer"
	"github.com/dotcreep/go-automate-deploy/internal/service/system"
	"github.com/dotcreep/go-automate-deploy/internal/utils"
)

type Secrets struct {
	Cloudflare   CloudflareSecrets
	Portainer    PortainerSecrets
	Jenkins      JenkinsSecrets
	UserRegister UserRegister
}

type CloudflareSecrets struct {
	Token string
	Email string
}

type PortainerSecrets struct {
	Username string
	Password string
}

type JenkinsSecrets struct {
	Username string
	Password string
}

type UserRegister struct {
	Username string
	Password string
}

// Response is the specific payload returned by this API.
//
//	@Cloudflare	contains the status of the Cloudflare service.
//	@Portainer	contains the status of the Portainer service.
//	@Jenkins	contains the status of the Jenkins service.
type Response struct {
	Cloudflare string `json:"cloudflare"`
	Portainer  string `json:"portainer"`
	Jenkins    string `json:"jenkins"`
}

type RequestInput struct {
	Domain       string `json:"domain" example:"example.com"`
	Username     string `json:"username" example:"exampleusername"`
	Email        string `json:"email" example:"sample@example.com"`
	MerchantName string `json:"merchant_name" example:"Example Name"`
}

func (s *Secrets) GetSecret() (*Secrets, error) {
	config, err := utils.Open()
	if err != nil {
		return nil, err
	}

	s.Cloudflare.Token = config.Cloudflare.Key
	s.Cloudflare.Email = config.Cloudflare.Email
	s.Portainer.Username = config.Portainer.Username
	s.Portainer.Password = config.Portainer.APIKey
	s.Jenkins.Username = config.Jenkins.Username
	s.Jenkins.Password = config.Jenkins.APIKey

	return s, nil
}

// @Summary		Deploy All Third Party Environment
// @Description	Deployment to Cloudflare, Portainer, and Jenkins (Support rollback action if failed)
// @Tags			Deploy
// @Accept			json
// @Produce		json
// @Security		X-Token
// @Param			body	body		RequestInput				true	"User Data"
// @Success		200		{object}	utils.SuccessDeploy			"Success"
// @Failure		400		{object}	utils.BadRequest			"Bad request"
// @Failure		500		{object}	utils.InternalServerError	"Internal server error"
// @Router			/api/v1/deploy/start [post]
func Deploy(w http.ResponseWriter, r *http.Request) {
	/**
	0. Must check of value first
	1. Check ram if more than 80% cannot create user
	2. Check exists in thord party
	3. If fail always delete or rollback
	4. Success will pass
	*/

	// 0. Variables and Initialize
	Json := utils.Json{}
	yamlConfig, err := utils.Open()
	if err != nil {
		Json.NewResponse(false, w, nil, "config.yml not found", http.StatusInternalServerError, err.Error())
		return
	}
	ctx, cancel := utils.Cfgx{}.LongTimeout()
	defer cancel()
	// 1. Check RAM
	percent, err := system.CheckFreeRam()
	if err != nil {
		Json.NewResponse(false, w, nil, "unable check ram", http.StatusInternalServerError, err.Error())
		return
	}
	if percent > 80 {
		Json.NewResponse(false, w, nil, "cannot create user", http.StatusInternalServerError, "ram usage more than 80%")
		return
	}
	// 2. Check input
	if r.Header.Get("Content-Type") != "application/json" {
		Json.NewResponse(false, w, nil, "Content-Type is not application/json", http.StatusBadRequest, nil)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		Json.NewResponse(false, w, nil, "unable read request", http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	var requestData RequestInput
	err = json.Unmarshal(body, &requestData)
	if err != nil {
		Json.NewResponse(false, w, nil, "unable unmarshal request", http.StatusBadRequest, err.Error())
		return
	}

	// 1. Always check input
	/*
		1. username
		2. email
		3. merchant_name
		4. domain
	*/
	usernameClient := requestData.Username
	if usernameClient == "" {
		Json.NewResponse(false, w, nil, "username is required", http.StatusBadRequest, nil)
		return
	}
	emailClient := requestData.Email
	if emailClient == "" {
		Json.NewResponse(false, w, nil, "email is required", http.StatusBadRequest, nil)
		return
	}
	merchantName := requestData.MerchantName
	if merchantName == "" {
		Json.NewResponse(false, w, nil, "merchant_name is required", http.StatusBadRequest, nil)
		return
	}
	domainClient := requestData.Domain
	if domainClient == "" {
		Json.NewResponse(false, w, nil, "domain is required", http.StatusBadRequest, nil)
		return
	}
	tunnelId := yamlConfig.Cloudflare.TunnelID
	if tunnelId == "" {
		Json.NewResponse(false, w, nil, "tunnel_id is required", http.StatusBadRequest, nil)
		return
	}
	domainPath := "" // requestData["domain_path"].(string)
	serviceClient := fmt.Sprintf("https://web_%s:443", usernameClient)
	loadBalancer := false
	apiURLClient := yamlConfig.Jenkins.APIURL
	if apiURLClient == "" {
		Json.NewResponse(false, w, nil, "api_url is required", http.StatusBadRequest, nil)
		return
	}

	// 2. Third party checker
	// 2.1. Check domain
	// Search by tunnel id configuration first
	// res, err := connectcf.GetDomainFromTunnelConfiguration(ctx, domainClient, yamlConfig.Cloudflare.TunnelID)
	// if err != nil {
	// 	Json.NewResponse(false, w, nil, err.Error(), http.StatusInternalServerError, res)
	// 	return
	// }

	// 2.2. Create connection thirdparty
	connectcf, err := connectCloudflare()
	if err != nil {
		Json.NewResponse(false, w, nil, "failed connect cloudflare", http.StatusInternalServerError, err.Error())
		return
	}
	connectpt, err := connectPortainer()
	if err != nil {
		Json.NewResponse(false, w, nil, "failed connect portainer", http.StatusInternalServerError, err.Error())
		return
	}
	connectjen, err := connectJenkins()
	if err != nil {
		Json.NewResponse(false, w, nil, "failed connect jenkins", http.StatusInternalServerError, err.Error())
		return
	}

	// -------------------------------------------------------------------------------
	// Create Directory
	// -------------------------------------------------------------------------------
	pathListEnvCreated := []string{
		"web",
		"api",
	}
	for _, path := range pathListEnvCreated {
		err = service.CreateDir("/nfs/environment", usernameClient, path)
		if err != nil {
			Json.NewResponse(false, w, nil, "failed create directory", http.StatusInternalServerError, err.Error())
			return
		}
	}
	pathListCreated := []string{
		"database",
		"web-merchant/logs",
		"web-merchant/framework/sessions",
		"web-merchant/framework/views",
		"web-merchant/framework/cache",
		"web-merchant/framework/storage",
	}
	for _, path := range pathListCreated {
		err = service.CreateDir("/nfs/client", usernameClient, path)
		if err != nil {
			Json.NewResponse(false, w, nil, "failed create directory", http.StatusInternalServerError, err.Error())
			return
		}
	}

	// ------------------------------------------------------------------
	// Deploy Cloudflare
	// ------------------------------------------------------------------
	domain := domainClient
	tunnelID := tunnelId
	cloudflare := &cloudflare.Subdomains{
		Domain:       domain,
		TunnelID:     tunnelID,
		Path:         domainPath,
		LoadBalancer: loadBalancer,
		Service:      serviceClient,
	}
	resCloudflare, err := DeployCloudflare(ctx, cloudflare)
	if err != nil {
		defer connectcf.RollbackAddDomain(ctx, cloudflare)
		Json.NewResponse(false, w, nil, "failed add domain", http.StatusBadRequest, err.Error())
		return
	}
	responseCloudflare := resCloudflare
	// ------------------------------------------------------------------
	// Deploy Portainer
	// ------------------------------------------------------------------

	reverseDomain := strings.Split(domainClient, ".")
	reverseDomain = reverseDomain[len(reverseDomain)-2:]

	merch := strings.Split(merchantName, " ")
	if len(merch) > 2 {
		merchantName = fmt.Sprintf("%s %s", merch[0], merch[1])
	}
	genEnv := &cli.EGenerate{
		BaseWebURL: fmt.Sprintf("https://%s", domainClient),
		BaseAPIURL: apiURLClient,
		ClientName: usernameClient,
		Email:      emailClient,
		EGenerateAndroid: cli.EGenerateAndroid{
			WebHost:     domainClient,
			PackageName: fmt.Sprintf("%s.%s.%s", reverseDomain[0], reverseDomain[1], usernameClient),
			LabelApps:   merchantName,
		},
		WebApi:   yamlConfig.Config.ChatAPI,
		AppTitle: merchantName,
	}

	// Generate ENV for Client
	dataGenerator, err := cli.Generator(genEnv)
	if err != nil {
		Json.NewResponse(false, w, nil, "failed generate env", http.StatusInternalServerError, err.Error())
		return
	}

	// Copy ENV api to folder client
	// err = service.CopyFile(fmt.Sprintf("storage/dist/%s/api/.env", usernameClient), fmt.Sprintf("/nfs/environment/%s/api/.env", usernameClient))
	// if err != nil {
	// Json.NewResponse(false, w, nil, err.Error(), http.StatusInternalServerError, nil)
	// return
	// }
	// Copy ENV web to folder client
	err = service.CopyFile(fmt.Sprintf("storage/dist/%s/web/.env", usernameClient), fmt.Sprintf("/nfs/environment/%s/web/.env", usernameClient))
	if err != nil {
		Json.NewResponse(false, w, nil, "failed copy env", http.StatusInternalServerError, err.Error())
		return
	}

	passDBRoot := dataGenerator.PasswordDB
	passDBWeb := dataGenerator.PasswordDB
	passDBAPI := generator.Password(10)
	custom := &portainer.CustomInput{
		Name:              usernameClient,
		WebImageContainer: yamlConfig.Config.ImageWeb,
		APIImageContainer: yamlConfig.Config.ImageAPI,
		APIURL:            apiURLClient,
		DBHost:            fmt.Sprintf("db_%s", usernameClient),
		DBPort:            "5432",
		DBRootUser:        "postgres",
		DBRootPass:        passDBRoot,
		DBWebName:         fmt.Sprintf("web_%s", usernameClient),
		DBWebUser:         fmt.Sprintf("web_%s", usernameClient),
		DBWebPass:         passDBWeb,
		DBAPIName:         fmt.Sprintf("api_%s", usernameClient),
		DBAPIUser:         fmt.Sprintf("api_%s", usernameClient),
		DBAPIPass:         passDBAPI,
		DockerPath: portainer.DockerPath{
			Source: "/storage/src/docker",
			Dist:   fmt.Sprintf("/storage/dist/%s", usernameClient),
		},
	}

	resPortainer, err := DeployPortainer(custom)
	if err != nil {
		defer connectcf.RollbackAddDomain(ctx, cloudflare)
		defer connectpt.RollbackAddStack(ctx, usernameClient)
		Json.NewResponse(false, w, nil, "failed add container", http.StatusBadRequest, err.Error())
		return
	}
	responsePortainer := resPortainer

	// ------------------------------------------------------------------
	// Deploy Add User Merchant and User Admin Merchant
	// ------------------------------------------------------------------
	// 1. Create user admin
	// passSuperUser := generator.Password(10)
	// passUserMerchant := generator.Password(10)
	// statusRegister := struct {
	// 	Admin    string
	// 	Merchant string
	// }{}
	// externalInput := &external.External{}
	// externalInput.Config.Token = yamlConfig.Config.TokenRegis

	// // 2. Create admin merchant
	// externalInput.User.Username = fmt.Sprintf("admin_%s", usernameClient)
	// externalInput.User.Password = passSuperUser
	// externalInput.Config.URLSuperuser = yamlConfig.Config.RegisAdmin
	// resp, err := RegisterUser(externalInput, "admin")
	// if err != nil {
	// 	statusRegister.Admin = err.Error()
	// } else {
	// 	statusRegister.Admin = resp
	// }

	// // 3. Create user merchant
	// externalInput.User.Username = usernameClient
	// externalInput.User.Password = passUserMerchant
	// externalInput.Config.URLMerchant = yamlConfig.Config.RegisMerch
	// resp, err = RegisterUser(externalInput, "merchant")
	// if err != nil {
	// 	statusRegister.Merchant = err.Error()
	// } else {
	// 	statusRegister.Merchant = resp
	// }

	// 4. Send Email

	//-------------------------------------------------------------------
	// Deploy Jenkins
	// ------------------------------------------------------------------
	JenInput := &jenkins.JenkinsData{}
	JenInput.ID = usernameClient
	JenInput.Description = fmt.Sprintf("Environment for %s", usernameClient)
	JenInput.Name = usernameClient
	JenInput.DomainCredentials = yamlConfig.Jenkins.DomainCredentials
	JenInput.MerchantName = merchantName
	// JenInput.MerchantID = merchantId
	JenInput.APIURL = yamlConfig.Jenkins.APIURL
	resJenkins, err := DeployJenkins(ctx, JenInput, domainClient)
	if err != nil || resJenkins == "" {
		defer connectcf.RollbackAddDomain(ctx, cloudflare)
		defer connectpt.RollbackAddStack(ctx, usernameClient)
		defer connectjen.RollbackAddItem(ctx, JenInput)
		Json.NewResponse(false, w, nil, "failed add item", http.StatusBadRequest, err.Error())
		return
	}
	responseJenkins := resJenkins
	//------------------------------------------------------------------

	//------------------------------------------------------------------
	// Response
	allResponse := &Response{
		Cloudflare: responseCloudflare,
		Portainer:  responsePortainer,
		Jenkins:    responseJenkins,
	}
	Json.NewResponse(true, w, allResponse, "success deploy", http.StatusOK, nil)
}
