package jenkins_api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/dotcreep/go-automate-deploy/internal/utils"
)

// @Summary		Check status of mobile builder
// @Description	Check status of mobile builder with return 'success', 'no build', 'failed', 'unknown'
// @Tags			Mobile
// @Accept			json
// @Produce		json
// @Security		X-Token
// @Param			body	body		RequestInput				true	"Body Input"
// @Success		200		{object}	utils.Success				"Success"
// @Failure		400		{object}	utils.BadRequest			"Bad request"
// @Failure		500		{object}	utils.InternalServerError	"Internal server error"
// @Router			/api/v1/mobile/status [post]
func GetStatusOfItem(w http.ResponseWriter, r *http.Request) {
	// 1. Create new connection
	connect := newJenkins()
	Json := utils.Json{}
	ctx, cancel := utils.Cfgx{}.DefaultTimeout()
	defer cancel()

	// 2. Check input user
	if r.Header.Get("Content-Type") != "application/json" {
		Json.NewResponse(false, w, nil, "Content-Type is not application/json", http.StatusBadRequest, nil)
		return
	}

	// 3. Get input user
	body, err := io.ReadAll(r.Body)
	if err != nil {
		Json.NewResponse(false, w, nil, err.Error(), http.StatusBadRequest, nil)
		return
	}
	defer r.Body.Close()

	// 4. Unmarshal input user
	var requestData map[string]interface{}
	err = json.Unmarshal(body, &requestData)
	if err != nil {
		Json.NewResponse(false, w, nil, err.Error(), http.StatusBadRequest, "failed parse body data")
		return
	}

	if requestData["username"] == nil {
		Json.NewResponse(false, w, nil, "username is required", http.StatusBadRequest, nil)
		return
	}

	// 5. Get status of item
	status, err := connect.GetStatusOfItem(ctx, requestData["username"].(string))
	if err != nil {
		Json.NewResponse(false, w, nil, err.Error(), http.StatusNotFound, "fail")
		return
	}
	Json.NewResponse(true, w, status, "success get status of item", http.StatusOK, nil)
}
