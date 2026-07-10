package k8s

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// K8sManifestVars represents the options passed from the CLI/orchestrator
type K8sManifestVars struct {
	ServiceName   string
	ImageName     string
	ContainerPort int
	ServicePort   int 
	ServiceType   string
	Replicas      int
	CpuRequest    string
	MemoryRequest string
}

// GenerateK8sManifestes writes out deployment and service configuration blocks dynamically
func GenerateK8sManifestes(targetDir string, vars K8sManifestVars) error {
	k8sOutDir := filepath.Join(targetDir, "k8s")
	if err := os.MkdirAll(k8sOutDir, 0755); err != nil {
		return fmt.Errorf("failed to create k8s directory: %w", err)
	}

	// Dynamic safe fallbacks for resource parameters if empty
	cpuReq := vars.CpuRequest
	if cpuReq == "" { cpuReq = "250m" }
	
	memReq := vars.MemoryRequest
	if memReq == "" { memReq = "256Mi" }

	// Calculate production limits cleanly based on units
	cpuLimit := "1"
	if strings.Contains(cpuReq, "500m") || cpuReq == "1" {
		cpuLimit = "2"
	} else if cpuReq == "250m" {
		cpuLimit = "500m" // Safely double the 250m request out to 500m limit
	}

	memLimit := "512Mi"
	if strings.Contains(memReq, "512Mi") {
		memLimit = "1Gi"
	} else if strings.Contains(memReq, "1Gi") || strings.Contains(memReq, "1024Mi") {
		memLimit = "2Gi"
	} else if memReq == "256Mi" {
		memLimit = "512Mi"
	}

	// 1. Dynamic Deployment Blueprint Configuration
	// NOTE: Removed structural double quotes around resource inputs to avoid string escaping formatting bugs
	deploymentContent := fmt.Sprintf(`apiVersion: apps/v1
kind: Deployment
metadata:
  name: %s-deployment
  labels:
    app: %s
spec:
  replicas: %d
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  selector:
    matchLabels:
      app: %s
  template:
    metadata:
      labels:
        app: %s
    spec:
      containers:
      - name: %s
        image: %s
        imagePullPolicy: IfNotPresent
        ports:
        - containerPort: %d
        readinessProbe:
          httpGet:
            path: /api/health
            port: %d
          initialDelaySeconds: 5
          periodSeconds: 10
        livenessProbe:
          httpGet:
            path: /api/health
            port: %d
          initialDelaySeconds: 15
          periodSeconds: 20
        resources:
          requests:
            cpu: %s
            memory: %s
          limits:
            cpu: %s
            memory: %s
`, vars.ServiceName, vars.ServiceName, vars.Replicas, vars.ServiceName, vars.ServiceName,
		vars.ServiceName, vars.ImageName, vars.ContainerPort, vars.ContainerPort, vars.ContainerPort,
		cpuReq, memReq, cpuLimit, memLimit)

	deployFile := filepath.Join(k8sOutDir, fmt.Sprintf("%s-deployment.yaml", vars.ServiceName))
	if err := os.WriteFile(deployFile, []byte(deploymentContent), 0644); err != nil {
		return fmt.Errorf("failed to write deployment file: %w", err)
	}

	// 2. Dynamic Service Blueprint Configuration
	serviceContent := fmt.Sprintf(`apiVersion: v1
kind: Service
metadata:
  name: %s-service
  labels:
    app: %s
spec:
  type: %s
  ports:
  - port: %d
    targetPort: %d
    protocol: TCP
  selector:
    app: %s
`, vars.ServiceName, vars.ServiceName, vars.ServiceType, vars.ServicePort, vars.ContainerPort, vars.ServiceName)

	serviceFile := filepath.Join(k8sOutDir, fmt.Sprintf("%s-service.yaml", vars.ServiceName))
	if err := os.WriteFile(serviceFile, []byte(serviceContent), 0644); err != nil {
		return fmt.Errorf("failed to write service file: %w", err)
	}

	fmt.Printf("Kubernetes manifests generated successfully for %s in /k8s\n", vars.ServiceName)
	return nil
}