package service

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Environment struct {
	Android Android
	DataAPI DataAPI
}

type Android struct {
	BaseAPI     string
	BaseWeb     string
	PackageName string
	LabelApps   string
	Host        string
}

type DataAPI struct {
	Django       Django
	Laravel      Laravel
	SuperUser    SuperUser
	Merchant     Merchant
	RootDatabase RootDatabase
	URL          URL
	Management   Management
	AppTitle     string
}

type Django struct {
	Name     string
	User     string
	Password string
	Host     string
}

type Laravel struct {
	Name     string
	User     string
	Password string
	Host     string
}

type SuperUser struct {
	Username string
	Password string
	Email    string
	Group    string
}

type Merchant struct {
	Username string
	Password string
}

type RootDatabase struct {
	Password string
}

type URL struct {
	Host   string
	Web    string
	API    string
	WebApi string
}

type SuperAdmin struct {
	Name     string
	Email    string
	Password string
}

type Management struct {
	PathSource string
	PathDist   string
	ClientName string
}

func (e *Environment) AndroidEnvironment(data *Management) (string, error) {
	if data.PathSource == "" {
		return "", errors.New("path source is required")
	}
	if data.PathDist == "" {
		return "", errors.New("path distribution is required")
	}
	if e.Android.PackageName == "" {
		return "", errors.New("package name is required")
	}
	if e.Android.BaseAPI == "" {
		return "", errors.New("base API is required")
	}
	if e.Android.BaseWeb == "" {
		return "", errors.New("base web is required")
	}
	if e.Android.LabelApps == "" {
		return "", errors.New("label apps is required")
	}
	if e.Android.Host == "" {
		return "", errors.New("host is required")
	}
	if data.ClientName == "" {
		return "", errors.New("client name is required")
	}

	file, err := os.Open(data.PathSource)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "APP_ID") {
			line = fmt.Sprintf("APP_ID=%s", e.Android.PackageName)
		} else if strings.HasPrefix(line, "BASE_URL_API") {
			line = fmt.Sprintf("BASE_URL_API=%s", e.Android.BaseAPI)
		} else if strings.HasPrefix(line, "BASE_URL_WEB") {
			line = fmt.Sprintf("BASE_URL_WEB=%s", e.Android.BaseWeb)
		} else if strings.HasPrefix(line, "LABEL_APPS") {
			line = fmt.Sprintf("LABEL_APPS=%s", e.Android.LabelApps)
		} else if strings.HasPrefix(line, "HOST") {
			line = fmt.Sprintf("HOST=%s", e.Android.Host)
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	directory := fmt.Sprintf("%s/%s/android", data.PathDist, data.ClientName)
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		os.MkdirAll(directory, 0755)
	}
	outputFile, err := os.Create(filepath.Join(directory, ".env"))
	if err != nil {
		return "", err
	}
	defer outputFile.Close()
	writer := bufio.NewWriter(outputFile)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return "", err
		}
	}
	writer.Flush()
	return "Android environment created successfully", nil
}

