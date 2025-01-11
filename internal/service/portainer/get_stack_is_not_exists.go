package portainer

import (
	"context"
	"errors"
)

func (p *Portainer) GetStackIsNotExists(ctx context.Context, username string) (string, error) {
	if username == "" {
		return "", errors.New("username is required")
	}

	res, err := p.GetStack(ctx)
	if err != nil {
		return "", err
	}

	var isExists bool
	isExists = false
	for _, v := range res.Stacks {
		if v.Name == username {
			isExists = true
			break
		}
	}
	if isExists {
		return "", errors.New("system is exists")
	}
	return "system already for create", nil
}
