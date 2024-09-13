package docker

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateDockerfile(directory, language, entryPoint string) error {
	// Ensure the directory ends with a slash (/) or backslash (\), depending on the OS
	dockerfilePath := filepath.Join(directory, "Dockerfile")

	var content string

	switch language {
	case "go":
		content = fmt.Sprintf(`FROM golang:latest
WORKDIR /app
COPY . .
RUN go build -o main %s
EXPOSE 8080
CMD ["/app/main"]`, entryPoint)

	case "python":
		content = `FROM python:3.11
WORKDIR /app
COPY . .
RUN pip install -r requirements.txt
EXPOSE 8080
CMD ["python", "app.py"]`

	case "nodejs":
		content = `FROM node:latest
WORKDIR /app
COPY . .
RUN npm install
EXPOSE 8080
CMD ["node", "app.js"]`

	case "c":
		content = fmt.Sprintf(`FROM gcc:latest
WORKDIR /app
COPY . .
RUN gcc -o main %s
EXPOSE 8080
CMD ["/app/main"]`, entryPoint)

	default:
		return fmt.Errorf("unsupported language: %s", language)
	}

	return os.WriteFile(dockerfilePath, []byte(content), 0644)
}
