/*
Copyright © 2025 Quentin ROBCIS

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package providers

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"qrobcis/pkgsmanager/internal/models"
)

type AptProvider struct {
	*AbstractProvider
}

func (apt *AptProvider) InstallPackage(pkgConfiguration *models.PackageConfiguration) (err error, cmdErr error) {
	name, args := apt.buildCommand(apt.InstallCommand, true, pkgConfiguration.Name)
	cmd := exec.Command(name, args...)

	errBuffer := new(bytes.Buffer)
	cmd.Stderr = errBuffer
	err = cmd.Run()
	if err != nil {
		err = errors.New(fmt.Sprintf("Failed to install %s", pkgConfiguration.Name))
		cmdErr = errors.New(errBuffer.String())
	}

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
