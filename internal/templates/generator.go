package templates

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
	"devspace/internal/templates/k8s"
)

type ProjectMetadata struct {
	TargetDir      string
	ServiceName    string
	Backend        string
	Frontend       string
	GitHubUser     string
	K8sReplicas    int
	K8sServiceType string
	K8sCpuRequest  string 
	K8sMemRequest  string 
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

	// ✅ CALCULATE BACKEND PORT (Only once!)
	var backendPort int
	switch meta.Backend {
	case "Rust (Actix-web)":
		backendPort = 8080
	case "Node.js (Express)":
		backendPort = 3000
	case "Python (Django)":
		backendPort = 8000
	default: // Go (Golang)
		backendPort = 8080
	}
	
	// ✅ DYNAMIC KUBERNETES MANIFEST MAPPING (Using your GitHub User update!)
	backendK8sVars := k8s.K8sManifestVars{
		ServiceName:   fmt.Sprintf("%s-backend", meta.ServiceName),
		ImageName:     fmt.Sprintf("%s/%s-backend", meta.GitHubUser, meta.ServiceName), 
		ContainerPort: backendPort,
		ServicePort:   backendPort,
		ServiceType:   meta.K8sServiceType, 
		Replicas:      meta.K8sReplicas,    
		CpuRequest:    meta.K8sCpuRequest,   
		MemoryRequest: meta.K8sMemRequest,   
	}
	_ = k8s.GenerateK8sManifestes(basePath, backendK8sVars)

	// 3. Let the frontend CLI handle its own folder creation cleanly
	if meta.Frontend != "None (Pure Backend API)" {
		frontendPath := filepath.Join(basePath, "frontend")
		GenerateFrontendFramework(frontendPath, "", meta)

		// BONUS DYNAMIC FRONTEND K8S PIPELINE CONFIGURATION
		var frontendPort int
		switch meta.Frontend {
		case "Next.js", "Next":
			frontendPort = 3000
		case "React":
			frontendPort = 80 // Nginx distribution production engine standard
		default:
			frontendPort = 3000
		}

		frontendK8sVars := k8s.K8sManifestVars{
			ServiceName:   fmt.Sprintf("%s-frontend", meta.ServiceName),
			ImageName:     fmt.Sprintf("%s/%s-frontend", meta.GitHubUser, meta.ServiceName),
			ContainerPort: frontendPort,
			ServicePort:   80,
			ServiceType:   "LoadBalancer", 
			Replicas:      meta.K8sReplicas,
			CpuRequest:    "100m",
			MemoryRequest: "128Mi",
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
	var dockerfileContent string

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

	// Add dynamic footer signature using text/template parsing
	footerBlueprint := "\n\n# Provisioned securely via DevSpace for github.com/{{.GitHubUser}}/{{.ServiceName}}"
	tmpl, err := template.New("footer").Parse(footerBlueprint)
	if err == nil {
		var processedFooter bytes.Buffer
		if err := tmpl.Execute(&processedFooter, meta); err == nil {
			dockerfileContent += processedFooter.String()
		}
	}

	targetPath := filepath.Join(basePath, "Dockerfile")
	return os.WriteFile(targetPath, []byte(dockerfileContent), 0644)
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