package main

import (
	"bufio"
	"fmt"
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
	if _, err := os.Stat(ignoreFileName); err == nil {
		fmt.Printf("%s already exists. Use 'omniserve ignore add' to add patterns.\n", ignoreFileName)
		return
	}

	file, err := os.Create(ignoreFileName)
	if err != nil {
		fmt.Printf("Error creating %s: %v\n", ignoreFileName, err)
		return
	}
	defer file.Close()

	if !noDefaultPatterns {
		for _, pattern := range variables.DefaultIgnorePatterns {
			if _, err := file.WriteString(pattern + "\n"); err != nil {
				fmt.Printf("Error writing default pattern: %v\n", err)
				return
			}
		}
		fmt.Printf("%s created successfully with default patterns.\n", ignoreFileName)
	} else {
		fmt.Printf("%s created successfully (empty).\n", ignoreFileName)
	}
}

func runIgnoreAdd(cmd *cobra.Command, args []string) {
	file, err := os.OpenFile(ignoreFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening %s: %v\n", ignoreFileName, err)
		return
	}
	defer file.Close()

	for _, pattern := range args {
		if _, err := file.WriteString(pattern + "\n"); err != nil {
			fmt.Printf("Error writing pattern '%s': %v\n", pattern, err)
		} else {
			fmt.Printf("Added pattern: %s\n", pattern)
		}
	}
}

func runIgnoreList(cmd *cobra.Command, args []string) {
	file, err := os.Open(ignoreFileName)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("%s does not exist. Use 'omniserve ignore init' to create it.\n", ignoreFileName)
		} else {
			fmt.Printf("Error opening %s: %v\n", ignoreFileName, err)
		}
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	fmt.Printf("Contents of %s:\n", ignoreFileName)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			fmt.Println(line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading %s: %v\n", ignoreFileName, err)
	}
}
