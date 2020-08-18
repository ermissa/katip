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
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit your saved command",
	Run: func(cmd *cobra.Command, args []string) {
		// check if app directory exists
		isAppDirExists, err := checkIfAppDirExists()
		if err != nil || isAppDirExists == false {
			// if app directory does not exist, call init command
			initCmd.Run(cmd, args)
			return
		}

		if !checkIfCommandsFileExists() {
			fmt.Println(warningCommandsFileNotExist)
			return
		}
		commandsFilePath, err := getCommandsFilePath()
		if err != nil {
			fmt.Println("error while getting commands file path:", err)
			return
		}
		editor := exec.Command("vim", commandsFilePath)
		editor.Stdin = os.Stdin
		editor.Stdout = os.Stdout
		err = editor.Run()
		if err != nil {
			fmt.Println("error occured while opening vim:", err)
			return
		}
		return
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