func (e *Environment) ApiEnvironment(data *Management) (string, error) {
	if e.DataAPI.Django.Name == "" {
		return "", errors.New("database api name is required")
	}
	if e.DataAPI.Django.User == "" {
		return "", errors.New("database api user is required")
	}
	if e.DataAPI.Django.Password == "" {
		return "", errors.New("database api password is required")
	}
	if e.DataAPI.Django.Host == "" {
		return "", errors.New("database api host is required")
	}
	if e.DataAPI.Laravel.Name == "" {
		return "", errors.New("database laravel name is required")
	}
	if e.DataAPI.Laravel.User == "" {
		return "", errors.New("database laravel user is required")
	}
	if e.DataAPI.Laravel.Password == "" {
		return "", errors.New("database laravel password is required")
	}
	if e.DataAPI.Laravel.Host == "" {
		return "", errors.New("database laravel host is required")
	}
	if e.DataAPI.SuperUser.Username == "" {
		return "", errors.New("superuser account username is required")
	}
	if e.DataAPI.SuperUser.Password == "" {
		return "", errors.New("superuser account password is required")
	}
	if e.DataAPI.SuperUser.Email == "" {
		return "", errors.New("superuser account email is required")
	}
	if e.DataAPI.Merchant.Username == "" {
		return "", errors.New("merchant username is required")
	}
	if e.DataAPI.Merchant.Password == "" {
		return "", errors.New("merchant password is required")
	}
	if e.DataAPI.URL.Web == "" {
		return "", errors.New("merchant web is required")
	}
	if data.PathSource == "" {
		return "", errors.New("path source is required")
	}
	if data.PathDist == "" {
		return "", errors.New("path distribution is required")
	}
	if data.ClientName == "" {
		return "", errors.New("client name is required")
	}
	file, err := os.Open(data.PathSource)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "DATABASE_NAME") {
			line = fmt.Sprintf("DATABASE_NAME=%s", e.DataAPI.Django.Name)
		} else if strings.HasPrefix(line, "DATABASE_USER") {
			line = fmt.Sprintf("DATABASE_USER=%s", e.DataAPI.Django.User)
		} else if strings.HasPrefix(line, "DATABASE_PASSWORD") {
			line = fmt.Sprintf("DATABASE_PASSWORD=%s", e.DataAPI.Django.Password)
		} else if strings.HasPrefix(line, "DATABASE_HOST") {
			line = fmt.Sprintf("DATABASE_HOST=%s", e.DataAPI.Django.Host)
		} else if strings.HasPrefix(line, "WEBMIN_DATABASE_NAME") {
			line = fmt.Sprintf("WEBMIN_DATABASE_NAME=%s", e.DataAPI.Laravel.Name)
		} else if strings.HasPrefix(line, "WEBMIN_DATABASE_USER") {
			line = fmt.Sprintf("WEBMIN_DATABASE_USER=%s", e.DataAPI.Laravel.User)
		} else if strings.HasPrefix(line, "WEBMIN_DATABASE_PASSWORD") {
			line = fmt.Sprintf("WEBMIN_DATABASE_PASSWORD=%s", e.DataAPI.Laravel.Password)
		} else if strings.HasPrefix(line, "WEBMIN_DATABASE_HOST") {
			line = fmt.Sprintf("WEBMIN_DATABASE_HOST=%s", e.DataAPI.Laravel.Host)
		} else if strings.HasPrefix(line, "DJANGO_SUPERUSER_USERNAME") {
			line = fmt.Sprintf("DJANGO_SUPERUSER_USERNAME=%s", e.DataAPI.SuperUser.Username)
		} else if strings.HasPrefix(line, "DJANGO_SUPERUSER_PASSWORD") {
			line = fmt.Sprintf("DJANGO_SUPERUSER_PASSWORD=%s", e.DataAPI.SuperUser.Password)
		} else if strings.HasPrefix(line, "DJANGO_SUPERUSER_EMAIL") {
			line = fmt.Sprintf("DJANGO_SUPERUSER_EMAIL=%s", e.DataAPI.SuperUser.Email)
		} else if strings.HasPrefix(line, "DJANGO_SUPERUSER_GROUP") {
			line = fmt.Sprintf("DJANGO_SUPERUSER_GROUP=%s", e.DataAPI.SuperUser.Group)
		} else if strings.HasPrefix(line, "ROOT_DB_PASS") {
			line = fmt.Sprintf("ROOT_DB_PASS=%s", e.DataAPI.RootDatabase.Password)
		} else if strings.HasPrefix(line, "MERCHANT_NAME") {
			line = fmt.Sprintf("MERCHANT_NAME=%s", e.DataAPI.Merchant.Username)
		} else if strings.HasPrefix(line, "MERCHANT_PASS") {
			line = fmt.Sprintf("MERCHANT_PASS=%s", e.DataAPI.Merchant.Password)
		} else if strings.HasPrefix(line, "WEBMIN_URL") {
			line = fmt.Sprintf("WEBMIN_URL=%s", e.DataAPI.URL.Web)
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	directory := fmt.Sprintf("%s/%s/api", data.PathDist, data.ClientName)
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		os.MkdirAll(directory, 0755)
	}
	outputFile, err := os.Create(filepath.Join(directory, ".env"))
	if err != nil {
		return "", err
	}
	defer outputFile.Close()
	writer := bufio.NewWriter(outputFile)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return "", err
		}
	}
	writer.Flush()
	return "API environment created successfully", nil
}

