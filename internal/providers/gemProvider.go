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

type GemProvider struct {
	*AbstractProvider
}

func (gem *GemProvider) InstallPackage(pkgConfiguration *models.PackageConfiguration) (err error, cmdErr error) {
	var versionArg []string

	if pkgConfiguration.Version != "" {
		versionArg = []string{"-v", pkgConfiguration.Version}
	}

	packageArgs := []string{pkgConfiguration.Name}
	packageArgs = append(packageArgs, versionArg...)

	name, args := gem.buildCommand(gem.InstallCommand, packageArgs...)

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

func (gem *GemProvider) UpdateRegistry() (err error, cmdErr error) {
	return
}

func (gem *GemProvider) CleanRegistry() (err error, cmdErr error) {
	return
}

func (gem *GemProvider) buildCommand(subCommand string, options ...string) (name string, args []string) {
	name = gem.Command
	if gem.RequiresRoot == true {
		name = "sudo"
		args = append(args, gem.Command)
	}

	args = append(args, subCommand)

	if len(options) > 0 {
		args = append(args, options...)
	}

	return
}

func NewGemProvider() *GemProvider {
	return &GemProvider{
		&AbstractProvider{
			Command:          "gem",
			InstallCommand:   "install",
			UpdateCommand:    "",
			CleanCommand:     "",
			RequiresRoot:     true,
			VersionSeparator: "",
		},
	}
}
