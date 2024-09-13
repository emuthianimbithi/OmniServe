package utils

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/emuthianimbithi/OmniServe/pkg/cliconfig"
	"github.com/emuthianimbithi/OmniServe/pkg/variables"
	"github.com/sabhiram/go-gitignore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var Verbose bool

var ignoreList *ignore.GitIgnore

func InitIgnoreList() error {
	var err error
	ignoreList, err = ignore.CompileIgnoreFile(".omniserve-ignore")
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func ShouldIgnore(path string) bool {
	if ignoreList == nil {
		return false
	}
	return ignoreList.MatchesPath(path)
}

func IsValidLanguage(lang string) bool {
	if langConfig, exists := cliconfig.CliConfig.Languages[lang]; exists {
		return langConfig.EntryPoint != ""
	}
	return variables.DefaultSupportedLanguages[lang]
}

func GetDefaultEntryPoint(language string) string {
	if langConfig, exists := cliconfig.CliConfig.Languages[language]; exists {
		return langConfig.EntryPoint
	}
	return variables.DefaultEntryPointTemplate[language]
}

// VerboseLog used to give more detailed output
func VerboseLog(message string) {
	if Verbose {
		fmt.Fprintln(os.Stderr, "VERBOSE:", message)
	}
}
func GetAllFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func GetFiles(path string) ([]string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return []string{path}, nil
	}
	return GetAllFiles(path)
}

type ExternalGRPCConnectionResponse struct {
	Error  string
	Status int64
	Conn   *grpc.ClientConn
}

func GetGRPCConnection(target string) ExternalGRPCConnectionResponse {
	// Check if the connection should be insecure
	_insecure := strings.Contains(target, "localhost") || strings.Contains(target, "0.0.0.0")

	// Remove the 'https://' or 'http://' prefix if present
	if strings.HasPrefix(target, "https://") {
		target = strings.TrimPrefix(target, "https://")
	} else if strings.HasPrefix(target, "http://") {
		target = strings.TrimPrefix(target, "http://")
	}

	// Ensure the target contains a port
	if !strings.Contains(target, ":") {
		target = target + ":443"
	}

	fmt.Println("Attempting to connect to:", target)

	var conn *grpc.ClientConn
	var err error

	maxMsgSize := 1024 * 1024 * 1024 // 1GB

	// Create gRPC connection with appropriate credentials
	if _insecure {
		// For insecure connections (e.g., localhost), use insecure credentials
		conn, err = grpc.Dial(target,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMsgSize)),
			grpc.WithReadBufferSize(maxMsgSize),
			grpc.WithWriteBufferSize(maxMsgSize),
			grpc.WithInitialWindowSize(int32(maxMsgSize)),
			grpc.WithInitialConnWindowSize(int32(maxMsgSize)))
	} else {
		// For secure connections, use TLS credentials
		systemRoots, sysErr := x509.SystemCertPool()
		if sysErr != nil {
			return ExternalGRPCConnectionResponse{
				Status: http.StatusInternalServerError,
				Error:  "Error fetching system root certificates: " + sysErr.Error(),
			}
		}
		cred := credentials.NewTLS(&tls.Config{
			RootCAs: systemRoots,
		})
		conn, err = grpc.Dial(target,
			grpc.WithTransportCredentials(cred),
			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMsgSize)),
			grpc.WithReadBufferSize(maxMsgSize),
			grpc.WithWriteBufferSize(maxMsgSize),
			grpc.WithInitialWindowSize(int32(maxMsgSize)),
			grpc.WithInitialConnWindowSize(int32(maxMsgSize)))
	}

	if err != nil {
		return ExternalGRPCConnectionResponse{
			Status: http.StatusInternalServerError,
			Error:  "Error connecting to gRPC server: " + err.Error(),
		}
	}

	return ExternalGRPCConnectionResponse{
		Status: http.StatusOK,
		Conn:   conn,
	}
}
