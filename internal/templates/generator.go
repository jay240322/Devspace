package templates

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

type ProjectMetadata struct {
	TargetDir   string
	ServiceName string
	Backend     string
	Frontend    string
	GitHubUser  string
}

func GenerateBoilerplate(meta ProjectMetadata) error {
	basePath := filepath.Join(meta.TargetDir,meta.ServiceName)
	backendPath := filepath.Join(basePath, "backend")
	frontendPath := filepath.Join(basePath, "frontend")

	// Create physical folders
	if err := os.MkdirAll(backendPath, 0755); err != nil {
		return err
	}
	if meta.Frontend != "None (Pure Backend API)" {
		if err := os.MkdirAll(frontendPath, 0755); err != nil {
			return err
		}
	}

	// Ready-made dockerfile blueprint using BACKTICKS (``) for multi-line support
	dockerfileBlueprint := `FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o main .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
# Provisioned securely via DevSpace for github.com/{{.GitHubUser}}/{{.ServiceName}}`

	// Compile the dockerfile blueprint and inject the real user variables
	tmpl, err := template.New("dockerfile").Parse(dockerfileBlueprint)
	if err != nil {
		return fmt.Errorf("failed to parse dockerfile template: %v", err)
	}

	var processedCode bytes.Buffer
	if err := tmpl.Execute(&processedCode, meta); err != nil {
		return fmt.Errorf("failed to inject variables: %v", err)
	}

	// Write out the target file
	targetFile := filepath.Join(backendPath, "Dockerfile")
	if err := os.WriteFile(targetFile, processedCode.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write generated file: %v", err)
	}

	fmt.Printf("✅ Generated personalized Dockerfile at %s\n", targetFile)
	return nil
}