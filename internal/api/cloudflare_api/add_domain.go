package cloudflare_api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/dotcreep/go-automate-deploy/internal/service/cloudflare"
	"github.com/dotcreep/go-automate-deploy/internal/utils"
)

type RequestInputAddDomain struct {
	Domain   string `json:"domain" example:"example.com"`
	Username string `json:"username" example:"exampleusername"`
}

// @Summary		Add domain to cloudflare
// @Description	Add domain to cloudflare tunnel and dns record
// @Tags			Domain
// @Accept			json
// @Produce		json
// @Security		X-Token
// @Param			body	body		RequestInputAddDomain		true	"Body"
// @Success		200		{object}	utils.Success				"Success"
// @Failure		400		{object}	utils.BadRequest			"Bad request"
// @Failure		500		{object}	utils.InternalServerError	"Internal server error"
// @Router			/api/v1/domain/add [post]
func AddDomain(w http.ResponseWriter, r *http.Request) {
	connect := newCloudflare()
	Json := utils.Json{}
	yamlConfig, err := utils.Open()
	if err != nil {
		Json.NewResponse(false, w, nil, "config.yml not found", http.StatusInternalServerError, err.Error())
		return
	}
	ctx, cancel := utils.Cfgx{}.DefaultTimeout()
	defer cancel()
	if r.Header.Get("Content-Type") != "application/json" {
		Json.NewResponse(false, w, nil, "Content-Type is not application/json", http.StatusBadRequest, nil)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		Json.NewResponse(false, w, nil, err.Error(), http.StatusBadRequest, nil)
		return
	}
	defer r.Body.Close()

	var requestData RequestInputAddDomain
	err = json.Unmarshal(body, &requestData)
	if err != nil {
		Json.NewResponse(false, w, nil, err.Error(), http.StatusBadRequest, nil)
		return
	}
	if requestData.Domain == "" {
		Json.NewResponse(false, w, nil, "domain is required", http.StatusBadRequest, "domain input is empty")
		return
	}
	if requestData.Username == "" {
		Json.NewResponse(false, w, nil, "username is required", http.StatusBadRequest, "username input is empty")
		return
	}
	data := &cloudflare.Subdomains{
		Service:      fmt.Sprintf("https://web_%s:443", requestData.Username),
		Domain:       requestData.Domain,
		TunnelID:     yamlConfig.Cloudflare.TunnelID,
		Path:         "",
		LoadBalancer: false,
	}
	// Parse domain if subdomain or domain
	parts := strings.Split(data.Domain, ".")
	baseDomain := data.Domain
	var domainTypes string
	if len(parts) > 2 {
		domainTypes = "Subdomain"
		baseDomain = strings.Join(parts[len(parts)-2:], ".")
	} else {
		domainTypes = "Domain"
	}
	// For new Domain
	if domainTypes == "Domain" {
		_, err := connect.GetZone(ctx, baseDomain)
		if err != nil {
			_, errData := connect.Register(ctx, data)
			if errData != nil {
				Json.NewResponse(false, w, nil, "failed register domain", http.StatusInternalServerError, errData.Error())
				return
			}
		}
		zone, err := connect.GetZone(ctx, baseDomain)
		if err != nil {
			Json.NewResponse(false, w, nil, "failed get zone", http.StatusInternalServerError, err.Error())
			return
		}
		if zone.Status == "pending" {
			Json.NewResponse(false, w, nil, "domain is pending", http.StatusInternalServerError, "pending")
			return
		}
	}
	newResp := struct {
		Tunnel    string `json:"tunnel"`
		DNSRecord string `json:"dns-record"`
	}{}
	respTunnel, err := connect.AddDomainToTunnelConfiguration(ctx, data)
	if err != nil {
		Json.NewResponse(false, w, nil, err.Error(), http.StatusInternalServerError, err)
		return
	}
	respDNSRecord, err := connect.AddDNSRecord(ctx, data)
	if err != nil {
		Json.NewResponse(false, w, nil, err.Error(), http.StatusInternalServerError, err)
		return
	}
	newResp.Tunnel = respTunnel
	newResp.DNSRecord = respDNSRecord
	Json.NewResponse(true, w, newResp, "success add domain", http.StatusOK, nil)
}
