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
	if meta.TargetDir == "" {
		meta.TargetDir = "."
	}

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

	if meta.BackendPort == 0 {
		switch meta.Backend {
		case "Node.js (Express)":
			meta.BackendPort = 8080
		case "Rust (Actix-web)":
			meta.BackendPort = 8080
		case "Python (Django)":
			meta.BackendPort = 8000
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

	// ✅ DYNAMIC KUBERNETES MANIFEST MAPPING (BACKEND)
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

		// ✅ DYNAMIC KUBERNETES MANIFEST MAPPING (FRONTEND)
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

	// 4. Generate ONE unified multi-stage Dockerfile at the root directory level
	err := generateDynamicDockerfile(basePath, meta)
	if err != nil {
		fmt.Printf("⚠️  Dockerfile compilation skipped: %v\n", err)
	}

	// 5. Generate a root .dockerignore file to bypass heavy node installation contexts
	ignoreContent := `node_modules
.git
frontend/node_modules
backend/venv
*.exe
.idea
.vscode
`
	ignorePath := filepath.Join(basePath, ".dockerignore")
	_ = os.WriteFile(ignorePath, []byte(ignoreContent), 0644)

	fmt.Printf("\n✨ Successfully populated dynamic structure inside: %s\n", basePath)
	return nil
}

// Unified Multi-Stage Dockerfile Builder Factory
func generateDynamicDockerfile(basePath string, meta ProjectMetadata) error {
	var dockerfileContent string

	// A. Inject FRONTEND Stages if a user layout option is provided
	if meta.Frontend != "None (Pure Backend API)" {
		dockerfileContent += `# --- Stage 1: Frontend Dependency & Build Pipeline ---
FROM node:20-alpine AS frontend-builder
ENV NODE_ENV=production
ENV NODE_OPTIONS="--max-old-space-size=2048"
WORKDIR /app
COPY frontend/package*.json ./
RUN npm ci --only=production --quiet --no-audit --no-fund
COPY frontend/ ./
RUN npm run build --if-present

# --- Stage 2: Frontend Production Web Server Runtime Target ---
FROM nginx:alpine AS frontend-runtime
COPY --from=frontend-builder /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]

`
	}

	// B. Append BACKEND Runtime Stages dynamically at the bottom
	switch meta.Backend {
	case "Python (Django)":
		dockerfileContent += fmt.Sprintf(`# --- Stage 3: Python Django Backend Engine ---
FROM python:3.11-slim AS backend-runtime
WORKDIR /app
RUN pip install django django-cors-headers
COPY backend/ ./backend/
EXPOSE %d
CMD ["python", "backend/manage.py", "runserver", "0.0.0.0:%d"]
`, meta.BackendPort, meta.BackendPort)

	case "Node.js (Express)":
		dockerfileContent += `# --- Stage 3: Node.js Express Backend Engine ---
FROM node:20-alpine AS backend-runtime
WORKDIR /app
COPY backend/package*.json ./backend/
RUN cd backend && npm install
COPY backend/ ./backend/
EXPOSE 8080
CMD ["node", "backend/index.js"]
`
	case "Rust (Actix-web)":
		dockerfileContent += `# --- Stage 3: Rust Compilation Pipe ---
FROM rust:1.75 AS rust-builder
WORKDIR /app
COPY backend/ ./backend/
WORKDIR /app/backend
RUN cargo build --release

# --- Stage 4: Rust Production Binary Target ---
FROM debian:bookworm-slim AS backend-runtime
WORKDIR /root/
COPY --from=rust-builder /app/backend/target/release/backend ./main
EXPOSE 8080
CMD ["./main"]
`
	default: // Go (Golang)
		dockerfileContent += `# --- Stage 3: Go Compiler Core ---
FROM golang:1.22-alpine AS go-builder
WORKDIR /app
COPY backend/ ./backend/
WORKDIR /app/backend
RUN go build -o main .

# --- Stage 4: Go Production Binary Target ---
FROM alpine:latest AS backend-runtime
WORKDIR /root/
COPY --from=go-builder /app/backend/main .
EXPOSE 8080
CMD ["./main"]
`
	}

	dockerfilePath := filepath.Join(basePath, "Dockerfile")
	return os.WriteFile(dockerfilePath, []byte(dockerfileContent), 0644)
}

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