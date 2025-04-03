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

type GoProvider struct {
	*AbstractProvider
}

func (golang *GoProvider) InstallPackage(pkgConfiguration *models.PackageConfiguration) (err error, cmdErr error) {
	packageNameVersionned := ""
	if pkgConfiguration.Version != "" {
		packageNameVersionned = pkgConfiguration.Name + "@" + pkgConfiguration.Version
	} else {
		packageNameVersionned = pkgConfiguration.Name
	}

	name, args := golang.buildCommand(golang.InstallCommand, packageNameVersionned)
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

func (golang *GoProvider) UpdateRegistry() (err error, cmdErr error) {
	return
}

func (golang *GoProvider) UpgradePackages() (err error, cmdErr error) {
	name, args := golang.buildCommand(golang.UpdateCommand)
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

func (golang *GoProvider) CleanRegistry() (err error, cmdErr error) {
	name, args := golang.buildCommand(golang.CleanCommand)
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

func (golang *GoProvider) buildCommand(subCommand string, options ...string) (name string, args []string) {
	name = golang.Command
	if golang.RequiresRoot == true {
		name = "sudo"
		args = append(args, golang.Command)
	}

	args = append(args, subCommand)

	if subCommand == golang.CleanCommand {
		args = append(args, "-cache")
	}

	if len(options) > 0 {
		args = append(args, options...)
	}

	return
}

func NewGoProvider() *GoProvider {
	return &GoProvider{
		&AbstractProvider{
			Command:          "go",
			InstallCommand:   "install",
			UpdateCommand:    "",
			CleanCommand:     "clean",
			RequiresRoot:     false,
			VersionSeparator: "@",
		},
	}
}
