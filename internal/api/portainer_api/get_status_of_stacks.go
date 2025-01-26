package portainer_api

import (
	"encoding/json"
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
			Json.NewResponse(false, w, nil, "failed parse body data", http.StatusInternalServerError, err.Error())
			return
		}
		u, ok := data["username"].(string)
		if !ok {
			Json.NewResponse(false, w, nil, "username is required", http.StatusBadRequest, nil)
			return
		}
		username = u
	}
	if username == "" {
		Json.NewResponse(false, w, nil, "username is required", http.StatusBadRequest, nil)
		return
	}
	yamlConf, err := utils.Open()
	if err != nil {
		Json.NewResponse(false, w, nil, "failed open config.yml", http.StatusInternalServerError, err.Error())
		return
	}
	resp, err := connect.Endpoint(yamlConf.Portainer.EndpointId)
	if err != nil {
		Json.NewResponse(false, w, nil, "failed get endpoint", http.StatusInternalServerError, err.Error())
		return
	}
	var responseEndpoint portainer.ResponseContainer
	err = json.NewDecoder(resp.Body).Decode(&responseEndpoint)
	if err != nil {
		Json.NewResponse(false, w, nil, "failed decode response", http.StatusInternalServerError, err.Error())
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

	if respData.WebActive == 0 && respData.DBActive == 0 {
		Json.NewResponse(false, w, "failed", "berhasil mendapatkan status", http.StatusOK, "system not running")
		return
	}
	if respData.WebActive == yamlConf.Portainer.MaxClientWeb && respData.DBActive == yamlConf.Portainer.MaxClientDB {
		Json.NewResponse(true, w, "ready", "berhasil mendapatkan status", http.StatusOK, nil)
		return
	}

	Json.NewResponse(true, w, "process", "berhasil mendapatkan status", http.StatusOK, nil)
}
