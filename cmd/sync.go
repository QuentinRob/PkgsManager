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

package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"qrobcis/pkgsmanager/internal"
	"strings"

	"github.com/spf13/cobra"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Install/Remove packages based on the configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		pterm.Info.Println("Synchronizing packages...")
		totalRequestedPackages := 0
		totalInstalledPackages := 0
		groups := viper.AllKeys()
		for _, groupName := range groups {
			intalled, requested := installGroup(groupName)
			totalInstalledPackages += intalled
			totalRequestedPackages += requested
		}

		pterm.Println()
		pterm.Info.Println("Installed ", totalInstalledPackages, "/", totalRequestedPackages, " packages.")
	},
}

func installGroup(groupName string) (success int, requested int) {
	success = 0
	var pkgsConfigurations []internal.PackageConfiguration
	if err := viper.UnmarshalKey(groupName, &pkgsConfigurations); err != nil {
		panic(err)
	}
	requested = len(pkgsConfigurations)
	pterm.DefaultSection.Println("Installing group: " + groupName)
	progress, _ := pterm.DefaultProgressbar.WithRemoveWhenDone(true).WithTotal(len(pkgsConfigurations)).WithTitle(fmt.Sprint("Installing packages for group:", pterm.Blue(" ", groupName))).Start()
	for _, pkgConfiguration := range pkgsConfigurations {
		isSuccessful := installPackage(pkgConfiguration, progress)
		if isSuccessful == true {
			success += 1
		}
	}
	pterm.Println()
	pterm.Info.Println("Successfully installed ", success, "/", len(pkgsConfigurations), " packages.")
	pterm.Println()

	return
}

func installPackage(pkgConfiguration internal.PackageConfiguration, progress *pterm.ProgressbarPrinter) (successful bool) {

	progress.UpdateTitle("Installing package " + pkgConfiguration.Name)

	var cmd *exec.Cmd
	if pkgConfiguration.Provider == "" || pkgConfiguration.Provider == "apt" {
		if pkgConfiguration.SourceList != "" {
			addSourceList(pkgConfiguration)
		}
		cmd = exec.Command("sudo", "apt-get", "install", "-y", pkgConfiguration.Name)
	} else if pkgConfiguration.Provider == "go" {
		cmd = buildGoCommand(pkgConfiguration)
	} else {
		pterm.Warning.Println("Provider not supported: " + pkgConfiguration.Provider)
		successful = false
		return
	}

	errBuffer := new(bytes.Buffer)
	cmd.Stderr = errBuffer
	err := cmd.Run()
	if err != nil {
		pterm.Error.Println("Failed to install " + pkgConfiguration.Name)
		pterm.DefaultParagraph.WithMaxWidth(60).Println(errBuffer.String())
		successful = false
	} else {
		pterm.Success.Println("Installed package " + pkgConfiguration.Name)
		successful = true
	}
	progress.Increment()

	return
}

func buildGoCommand(pkgConfiguration internal.PackageConfiguration) (cmd *exec.Cmd) {
	packageNameVersionned := ""
	if pkgConfiguration.Version != "" {
		packageNameVersionned = pkgConfiguration.Name + "@" + pkgConfiguration.Version
	} else {
		packageNameVersionned = pkgConfiguration.Name
	}
	cmd = exec.Command("go", "install", packageNameVersionned)

	return
}

func installGPGKey(GPGKey string, packageName string) (keyPath string) {
	keyPath = "/etc/apt/keyrings/" + packageName + "-apt-keyring.gpg"
	if _, err := os.Stat(keyPath); errors.Is(err, os.ErrNotExist) {
		cmdCurl := exec.Command("curl", "-fsSL", GPGKey)
		cmd := exec.Command("sudo", "gpg", "--dearmor", "-o", keyPath)
		cmd.Stdin, _ = cmdCurl.StdoutPipe()
		errBuffer := new(bytes.Buffer)
		cmd.Stderr = errBuffer
		err := cmd.Start()
		_ = cmdCurl.Run()
		err = cmd.Wait()
		if err != nil {
			pterm.Error.Println("Failed to add GPG Key: " + GPGKey)
			pterm.DefaultParagraph.WithMaxWidth(60).Println(errBuffer.String())
		} else {
			pterm.Success.Println("Installed GPG Key: " + GPGKey)
		}
	}
	return
}

func addSourceList(pkgConfiguration internal.PackageConfiguration) {
	sourceListPath := "/etc/apt/sources.list.d/" + pkgConfiguration.Name + ".list"
	if _, err := os.Stat(sourceListPath); errors.Is(err, os.ErrNotExist) {
		sourceListSignature := ""
		if pkgConfiguration.GPGKey != "" {
			keyPath := installGPGKey(pkgConfiguration.GPGKey, pkgConfiguration.Name)
			sourceListSignature = "[arch=amd64 signed-by=" + keyPath + "]"
		}

		sourceList := "deb " + sourceListSignature + " " + pkgConfiguration.SourceList

		cmd := exec.Command("sudo", "tee", "-a", sourceListPath)
		cmd.Stdin = strings.NewReader(sourceList)
		errBuffer := new(bytes.Buffer)
		cmd.Stderr = errBuffer
		err := cmd.Run()
		if err != nil {
			pterm.Error.Println("Failed to add source list: " + sourceList)
			pterm.DefaultParagraph.WithMaxWidth(60).Println(errBuffer.String())
		} else {
			pterm.Success.Println("Added source list: " + sourceList)
		}

		cmdUpdate := exec.Command("sudo", "apt", "update")
		errBufferUpdate := new(bytes.Buffer)
		cmdUpdate.Stderr = errBufferUpdate
		errUpdate := cmdUpdate.Run()
		if errUpdate != nil {
			pterm.Error.Println("Failed to update apt sources")
			pterm.DefaultParagraph.WithMaxWidth(60).Println(errBuffer.String())
		} else {
			pterm.Success.Println("Updated apt sources")
		}
	}
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
