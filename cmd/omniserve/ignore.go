package main

import (
	"bufio"
	"fmt"
	"github.com/emuthianimbithi/OmniServe/pkg/utils"
	"github.com/emuthianimbithi/OmniServe/pkg/variables"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

const ignoreFileName = ".omniserve-ignore"

var (
	ignoreCmd = &cobra.Command{
		Use:   "ignore",
		Short: "Manage ignore patterns for OmniServe",
	}

	ignoreInitCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize the .omniserve-ignore file",
		Run:   runIgnoreInit,
	}

	ignoreAddCmd = &cobra.Command{
		Use:   "add [patterns...]",
		Short: "Add patterns to the .omniserve-ignore file",
		Args:  cobra.MinimumNArgs(1),
		Run:   runIgnoreAdd,
	}

	ignoreListCmd = &cobra.Command{
		Use:   "list",
		Short: "List the contents of the .omniserve-ignore file",
		Run:   runIgnoreList,
	}

	noDefaultPatterns bool
)

func init() {
	ignoreInitCmd.Flags().BoolVar(&noDefaultPatterns, "no-defaults", false, "Initialize without default ignore patterns")

	ignoreCmd.AddCommand(ignoreInitCmd)
	ignoreCmd.AddCommand(ignoreAddCmd)
	ignoreCmd.AddCommand(ignoreListCmd)
	rootCmd.AddCommand(ignoreCmd)
}

func runIgnoreInit(cmd *cobra.Command, args []string) {
	InitProjectConfig()
	cmd.Print("Initializing ignore file\n.Number o")

	utils.VerboseLog("Starting runIgnoreInit function")
	utils.VerboseLog(fmt.Sprintf("Ignore file name: %s", ignoreFileName))

	if _, err := os.Stat(ignoreFileName); err == nil {
		utils.VerboseLog(fmt.Sprintf("%s already exists", ignoreFileName))
		fmt.Printf("%s already exists. Use 'omniserve ignore add' to add patterns.\n", ignoreFileName)
		return
	}
	utils.VerboseLog(fmt.Sprintf("%s does not exist, proceeding with creation", ignoreFileName))

	utils.VerboseLog(fmt.Sprintf("Attempting to create %s", ignoreFileName))
	file, err := os.Create(ignoreFileName)
	if err != nil {
		utils.VerboseLog(fmt.Sprintf("Error creating %s: %v", ignoreFileName, err))
		fmt.Printf("Error creating %s: %v\n", ignoreFileName, err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			utils.VerboseLog(fmt.Sprintf("Error closing %s: %v", ignoreFileName, err))
			fmt.Printf("Error closing %s: %v\n", ignoreFileName, err)
		}
	}(file)
	utils.VerboseLog(fmt.Sprintf("%s created successfully", ignoreFileName))

	if !noDefaultPatterns {
		utils.VerboseLog("Adding default patterns to ignore file")
		for _, pattern := range variables.DefaultIgnorePatterns {
			utils.VerboseLog(fmt.Sprintf("Writing pattern: %s", pattern))
			if _, err := file.WriteString(pattern + "\n"); err != nil {
				utils.VerboseLog(fmt.Sprintf("Error writing default pattern: %v", err))
				fmt.Printf("Error writing default pattern: %v\n", err)
				return
			}
		}
		utils.VerboseLog("Default patterns added successfully")
		fmt.Printf("%s created successfully with default patterns.\n", ignoreFileName)
	} else {
		utils.VerboseLog("No default patterns flag set, creating empty ignore file")
		fmt.Printf("%s created successfully (empty).\n", ignoreFileName)
	}
}

func runIgnoreAdd(cmd *cobra.Command, args []string) {
	InitProjectConfig()
	utils.VerboseLog("Starting runIgnoreAdd function")
	utils.VerboseLog(fmt.Sprintf("Ignore file name: %s", ignoreFileName))

	utils.VerboseLog(fmt.Sprintf("Attempting to open %s in append mode", ignoreFileName))
	file, err := os.OpenFile(ignoreFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		utils.VerboseLog(fmt.Sprintf("Error opening %s: %v", ignoreFileName, err))
		fmt.Printf("Error opening %s: %v\n", ignoreFileName, err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			utils.VerboseLog(fmt.Sprintf("Error closing %s: %v", ignoreFileName, err))
			fmt.Printf("Error closing %s: %v\n", ignoreFileName, err)
		}
	}(file)
	utils.VerboseLog(fmt.Sprintf("%s opened successfully", ignoreFileName))

	for _, pattern := range args {
		utils.VerboseLog(fmt.Sprintf("Attempting to write pattern: %s", pattern))
		if _, err := file.WriteString(pattern + "\n"); err != nil {
			utils.VerboseLog(fmt.Sprintf("Error writing pattern '%s': %v", pattern, err))
			fmt.Printf("Error writing pattern '%s': %v\n", pattern, err)
		} else {
			utils.VerboseLog(fmt.Sprintf("Successfully added pattern: %s", pattern))
			fmt.Printf("Added pattern: %s\n", pattern)
		}
	}
}

func runIgnoreList(cmd *cobra.Command, args []string) {
	InitProjectConfig()
	utils.VerboseLog("Starting runIgnoreList function")
	utils.VerboseLog(fmt.Sprintf("Ignore file name: %s", ignoreFileName))

	utils.VerboseLog(fmt.Sprintf("Attempting to open %s", ignoreFileName))
	file, err := os.Open(ignoreFileName)
	if err != nil {
		if os.IsNotExist(err) {
			utils.VerboseLog(fmt.Sprintf("%s does not exist", ignoreFileName))
			fmt.Printf("%s does not exist. Use 'omniserve ignore init' to create it.\n", ignoreFileName)
		} else {
			utils.VerboseLog(fmt.Sprintf("Error opening %s: %v", ignoreFileName, err))
			fmt.Printf("Error opening %s: %v\n", ignoreFileName, err)
		}
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			utils.VerboseLog(fmt.Sprintf("Error closing %s: %v", ignoreFileName, err))
			fmt.Printf("Error closing %s: %v\n", ignoreFileName, err)
		}
	}(file)
	utils.VerboseLog(fmt.Sprintf("%s opened successfully", ignoreFileName))

	scanner := bufio.NewScanner(file)
	fmt.Printf("Contents of %s:\n", ignoreFileName)
	utils.VerboseLog("Starting to read file contents")
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		utils.VerboseLog(fmt.Sprintf("Read line: %s", line))
		if line != "" && !strings.HasPrefix(line, "#") {
			utils.VerboseLog(fmt.Sprintf("Printing non-empty, non-comment line: %s", line))
			fmt.Println(line)
		} else {
			utils.VerboseLog(fmt.Sprintf("Skipping empty or comment line: %s", line))
		}
	}

	if err := scanner.Err(); err != nil {
		utils.VerboseLog(fmt.Sprintf("Error reading %s: %v", ignoreFileName, err))
		fmt.Printf("Error reading %s: %v\n", ignoreFileName, err)
	} else {
		utils.VerboseLog("File read successfully")
	}
}
