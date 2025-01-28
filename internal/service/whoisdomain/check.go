package whoisdomain

import (
	"context"
	"errors"
)

func CheckDomainIsUsed(domain string) (bool, error) {
	if domain == "" {
		return false, errors.New("domain is required")
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	data, err := WhoisDomain(ctx, domain)
	if err != nil {
		return false, err
	}
	if data == "unavailable" {
		return false, nil
	} else if data == "available" {
		return true, nil
	}
	return false, errors.New("unexpected response")
}
