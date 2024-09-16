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
		_, err := fmt.Fprintln(os.Stderr, "VERBOSE:", message)
		if err != nil {
			return
		}
	}
}

type ExternalGRPCConnectionResponse struct {
	Error  string
	Status int64
	Conn   *grpc.ClientConn
}

func GetGRPCConnection(target string) ExternalGRPCConnectionResponse {
	VerboseLog(fmt.Sprintf("Starting GetGRPCConnection with target: %s", target))

	isInsecure := strings.Contains(target, "localhost") || strings.Contains(target, "0.0.0.0") || strings.Contains(target, "ngrok.io")
	VerboseLog(fmt.Sprintf("Connection is insecure: %v", isInsecure))

	// Remove the 'https://' or 'http://' prefix if present
	target = strings.TrimPrefix(strings.TrimPrefix(target, "https://"), "http://")
	VerboseLog(fmt.Sprintf("Target after prefix removal: %s", target))

	// Only add default port if there isn't one already
	if !strings.Contains(target, ":") {
		target += ":443"
		VerboseLog("Default port :443 added to target")
	}

	fmt.Println("Attempting to connect to:", target)
	VerboseLog(fmt.Sprintf("Final target for connection: %s", target))

	maxMsgSize := 1024 * 1024 * 1024 // 1GB
	VerboseLog(fmt.Sprintf("Setting max message size to %d bytes", maxMsgSize))

	dialOptions := []grpc.DialOption{
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMsgSize)),
		grpc.WithReadBufferSize(maxMsgSize),
		grpc.WithWriteBufferSize(maxMsgSize),
		grpc.WithInitialWindowSize(int32(maxMsgSize)),
		grpc.WithInitialConnWindowSize(int32(maxMsgSize)),
	}
	VerboseLog("Basic dial options set")

	if isInsecure {
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))
		VerboseLog("Using insecure credentials")
	} else {
		VerboseLog("Attempting to use secure credentials")
		systemRoots, err := x509.SystemCertPool()
		if err != nil {
			errMsg := fmt.Sprintf("Error fetching system root certificates: %v", err)
			VerboseLog(errMsg)
			return ExternalGRPCConnectionResponse{
				Status: http.StatusInternalServerError,
				Error:  errMsg,
			}
		}
		VerboseLog("System root certificates fetched successfully")

		cred := credentials.NewTLS(&tls.Config{
			RootCAs:            systemRoots,
			InsecureSkipVerify: strings.Contains(target, "ngrok.io"),
		})
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(cred))
		VerboseLog(fmt.Sprintf("TLS credentials created. InsecureSkipVerify: %v", strings.Contains(target, "ngrok.io")))
	}

	VerboseLog("Attempting to establish gRPC connection")
	conn, err := grpc.Dial(target, dialOptions...)
	if err != nil {
		errMsg := fmt.Sprintf("Error connecting to gRPC server: %v", err)
		VerboseLog(errMsg)
		return ExternalGRPCConnectionResponse{
			Status: http.StatusInternalServerError,
			Error:  errMsg,
		}
	}
	VerboseLog("gRPC connection established successfully")

	VerboseLog("Returning successful ExternalGRPCConnectionResponse")
	return ExternalGRPCConnectionResponse{
		Status: http.StatusOK,
		Conn:   conn,
	}
}
