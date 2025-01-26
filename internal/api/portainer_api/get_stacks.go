package portainer_api

import (
	"net/http"

	"github.com/dotcreep/go-automate-deploy/internal/utils"
)

// @Summary		Get all system stack
// @Description	Get all stack from portainer
// @Tags			System
// @Accept			json
// @Produce		json
// @Security		X-Token
// @Success		200	{object}	utils.Success				"Success"
// @Failure		400	{object}	utils.BadRequest			"Bad request"
// @Failure		500	{object}	utils.InternalServerError	"Internal server error"
// @Router			/api/v1/system/stack [get]
func GetStack(w http.ResponseWriter, r *http.Request) {
	connect := newPortainer()
	Json := utils.Json{}
	ctx, cancel := utils.Cfgx{}.LongTimeout()
	defer cancel()
	res, err := connect.GetStack(ctx)
	if err != nil {
		Json.NewResponse(false, w, res, "gagal mendapatkan stack", http.StatusOK, err.Error())
		return
	}
	Json.NewResponse(true, w, res, "berhasil mendapatkan stack", http.StatusOK, nil)
}
