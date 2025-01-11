package cloudflare_api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dotcreep/go-automate-deploy/internal/utils"
)

//	@Summary		Get domain is available
//	@Description	Get domain is available or unavailable, only support for cloudflare providers
//	@Tags			Domain
//	@Accept			json
//	@Produce		json
//	@Security		X-Token
//	@Param			body	body		RequestInput				true	"User Data"
//	@Success		200		{object}	utils.Success				"Success"
//	@Success		302		{object}	utils.FoundSuccess			"Found"
//	@Failure		400		{object}	utils.BadRequest			"Bad request"
//	@Failure		500		{object}	utils.InternalServerError	"Internal server error"
//	@Router			/api/v1/domain/check [post]
func GetBasedomainRegisteredStatus(w http.ResponseWriter, r *http.Request) {
	connect := newCloudflare()
	ctx, cancel := utils.Cfgx{}.LongTimeout()
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
		d, ok := data["domain"].(string)
		if !ok {
			log.Println("domain is required")
			Json.NewResponse(false, w, nil, "domain is required", http.StatusBadRequest, nil)
			return
		}
		domain = d
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

	zone, err := connect.GetZone(ctx, baseDomain)
	if err != nil {
		log.Println(err)
		Json.NewResponse(false, w, nil, "failed get domain zone", http.StatusNotFound, err.Error())
		return
	}
	res, err := connect.GetDNSRecord(ctx, zone.ID)
	if err != nil {
		log.Println(err)
		Json.NewResponse(false, w, res, "failed get dns record", http.StatusNotFound, err.Error())
		return
	}
	for _, v := range res {
		if v == domain {
			log.Printf("%s already used", v)
			Json.NewResponse(false, w, nil, "domain is unavailable", http.StatusFound, "unavailable")
			return
		}
	}
	Json.NewResponse(true, w, "available", "domain is available", http.StatusOK, nil)
}
