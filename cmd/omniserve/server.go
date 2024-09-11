package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/emuthianimbithi/OmniServe/pkg/cliconfig"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Manage the OmniServe server service",
}

var startServerCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the OmniServe server",
	Run: func(cmd *cobra.Command, args []string) {
		port := cliconfig.Config.Server.Port
		if port == "" {
			port = "8765"
		}
		if err := startServer(port); err != nil {
			fmt.Printf("Failed to start server: %v\n", err)
		}
	},
}

var stopServerCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the OmniServe server",
	Run: func(cmd *cobra.Command, args []string) {
		if err := stopServer(); err != nil {
			fmt.Printf("Failed to stop server: %v\n", err)
		}
	},
}

var restartServerCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart the OmniServe server",
	Run: func(cmd *cobra.Command, args []string) {
		if err := restartServer(); err != nil {
			fmt.Printf("Failed to restart server: %v\n", err)
		}
	},
}

func init() {
	serverCmd.AddCommand(startServerCmd, stopServerCmd, restartServerCmd)
	rootCmd.AddCommand(serverCmd)
}

func startServer(port string) error {
	fmt.Printf("Starting server on port %s\n", port)
	http.HandleFunc("/push", handlePush)
	return http.ListenAndServe(":"+port, nil)
}

func stopServer() error {
	fmt.Println("Stopping server")
	// Implement actual server stop logic here
	return nil
}

func restartServer() error {
	if err := stopServer(); err != nil {
		return err
	}
	port := cliconfig.Config.Server.Port
	if port == "" {
		port = "8765"
	}
	return startServer(port)
}

func handlePush(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the multipart form data
	err := r.ParseMultipartForm(10 << 20) // 10 MB max memory
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the project details
	savedProject := r.FormValue("project")
	if savedProject == "" {
		http.Error(w, "Project name is required", http.StatusBadRequest)
		return
	}

	// Create the project directory if it doesn't exist
	projectDir := filepath.Join("./projects", savedProject)
	if err := os.MkdirAll(projectDir, os.ModePerm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Process each file in the form
	for _, fileHeaders := range r.MultipartForm.File {
		for _, fileHeader := range fileHeaders {
			file, err := fileHeader.Open()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer file.Close()

			// Create the destination file
			dst, err := os.Create(filepath.Join(projectDir, fileHeader.Filename))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer dst.Close()

			// Copy the uploaded file to the destination file
			_, err = io.Copy(dst, file)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	response := map[string]string{"message": fmt.Sprintf("Files for project '%s' pushed successfully", savedProject)}
	json.NewEncoder(w).Encode(response)
}
