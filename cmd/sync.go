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
	"context"
	"errors"
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"qrobcis/pkgsmanager/internal/models"
	"qrobcis/pkgsmanager/internal/providers"
	"qrobcis/pkgsmanager/internal/types/provider"
	"strings"

	"github.com/spf13/cobra"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Install/Remove packages based on the configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		pterm.Info.Println("Synchronizing packages...")
		pterm.Println()

		providersMap := initProviders()
		configuration := initConfiguration()
		ctx = context.WithValue(ctx, "providers", providersMap)

		err, cmdErr := providersMap[provider.APT].UpdateRegistry()
		if err != nil {
			pterm.Error.Println(err)
			pterm.DefaultParagraph.WithMaxWidth(60).Println(cmdErr)
			os.Exit(1)
			return
		}
		totalRequestedPackages := 0
		totalInstalledPackages := 0

		for _, group := range configuration {
			intalled, requested := installGroup(ctx, group)
			totalInstalledPackages += intalled
			totalRequestedPackages += requested
		}

		err, cmdErr = providersMap[provider.APT].CleanRegistry()
		if err != nil {
			pterm.Error.Println(err)
			pterm.DefaultParagraph.WithMaxWidth(60).Println(cmdErr)
			os.Exit(1)
			return
		}

		pterm.Println()
		pterm.Info.Println("Installed ", totalInstalledPackages, "/", totalRequestedPackages, " packages.")
	},
}

func initProviders() (providersMap map[provider.Provider]providers.PackageProvider) {
	providersMap = make(map[provider.Provider]providers.PackageProvider)
	providersMap[provider.APT] = providers.NewAptProvider()

	return
}

func initConfiguration() (configuration map[string]*models.GroupConfiguration) {
	groups := viper.AllKeys()
	configuration = make(map[string]*models.GroupConfiguration)

	for _, groupName := range groups {
		groupConfiguration := models.NewGroupConfiguration(groupName)
		configuration[groupName] = groupConfiguration

		var packagesConfigurations []models.RawPackageConfiguration
		if err := viper.UnmarshalKey(groupName, &packagesConfigurations); err != nil {
			panic(err)
		}
		for _, pkgConfiguration := range packagesConfigurations {
			groupConfiguration.AddPackage(models.NewPackageConfiguration(
				pkgConfiguration.Name,
				pkgConfiguration.GPGKey,
				pkgConfiguration.SourceList,
				pkgConfiguration.Provider,
				pkgConfiguration.Version,
			))
		}
	}

	return
}

func installGroup(ctx context.Context, group *models.GroupConfiguration) (success int, requested int) {
	success = 0
	requested = len(group.Packages)

	pterm.DefaultSection.Println("Installing group: " + group.Name)
	progress, _ := pterm.DefaultProgressbar.WithRemoveWhenDone(true).WithTotal(len(group.Packages)).WithTitle(fmt.Sprint("Installing packages for group:", pterm.Blue(" ", group.Name))).Start()

	for _, packageConfiguration := range group.Packages {
		err, cmdErr := installPackage(ctx, packageConfiguration, progress)
		if err != nil {
			pterm.Error.Println(err)
			if cmdErr != nil {
				pterm.DefaultParagraph.Printfln(cmdErr.Error())
			}
		} else {
			paddedProvider := formatProvider(packageConfiguration)
			pterm.FgGreen.Println("| " + paddedProvider + "| Installed package " + packageConfiguration.Name)
		}
	}
	pterm.Println()
	pterm.Info.Println("Successfully installed ", success, "/", len(group.Packages), " packages.")
	pterm.Println()

	return
}

func installPackage(ctx context.Context, pkgConfiguration *models.PackageConfiguration, progress *pterm.ProgressbarPrinter) (err error, cmdErr error) {

	progress.UpdateTitle("Installing package " + pkgConfiguration.Name)

	var providersMap map[provider.Provider]providers.PackageProvider
	providersMap = ctx.Value("providers").(map[provider.Provider]providers.PackageProvider)
	if pkgConfiguration.Provider != provider.Unknown {
		err, cmdErr = providersMap[pkgConfiguration.Provider].InstallPackage(pkgConfiguration)
	} else {
		err = errors.New(fmt.Sprintf("Provider not supported: %s", pkgConfiguration.Provider))
	}
	progress.Increment()

	return
}

func formatProvider(pkgConfiguration *models.PackageConfiguration) (paddedProvider string) {
	var providerStyle *pterm.Style
	if pkgConfiguration.Provider == provider.APT || pkgConfiguration.Provider == provider.Unset {
		pkgConfiguration.Provider = provider.APT
		providerStyle = pterm.NewStyle(pterm.Bold, pterm.FgYellow)
	} else if pkgConfiguration.Provider == provider.Golang {
		providerStyle = pterm.NewStyle(pterm.Bold, pterm.FgBlue)
	} else if pkgConfiguration.Provider == provider.NPM {
		providerStyle = pterm.NewStyle(pterm.Bold, pterm.FgGreen)
	} else if pkgConfiguration.Provider == provider.Gem {
		providerStyle = pterm.NewStyle(pterm.Bold, pterm.FgRed)
	} else if pkgConfiguration.Provider == provider.Pip {
		providerStyle = pterm.NewStyle(pterm.Bold, pterm.FgCyan)
	} else {
		providerStyle = pterm.NewStyle(pterm.Bold, pterm.FgDefault)
	}
	paddedProvider = providerStyle.Sprintf("%-5s", pkgConfiguration.Provider)

	return
}

func buildPipCommand(pkgConfiguration models.PackageConfiguration) (cmd *exec.Cmd) {
	versionnedName := pkgConfiguration.Name

	if pkgConfiguration.Version != "" {
		versionnedName = "'" + versionnedName + "==" + pkgConfiguration.Version + "'"
	}

	args := []string{"install", versionnedName}

	cmd = exec.Command("pipx", args...)

	return
}

func buildGemCommand(pkgConfiguration models.PackageConfiguration) (cmd *exec.Cmd) {
	var versionArg []string

	if pkgConfiguration.Version != "" {
		versionArg = []string{"-v", pkgConfiguration.Version}
	}

	args := []string{"gem", "install", pkgConfiguration.Name}
	args = append(args, versionArg...)

	cmd = exec.Command("sudo", args...)

	return
}

func buildNpmCommand(pkgConfiguration models.PackageConfiguration) (cmd *exec.Cmd) {
	packageNameVersionned := ""
	if pkgConfiguration.Version != "" {
		packageNameVersionned = pkgConfiguration.Name + "@" + pkgConfiguration.Version
	} else {
		packageNameVersionned = pkgConfiguration.Name
	}
	cmd = exec.Command("npm", "install", "-g", packageNameVersionned)

	return
}

func buildGoCommand(pkgConfiguration models.PackageConfiguration) (cmd *exec.Cmd) {
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

func addSourceList(pkgConfiguration models.PackageConfiguration) {
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

		//		updateApt()
	}
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
