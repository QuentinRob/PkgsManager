package providers

import (
	"bytes"
	"errors"
	"os/exec"
	"qrobcis/pkgsmanager/internal/models"
)

type AptProvider struct {
	*AbstractProvider
}

func (apt *AptProvider) InstallPackage(pkgConfiguration models.PackageConfiguration) (err error, cmdErr error) {
	name, args := apt.buildCommand(apt.InstallCommand, true, pkgConfiguration.Name)
	cmd = exec.Command(name, args...)
	
	

	return
}


func (apt *AptProvider) UpdateRegistry() (err error, cmdErr error) {
	name, args := apt.buildCommand(apt.UpdateCommand, false)
	cmd := exec.Command(name, args...)
	errBuffer := new(bytes.Buffer)
	cmd.Stderr = errBuffer
	err = cmd.Run()
	if err != nil {
		err = errors.New("failed to update apt sources")
		cmdErr = errors.New(errBuffer.String())
	}

	return
}

func (apt *AptProvider) CleanRegistry() (err error, cmdErr error) {
	name, args := apt.buildCommand(apt.CleanCommand, true)
	cmd := exec.Command(name, args...)
	errBuffer := new(bytes.Buffer)
	cmd.Stderr = errBuffer
	err = cmd.Run()
	if err != nil {
		err = errors.New("failed to update apt sources")
		cmdErr = errors.New(errBuffer.String())
	}

	return
}

func (apt *AptProvider) buildCommand(subCommand string, autoApprove bool, options ...string) (name string, args []string) {
	name = apt.Command
	if apt.RequiresRoot == true {
		name = "sudo"
	}
	args = append(args, apt.Command, subCommand)

	if autoApprove == true {
		args = append(args, "-y")
	}

	if len(options) > 0 {
		args = append(args, options...)
	}

	return
}

func NewAptProvider() *AptProvider {
	return &AptProvider{
		&AbstractProvider{
			Command:          "apt-get",
			InstallCommand:   "install",
			UpdateCommand:    "update",
			CleanCommand:     "clean",
			RequiresRoot:     true,
			VersionSeparator: "=",
		},
	}
}
