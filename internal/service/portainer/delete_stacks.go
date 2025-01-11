package portainer

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

func (p *Portainer) DeleteStackByName(ctx context.Context, name string) (string, error) {
	// 1. Input and Variables
	// 1.1. Input
	if name == "" {
		return "", errors.New("username is required")
	}

	// 2. Get all stack
	stackLists, err := p.GetStack(ctx)
	if err != nil {
		return "", err
	}

	// 3. Filter
	var isExist bool
	var id int
	var endpointId int
	isExist = false
	for _, v := range stackLists.Stacks {
		if v.Name == name {
			isExist = true
			id = v.Id
			endpointId = v.EndpointId
			break
		}
	}

	if !isExist {
		return "", errors.New("stack already deleted")
	}

	// 4. Delete stack
	url := fmt.Sprintf("%s/stacks/%d?endpointId=%d", p.BaseURL, id, endpointId)

	resp, err := p.DeletePortainer(ctx, url)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusNoContent {
		return "", errors.New("delete stack failed")
	}

	return fmt.Sprintf("success delete stack %s", name), nil
}
