package templates

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"devspace/internal/templates/k8s"
)

type ProjectMetadata struct {
	TargetDir           string
	ServiceName         string
	Backend             string
	Frontend            string
	GitHubUser          string
	K8sReplicas         int
	K8sServiceType      string
	K8sCpuRequest       string 
	K8sMemRequest       string 
	BackendPort         int
	BackendHealthPath   string
	FrontendPort        int
	FrontendServicePort int
	FrontendHealthPath  string
}

func GenerateBoilerplate(meta ProjectMetadata) error {
	// If TargetDir is empty, default to current directory safely
	if meta.TargetDir == "" {
		meta.TargetDir = "."
	}

	// This cleanly joins your custom target path and service folder name
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

	// Fallback/defaults in GenerateBoilerplate in case meta wasn't initialized via CLI
	if meta.BackendPort == 0 {
		switch meta.Backend {
		case "Node.js (Express)":
			meta.BackendPort = 8080
		case "Rust (Actix-web)":
			meta.BackendPort = 8080
		case "Python (Django)":
			meta.BackendPort = 8080
		default: // Go (Golang)
			meta.BackendPort = 8080
		}
	}
	if meta.BackendHealthPath == "" {
		meta.BackendHealthPath = "/api/health"
	}
	if meta.K8sServiceType == "" {
		meta.K8sServiceType = "ClusterIP"
	}
	if meta.K8sReplicas <= 0 {
		meta.K8sReplicas = 1
	}

	// ✅ DYNAMIC KUBERNETES MANIFEST MAPPING
	backendK8sVars := k8s.K8sManifestVars{
		ServiceName:   fmt.Sprintf("%s-backend", meta.ServiceName),
		ImageName:     fmt.Sprintf("%s/%s-backend", meta.GitHubUser, meta.ServiceName), 
		ContainerPort: meta.BackendPort,
		ServicePort:   meta.BackendPort,
		ServiceType:   meta.K8sServiceType, 
		Replicas:      meta.K8sReplicas,    
		CpuRequest:    meta.K8sCpuRequest,   
		MemoryRequest: meta.K8sMemRequest,
		HealthPath:    meta.BackendHealthPath,
	}
	_ = k8s.GenerateK8sManifestes(basePath, backendK8sVars)

	// 3. Let the frontend CLI handle its own folder creation cleanly
	if meta.Frontend != "None (Pure Backend API)" {
		frontendPath := filepath.Join(basePath, "frontend")
		GenerateFrontendFramework(frontendPath, "", meta)

		if meta.FrontendPort == 0 {
			if strings.Contains(meta.Frontend, "Next.js") {
				meta.FrontendPort = 3000
			} else {
				meta.FrontendPort = 80
			}
		}
		if meta.FrontendServicePort == 0 {
			if strings.Contains(meta.Frontend, "Next.js") {
				meta.FrontendServicePort = 3000
			} else {
				meta.FrontendServicePort = 80
			}
		}
		if meta.FrontendHealthPath == "" {
			meta.FrontendHealthPath = "/"
		}

		frontendK8sVars := k8s.K8sManifestVars{
			ServiceName:   fmt.Sprintf("%s-frontend", meta.ServiceName),
			ImageName:     fmt.Sprintf("%s/%s-frontend", meta.GitHubUser, meta.ServiceName),
			ContainerPort: meta.FrontendPort,
			ServicePort:   meta.FrontendServicePort,
			ServiceType:   "LoadBalancer", 
			Replicas:      meta.K8sReplicas,
			CpuRequest:    "100m",
			MemoryRequest: "128Mi",
			HealthPath:    meta.FrontendHealthPath,
		}
		_ = k8s.GenerateK8sManifestes(basePath, frontendK8sVars)
	}

	// 4. Generate custom full-stack Dockerfile on the fly at the ROOT level
	err := generateDynamicDockerfile(basePath, meta)
	if err != nil {
		fmt.Printf("⚠️  Dockerfile compilation skipped: %v\n", err)
	}

	fmt.Printf("\n✨ Successfully populated dynamic structure inside: %s\n", basePath)
	return nil
}

// Dynamic Dockerfile Builder Factory
func generateDynamicDockerfile(basePath string, meta ProjectMetadata) error {
	// 1. Generate FRONTEND Dockerfile (If frontend exists)
	if meta.Frontend != "None (Pure Backend API)" {
		frontendContent := fmt.Sprintf(`# --- Frontend Production Nginx Stack ---
FROM node:20-alpine AS builder
WORKDIR /app
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ ./
RUN npm run build --if-present

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
`)
		
		frontPath := filepath.Join(basePath, "frontend.Dockerfile")
		if err := os.WriteFile(frontPath, []byte(frontendContent), 0644); err != nil {
			return err
		}
	}

	// 2. Generate BACKEND Dockerfile
	var backendContent string
	switch meta.Backend {
	case "Python (Django)":
		backendContent = `# --- Python Django Environment ---
FROM python:3.11-slim
WORKDIR /app
RUN pip install django django-cors-headers
COPY backend/ ./backend/
EXPOSE 8080
CMD ["python", "backend/manage.py", "runserver", "0.0.0.0:8080"]
`
	case "Node.js (Express)":
		backendContent = `# --- Node.js Express Runtime ---
FROM node:20-alpine
WORKDIR /app
COPY backend/package*.json ./backend/
RUN cd backend && npm install
COPY backend/ ./backend/
EXPOSE 8080
CMD ["node", "backend/index.js"]
`
	case "Rust (Actix-web)":
		backendContent = `# --- Rust Actix Pipeline ---
FROM rust:1.75 as builder
WORKDIR /app
COPY backend/ ./backend/
WORKDIR /app/backend
RUN cargo build --release

FROM debian:bookworm-slim
WORKDIR /root/
COPY --from=builder /app/backend/target/release/backend ./main
EXPOSE 8080
CMD ["./main"]
`
	default: // Go (Golang)
		backendContent = `# --- Go Deployment Core ---
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY backend/ ./backend/
WORKDIR /app/backend
RUN go build -o main .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/backend/main .
EXPOSE 8080
CMD ["./main"]
`
	}

	backPath := filepath.Join(basePath, "backend.Dockerfile")
	return os.WriteFile(backPath, []byte(backendContent), 0644)
}

// generate serviceworkspace into custom folder structures cleanly
func GenerateServiceWorkspace(customPath string, meta ProjectMetadata) error {
	if err := os.MkdirAll(customPath, 0755); err != nil {
		return fmt.Errorf("failed to provision target directory root: %w", err)
	}

	fmt.Printf("📂 Initializing target workspace root at: %s\n", customPath)
	meta.TargetDir = filepath.Dir(customPath)
	meta.ServiceName = filepath.Base(customPath)
	return GenerateBoilerplate(meta)
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