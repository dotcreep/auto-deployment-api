package external

import (
	"context"
	"io"
	"net/http"
)

type UserData struct {
	Superuser         string `json:"superuser"`
	PasswordSuperuser string `json:"password-superuser"`
	Merchant          string `json:"merchant"`
	PasswordMerchant  string `json:"password-merchant"`
}

func SendEmail(ctx context.Context, url string, data io.Reader) (*http.Response, error) {
	return nil, nil
}
