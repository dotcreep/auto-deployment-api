package portainer_api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/dotcreep/go-automate-deploy/internal/utils"
)

//	@Summary		True if stack is not exist
//	@Description	True if stack is not exist
//	@Tags			System
//	@Accept			json
//	@Produce		json
//	@Security		X-Token
//	@Param			body	body		RequestInput				true	"Body Input"
//	@Success		200		{object}	utils.Success				"Success"
//	@Success		302		{object}	utils.FoundFail				"Found"
//	@Failure		400		{object}	utils.BadRequest			"Bad request"
//	@Failure		500		{object}	utils.InternalServerError	"Internal server error"
//	@Router			/api/v1/system/is-not-exists [post]
func GetStackIsNotExists(w http.ResponseWriter, r *http.Request) {
	connect := newPortainer()
	Json := utils.Json{}
	ctx, cancel := utils.Cfgx{}.LongTimeout()
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

	res, err := connect.GetStack(ctx)
	if err != nil {
		Json.NewResponse(false, w, nil, "failed get stacks", http.StatusInternalServerError, err)
		return
	}

	var isExists bool = false
	for _, v := range res.Stacks {
		if v.Name == name {
			isExists = true
			break
		}
	}
	if isExists {
		Json.NewResponse(false, w, nil, "stack is exists", http.StatusFound, true)
		return
	}
	Json.NewResponse(true, w, true, "stack is ready for create", http.StatusOK, nil)
}
