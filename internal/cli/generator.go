package cli

import (
	"fmt"

	"github.com/dotcreep/go-automate-deploy/internal/generator"
	"github.com/dotcreep/go-automate-deploy/internal/service"
	"github.com/dotcreep/go-automate-deploy/internal/utils"
)

type ReturnGenerate struct {
	Android    string
	API        string
	Web        string
	PasswordDB string
}

type EGenerate struct {
	BaseWebURL       string
	BaseAPIURL       string
	EGenerateAndroid EGenerateAndroid
	ClientName       string
	Email            string
	WebApi           string
	AppTitle         string
	PaketMerchant    string
}

type EGenerateAndroid struct {
	WebHost     string
	LabelApps   string
	PackageName string
}

func Generator(g *EGenerate) (*ReturnGenerate, error) {
	// This is area is copying file src to dist, always check with Deploy function
	data := service.Environment{}
	var Android string
	var Web string
	var API string
	m := &data.DataAPI.Management
	m.PathDist = "storage/dist"
	m.ClientName = g.ClientName

	a := &data.Android
	a.BaseAPI = g.BaseAPIURL
	a.BaseWeb = g.BaseWebURL
	a.Host = g.EGenerateAndroid.WebHost
	a.LabelApps = g.EGenerateAndroid.LabelApps
	a.PackageName = utils.GeneratePackageName(g.ClientName, g.EGenerateAndroid.WebHost)

	// Web
	e := &data.DataAPI
	e.URL.Web = g.BaseWebURL
	e.URL.API = g.BaseAPIURL
	e.URL.WebApi = g.WebApi
	passDB := generator.Password(10)
	e.Django.Name = fmt.Sprintf("api_%s", g.ClientName)
	e.Django.User = fmt.Sprintf("api_%s", g.ClientName)
	e.Django.Password = generator.Password(10)
	e.Django.Host = fmt.Sprintf("db_%s", g.ClientName)
	e.Laravel.Name = fmt.Sprintf("web_%s", g.ClientName)
	e.Laravel.User = fmt.Sprintf("web_%s", g.ClientName)
	e.Laravel.Password = passDB
	e.Laravel.Host = fmt.Sprintf("db_%s", g.ClientName)
	e.RootDatabase.Password = passDB
	e.SuperUser.Username = fmt.Sprintf("admin_%s", g.ClientName)
	e.SuperUser.Email = g.Email
	e.SuperUser.Password = generator.Password(10)
	e.Merchant.Username = g.ClientName
	e.Merchant.Password = generator.Password(10)
	e.AppTitle = fmt.Sprintf("\"%s\"", g.AppTitle)
	e.Username = fmt.Sprintf("\"%s\"", g.ClientName)
	e.PacketMerchant = fmt.Sprintf("\"%s\"", g.PaketMerchant)
	// Return

	m.PathSource = "storage/src/android/.env"
	result, err := data.AndroidEnvironment(m)
	if err != nil {
		return nil, err
	}
	Android = result
	m.PathSource = "storage/src/api/.env"
	result, err = data.ApiEnvironment(m)
	if err != nil {
		return nil, err
	}
	API = result
	m.PathSource = "storage/src/web/.env"
	result, err = data.WebEnvironment(m)
	if err != nil {
		return nil, err
	}
	Web = result

	fallback := &ReturnGenerate{
		Android:    Android,
		API:        API,
		Web:        Web,
		PasswordDB: passDB,
	}

	return fallback, nil
}
