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
	basePath := filepath.Join(meta.TargetDir, meta.ServiceName)

	// 1. Physically construct the root service directory
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return err
	}

	// 2. Launch the Automated Backend CLIs
	switch meta.Backend {
	case "Python (Django)":
		GeneratePythonBackend(basePath, meta)
	case "Node.js (Express)":
		GenerateNodeBackend(basePath, meta)
	case "Rust (Actix-web)":
		GenerateRustBackend(basePath, meta)
	default:
		GenerateGoBackend(basePath, meta)
	}

	// 3. Let the frontend CLI handle its own folder creation cleanly
	if meta.Frontend != "None (Pure Backend API)" {
		frontendPath := filepath.Join(basePath, "frontend")
		GenerateFrontendFramework(frontendPath, "", meta)
	}

	// 4. Generate custom full-stack Dockerfile on the fly
	err := generateDynamicDockerfile(basePath, meta)
	if err != nil {
		fmt.Printf("⚠️  Dockerfile compilation skipped: %v\n", err)
	}

	fmt.Printf("\n✨ Successfully populated dynamic structure inside: %s\n", basePath)
	return nil
}

// Dynamic Dockerfile Builder Factory
func generateDynamicDockerfile(basePath string, meta ProjectMetadata) error {
	var dockerfileContent string

	// Step A: Append Frontend Builder Stage if a frontend exists (Explicitly stating the language name)
	if meta.Frontend != "None (Pure Backend API)" {
		dockerfileContent += fmt.Sprintf(`# --- Stage 1: Dynamic Frontend Builder Layer (%s) ---
FROM node:20-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ ./
RUN npm run build --if-present

`, meta.Frontend)
	}

	// Step B: Append the exact tailored Backend deployment environment stage
	switch meta.Backend {
	case "Python (Django)":
		dockerfileContent += `# --- Stage 2: Python Django Production Environment ---
FROM python:3.11-slim
WORKDIR /app
RUN pip install django django-cors-headers
COPY backend/ ./backend/
`
		if meta.Frontend != "None (Pure Backend API)" {
			dockerfileContent += "COPY --from=frontend-builder /app/frontend/dist ./frontend/dist\n"
		}
		dockerfileContent += "EXPOSE 8080\nCMD [\"python\", \"backend/manage.py\", \"runserver\", \"0.0.0.0:8080\"]"

	case "Node.js (Express)":
		dockerfileContent += `# --- Stage 2: Nodejs Express Production Runtime ---
FROM node:20-alpine
WORKDIR /app
COPY backend/package*.json ./backend/
RUN cd backend && npm install
COPY backend/ ./backend/
`
		if meta.Frontend != "None (Pure Backend API)" {
			dockerfileContent += "COPY --from=frontend-builder /app/frontend/dist ./backend/public\n"
		}
		dockerfileContent += "EXPOSE 8080\nCMD [\"node\", \"backend/index.js\"]"

	case "Rust (Actix-web)":
		dockerfileContent += `# --- Stage 2: Rust Actix Compiled Binary Stage ---
FROM rust:1.75 as backend-builder
WORKDIR /app
COPY backend/ ./backend/
WORKDIR /app/backend
RUN cargo build --release

# --- Stage 3: Minimal Linux Execution Core ---
FROM debian:bookworm-slim
WORKDIR /root/
COPY --from=backend-builder /app/backend/target/release/backend ./main
`
		if meta.Frontend != "None (Pure Backend API)" {
			dockerfileContent += "COPY --from=frontend-builder /app/frontend/dist ./public\n"
		}
		dockerfileContent += "EXPOSE 8080\nCMD [\"./main\"]"

	default: // Go (Golang)
		dockerfileContent += `# --- Stage 2: Go Binary Builder Stage ---
FROM golang:1.22-alpine AS backend-builder
WORKDIR /app
COPY backend/ ./backend/
WORKDIR /app/backend
RUN go build -o main .

# --- Stage 3: Lightweight Alpine Deployment Core ---
FROM alpine:latest
WORKDIR /root/
COPY --from=backend-builder /app/backend/main .
`
		if meta.Frontend != "None (Pure Backend API)" {
			dockerfileContent += "COPY --from=frontend-builder /app/frontend/dist ./public\n"
		}
		dockerfileContent += "EXPOSE 8080\nCMD [\"./main\"]"
	}

	// Write out the assembled layers into the root workspace folder
	targetPath := filepath.Join(basePath, "Dockerfile")
	return os.WriteFile(targetPath, []byte(dockerfileContent), 0644)
}

func writeTemplate(targetFilePath string, blueprint string, meta ProjectMetadata) error {
	tmpl, err := template.New(filepath.Base(targetFilePath)).Parse(blueprint)
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}
	var processedCode bytes.Buffer
	if err := tmpl.Execute(&processedCode, meta); err != nil {
		return err
	}
	_ = os.MkdirAll(filepath.Dir(targetFilePath), 0755)
	return os.WriteFile(targetFilePath, processedCode.Bytes(), 0644)
}