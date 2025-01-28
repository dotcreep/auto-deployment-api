package whoisdomain

import (
	"context"
	"errors"
	"strings"

	"github.com/domainr/whois"
)

func WhoisDomain(ctx context.Context, domain string) (string, error) {
	if domain == "" {
		return "", errors.New("domain is required")
	}

	is, err := isDomain(domain)
	if err != nil {
		return "", err
	}
	if !is {
		return "", errors.New("domain is invalid")
	}

	request, err := whois.NewRequest(domain)
	if err != nil {
		return "", err
	}
	response, err := whois.DefaultClient.FetchContext(ctx, request)
	if err != nil {
		return "", err
	}
	if strings.Contains(string(response.Body), "redemptionPeriod") ||
		strings.Contains(string(response.Body), "autoRenewPeriod") ||
		strings.Contains(string(response.Body), "autoRenew") ||
		strings.Contains(string(response.Body), "pendingDelete") ||
		strings.Contains(string(response.Body), "pendingTransfer") ||
		strings.Contains(string(response.Body), "pendingUpdate") ||
		strings.Contains(string(response.Body), "clientHold") ||
		strings.Contains(string(response.Body), "clientTransferProhibited") ||
		strings.Contains(string(response.Body), "serverHold") ||
		strings.Contains(string(response.Body), "ok") {
		return "unavailable", nil
	} else {
		return "available", nil
	}
}

func isDomain(domain string) (bool, error) {
	if domain == "" {
		return false, errors.New("domain is required")
	}

	parts := strings.Split(domain, ".")
	if len(parts) < 2 {
		return false, errors.New("domain is invalid")
	}
	if len(parts) > 2 {
		if parts[len(parts)-2] == "co" || parts[len(parts)-2] == "biz" || parts[len(parts)-2] == "com" {
			return true, nil
		}
	} else if len(parts) == 2 {
		return true, nil
	}
	return false, errors.New("domain is invalid")
}
