package deploy_api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/dotcreep/go-automate-deploy/internal/service/cloudflare"
	"github.com/dotcreep/go-automate-deploy/internal/service/jenkins"
	"github.com/dotcreep/go-automate-deploy/internal/utils"
)

type RequestInputAddDomain struct {
	Domain   string `json:"domain" example:"example.com"`
	Username string `json:"username" example:"exampleusername"`
}

// @Summary		Undeploy user data
// @Description	Remove all data of user by username and domain used
// @Tags			Deploy
// @Accept			json
// @Produce		json
// @Security		X-Token
// @Param			body	body		RequestInputAddDomain		true	"Body"
// @Success		200		{object}	utils.Success				"Success"
// @Failure		400		{object}	utils.BadRequest			"Bad request"
// @Failure		500		{object}	utils.InternalServerError	"Internal server error"
// @Router			/api/v1/deploy/remove [delete]
func Undeploy(w http.ResponseWriter, r *http.Request) {
	Json := utils.Json{}
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cfg, err := utils.Open()
	if err != nil {
		Json.NewResponse(false, w, nil, "internal server error", http.StatusInternalServerError, err.Error())
		return
	}
	data := struct {
		Domain   string `json:"domain"`
		Username string `json:"username"`
	}{}
	if r.Header.Get("Content-Type") != "application/json" {
		data.Domain = r.FormValue("domain")
		data.Username = r.FormValue("username")
	} else {
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			Json.NewResponse(false, w, nil, "failed parse body data", http.StatusInternalServerError, err.Error())
			return
		}
	}

	if data.Domain == "" {
		Json.NewResponse(false, w, nil, "domain is required", http.StatusBadRequest, nil)
		return
	}
	if data.Username == "" {
		Json.NewResponse(false, w, nil, "username is required", http.StatusBadRequest, nil)
		return
	}

	removal := struct {
		Cloudflare struct {
			DNSRecords          bool `json:"dns_records"`
			TunnelConfiguration bool `json:"tunnel_configuration"`
		} `json:"cloudflare"`
		Portainer struct {
			Stacks        bool `json:"stacks"`
			DataDirectory bool `json:"data_directory"`
		} `json:"portainer"`
		Jenkins struct {
			Job         bool `json:"job"`
			Credentials bool `json:"credentials"`
		} `json:"jenkins"`
	}{}
	cfinput := &cloudflare.Subdomains{
		Domain:   data.Domain,
		TunnelID: cfg.Cloudflare.TunnelID,
	}
	_, err = connectcf.DeleteDomainDNSRecords(ctx, cfinput)
	if err != nil {
		log.Println(err)
		removal.Cloudflare.DNSRecords = false
	} else {
		removal.Cloudflare.DNSRecords = true
	}
	_, err = connectcf.DeleteDomainFromTunnelConfiguration(ctx, cfinput)
	if err != nil {
		log.Println(err)
		removal.Cloudflare.TunnelConfiguration = false
	} else {
		removal.Cloudflare.TunnelConfiguration = true
	}

	_, err = connectpt.DeleteStackByName(ctx, data.Username)
	if err != nil {
		log.Println(err)
		removal.Portainer.Stacks = false
	} else {
		removal.Portainer.Stacks = true
	}

	err = connectpt.RemoveClientDirectory(data.Username)
	if err != nil {
		log.Println(err)
		removal.Portainer.DataDirectory = false
	} else {
		removal.Portainer.DataDirectory = true
	}
	jkdata := jenkins.JenkinsData{}
	jkdata.Name = data.Username
	jkdata.DomainCredentials = cfg.Jenkins.DomainCredentials
	_, err = connectjen.DeleteJob(ctx, &jkdata)
	if err != nil {
		log.Println(err)
		removal.Jenkins.Job = false
	} else {
		removal.Jenkins.Job = true
	}

	_, err = connectjen.DeleteCredential(ctx, &jkdata)
	if err != nil {
		log.Println(err)
		removal.Jenkins.Credentials = false
	} else {
		removal.Jenkins.Credentials = true
	}
	Json.NewResponse(true, w, removal, "removal user success", http.StatusOK, nil)
}
