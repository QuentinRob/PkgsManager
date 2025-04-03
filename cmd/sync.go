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
	"context"
	"errors"
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"qrobcis/pkgsmanager/internal/models"
	"qrobcis/pkgsmanager/internal/providers"
	"qrobcis/pkgsmanager/internal/types/provider"
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
	providersMap[provider.NPM] = providers.NewNpmProvider()
	providersMap[provider.Gem] = providers.NewGemProvider()
	providersMap[provider.Golang] = providers.NewGoProvider()

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
			success += 1
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

func init() {
	rootCmd.AddCommand(syncCmd)
}
