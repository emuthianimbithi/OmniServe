package main

import (
	"bytes"
	"fmt"
	"github.com/emuthianimbithi/OmniServe/pkg/cliconfig"
	"github.com/emuthianimbithi/OmniServe/pkg/config"
	"github.com/emuthianimbithi/OmniServe/pkg/stagedfiles"
	"github.com/emuthianimbithi/OmniServe/pkg/utils"
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
	err := utils.InitIgnoreList()
	if err != nil {
		fmt.Printf("Error initializing ignore list: %v\n", err)
		return
	}

	stagedFiles, err := stagedfiles.LoadStagedFiles()
	if err != nil {
		fmt.Printf("Error loading staged files: %v\n", err)
		return
	}

	if len(args) == 0 {
		fmt.Println("Please specify files or directories to add, or use '.' for the current directory.")
		return
	}

	addedFiles := 0
	baseDir, _ := os.Getwd()

	for _, arg := range args {
		if arg == "." {
			arg = baseDir
		}

		err := filepath.Walk(arg, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("Error accessing path %q: %v\n", path, err)
				return err
			}

			relPath, err := filepath.Rel(baseDir, path)
			if err != nil {
				fmt.Printf("Error getting relative path for %s: %v\n", path, err)
				return nil
			}

			if !info.IsDir() && !utils.ShouldIgnore(relPath) {
				if !contains(stagedFiles, relPath) {
					stagedFiles = append(stagedFiles, relPath)
					addedFiles++
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

	fmt.Printf("Added %d files to be pushed\n", addedFiles)
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

	serverHost := cliconfig.CliConfig.Server.Host
	serverPort := cliconfig.CliConfig.Server.Port

	if serverHost == "" {
		serverHost = "localhost"
	}
	if serverPort == "" {
		serverPort = "8765"
	}

	url := fmt.Sprintf("http://%s:%s/push", serverHost, serverPort)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	currentDir, _ := os.Getwd()
	// create an uuid for the project. this will be used to know what project to start
	projectCode := config.ProjectConfig.Code
	err = writer.WriteField("project", projectCode)
	if err != nil {
		return
	}

	// Iterate through the staged files and send them with their full relative paths
	for i, relPath := range stagedFiles {
		absPath := filepath.Join(currentDir, relPath)

		fileInfo, err := os.Stat(absPath)
		if err != nil {
			fmt.Printf("Error getting file info for %s: %v\n", absPath, err)
			continue
		}

		if fileInfo.IsDir() {
			fmt.Printf("Skipping directory: %s\n", absPath)
			continue
		}

		file, err := os.Open(absPath)
		if err != nil {
			fmt.Printf("Error opening file %s: %v\n", absPath, err)
			continue
		}

		// Create a unique key for this file
		fileKey := fmt.Sprintf("file_%d", i)

		// Add the file path as a separate form field
		err = writer.WriteField(fmt.Sprintf("%s_path", fileKey), relPath)
		if err != nil {
			return
		}

		// Add the file content
		part, err := writer.CreateFormFile(fileKey, filepath.Base(relPath))
		if err != nil {
			fmt.Printf("Error creating form file for %s: %v\n", relPath, err)
			err = file.Close()
			if err != nil {
				return
			}
			continue
		}

		_, err = io.Copy(part, file)
		err = file.Close()
		if err != nil {
			return
		}

		fmt.Printf("Pushing file %d of %d: %s\n", i+1, len(stagedFiles), relPath)
	}
	err = writer.Close()
	if err != nil {
		return
	}

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
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}(resp.Body)

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Println(string(respBody))

	err = stagedfiles.ClearStagedFiles()
	if err != nil {
		fmt.Printf("Error clearing staged files: %v\n", err)
	} else {
		fmt.Println("All staged files have been pushed and cleared.")
	}

	fmt.Println("Project saved. Your project unique code is: ", projectCode)
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
