// Copyright © 2018 Johannes Mitlmeier <dev.jojomi@yahoo.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/sahilm/fuzzy"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	flagVerbose   bool
	flagCfgFile   string
	flagFirstOnly bool
)

// This represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "ff",
	Short: "find folder that matches the given fuzzy search",
	Run:   rootCmdHandler,
}

const highlight = "\033[1;92m%s\033[0m"

func rootCmdHandler(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		log.Fatal("no search given")
	}
	search := args[0]
	if flagVerbose {
		fmt.Printf("Searching for: %s\n", search)
	}

	paths := viper.GetStringSlice("paths")
	if flagVerbose {
		fmt.Println("Search base paths:")
		for _, p := range paths {
			fmt.Println(p)
		}
	}

	subDirs := getSubDirs(paths)
	if flagVerbose {
		fmt.Println("Folders in search paths:")
		for _, s := range subDirs {
			fmt.Println(s)
		}
	}

	matches := fuzzy.Find(search, subDirs)

	interactive := terminal.IsTerminal(int(os.Stdout.Fd()))
	if flagVerbose {
		fmt.Println("Results:")
	}
	for _, match := range matches {
		for i, m := range match.Str {
			if interactive && contains(i, match.MatchedIndexes) {
				fmt.Print(fmt.Sprintf(highlight, string(m)))
			} else {
				fmt.Print(string(m))
			}
		}
		fmt.Println()
		if flagFirstOnly {
			break
		}
	}
}

func contains(needle int, haystack []int) bool {
	for _, i := range haystack {
		if needle == i {
			return true
		}
	}
	return false
}

func getSubDirs(paths []string) []string {
	result := []string{}
	var (
		files []os.FileInfo
		err   error
	)
	for _, path := range paths {
		path, err = homedir.Expand(path)
		if err != nil {
			// TODO print on stderr
			continue
		}
		files, err = ioutil.ReadDir(path)
		if err != nil {
			// TODO print on stderr
			continue
		}

		for _, f := range files {
			if !f.IsDir() {
				continue
			}
			result = append(result, filepath.Join(path, f.Name()))
		}
	}
	return result
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	f := RootCmd.PersistentFlags()
	f.StringVarP(&flagCfgFile, "config", "c", "", "config file (default is $HOME/.ff/config.yml)")
	f.BoolVarP(&flagFirstOnly, "first-only", "f", false, "print only first element")
	f.BoolVarP(&flagVerbose, "verbose", "v", false, "print search details")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if flagCfgFile != "" {
		viper.SetConfigFile(flagCfgFile)
	}

	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.ff")
	viper.AddConfigPath("/etc/ff")
	viper.AutomaticEnv()

	// If a config file is found, parse it
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if flagVerbose {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
