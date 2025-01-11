package service

import (
	"errors"
	"fmt"
	"os/exec"
)

type SSHAccount struct {
	Username string
	Hostname string
	Password string
}

func (s *SSHAccount) SendEnvinvorment(dir string, d *Environment) (string, error) {
	path := d.DataAPI.Management.PathSource
	if path == "" {
		return "", errors.New("path is required")
	}
	if dir == "" {
		return "", errors.New("directory type is required")
	}
	var cmd *exec.Cmd
	command := fmt.Sprintf("scp %s %s@%s:/nfs/environment/%s/%s", path, s.Username, s.Hostname, d.DataAPI.Merchant.Username, dir)
	cmd = exec.Command(command)
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return "file transferred successfully", nil
}
