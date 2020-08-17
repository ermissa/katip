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
	"strconv"

	"github.com/spf13/cobra"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Deletes selected command or commands",
	Run: func(cmd *cobra.Command, args []string) {
		if !checkIfCommandsFileExists() {
			fmt.Println(warningCommandsFileNotExist)
			return
		}

		// get commands from commands file
		var commands *Commands
		commands, err := getCommands()
		if err != nil {
			fmt.Println("get commands error:", err)
			return
		}
		if len(commands.Commands) == 0 {
			fmt.Println(warningCommandsFileNotExist)
			return
		}

		printCommandsAsTableWithIndexes(commands)

		fmt.Println("Which one do you want to delete? :")
		var rmIndexStr string
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		rmIndexStr = scanner.Text()
		rmIndex, err := strconv.Atoi(rmIndexStr)
		if err != nil {
			fmt.Println("Invalid input")
			return
		}
		rmIndex--
		if rmIndex >= 0 && rmIndex <= len(commands.Commands) {
			// ask for delete confirmation
			fmt.Println(commands.Commands[rmIndex].Command + " :: " + commands.Commands[rmIndex].Description + " :: " + commands.Commands[rmIndex].Alias)
			if askForConfirmation(confirmationTextForDeleteCommand) {
				commands.Commands = append(commands.Commands[:rmIndex], commands.Commands[rmIndex+1:]...)
				err = writeCommandsToFile(*commands)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Println("Command is successfully removed")
				return
			}
			fmt.Println("Aborted")
			return
		}
		fmt.Println("There is no command with this index")
		return
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}
