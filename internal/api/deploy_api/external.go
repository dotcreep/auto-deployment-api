package deploy_api

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dotcreep/go-automate-deploy/internal/service/external"
)

func connectExternal() (*external.External, error) {
	e := external.External{}
	secrets := &Secrets{}
	secrets.GetSecret()
	connect, err := e.NewExternal()
	if err != nil {
		return nil, err
	}
	return connect, nil
}

func RegisterUser(data *external.External, role string) (string, error) {
	connect, err := connectExternal()
	if err != nil {
		return "", err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if role == "" {
		return "", errors.New("role is required")
	}
	if data.User.Username == "" {
		return "", errors.New("username is required")
	}
	if data.User.Password == "" {
		return "", errors.New("password is required")
	}
	_, err = connect.AddUser(ctx, role)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("success register %s", role), nil
}
