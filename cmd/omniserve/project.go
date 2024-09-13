package main

import (
	"context"
	"fmt"
	"github.com/emuthianimbithi/OmniServe/pkg/cliconfig"
	"github.com/emuthianimbithi/OmniServe/pkg/config"
	"github.com/emuthianimbithi/OmniServe/pkg/pb/omniserve_proto"
	"github.com/emuthianimbithi/OmniServe/pkg/stagedfiles"
	"github.com/emuthianimbithi/OmniServe/pkg/utils"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
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

	conn := utils.GetGRPCConnection(cliconfig.CliConfig.Server.Host)

	client := omniserve_proto.NewOmniServeClient(conn.Conn)
	stream, err := client.PushFiles(context.Background())
	if err != nil {
		log.Fatalf("Error creating stream: %v", err)
	}

	currentDir, _ := os.Getwd()
	projectCode := config.ProjectConfig.Code

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

		content, err := ioutil.ReadFile(absPath)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", absPath, err)
			continue
		}

		chunk := &omniserve_proto.FileChunk{
			ProjectCode: projectCode,
			FilePath:    relPath,
			Content:     content,
		}

		if err := stream.Send(chunk); err != nil {
			log.Printf("Error sending chunk: %v\n", err)
		}

		fmt.Printf("Pushing file %d of %d: %s\n", i+1, len(stagedFiles), relPath)
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving response: %v", err)
	}

	fmt.Println(resp.Message)

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
