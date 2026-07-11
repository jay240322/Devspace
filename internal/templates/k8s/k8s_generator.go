package k8s

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
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
	CpuLimit      string
	MemoryLimit   string
	HealthPath    string
}

// GenerateK8sManifestes writes out deployment and service configuration blocks dynamically
func GenerateK8sManifestes(targetDir string, vars K8sManifestVars) error {
	k8sOutDir := filepath.Join(targetDir, "k8s")
	if err := os.MkdirAll(k8sOutDir, 0755); err != nil {
		return fmt.Errorf("failed to create k8s directory: %w", err)
	}

	// Dynamic safe fallbacks for resource parameters if empty
	if vars.CpuRequest == "" {
		vars.CpuRequest = "250m"
	}
	if vars.MemoryRequest == "" {
		vars.MemoryRequest = "256Mi"
	}

	// Calculate production limits cleanly based on units
	if vars.CpuLimit == "" {
		cpuLimit := "1"
		if strings.Contains(vars.CpuRequest, "500m") || vars.CpuRequest == "1" {
			cpuLimit = "2"
		} else if vars.CpuRequest == "250m" {
			cpuLimit = "500m"
		}
		vars.CpuLimit = cpuLimit
	}

	if vars.MemoryLimit == "" {
		memLimit := "512Mi"
		if strings.Contains(vars.MemoryRequest, "512Mi") {
			memLimit = "1Gi"
		} else if strings.Contains(vars.MemoryRequest, "1Gi") || strings.Contains(vars.MemoryRequest, "1024Mi") {
			memLimit = "2Gi"
		} else if vars.MemoryRequest == "256Mi" {
			memLimit = "512Mi"
		}
		vars.MemoryLimit = memLimit
	}

	if vars.Replicas <= 0 {
		vars.Replicas = 1
	}

	if vars.HealthPath == "" {
		vars.HealthPath = "/api/health"
	}

	if vars.ServiceType == "" {
		vars.ServiceType = "ClusterIP"
	}

	// 1. Dynamic Deployment Blueprint Configuration using text/template
	deploymentTemplateStr := `apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.ServiceName}}-deployment
  labels:
    app: {{.ServiceName}}
spec:
  replicas: {{.Replicas}}
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  selector:
    matchLabels:
      app: {{.ServiceName}}
  template:
    metadata:
      labels:
        app: {{.ServiceName}}
    spec:
      containers:
      - name: {{.ServiceName}}
        image: {{.ImageName}}
        imagePullPolicy: IfNotPresent
        ports:
        - containerPort: {{.ContainerPort}}
        readinessProbe:
          httpGet:
            path: {{.HealthPath}}
            port: {{.ContainerPort}}
          initialDelaySeconds: 5
          periodSeconds: 10
        livenessProbe:
          httpGet:
            path: {{.HealthPath}}
            port: {{.ContainerPort}}
          initialDelaySeconds: 15
          periodSeconds: 20
        resources:
          requests:
            cpu: {{.CpuRequest}}
            memory: {{.MemoryRequest}}
          limits:
            cpu: {{.CpuLimit}}
            memory: {{.MemoryLimit}}
`

	tmplDeploy, err := template.New("deployment").Parse(deploymentTemplateStr)
	if err != nil {
		return fmt.Errorf("failed to parse deployment template: %w", err)
	}

	var deployBuf bytes.Buffer
	if err := tmplDeploy.Execute(&deployBuf, vars); err != nil {
		return fmt.Errorf("failed to execute deployment template: %w", err)
	}

	deployFile := filepath.Join(k8sOutDir, fmt.Sprintf("%s-deployment.yaml", vars.ServiceName))
	if err := os.WriteFile(deployFile, deployBuf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write deployment file: %w", err)
	}

	// 2. Dynamic Service Blueprint Configuration using text/template
	serviceTemplateStr := `apiVersion: v1
kind: Service
metadata:
  name: {{.ServiceName}}-service
  labels:
    app: {{.ServiceName}}
spec:
  type: {{.ServiceType}}
  ports:
  - port: {{.ServicePort}}
    targetPort: {{.ContainerPort}}
    protocol: TCP
  selector:
    app: {{.ServiceName}}
`

	tmplService, err := template.New("service").Parse(serviceTemplateStr)
	if err != nil {
		return fmt.Errorf("failed to parse service template: %w", err)
	}

	var serviceBuf bytes.Buffer
	if err := tmplService.Execute(&serviceBuf, vars); err != nil {
		return fmt.Errorf("failed to execute service template: %w", err)
	}

	serviceFile := filepath.Join(k8sOutDir, fmt.Sprintf("%s-service.yaml", vars.ServiceName))
	if err := os.WriteFile(serviceFile, serviceBuf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write service file: %w", err)
	}

	fmt.Printf("Kubernetes manifests generated successfully for %s in /k8s\n", vars.ServiceName)
	return nil
}