package cloudflare_api

import (
	"encoding/json"
	"net/http"

	"github.com/dotcreep/go-automate-deploy/internal/utils"
)

// @Summary		Get status domain
// @Description	Are domain is accessible or still cannot access
// @Tags			Domain
// @Accept			json
// @Produce		json
// @Security		X-Token
// @Param			body	body		RequestInput				true	"User Data"
// @Success		200		{object}	utils.Success				"Success"
// @Failure		400		{object}	utils.BadRequest			"Bad request"
// @Failure		500		{object}	utils.InternalServerError	"Internal server error"
// @Router			/api/v1/domain/status [post]
func StatusDomain(w http.ResponseWriter, r *http.Request) {
	getStatus := utils.GetStatusOK
	Json := utils.Json{}
	var domain string
	if r.Header.Get("Content-Type") != "application/json" {
		domain = r.FormValue("domain")
	} else {
		var data map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			Json.NewResponse(false, w, nil, "failed parse body data", http.StatusInternalServerError, err.Error())
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
	ctx, cancel := utils.Cfgx{}.DefaultTimeout()
	defer cancel()
	resp, _ := getStatus(ctx, domain)
	if resp != http.StatusOK {
		Json.NewResponse(false, w, "fail", "domain tidak valid", http.StatusOK, "no such host")
		return
	}
	Json.NewResponse(true, w, "ready", "domain dapat diakses", http.StatusOK, nil)
}
