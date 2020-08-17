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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Saves a new command",
	Run: func(cmd *cobra.Command, args []string) {
		isAppDirExists, err := checkIfAppDirExists()
		if err != nil || isAppDirExists == false {
			// if app directory does not exist, call init command
			initCmd.Run(cmd, args)
			return
		}
		// get command and description
		var commandInput, descriptionInput, aliasInput string
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Printf("Command: ")
		scanner.Scan()
		commandInput = scanner.Text()

		fmt.Printf("Description: ")
		scanner.Scan()
		descriptionInput = scanner.Text()

		fmt.Printf("Alias: ")
		scanner.Scan()
		aliasInput = scanner.Text()

		// create commands file if not exist
		if !checkIfCommandsFileExists() {
			err := createCommandsFile()
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		// write command to file as json
		// check if there is a record on file.
		newCommands := Commands{
			Commands: []Command{
				{
					Command:     commandInput,
					Description: descriptionInput,
					Alias:       aliasInput,
				},
			},
		}
		commandsFilePath, err := getCommandsFilePath()
		if err != nil {
			fmt.Println("error while getting commands file path:", err)
			return
		}
		fileStat, err := os.Stat(commandsFilePath)
		if err != nil {
			fmt.Println("error while getting commands file size:", err)
			return
		}
		// if file is empty, write command to file as JSON
		if fileStat.Size() == 0 {
			err = writeCommandsToFile(newCommands)
			if err != nil {
				fmt.Println(err)
				return
			}
			return
		}

		// if there is a record on file, add new command to end of existing JSON
		var existingCommands Commands
		file, _ := ioutil.ReadFile(commandsFilePath)
		err = json.Unmarshal(file, &existingCommands)
		if err != nil {
			fmt.Println("unmarshall error:", err)
			return
		}
		newCommand := Command{Command: commandInput, Description: descriptionInput, Alias: aliasInput}
		existingCommands.Commands = append(existingCommands.Commands, newCommand)

		err = writeCommandsToFile(existingCommands)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Command is successfully saved")
		return

	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
