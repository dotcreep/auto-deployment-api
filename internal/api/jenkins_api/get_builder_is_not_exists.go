package jenkins_api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/dotcreep/go-automate-deploy/internal/utils"
)

//	@Summary		True if item is not exist
//	@Description	True if item is not exist
//	@Tags			Mobile
//	@Accept			json
//	@Produce		json
//	@Security		X-Token
//	@Param			body	body		RequestInput				true	"Body Input"
//	@Success		200		{object}	utils.Success				"Success"
//	@Success		302		{object}	utils.FoundFail				"Found"
//	@Failure		400		{object}	utils.BadRequest			"Bad request"
//	@Failure		500		{object}	utils.InternalServerError	"Internal server error"
//	@Router			/api/v1/mobile/is-not-exists [post]
func GetBuilderIsNotExists(w http.ResponseWriter, r *http.Request) {
	connect := newJenkins()
	Json := utils.Json{}
	ctx, cancel := utils.Cfgx{}.DefaultTimeout()
	defer cancel()
	// Get input user
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

	var requestData map[string]interface{}
	err = json.Unmarshal(body, &requestData)
	if err != nil {
		Json.NewResponse(false, w, nil, err.Error(), http.StatusBadRequest, "failed parse body data")
		return
	}

	name, ok := requestData["username"].(string)
	if !ok {
		Json.NewResponse(false, w, nil, "username is required", http.StatusBadRequest, nil)
		return
	}
	_, err = connect.GetItemIsNotExists(ctx, name)
	if err != nil {
		Json.NewResponse(false, w, nil, "item name is exist", http.StatusBadRequest, true)
		return
	}
	_, err = connect.GetCredentialsIsNotExists(ctx, name)
	if err != nil {
		Json.NewResponse(false, w, nil, "credential id is exist", http.StatusBadRequest, true)
		return
	}
	Json.NewResponse(true, w, true, "item is ready for create", http.StatusOK, nil)
}
