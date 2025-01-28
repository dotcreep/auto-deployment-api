package test

import (
	"fmt"
	"testing"

	"github.com/dotcreep/go-automate-deploy/internal/service/whoisdomain"
)

func TestCheckDomain(t *testing.T) {
	userDomain := "sella.com"

	data, err := whoisdomain.CheckDomainIsUsed(userDomain)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(data)

}
