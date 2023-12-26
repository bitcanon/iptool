/*
Copyright © 2024 Mikael Schultz <mikael@conf-t.se>

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
	"os"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Version: "1.1.0",
	Use:     "iptool",
	Short:   "Simplify IP address and subnet tasks on the command line",
	Long: `Simplify IP address and subnet tasks on the command line.

The iptool is useful when detailed information about an IP address or various
subnetting information is required. The tool takes input from the user,
parses the input and presents the user with detailed information.

Author: Mikael Schultz <mikael@conf-t.se>
GitHub: https://github.com/bitcanon/iptool
`,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Set default config file path for the flag help text
	var defaultConfigPath string
	if runtime.GOOS == "windows" {
		defaultConfigPath = "%USERPROFILE%\\.iptool.yaml"
	} else {
		defaultConfigPath = "~/.iptool.yaml"
	}

	// Add flag for custom config file path
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is "+defaultConfigPath+")")

	// Add persistent flag for debug mode
	rootCmd.PersistentFlags().Bool("debug", false, "show debug info")
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	rootCmd.PersistentFlags().Lookup("debug").Hidden = true

	// Set a custom version template
	rootCmd.SetVersionTemplate(`{{ printf "%s %s" .Name .Version }}`)

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".iptool" (without extension)
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".iptool")
	}

	// Check for environment variables prefixed with MACTOOL
	replacer := strings.NewReplacer("-", "_", ".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("IPTOOL")

	// Load any environment variable that match an existing config key
	viper.AutomaticEnv()

	// Print all environment variables loaded in viper
	// viper.Debug()

	// If a config file is found, read it in
	viper.ReadInConfig()
}
