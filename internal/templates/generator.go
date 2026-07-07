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
	backendPath := filepath.Join(basePath, "backend")
	frontendPath := filepath.Join(basePath, "frontend")

	// 1. Physically construct directories
	if err := os.MkdirAll(backendPath, 0755); err != nil {
		return err
	}

	// 2. Dynamically route to the correct isolated Backend File
	switch meta.Backend {
	case "Python (Django)":
		GeneratePythonBackend(backendPath, meta)
	case "Node.js (Express)":
		GenerateNodeBackend(backendPath, meta)
	case "Rust (Actix-web)":
		GenerateRustBackend(backendPath, meta)
	default:
		GenerateGoBackend(backendPath, meta)
	}

	// 3. Dynamically route to the correct isolated Frontend File
	if meta.Frontend != "None (Pure Backend API)" {
		frontendSrcPath := filepath.Join(frontendPath, "src")
		_ = os.MkdirAll(frontendSrcPath, 0755)
		GenerateFrontendFramework(frontendPath, frontendSrcPath, meta)
	}

	fmt.Printf("\n✨ Successfully populated dynamic structure inside: %s\n", basePath)
	return nil
}

// Shared helper utility to process dynamic injection maps
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