package cloudflare_api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dotcreep/go-automate-deploy/internal/service/cloudflare"
	"github.com/dotcreep/go-automate-deploy/internal/utils"
)

//	@Summary		Get status register domain
//	@Description	Check status register domain with status `pending` and `active`
//	@Tags			Domain
//	@Accept			json
//	@Produce		json
//	@Security		X-Token
//	@Param			body	body		RequestInput				true	"User Data"
//	@Success		200		{object}	utils.SuccessDeploy			"Success"
//	@Success		302		{object}	utils.FoundSuccess			"Found"
//	@Failure		400		{object}	utils.BadRequest			"Bad request"
//	@Failure		500		{object}	utils.InternalServerError	"Internal server error"
//	@Router			/api/v1/domain/register-status [post]
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
			log.Println(err)
			Json.NewResponse(false, w, nil, "failed parse body data", http.StatusInternalServerError, err)
			return
		}
		domain = data["domain"].(string)
	}
	if domain == "" {
		log.Println("domain is required")
		Json.NewResponse(false, w, nil, "domain is required", http.StatusBadRequest, nil)
		return
	}
	slDom := strings.Split(domain, ".")
	var baseDomain string
	if len(slDom) > 2 {
		domainsub := []string{"co", "biz,"}
		for _, v := range domainsub {
			if slDom[len(slDom)-2] == v {
				baseDomain = fmt.Sprintf("%s.%s.%s", slDom[len(slDom)-3], slDom[len(slDom)-2], slDom[len(slDom)-1])
				break
			}
			baseDomain = fmt.Sprintf("%s.%s", slDom[len(slDom)-2], slDom[len(slDom)-1])
		}
		//baseDomain = fmt.Sprintf("%s.%s", slDom[len(slDom)-2], slDom[len(slDom)-1])
	} else {
		baseDomain = domain
	}
	data := &cloudflare.Subdomains{
		Domain: baseDomain,
	}
	res, err := connect.StatusDomain(ctx, data)
	if err != nil {
		Json.NewResponse(false, w, nil, err.Error(), http.StatusNotFound, res)
		return
	}

	Json.NewResponse(true, w, res, "success get status domain", http.StatusOK, nil)
	// result is "pending" and "active"
}
