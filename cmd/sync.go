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
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/viper"
	"os/exec"
	"time"

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
	pkgs := viper.GetStringSlice(groupName)
	requested = len(pkgs)
	pterm.DefaultSection.Println("Installing group: " + groupName)
	progress, _ := pterm.DefaultProgressbar.WithRemoveWhenDone(true).WithTotal(len(pkgs)).WithTitle(fmt.Sprint("Installing packages for group:", pterm.Blue(" ", groupName))).Start()
	for _, pkg := range pkgs {
		isSuccessful := installPackage(pkg, progress)
		if isSuccessful == true {
			success += 1
		}
	}
	pterm.Println()
	pterm.Info.Println("Successfully installed ", success, "/", len(pkgs), " packages.")
	pterm.Println()

	return
}

func installPackage(packageName string, progress *pterm.ProgressbarPrinter) (successful bool) {
	progress.UpdateTitle("Installing package " + packageName)
	time.Sleep(time.Millisecond * 350)
	cmd := exec.Command("sudo", "apt-get", "install", "-y", packageName)
	errBuffer := new(bytes.Buffer)
	cmd.Stderr = errBuffer
	err := cmd.Run()
	if err != nil {
		pterm.Error.Println("Failed to install " + packageName)
		pterm.DefaultParagraph.WithMaxWidth(60).Println(errBuffer.String())
		successful = false
	} else {
		pterm.Success.Println("Installed package " + packageName)
		successful = true
	}
	progress.Increment()

	return
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
