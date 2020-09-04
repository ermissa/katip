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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/jedib0t/go-pretty/table"
)

var (
	version                               = "0.1.0"
	appDirName                            = ".katip"
	commandsFileName                      = "commands.json"
	warningCommandsFileNotExist           = "No commands to show. Add one by 'katip new'"
	confirmationTextForDeleteCommand      = "Command will be removed"
	confirmationTextForRunCommand         = "Execute?"
	confirmationTextForDeleteAppDirectory = "[CRITICAL] Remove everything inside ~/.katip ?"
)

type Command struct {
	Command     string `json:"command"`
	Description string `json:"description"`
	Alias       string `json:"alias"`
}
type Commands struct {
	Commands []Command `json:"commands"`
}

// Checks if app directory is exists
func checkIfAppDirExists() (bool, error) {
	appDirPath, err := getAppDirPath()
	if err != nil {
		return false, err
	}
	if _, err := os.Stat(appDirPath); !os.IsNotExist(err) {
		return true, nil
	}
	return false, nil
}

// Checks if commands file exists
func checkIfCommandsFileExists() bool {
	commandsFilePath, err := getCommandsFilePath()
	if err != nil {
		return false
	}
	if _, err := os.Stat(commandsFilePath); !os.IsNotExist(err) {
		return true
	}
	return false
}

// Returns home directory path
func getHomeDirPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return home, nil
}

// Returns application directory path
func getAppDirPath() (string, error) {
	homeDir, err := getHomeDirPath()
	if err != nil {
		return "", err
	}
	return homeDir + "/" + appDirName, nil
}

// Returns commands file path
func getCommandsFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return homeDir + "/" + appDirName + "/" + commandsFileName, nil
}

// Returns saved commands
func getCommands() (*Commands, error) {
	var existingCommands *Commands
	commandsFilePath, err := getCommandsFilePath()
	if err != nil {
		return nil, err
	}
	file, _ := ioutil.ReadFile(commandsFilePath)
	err = json.Unmarshal(file, &existingCommands)
	if err != nil {
		return nil, err
	}
	return existingCommands, nil
}

// Creates app directory
func createAppDirectory(dirPath string) error {
	fmt.Printf("Initializing environment in ~/%s directory\n", appDirName)
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()
	time.Sleep(2 * time.Second)
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		s.Stop()
		return err
	}
	s.Stop()
	return nil
}

// Creates commands file in app directory
func createCommandsFile() error {
	commandsFilePath, err := getCommandsFilePath()
	if err != nil {
		return err
	}
	f, err := os.OpenFile(commandsFilePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

// Writes commands to commands file
func writeCommandsToFile(commands Commands) error {
	commandsJSON, err := json.MarshalIndent(commands, "", "")
	if err != nil {
		return err
	}
	commandsFilePath, err := getCommandsFilePath()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(commandsFilePath, commandsJSON, 0644)
	if err != nil {
		return err
	}
	return nil
}

// Search commands that contains concatenated string of args
func searchCommands(args []string) ([]string, error) {
	// concatenate args into the single string
	concatenatedArgs := strings.Join(args[:], " ")
	commands, err := getCommands()
	if err != nil {
		return nil, err
	}
	var concatenatedCommands, matches []string
	for _, command := range commands.Commands {
		concatenatedCommands = append(concatenatedCommands, command.Command+" :: "+command.Description+" :: "+command.Alias)
	}
	for _, concatenatedCommand := range concatenatedCommands {
		isMatched, err := regexp.MatchString(concatenatedArgs, concatenatedCommand)
		if err != nil {
			return nil, err
		} else if isMatched {
			matches = append(matches, concatenatedCommand)
		}
	}
	if len(matches) == 0 {
		return nil, err
	}
	return matches, nil
}

// Prints commands as table
func printCommandsAsTable(commands *Commands) {

	// convert commands to table.Row type
	var commandRow []table.Row
	for _, command := range commands.Commands {
		commandRow = append(commandRow, table.Row{command.Command, command.Description, command.Alias})
	}
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.Style().Options.SeparateRows = true
	t.AppendHeader(table.Row{"Command", "Description", "Alias"})
	t.AppendRows(commandRow)
	t.Render()
	return
}

// Prints commands with indexes
func printCommandsAsTableWithIndexes(commands *Commands) {

	// convert commands to table.Row type
	var commandRow []table.Row
	for index, command := range commands.Commands {
		commandRow = append(commandRow, table.Row{index + 1, command.Command, command.Description, command.Alias})
	}
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.Style().Options.SeparateRows = true
	t.AppendHeader(table.Row{"#", "Command", "Description", "Alias"})
	t.AppendRows(commandRow)
	t.Render()
	return
}

func isIntInSlice(i int, slice []int) bool {
	for _, a := range slice {
		if a == i {
			return true
		}
	}
	return false
}