func (e *Environment) WebEnvironment(data *Management) (string, error) {
	if e.DataAPI.URL.Web == "" {
		return "", errors.New("web url is required")
	}
	if e.DataAPI.Laravel.Host == "" {
		return "", errors.New("database laravel host is required")
	}
	if e.DataAPI.Laravel.Name == "" {
		return "", errors.New("database laravel name is required")
	}
	if e.DataAPI.Laravel.User == "" {
		return "", errors.New("database laravel user is required")
	}
	if e.DataAPI.Laravel.Password == "" {
		return "", errors.New("database laravel password is required")
	}
	if e.DataAPI.URL.API == "" {
		return "", errors.New("api url is required")
	}
	if e.DataAPI.URL.WebApi == "" {
		return "", errors.New("web api url is required")
	}
	if e.DataAPI.AppTitle == "" {
		return "", errors.New("app title is required")
	}
	if data.PathSource == "" {
		return "", errors.New("path source is required")
	}
	if data.PathDist == "" {
		return "", errors.New("path distribution is required")
	}
	if data.ClientName == "" {
		return "", errors.New("client name is required")
	}
	file, err := os.Open(data.PathSource)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "APP_URL") {
			line = fmt.Sprintf("APP_URL=%s", e.DataAPI.URL.Web)
		} else if strings.HasPrefix(line, "DB_HOST") {
			line = fmt.Sprintf("DB_HOST=%s", e.DataAPI.Laravel.Host)
		} else if strings.HasPrefix(line, "DB_DATABASE") {
			line = fmt.Sprintf("DB_DATABASE=%s", e.DataAPI.Laravel.Name)
		} else if strings.HasPrefix(line, "DB_USERNAME") {
			line = fmt.Sprintf("DB_USERNAME=%s", e.DataAPI.Laravel.User)
		} else if strings.HasPrefix(line, "DB_PASSWORD") {
			line = fmt.Sprintf("DB_PASSWORD=%s", e.DataAPI.Laravel.Password)
		} else if strings.HasPrefix(line, "BASE_API_WIL") {
			line = fmt.Sprintf("BASE_API_WIL=%s", e.DataAPI.URL.API)
		} else if strings.HasPrefix(line, "BASE_API_WEB") {
			line = fmt.Sprintf("BASE_API_WEB=%s", e.DataAPI.URL.WebApi)
		} else if strings.HasPrefix(line, "BASE_API") {
			line = fmt.Sprintf("BASE_API=%s", e.DataAPI.URL.API)
		} else if strings.HasPrefix(line, "CHAT_API") {
			line = fmt.Sprintf("CHAT_API=%s", e.DataAPI.URL.WebApi)
		} else if strings.HasPrefix(line, "APP_TITLE") {
			line = fmt.Sprintf("APP_TITLE=%s", e.DataAPI.AppTitle)
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	directory := fmt.Sprintf("%s/%s/web", data.PathDist, data.ClientName)
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		os.MkdirAll(directory, 0755)
	}
	outputFile, err := os.Create(filepath.Join(directory, ".env"))
	if err != nil {
		return "", err
	}
	defer outputFile.Close()
	writer := bufio.NewWriter(outputFile)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return "", err
		}
	}
	writer.Flush()
	return "Web environment created successfully", nil
}
