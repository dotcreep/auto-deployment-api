package external

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/dotcreep/go-automate-deploy/internal/utils"
)

type External struct {
	User                  User
	Config                configExternal
	NewExternalConnection NewExternalConnection
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type configExternal struct {
	URLSuperuser string
	URLMerchant  string
	Token        string
}

type NewExternalConnection struct {
	URLAdmin    string
	URLMerchant string
	Headers     map[string]string
}

func (e *External) ExternalConfig() (*utils.YamlStruct, error) {
	config, err := utils.Open()
	if config == nil {
		return nil, err
	}
	return config, nil
}

func (e *External) NewExternal() (*External, error) {
	if e.Config.URLSuperuser == "" {
		return nil, errors.New("url superuser is required")
	}
	if e.Config.URLMerchant == "" {
		return nil, errors.New("url merchant is required")
	}
	if e.Config.Token == "" {
		return nil, errors.New("token is required")
	}

	return &External{
		NewExternalConnection: NewExternalConnection{
			URLAdmin:    e.Config.URLSuperuser,
			URLMerchant: e.Config.URLMerchant,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"User-Agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36",
				"X-Token":      e.Config.Token,
			},
		},
	}, nil
}

func (e *External) PostExternal(ctx context.Context, url string, data io.Reader) (*http.Response, error) {
	if url == "" {
		return nil, errors.New("url is required")
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, data)
	if err != nil {
		return nil, err
	}

	for k, v := range e.NewExternalConnection.Headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Println(string(b))
		resp.Body = io.NopCloser(bytes.NewReader(b))
		return nil, fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}
	return resp, nil
}

func (e *External) AddUser(ctx context.Context, role string) (*http.Response, error) {
	if role == "" {
		return nil, errors.New("role is required")
	}
	if e.Config.Token == "" {
		return nil, errors.New("token is required")
	}
	var url string
	if role == "superuser" {
		if e.Config.URLSuperuser == "" {
			return nil, errors.New("url superuser is required")
		}
		if e.User.Username == "" {
			return nil, errors.New("superuser username is required")
		}
		if e.User.Password == "" {
			return nil, errors.New("superuser password is required")
		}

		url = e.NewExternalConnection.URLAdmin
	} else if role == "merchant" {
		if e.Config.URLMerchant == "" {
			return nil, errors.New("url merchant is required")
		}
		if e.User.Username == "" {
			return nil, errors.New("merchant username is required")
		}
		if e.User.Password == "" {
			return nil, errors.New("merchant password is required")
		}
		url = e.NewExternalConnection.URLMerchant

	} else {
		return nil, errors.New("role not found")
	}
	jsonData := &User{
		Username: e.User.Username,
		Password: e.User.Password,
	}
	data, err := json.Marshal(jsonData)
	if err != nil {
		return nil, err
	}
	resp, _ := e.PostExternal(ctx, url, bytes.NewReader(data))
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	jsonStruct := struct {
		Message string `json:"message"`
		Errors  struct {
			Username string `json:"username"`
		} `json:"errors"`
	}{}
	err = json.Unmarshal(b, &jsonStruct)
	if err != nil {
		return nil, err
	}
	if jsonStruct.Errors.Username != "" {
		return nil, errors.New(jsonStruct.Errors.Username)
	}
	return resp, nil
}
