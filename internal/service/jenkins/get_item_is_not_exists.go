package jenkins

import (
	"context"
	"errors"
)

func (j *Jenkins) GetItemIsNotExists(ctx context.Context, username string) (string, error) {
	if username == "" {
		return "", errors.New("username is required")
	}

	res, err := j.GetAllJobItem(ctx)
	if err != nil {
		return "", err
	}
	var isExists bool = false
	for _, v := range res {
		if v == username {
			isExists = true
			break
		}
	}
	if isExists {
		return "", errors.New("builder is exists")
	}
	return "builder already for create", nil
}
