package main

import (
	"fmt"
	"github.com/emuthianimbithi/OmniServe/pkg/cliconfig"
	"github.com/emuthianimbithi/OmniServe/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	Version    string // Will be set at build time
	configPath string
)

var rootCmd = &cobra.Command{
	Use:   "omniserve",
	Short: "OmniServe - A multi-language serverless platform CLI",
	Long: `OmniServe is a powerful CLI tool for creating and managing serverless projects across multiple programming languages.
Complete documentation is available at https://github.com/emuthianimbithi/OmniServe`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if version flag is used
		versionFlag, _ := cmd.Flags().GetBool("version")
		if versionFlag {
			fmt.Printf("OmniServe version %s\n", Version)
			return
		}

		// Check if info flag is used
		infoFlag, _ := cmd.Flags().GetBool("info")
		if infoFlag {
			printInfo()
			return
		}

		// If no flags are used, print help
		cmd.Help()
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolVarP(&utils.Verbose, "verbose", "v", false, "enable verbose output")
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "config file (default is $HOME/.omniserve.yaml)")
	rootCmd.Flags().Bool("version", false, "print the version number of OmniServe")
	rootCmd.Flags().Bool("info", false, "print information about OmniServe")
}

func Execute() error {
	return rootCmd.Execute()
}

func printInfo() {
	fmt.Println("OmniServe - A multi-language serverless platform CLI")
	fmt.Printf("Version: %s\n", Version)
	fmt.Println("Author: Emmanuel Muthiani Mbithi")
	fmt.Println("License: emmanuel.ke")
	fmt.Println("Repository: https://github.com/emuthianimbithi/OmniServe")
	fmt.Println("Supported Languages: Go, C, Python, JavaScript")
}

func initConfig() {
	err := cliconfig.LoadConfig(configPath)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
	}
}
