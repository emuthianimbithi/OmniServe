package main

import (
	"bytes"
	"fmt"
	"github.com/emuthianimbithi/OmniServe/pkg/cliconfig"
	"github.com/emuthianimbithi/OmniServe/pkg/stagedfiles"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
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

	for _, arg := range args {
		absPath, err := filepath.Abs(arg)
		if err != nil {
			fmt.Printf("Error getting absolute path for %s: %v\n", arg, err)
			continue
		}
		stagedFiles = append(stagedFiles, absPath)
	}

	err = stagedfiles.SaveStagedFiles(stagedFiles)
	if err != nil {
		fmt.Printf("Error saving staged files: %v\n", err)
		return
	}

	fmt.Printf("Added %d files to be pushed\n", len(args))
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

	// Add the project name
	writer.WriteField("project", filepath.Base(getCurrentDir()))

	// Add each file to the multipart writer
	for _, file := range stagedFiles {
		part, err := writer.CreateFormFile("file", filepath.Base(file))
		if err != nil {
			fmt.Printf("Error creating form file: %v\n", err)
			return
		}

		// Open and copy the file content
		f, err := os.Open(file)
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
			return
		}
		_, err = io.Copy(part, f)
		f.Close()
		if err != nil {
			fmt.Printf("Error copying file content: %v\n", err)
			return
		}
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
	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(respBody))

	// Clear staged files after successful push
	err = stagedfiles.ClearStagedFiles()
	if err != nil {
		fmt.Printf("Error clearing staged files: %v\n", err)
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
