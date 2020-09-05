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
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run [ARGS]",
	Short: "Executes a saved command",
	Long:  `Give this command a hint about your saved command (alias, description or command itself (it is not logical to run the command you know through katip instead of writing it directly btw)) and your command will be executed.`,
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
		var cmdIndexes []int
		for i, command := range commands.Commands {
			concatenatedCommands = append(concatenatedCommands, strconv.Itoa(i)+" - "+command.Command+" :: "+command.Description+" :: "+command.Alias)
		}
		for i, concatenatedCommand := range concatenatedCommands {
			isMatched, err := regexp.MatchString(concatenatedArgs, concatenatedCommand)
			if err != nil {
				fmt.Println("error while searching argument:", err)
				return
			} else if isMatched {
				cmdIndexes = append(cmdIndexes, i)
				matches = append(matches, concatenatedCommand)
			}
		}
		if len(matches) == 0 {
			fmt.Println("No saved commands matches the pattern: ", concatenatedArgs)
			return
		}
		// check if a single or multiple commands are found
		if len(matches) > 1 {
			// if there are more than one possible commands to execute
			fmt.Printf("Multiple commands found. Please enter the index of command you want to execute: \n\n")
			for _, match := range matches {
				fmt.Println(match)
			}
			var cmdIndexStr string
			scanner := bufio.NewScanner(os.Stdin)
			fmt.Printf("\nIndex of command you want to execute: ")
			scanner.Scan()
			cmdIndexStr = scanner.Text()
			cmdIndex, err := strconv.Atoi(cmdIndexStr)
			if err != nil {
				fmt.Println("Invalid input")
				return
			}
			if isIntInSlice(cmdIndex, cmdIndexes) {
				// ask for confirmation to execute
				fmt.Println("\n" + commands.Commands[cmdIndex].Command + " :: " + commands.Commands[cmdIndex].Description + " :: " + commands.Commands[cmdIndex].Alias)
				if askForConfirmation(confirmationTextForRunCommand) {
					out, err := exec.Command("bash", "-c", commands.Commands[cmdIndex].Command).Output()
					if err != nil {
						fmt.Println("error : ", err)
					}
					fmt.Printf("\n%s\n", out)
					return
				}
				fmt.Println("Aborted")
				return
			}
			fmt.Println("There is no command with this index")
			return
		}
		// if there is one possible command to execute
		fmt.Printf("%s\n", matches[0])
		if askForConfirmation(confirmationTextForRunCommand) {
			out, err := exec.Command("bash", "-c", commands.Commands[cmdIndexes[0]].Command).Output()
			if err != nil {
				fmt.Println("error : ", err)
			}
			fmt.Printf("\n%s\n", out)
			return
		}
		fmt.Println("Aborted")
		return
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
