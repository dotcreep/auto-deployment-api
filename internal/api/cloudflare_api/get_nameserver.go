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

// @Summary		Get cloudflare nameserver
// @Description	Get cloudflare nameserver of base domain
// @Tags			Domain
// @Accept			json
// @Produce		json
// @Security		X-Token
// @Param			body	body		RequestInput				true	"Body Input"
// @Success		200		{object}	utils.Success				"Success"
// @Success		302		{object}	utils.FoundSuccess			"Found"
// @Failure		400		{object}	utils.BadRequest			"Bad request"
// @Failure		500		{object}	utils.InternalServerError	"Internal server error"
// @Router			/api/v1/domain/nameserver [post]
func GetNameserver(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := utils.Cfgx{}.LongTimeout()
	defer cancel()
	connect := newCloudflare()
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

	_, err := connect.GetZone(ctx, baseDomain)
	if err != nil {
		log.Println(err)
		Json.NewResponse(false, w, nil, "failed get domain zone", http.StatusNotFound, err.Error())
		return
	}
	data := &cloudflare.Subdomains{
		Domain: baseDomain,
	}
	res, err := connect.GetNameserver(ctx, data)
	if err != nil {
		log.Println(err)
		Json.NewResponse(false, w, err.Error(), "nameserver is up to date", http.StatusOK, nil)
		return
	}
	Json.NewResponse(true, w, res, "success get nameserver", http.StatusOK, nil)
}
