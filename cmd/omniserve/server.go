package main

import (
	"fmt"
	"github.com/emuthianimbithi/OmniServe/pkg/cliconfig"
	"github.com/spf13/cobra"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Manage the OmniServe server service",
}

var startServerCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the OmniServe server",
	Run: func(cmd *cobra.Command, args []string) {
		port := cliconfig.CliConfig.Server.Port
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
	port := cliconfig.CliConfig.Server.Port
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

	err := r.ParseMultipartForm(32 << 20) // 32MB max memory
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	projectName := r.FormValue("project")
	if projectName == "" {
		http.Error(w, "Project name is required", http.StatusBadRequest)
		return
	}

	projectDir := filepath.Join("./projects", projectName)

	// Ensure the base project directory exists
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Iterate through the form data
	for key, _ := range r.MultipartForm.File {
		if !strings.HasPrefix(key, "file_") {
			continue // Skip non-file fields
		}

		// Get the corresponding path
		pathKey := key + "_path"
		paths, exists := r.MultipartForm.Value[pathKey]
		if !exists || len(paths) == 0 {
			http.Error(w, fmt.Sprintf("Missing path for %s", key), http.StatusBadRequest)
			return
		}
		relPath := paths[0]

		// Get the file
		files := r.MultipartForm.File[key]
		if len(files) == 0 {
			http.Error(w, fmt.Sprintf("Missing file for %s", key), http.StatusBadRequest)
			return
		}
		fileHeader := files[0]

		// Open the file
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer func(file multipart.File) {
			err := file.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}(file)

		// Construct the full path
		fullPath := filepath.Join(projectDir, relPath)

		// Ensure the relative path doesn't try to escape the project directory
		if !strings.HasPrefix(filepath.Clean(fullPath), projectDir) {
			http.Error(w, "Invalid file path", http.StatusBadRequest)
			return
		}

		// Ensure necessary parent directories are created
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Create and save the file
		dst, err := os.Create(fullPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer func(dst *os.File) {
			err := dst.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}(dst)

		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Printf("File saved: %s\n", fullPath)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Files for project '%s' pushed successfully", projectName)
}
