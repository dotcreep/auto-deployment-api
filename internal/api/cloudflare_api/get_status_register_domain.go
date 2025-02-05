package cloudflare_api

import (
	"encoding/json"
	"net/http"

	"github.com/dotcreep/go-automate-deploy/internal/utils"
)

// @Summary		Get status register domain
// @Description	Check status register domain with status `pending` and `active`
// @Tags			Domain
// @Accept			json
// @Produce		json
// @Security		X-Token
// @Param			body	body		RequestInput				true	"User Data"
// @Success		200		{object}	utils.Success				"Success"
// @Failure		400		{object}	utils.BadRequest			"Bad request"
// @Failure		500		{object}	utils.InternalServerError	"Internal server error"
// @Router			/api/v1/domain/register-status [post]
func StatusRegisterDomain(w http.ResponseWriter, r *http.Request) {
	connect := newCloudflare()
	ctx, cancel := utils.Cfgx{}.DefaultTimeout()
	defer cancel()
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
		domain = data["domain"].(string)
	}
	if domain == "" {
		Json.NewResponse(false, w, nil, "domain is required", http.StatusBadRequest, nil)
		return
	}
	baseDomain := utils.GetBaseDomain(domain)
	if baseDomain == "" {
		Json.NewResponse(false, w, nil, "domain is invalid", http.StatusInternalServerError, nil)
		return
	}
	res, err := connect.GetZone(ctx, baseDomain)
	if err != nil {
		Json.NewResponse(true, w, "unregistered", "domain belum terdaftar", http.StatusOK, err.Error())
		return
	}

	Json.NewResponse(true, w, res.Status, "berhasil mendapatkan status domain", http.StatusOK, nil)
	// result is "pending" and "active"
}
