package cloudflare_api

import (
	"encoding/json"
	"net/http"

	"github.com/dotcreep/go-automate-deploy/internal/utils"
)

// @Summary		True if domain is not exist
// @Description	True if domain is not exist
// @Tags			Domain
// @Accept			json
// @Produce		json
// @Security		X-Token
// @Param			body	body		RequestInput				true	"Body Input"
// @Success		200		{object}	utils.Success				"Success"
// @Success		302		{object}	utils.FoundFail				"Found"
// @Failure		400		{object}	utils.BadRequest			"Bad request"
// @Failure		500		{object}	utils.InternalServerError	"Internal server error"
// @Router			/api/v1/domain/is-not-exists [post]
func GetDomainIsNotExists(w http.ResponseWriter, r *http.Request) {
	connect := newCloudflare()
	Json := utils.Json{}
	yamlConfig, err := utils.Open()
	if err != nil {
		Json.NewResponse(false, w, nil, "failed load config", http.StatusInternalServerError, err.Error())
		return
	}
	ctx, cancel := utils.Cfgx{}.DefaultTimeout()
	defer cancel()
	if yamlConfig == nil {
		Json.NewResponse(false, w, nil, "config.yml not found", http.StatusInternalServerError, nil)
		return
	}
	var domain string
	if r.Header.Get("Content-Type") != "application/json" {
		domain = r.FormValue("domain")
	} else {
		var data map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			Json.NewResponse(false, w, nil, "failed parse body data", http.StatusInternalServerError, err)
			return
		}
		d, ok := data["domain"].(string)
		if !ok {
			Json.NewResponse(false, w, nil, "domain is required", http.StatusBadRequest, nil)
			return
		}
		domain = d
	}
	if domain == "" {
		Json.NewResponse(false, w, nil, "domain is required", http.StatusBadRequest, nil)
		return
	}
	// 1. Search from Tunnel Configuration
	res, err := connect.GetDomainFromTunnelConfiguration(ctx, domain, yamlConfig.Cloudflare.TunnelID)
	if err != nil {
		Json.NewResponse(false, w, res, err.Error(), http.StatusInternalServerError, err)
		return
	}
	var domainFound bool
	domainFound = false
	for _, v := range res.Result.Config.Ingress {
		if v.Hostname == domain {
			domainFound = true
			break
		}
	}
	if domainFound {
		Json.NewResponse(false, w, false, "domain sudah terpakai", http.StatusOK, "telah dipakai")
		return
	}

	// 2. Search from Zone
	baseDomain := utils.GetBaseDomain(domain)
	zone, err := connect.GetZone(ctx, baseDomain)
	if err != nil {
		Json.NewResponse(false, w, nil, err.Error(), http.StatusInternalServerError, err)
		return
	}

	dnsrecords, err := connect.GetDNSRecord(ctx, zone.ID)
	if err != nil {
		Json.NewResponse(false, w, nil, err.Error(), http.StatusInternalServerError, err)
		return
	}

	for _, v := range dnsrecords {
		if v == domain {
			domainFound = true
			break
		}
	}
	if domainFound {
		Json.NewResponse(false, w, false, "domain sudah terpakai", http.StatusOK, "telah dipakai")
		return
	}
	Json.NewResponse(true, w, true, "domain dapat dibuat", http.StatusOK, nil)
}
