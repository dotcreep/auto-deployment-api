package portainer_api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dotcreep/go-automate-deploy/internal/service/portainer"
	"github.com/dotcreep/go-automate-deploy/internal/utils"
)

// @Summary		Get status of stack
// @Description	Get status of stack
// @Tags			System
// @Accept			json
// @Produce		json
//
// @Security		X-Token
//
// @Param			body	body		RequestInput				true	"Body Input"
// @Success		200		{object}	utils.Success				"Success"
// @Success		302		{object}	utils.FoundSuccess			"Found"
// @Failure		400		{object}	utils.BadRequest			"Bad request"
// @Failure		500		{object}	utils.InternalServerError	"Internal server error"
// @Router			/api/v1/system/status [post]
func GetStatusOfStack(w http.ResponseWriter, r *http.Request) {
	connect := newPortainer()
	Json := utils.Json{}
	var username string
	if r.Header.Get("Content-Type") != "application/json" {
		username = r.FormValue("username")
	} else {
		var data map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			log.Println(err)
			Json.NewResponse(false, w, nil, "failed parse body data", http.StatusInternalServerError, err)
			return
		}
		u, ok := data["username"].(string)
		if !ok {
			log.Println("username is required")
			Json.NewResponse(false, w, nil, "username is required", http.StatusBadRequest, nil)
			return
		}
		username = u
	}
	if username == "" {
		log.Println("username is required")
		Json.NewResponse(false, w, nil, "username is required", http.StatusBadRequest, nil)
		return
	}
	yamlConf, err := utils.Open()
	if err != nil {
		log.Println(err)
		Json.NewResponse(false, w, nil, "failed open config.yml", http.StatusInternalServerError, nil)
		return
	}
	resp, err := connect.Endpoint(yamlConf.Portainer.EndpointId)
	if err != nil {
		log.Println(err)
		Json.NewResponse(false, w, nil, "failed get endpoint", http.StatusInternalServerError, err)
		return
	}
	var responseEndpoint portainer.ResponseContainer
	err = json.NewDecoder(resp.Body).Decode(&responseEndpoint)
	if err != nil {
		log.Println(err)
		Json.NewResponse(false, w, nil, "failed decode response", http.StatusInternalServerError, err)
		return
	}
	respData := struct {
		Namespace string
		State     string
		Status    string
		WebActive int
		DBActive  int
	}{}
	for _, v := range responseEndpoint.Snapshots {
		for _, x := range v.DockerSnapshotRaw.Containers {
			if x.State == "running" {
				if x.Labels.Namespace == username && strings.Contains(x.Labels.ServiceName, "web") {
					respData.WebActive++
				}
				if x.Labels.Namespace == username && strings.Contains(x.Labels.ServiceName, "db") {
					respData.DBActive++
				}
			}
		}
	}
	if respData.WebActive == 0 || respData.DBActive == 0 {
		newResp := fmt.Sprintf("container is not running, web running %d, db running %d", respData.WebActive, respData.DBActive)
		// Request by Dwi
		Json.NewResponse(false, w, "failed", newResp, http.StatusOK, respData)
		return
	}

	newResp := fmt.Sprintf("success get status, web running %d, db running %d", respData.WebActive, respData.DBActive)
	Json.NewResponse(true, w, "ready", newResp, http.StatusOK, nil)
}
