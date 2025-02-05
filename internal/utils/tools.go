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
		domainsub := []string{"co", "biz,", "my", "ac", "gov", "edu", "org", "web", "sch", "me", "com", "net", "go", "govt"}
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

func GeneratePackageName(username, domain string) string {
	baseDomain := GetBaseDomain(domain)
	reverseDomain := strings.Split(baseDomain, ".")
	var idApps string
	if len(reverseDomain) > 2 {
		reverseDomain = reverseDomain[len(reverseDomain)-3:]
		idApps = fmt.Sprintf("%s.%s.%s", reverseDomain[2], reverseDomain[1], reverseDomain[0])
	} else {
		reverseDomain = reverseDomain[len(reverseDomain)-2:]
		idApps = fmt.Sprintf("%s.%s", reverseDomain[1], reverseDomain[0])
	}
	return fmt.Sprintf("%s.%s", idApps, username)
}
