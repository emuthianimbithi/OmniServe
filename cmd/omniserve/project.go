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
	"log"
	"os"
	"path/filepath"
)

//var (
//	stagedFiles []string
//)

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
	utils.VerboseLog("Starting runAdd function")

	utils.VerboseLog("Initializing ignore list")
	err := utils.InitIgnoreList()
	if err != nil {
		utils.VerboseLog(fmt.Sprintf("Error initializing ignore list: %v", err))
		fmt.Printf("Error initializing ignore list: %v\n", err)
		return
	}

	utils.VerboseLog("Loading staged files")
	stagedFiles, err := stagedfiles.LoadStagedFiles()
	if err != nil {
		utils.VerboseLog(fmt.Sprintf("Error loading staged files: %v", err))
		fmt.Printf("Error loading staged files: %v\n", err)
		return
	}

	if len(args) == 0 {
		utils.VerboseLog("No arguments provided")
		fmt.Println("Please specify files or directories to add, or use '.' for the current directory.")
		return
	}

	addedFiles := 0
	baseDir, _ := os.Getwd()
	utils.VerboseLog(fmt.Sprintf("Base directory: %s", baseDir))

	for _, arg := range args {
		utils.VerboseLog(fmt.Sprintf("Processing argument: %s", arg))
		if arg == "." {
			arg = baseDir
			utils.VerboseLog("Argument is '.', using base directory")
		}

		err := filepath.Walk(arg, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				utils.VerboseLog(fmt.Sprintf("Error accessing path %q: %v", path, err))
				fmt.Printf("Error accessing path %q: %v\n", path, err)
				return err
			}

			relPath, err := filepath.Rel(baseDir, path)
			if err != nil {
				utils.VerboseLog(fmt.Sprintf("Error getting relative path for %s: %v", path, err))
				fmt.Printf("Error getting relative path for %s: %v\n", path, err)
				return nil
			}

			if !info.IsDir() && !utils.ShouldIgnore(relPath) {
				utils.VerboseLog(fmt.Sprintf("Checking file: %s", relPath))
				if !contains(stagedFiles, relPath) {
					utils.VerboseLog(fmt.Sprintf("Adding file to staged files: %s", relPath))
					stagedFiles = append(stagedFiles, relPath)
					addedFiles++
				} else {
					utils.VerboseLog(fmt.Sprintf("File already staged: %s", relPath))
				}
			} else if info.IsDir() {
				utils.VerboseLog(fmt.Sprintf("Skipping directory: %s", relPath))
			} else {
				utils.VerboseLog(fmt.Sprintf("Ignoring file: %s", relPath))
			}
			return nil
		})
		if err != nil {
			utils.VerboseLog(fmt.Sprintf("Error walking the path %q: %v", arg, err))
			fmt.Printf("Error walking the path %q: %v\n", arg, err)
		}
	}

	utils.VerboseLog("Saving staged files")
	err = stagedfiles.SaveStagedFiles(stagedFiles)
	if err != nil {
		utils.VerboseLog(fmt.Sprintf("Error saving staged files: %v", err))
		fmt.Printf("Error saving staged files: %v\n", err)
		return
	}

	utils.VerboseLog(fmt.Sprintf("Added %d files to be pushed", addedFiles))
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
	utils.VerboseLog("Starting runPush function")

	utils.VerboseLog("Loading staged files")
	stagedFiles, err := stagedfiles.LoadStagedFiles()
	if err != nil {
		utils.VerboseLog(fmt.Sprintf("Error loading staged files: %v", err))
		fmt.Printf("Error loading staged files: %v\n", err)
		return
	}

	if len(stagedFiles) == 0 {
		utils.VerboseLog("No files staged for pushing")
		fmt.Println("No files staged for pushing. Use 'omniserve add' to stage files.")
		return
	}

	utils.VerboseLog("Getting gRPC connection")
	conn := utils.GetGRPCConnection(cliconfig.CliConfig.Server.Host)

	utils.VerboseLog("Creating OmniServe client")
	client := omniserve_proto.NewOmniServeClient(conn.Conn)
	utils.VerboseLog("Creating PushFiles stream")
	stream, err := client.PushFiles(context.Background())
	if err != nil {
		utils.VerboseLog(fmt.Sprintf("Error creating stream: %v", err))
		log.Fatalf("Error creating stream: %v", err)
	}

	currentDir, _ := os.Getwd()
	projectCode := config.ProjectConfig.Code
	utils.VerboseLog(fmt.Sprintf("Current directory: %s, Project code: %s", currentDir, projectCode))

	for i, relPath := range stagedFiles {
		absPath := filepath.Join(currentDir, relPath)
		utils.VerboseLog(fmt.Sprintf("Processing file %d of %d: %s", i+1, len(stagedFiles), absPath))

		fileInfo, err := os.Stat(absPath)
		if err != nil {
			utils.VerboseLog(fmt.Sprintf("Error getting file info for %s: %v", absPath, err))
			fmt.Printf("Error getting file info for %s: %v\n", absPath, err)
			continue
		}

		if fileInfo.IsDir() {
			utils.VerboseLog(fmt.Sprintf("Skipping directory: %s", absPath))
			fmt.Printf("Skipping directory: %s\n", absPath)
			continue
		}

		utils.VerboseLog(fmt.Sprintf("Reading file: %s", absPath))
		content, err := os.ReadFile(absPath)
		if err != nil {
			utils.VerboseLog(fmt.Sprintf("Error reading file %s: %v", absPath, err))
			fmt.Printf("Error reading file %s: %v\n", absPath, err)
			continue
		}

		chunk := &omniserve_proto.FileChunk{
			ProjectCode: projectCode,
			FilePath:    relPath,
			Content:     content,
		}

		utils.VerboseLog(fmt.Sprintf("Sending chunk for file: %s", relPath))
		if err := stream.Send(chunk); err != nil {
			utils.VerboseLog(fmt.Sprintf("Error sending chunk: %v", err))
			log.Printf("Error sending chunk: %v\n", err)
		}

		fmt.Printf("Pushing file %d of %d: %s\n", i+1, len(stagedFiles), relPath)
	}

	utils.VerboseLog("Closing stream and receiving response")
	resp, err := stream.CloseAndRecv()
	if err != nil {
		utils.VerboseLog(fmt.Sprintf("Error receiving response: %v", err))
		log.Fatalf("Error receiving response: %v", err)
	}

	utils.VerboseLog(fmt.Sprintf("Server response: %s", resp.Message))
	fmt.Println(resp.Message)

	utils.VerboseLog("Clearing staged files")
	err = stagedfiles.ClearStagedFiles()
	if err != nil {
		utils.VerboseLog(fmt.Sprintf("Error clearing staged files: %v", err))
		fmt.Printf("Error clearing staged files: %v\n", err)
	} else {
		utils.VerboseLog("All staged files have been pushed and cleared")
		fmt.Println("All staged files have been pushed and cleared.")
	}

	utils.VerboseLog(fmt.Sprintf("Project saved. Project unique code: %s", projectCode))
	fmt.Println("Project saved. Your project unique code is: ", projectCode)
}

func runStatus(cmd *cobra.Command, args []string) {
	utils.VerboseLog("Starting runStatus function")

	utils.VerboseLog("Loading staged files")
	stagedFiles, err := stagedfiles.LoadStagedFiles()
	if err != nil {
		utils.VerboseLog(fmt.Sprintf("Error loading staged files: %v", err))
		fmt.Printf("Error loading staged files: %v\n", err)
		return
	}

	if len(stagedFiles) == 0 {
		utils.VerboseLog("No files staged for pushing")
		fmt.Println("No files staged for pushing.")
	} else {
		utils.VerboseLog(fmt.Sprintf("Found %d staged files", len(stagedFiles)))
		fmt.Println("Files staged for pushing:")
		for _, file := range stagedFiles {
			utils.VerboseLog(fmt.Sprintf("Staged file: %s", file))
			fmt.Println(file)
		}
	}
}
