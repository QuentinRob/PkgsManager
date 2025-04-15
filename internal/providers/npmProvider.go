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

type NpmProvider struct {
	*AbstractProvider
}

func (npm *NpmProvider) InstallPackage(pkgConfiguration *models.PackageConfiguration) (err error, cmdErr error) {
	packageNameVersionned := ""
	if pkgConfiguration.Version != "" {
		packageNameVersionned = pkgConfiguration.Name + "@" + pkgConfiguration.Version
	} else {
		packageNameVersionned = pkgConfiguration.Name
	}
	name, args := npm.buildCommand(npm.InstallCommand, packageNameVersionned)
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

func (npm *NpmProvider) UpdateRegistry() (err error, cmdErr error) {
	return
}

func (npm *NpmProvider) UpgradePackages() (err error, cmdErr error) {
	name, args := npm.buildCommand(npm.UpdateCommand)
	cmd := exec.Command(name, args...)
	errBuffer := new(bytes.Buffer)
	cmd.Stderr = errBuffer
	err = cmd.Run()
	if err != nil {
		err = errors.New("failed to update npm sources")
		cmdErr = errors.New(errBuffer.String())
	}

	return
}

func (npm *NpmProvider) CleanRegistry() (err error, cmdErr error) {
	name, args := npm.buildCommand(npm.CleanCommand)
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

func (npm *NpmProvider) buildCommand(subCommand string, options ...string) (name string, args []string) {
	name = npm.Command
	if npm.RequiresRoot == true {
		name = "sudo"
		args = append(args, npm.Command)
	}

	if subCommand == npm.InstallCommand || subCommand == npm.UpdateCommand {
		args = append(args, subCommand, "-g")
	} else if subCommand == npm.CleanCommand {
		args = append(args, "cache", npm.CleanCommand)
	}

	if len(options) > 0 {
		args = append(args, options...)
	}

	return
}

func NewNpmProvider() *NpmProvider {
	return &NpmProvider{
		&AbstractProvider{
			Command:          "npm",
			InstallCommand:   "install",
			UpdateCommand:    "update",
			CleanCommand:     "clean",
			RequiresRoot:     false,
			VersionSeparator: "@",
		},
	}
}
