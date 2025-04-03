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
	"os"
	"os/exec"
	"qrobcis/pkgsmanager/internal/models"
	"strings"
)

type AptProvider struct {
	*AbstractProvider
}

func (apt *AptProvider) InstallPackage(pkgConfiguration *models.PackageConfiguration) (err error, cmdErr error) {
	if pkgConfiguration.SourceList != "" {
		err, cmdErr = apt.addSourceList(pkgConfiguration)
		if err != nil || cmdErr != nil {
			return
		}
	}

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

func (apt *AptProvider) installGPGKey(GPGKey string, packageName string) (keyPath string, err error) {
	keyPath = "/etc/apt/keyrings/" + packageName + "-apt-keyring.gpg"
	if _, err := os.Stat(keyPath); errors.Is(err, os.ErrNotExist) {
		cmdCurl := exec.Command("curl", "-fsSL", GPGKey)
		cmd := exec.Command("sudo", "gpg", "--dearmor", "-o", keyPath)
		cmd.Stdin, _ = cmdCurl.StdoutPipe()
		errBuffer := new(bytes.Buffer)
		cmd.Stderr = errBuffer
		err = cmd.Start()
		_ = cmdCurl.Run()
		err = cmd.Wait()
	}
	return
}

func (apt *AptProvider) addSourceList(pkgConfiguration *models.PackageConfiguration) (err error, cmdErr error) {
	sourceListPath := "/etc/apt/sources.list.d/" + pkgConfiguration.Name + ".list"
	if _, err = os.Stat(sourceListPath); errors.Is(err, os.ErrNotExist) {
		sourceListSignature := ""
		if pkgConfiguration.GPGKey != "" {
			var keyPath string
			keyPath, err = apt.installGPGKey(pkgConfiguration.GPGKey, pkgConfiguration.Name)
			if err != nil {
				return
			}
			sourceListSignature = "[arch=amd64 signed-by=" + keyPath + "]"
		}

		sourceList := "deb " + sourceListSignature + " " + pkgConfiguration.SourceList

		cmd := exec.Command("sudo", "tee", "-a", sourceListPath)
		cmd.Stdin = strings.NewReader(sourceList)
		errBuffer := new(bytes.Buffer)
		cmd.Stderr = errBuffer
		err = cmd.Run()
		if err != nil {
			cmdErr = errors.New(errBuffer.String())
		}

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
