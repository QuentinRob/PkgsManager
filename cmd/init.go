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
	"github.com/pterm/pterm"
	"github.com/spf13/viper"
	"os"
	"qrobcis/pkgsmanager/internal/models"
	"qrobcis/pkgsmanager/internal/types/provider"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new configuration file if not present.",
	Run: func(cmd *cobra.Command, args []string) {
		home, err := os.UserHomeDir()
		spinner, _ := pterm.DefaultSpinner.Start("Initializing configuration file at: " + pterm.Red(" ", home, "/.pkgsmanager.yaml"))
		cobra.CheckErr(err)
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".pkgsmanager")
		viper.Set("default", [...]models.PackageConfiguration{{Name: "git", Provider: provider.APT}, {Name: "vim", Provider: provider.APT}})
		err = viper.SafeWriteConfig()
		if err != nil {
			spinner.Info()
			pterm.Info.Println(err)
		} else {
			spinner.Success()
			pterm.Success.Println("Configuration file initialized")
		}

	},
}

func init() {
	configCmd.AddCommand(initCmd)
}
