package main

import (
	"bytes"
	"fmt"
	"github.com/emuthianimbithi/OmniServe/pkg/cliconfig"
	"github.com/emuthianimbithi/OmniServe/pkg/stagedfiles"
	"github.com/spf13/cobra"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

var (
	stagedFiles []string
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the status of staged files",
	Run:   runStatus,
}

var addCmd = &cobra.Command{
	Use:   "add [files...]",
	Short: "Add files to be pushed to the server",
	Run:   runAdd,
}

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push staged files to the server",
	Run:   runPush,
}

func init() {
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(pushCmd)
	rootCmd.AddCommand(statusCmd)
}

func runAdd(cmd *cobra.Command, args []string) {
	stagedFiles, err := stagedfiles.LoadStagedFiles()
	if err != nil {
		fmt.Printf("Error loading staged files: %v\n", err)
		return
	}

	if len(args) == 0 {
		fmt.Println("Please specify files or directories to add, or use '.' for the current directory.")
		return
	}

	for _, arg := range args {
		if arg == "." {
			// Handle the special case of adding the current directory
			arg, _ = os.Getwd()
		}

		err := filepath.Walk(arg, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("Error accessing path %q: %v\n", path, err)
				return err
			}
			if !info.IsDir() {
				absPath, err := filepath.Abs(path)
				if err != nil {
					fmt.Printf("Error getting absolute path for %s: %v\n", path, err)
					return nil
				}
				// Check if the file is already staged
				if !contains(stagedFiles, absPath) {
					stagedFiles = append(stagedFiles, absPath)
				}
			}
			return nil
		})
		if err != nil {
			fmt.Printf("Error walking the path %q: %v\n", arg, err)
		}
	}

	err = stagedfiles.SaveStagedFiles(stagedFiles)
	if err != nil {
		fmt.Printf("Error saving staged files: %v\n", err)
		return
	}

	fmt.Printf("Added %d files to be pushed\n", len(stagedFiles))
}

// Helper function to check if a slice contains a string
func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

func runPush(cmd *cobra.Command, args []string) {
	stagedFiles, err := stagedfiles.LoadStagedFiles()
	if err != nil {
		fmt.Printf("Error loading staged files: %v\n", err)
		return
	}

	if len(stagedFiles) == 0 {
		fmt.Println("No files staged for pushing. Use 'omniserve add' to stage files.")
		return
	}

	serverHost := cliconfig.Config.Server.Host
	serverPort := cliconfig.Config.Server.Port

	if serverHost == "" {
		serverHost = "localhost"
	}
	if serverPort == "" {
		serverPort = "8765"
	}

	url := fmt.Sprintf("http://%s:%s/push", serverHost, serverPort)

	// Create a new multipart writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the project name (using current directory name as project name)
	currentDir, _ := os.Getwd()
	projectName := filepath.Base(currentDir)
	writer.WriteField("project", projectName)

	// Add each file to the form
	for i, filePath := range stagedFiles {
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Printf("Error opening file %s: %v\n", filePath, err)
			continue
		}

		part, err := writer.CreateFormFile("file", filePath)
		if err != nil {
			fmt.Printf("Error creating form file for %s: %v\n", filePath, err)
			file.Close()
			continue
		}

		_, err = io.Copy(part, file)
		file.Close()
		if err != nil {
			fmt.Printf("Error copying content of file %s: %v\n", filePath, err)
			continue
		}

		fmt.Printf("Pushing file %d of %d: %s\n", i+1, len(stagedFiles), filePath)
	}

	writer.Close()

	// Create and send the request
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Read and print the response
	respBody, _ := io.ReadAll(resp.Body)
	fmt.Println(string(respBody))

	// Clear staged files after successful push
	err = stagedfiles.ClearStagedFiles()
	if err != nil {
		fmt.Printf("Error clearing staged files: %v\n", err)
	} else {
		fmt.Println("All staged files have been pushed and cleared.")
	}
}
func runStatus(cmd *cobra.Command, args []string) {
	stagedFiles, err := stagedfiles.LoadStagedFiles()
	if err != nil {
		fmt.Printf("Error loading staged files: %v\n", err)
		return
	}

	if len(stagedFiles) == 0 {
		fmt.Println("No files staged for pushing.")
	} else {
		fmt.Println("Files staged for pushing:")
		for _, file := range stagedFiles {
			fmt.Println(file)
		}
	}
}

func getCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		return ""
	}
	return dir
}
