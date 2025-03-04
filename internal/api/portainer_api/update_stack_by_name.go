package portainer_api

import (
	"encoding/json"
	"net/http"

	"github.com/dotcreep/go-automate-deploy/internal/utils"
)

type RequestInputUpdateStack struct {
	Name string `json:"name" form:"name" example:"mystack"`
}

// @Summary		Update stack by name
// @Description	You can update stack only using name of stack
// @Tags			System
// @Accept			json
// @Produce		json
//
// @Security		X-Token
//
// @Param			body	body		RequestInputUpdateStack		true	"Body Input"
// @Success		200		{object}	utils.Success				"Success"
// @Failure		400		{object}	utils.BadRequest			"Bad request"
// @Failure		500		{object}	utils.InternalServerError	"Internal server error"
// @Router			/api/v1/system/update [post]
func UpdateStackByName(w http.ResponseWriter, r *http.Request) {
	Json := utils.Json{}
	connect := newPortainer()
	ctx, cancel := utils.Cfgx{}.LongTimeout()
	defer cancel()
	input := &RequestInputUpdateStack{}
	if r.Header.Get("Content-Type") != "application/json" {
		input.Name = r.FormValue("name")
	} else {
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			Json.NewResponse(false, w, nil, "failed parsing body", http.StatusInternalServerError, err.Error())
			return
		}
	}
	if input.Name == "" {
		Json.NewResponse(false, w, nil, "name is required", http.StatusBadRequest, nil)
		return
	}
	res, err := connect.UpdateStackByName(ctx, input.Name)
	if err != nil {
		Json.NewResponse(false, w, res, "gagal mendapatkan stack", http.StatusOK, err.Error())
		return
	}
	Json.NewResponse(true, w, res, "berhasil mendapatkan stack", http.StatusOK, nil)
}
