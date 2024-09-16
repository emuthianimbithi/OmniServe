package main

import (
	"context"
	"fmt"
	"github.com/emuthianimbithi/OmniServe/pkg/cliconfig"
	"github.com/emuthianimbithi/OmniServe/pkg/pb/omniserve_proto"
	"github.com/emuthianimbithi/OmniServe/pkg/utils"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Manage the OmniServe server service",
}

var startServerCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the OmniServe server",
	Run:   runStartServer,
}

var serverStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check the status of the OmniServe server",
	Run:   runServerStatus,
}

func init() {
	serverCmd.AddCommand(startServerCmd)
	serverCmd.AddCommand(serverStatusCmd)
	rootCmd.AddCommand(serverCmd)
}

func runStartServer(cmd *cobra.Command, args []string) {
	utils.VerboseLog("Starting runStartServer function")

	port := cliconfig.CliConfig.Server.Port
	if port == "" {
		port = "8765"
		utils.VerboseLog("No port specified in config, using default port: 8765")
	} else {
		utils.VerboseLog(fmt.Sprintf("Using port from config: %s", port))
	}

	utils.VerboseLog(fmt.Sprintf("Attempting to listen on port %s", port))
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		utils.VerboseLog(fmt.Sprintf("Failed to listen: %v", err))
		log.Fatalf("failed to listen: %v", err)
	}
	utils.VerboseLog(fmt.Sprintf("Successfully listening on %v", lis.Addr()))

	// Set gRPC options to allow a max of 1GB for message size
	maxMessageSize := 1 << 30 // 1GB in bytes
	utils.VerboseLog(fmt.Sprintf("Setting max message size to %d bytes (1GB)", maxMessageSize))

	utils.VerboseLog("Creating new gRPC server with max message size options")
	s := grpc.NewServer(
		grpc.MaxRecvMsgSize(maxMessageSize),
		grpc.MaxSendMsgSize(maxMessageSize),
	)

	utils.VerboseLog("Registering OmniServeServer")
	omniserve_proto.RegisterOmniServeServer(s, &server{})

	// Register the health service
	utils.VerboseLog("Creating and registering health service")
	healthSrv := health.NewServer()
	healthSrv.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(s, healthSrv)

	utils.VerboseLog(fmt.Sprintf("Server initialization complete. Listening at %v with max message size of 1GB", lis.Addr()))
	log.Printf("Server listening at %v with max message size of 1GB", lis.Addr())

	utils.VerboseLog("Starting to serve gRPC requests")
	if err := s.Serve(lis); err != nil {
		utils.VerboseLog(fmt.Sprintf("Failed to serve: %v", err))
		log.Fatalf("failed to serve: %v", err)
	}
}

type server struct {
	omniserve_proto.UnimplementedOmniServeServer
}

func (s *server) PushFiles(stream omniserve_proto.OmniServe_PushFilesServer) error {
	var projectCode string
	var fileCount int

	utils.VerboseLog("Starting PushFiles method")

	for {
		utils.VerboseLog("Waiting to receive file chunk")
		chunk, err := stream.Recv()
		if err == io.EOF {
			utils.VerboseLog(fmt.Sprintf("Received EOF. Total files pushed: %d", fileCount))
			return stream.SendAndClose(&omniserve_proto.PushResponse{
				Message: fmt.Sprintf("Successfully pushed %d files for project %s", fileCount, projectCode),
			})
		}
		if err != nil {
			utils.VerboseLog(fmt.Sprintf("Error receiving chunk: %v", err))
			return err
		}

		if projectCode == "" {
			projectCode = chunk.ProjectCode
			utils.VerboseLog(fmt.Sprintf("Set project code to: %s", projectCode))
		}

		fullPath := filepath.Join("./projects", projectCode, chunk.FilePath)
		dir := filepath.Dir(fullPath)

		utils.VerboseLog(fmt.Sprintf("Creating directory: %s", dir))
		if err := os.MkdirAll(dir, 0755); err != nil {
			utils.VerboseLog(fmt.Sprintf("Failed to create directory: %v", err))
			return fmt.Errorf("failed to create directory: %v", err)
		}

		utils.VerboseLog(fmt.Sprintf("Creating file: %s", fullPath))
		f, err := os.Create(fullPath)
		if err != nil {
			utils.VerboseLog(fmt.Sprintf("Failed to create file: %v", err))
			return fmt.Errorf("failed to create file: %v", err)
		}

		utils.VerboseLog(fmt.Sprintf("Writing content to file: %s", fullPath))
		_, err = f.Write(chunk.Content)
		if err != nil {
			utils.VerboseLog(fmt.Sprintf("Failed to write file content: %v", err))
			err := f.Close()
			if err != nil {
				return err
			}
			return fmt.Errorf("failed to write file: %v", err)
		}

		err = f.Close()
		if err != nil {
			utils.VerboseLog(fmt.Sprintf("Failed to close file: %v", err))
			return fmt.Errorf("failed to close file: %v", err)
		}

		fileCount++
		utils.VerboseLog(fmt.Sprintf("Successfully processed file: %s. Total files: %d", chunk.FilePath, fileCount))
		fmt.Printf("Received file: %s\n", chunk.FilePath)
	}
}

func runServerStatus(cmd *cobra.Command, args []string) {
	serverAddress := cliconfig.CliConfig.Server.Host

	utils.VerboseLog(fmt.Sprintf("Starting server status check for address: %s", serverAddress))
	fmt.Println("Attempting to connect to:", serverAddress)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	utils.VerboseLog("Created context with 15-second timeout")

	utils.VerboseLog("Calling GetGRPCConnection")
	conn := utils.GetGRPCConnection(serverAddress)
	if conn.Error != "" {
		errMsg := fmt.Sprintf("Server connection failed: %v", conn.Error)
		utils.VerboseLog(errMsg)
		fmt.Println(errMsg)
		return
	}
	utils.VerboseLog("Successfully established gRPC connection")
	defer func(Conn *grpc.ClientConn) {
		err := Conn.Close()
		if err != nil {
			utils.VerboseLog(fmt.Sprintf("Failed to close connection: %v", err))
		}
	}(conn.Conn)

	// Use the context when checking the connection
	utils.VerboseLog("Getting initial connection state")
	state := conn.Conn.GetState()
	utils.VerboseLog(fmt.Sprintf("Initial connection state: %s", state))

	utils.VerboseLog("Waiting for state change")
	if !conn.Conn.WaitForStateChange(ctx, state) {
		utils.VerboseLog("Connection state did not change within the timeout period")
		fmt.Println("Connection state did not change within the timeout period.")
		return
	}
	utils.VerboseLog("State change detected")

	newState := conn.Conn.GetState()
	utils.VerboseLog(fmt.Sprintf("New connection state: %s", newState))

	if newState == connectivity.Ready {
		utils.VerboseLog("Server is running and reachable")
		fmt.Println("Server is running and reachable.")
	} else {
		utils.VerboseLog(fmt.Sprintf("Server is not ready. Current state: %s", newState))
		fmt.Printf("Server is not ready. Current state: %s\n", newState)
	}
}
