package portainer

import (
	"errors"
	"fmt"
	"os"

	"github.com/dotcreep/go-automate-deploy/internal/utils"
)

func RemoveClientDirectory(name string) error {
	// 1. Input and Variables
	// 1.1. Check input
	if name == "" {
		return errors.New("name is required")
	}

	// 1.2. Variables
	yamlConfig, err := utils.Open()
	if err != nil {
		return err
	}
	pathClient := yamlConfig.Config.PathClient
	pathEnv := yamlConfig.Config.PathEnvironment

	// 2. Check folder if not exists
	// 2.1. Definite the path with name
	fullPathClient := fmt.Sprintf("%s/%s", pathClient, name)
	fullPathEnv := fmt.Sprintf("%s/%s", pathEnv, name)

	// 2.2. Create array of folder
	allPath := []string{fullPathClient, fullPathEnv}

	// 2.3. Check if folder exists
	for _, path := range allPath {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return fmt.Errorf("folder %s not found", path)
		}
	}

	// 3. Remove folder
	for _, path := range allPath {
		err := os.RemoveAll(path)
		if err != nil {
			return err
		}
	}

	return nil
}
