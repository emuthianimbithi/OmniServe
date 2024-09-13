package main

import (
	"context"
	"fmt"
	"github.com/emuthianimbithi/OmniServe/pkg/cliconfig"
	"github.com/emuthianimbithi/OmniServe/pkg/pb/omniserve_proto"
	"github.com/emuthianimbithi/OmniServe/pkg/utils"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
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
	port := cliconfig.CliConfig.Server.Port
	if port == "" {
		port = "50051"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Set gRPC options to allow a max of 1GB for message size
	maxMessageSize := 1 << 30 // 1GB in bytes (1GB = 1073741824 bytes)

	s := grpc.NewServer(
		grpc.MaxRecvMsgSize(maxMessageSize),
		grpc.MaxSendMsgSize(maxMessageSize),
	)

	omniserve_proto.RegisterOmniServeServer(s, &server{})

	log.Printf("Server listening at %v with max message size of 1GB", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type server struct {
	omniserve_proto.UnimplementedOmniServeServer
}

func (s *server) PushFiles(stream omniserve_proto.OmniServe_PushFilesServer) error {
	var projectCode string
	var fileCount int

	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&omniserve_proto.PushResponse{
				Message: fmt.Sprintf("Successfully pushed %d files for project %s", fileCount, projectCode),
			})
		}
		if err != nil {
			return err
		}

		if projectCode == "" {
			projectCode = chunk.ProjectCode
		}

		fullPath := filepath.Join("./projects", projectCode, chunk.FilePath)
		dir := filepath.Dir(fullPath)

		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}

		f, err := os.Create(fullPath)
		if err != nil {
			return fmt.Errorf("failed to create file: %v", err)
		}

		_, err = f.Write(chunk.Content)
		f.Close()
		if err != nil {
			return fmt.Errorf("failed to write file: %v", err)
		}

		fileCount++
		fmt.Printf("Received file: %s\n", chunk.FilePath)
	}
}

func runServerStatus(cmd *cobra.Command, args []string) {
	serverAddress := cliconfig.CliConfig.Server.Host

	// Output the server address to verify
	fmt.Println("Attempting to connect to:", serverAddress)

	// Set context timeout for the connection
	_, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	conn := utils.GetGRPCConnection(serverAddress)
	if conn.Error != "" {
		fmt.Printf("Server connection failed: %v\n", conn.Error)
		return
	}
	defer conn.Conn.Close()

	// If the connection succeeds
	fmt.Println("Server is running and reachable.")
}
