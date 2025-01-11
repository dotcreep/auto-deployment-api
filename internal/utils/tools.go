package utils

import (
	"fmt"
	"strings"
)

func GetBaseDomain(domain string) string {
	if domain == "" {
		return ""
	}
	slDom := strings.Split(domain, ".")
	var baseDomain string
	if len(slDom) > 2 {
		domainsub := []string{"co", "biz,"}
		for _, v := range domainsub {
			if slDom[len(slDom)-2] == v {
				baseDomain = fmt.Sprintf("%s.%s.%s", slDom[len(slDom)-3], slDom[len(slDom)-2], slDom[len(slDom)-1])
				break
			}
			baseDomain = fmt.Sprintf("%s.%s", slDom[len(slDom)-2], slDom[len(slDom)-1])
		}
		//baseDomain = fmt.Sprintf("%s.%s", slDom[len(slDom)-2], slDom[len(slDom)-1])
	} else {
		baseDomain = domain
	}
	return baseDomain
}
