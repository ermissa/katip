/*
Copyright © 2020 Fatih Ermiş <ermissaim@gmail.com>

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
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

// grepCmd represents the grep command
var grepCmd = &cobra.Command{
	Use:   "grep [ARGS]",
	Short: "Searches for your saved command",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// check if app directory exists
		isAppDirExists, err := checkIfAppDirExists()
		if err != nil || isAppDirExists == false {
			// if app directory does not exist, call init command
			initCmd.Run(cmd, args)
			return
		}
		// check if commands file exist
		if !checkIfCommandsFileExists() {
			err = createCommandsFile()
			if err != nil {
				fmt.Println("error occured while creating commands file: ", err)
				return
			}
			fmt.Println(warningCommandsFileNotExist)
			return
		}

		// concatenate args into the single string
		concatenatedArgs := strings.Join(args[:], " ")
		commands, err := getCommands()
		if err != nil {
			fmt.Println("error while getting commands:", err)
			return
		}

		if len(commands.Commands) == 0 {
			fmt.Println(warningCommandsFileNotExist)
			return
		}

		var concatenatedCommands, matches []string
		for _, command := range commands.Commands {
			concatenatedCommands = append(concatenatedCommands, command.Command+" :: "+command.Description+" :: "+command.Alias)
		}
		for _, concatenatedCommand := range concatenatedCommands {
			isMatched, err := regexp.MatchString(concatenatedArgs, concatenatedCommand)
			if err != nil {
				fmt.Println("error while searching argument:", err)
				return
			} else if isMatched {
				matches = append(matches, concatenatedCommand)
			}
		}
		if len(matches) == 0 {
			fmt.Println("No saved commands matches the pattern: ", concatenatedArgs)
			return
		}
		fmt.Printf("Command(s) found: \n\n")
		fmt.Println(strings.Join(matches[:], "\n"))
		return
	},
}

func init() {
	rootCmd.AddCommand(grepCmd)
}
