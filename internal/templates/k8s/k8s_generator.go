package k8s

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

type K8sManifestVars struct {
	ServiceName   string
	ImageName     string
	ContainerPort int
	ServicePort   int
	ServiceType   string
	Replicas      int
	CpuRequest    string
	MemoryRequest string
	HealthPath    string
}

// Global Deployment Blueprint Template String
const deploymentBlueprint = `apiVersion: apps/v1
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
        image: {{.ImageName}}:latest
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
            cpu: "{{.CpuRequest}}"
            memory: "{{.MemoryRequest}}"
          limits:
            cpu: "500m"
            memory: "512Mi"
`

// Global Service Blueprint Template String
const serviceBlueprint = `apiVersion: v1
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

func GenerateK8sManifestes(basePath string, vars K8sManifestVars) error {
	k8sDir := filepath.Join(basePath, "k8s")
	if err := os.MkdirAll(k8sDir, 0755); err != nil {
		return fmt.Errorf("failed to create k8s directory: %w", err)
	}

	// 1. Process and write the Deployment manifest file
	tmplDep, err := template.New("deployment").Parse(deploymentBlueprint)
	if err != nil {
		return fmt.Errorf("failed to parse deployment blueprint template: %w", err)
	}
	
	var depBuffer bytes.Buffer
	if err := tmplDep.Execute(&depBuffer, vars); err != nil {
		return fmt.Errorf("failed to execute deployment template mapping: %w", err)
	}

	depPath := filepath.Join(k8sDir, fmt.Sprintf("%s-deployment.yaml", vars.ServiceName))
	if err := os.WriteFile(depPath, depBuffer.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write deployment manifest: %w", err)
	}

	// 2. Process and write the Service manifest file
	tmplSvc, err := template.New("service").Parse(serviceBlueprint)
	if err != nil {
		return fmt.Errorf("failed to parse service blueprint template: %w", err)
	}

	var svcBuffer bytes.Buffer
	if err := tmplSvc.Execute(&svcBuffer, vars); err != nil {
		return fmt.Errorf("failed to execute service template mapping: %w", err)
	}

	svcPath := filepath.Join(k8sDir, fmt.Sprintf("%s-service.yaml", vars.ServiceName))
	if err := os.WriteFile(svcPath, svcBuffer.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write service manifest: %w", err)
	}

	return nil
}