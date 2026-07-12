package cmd

import (
	"devspace/internal/config"
	"devspace/internal/templates"
	"devspace/internal/templates/k8s" // Integrates the four-manifest generation engine
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func init(){
	RootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create and provision a new microservice instantly",
	Run: func(cmd *cobra.Command, args []string){
		fmt.Println("Initializing Devspace Service Provisioner..")
		
		// 1. Security validation
		_, err := config.LoadConfig()
		if err != nil {
			fmt.Printf("Security validation failed: %v\n", err)
			fmt.Println("Run 'set GITHUB_TOKEN=your_token' in terminal before running this tool")
			os.Exit(1)
		}

		// 2. Destination target directory path prompt
		pathPrompt := promptui.Prompt{
			Label:   "Enter full path destination where microservice should be built (use . for current directory)",
			Default: ".",
			Validate: func(input string) error {
				if input == "." {
					return nil
				}
				cleanPath := filepath.Clean(input)
				info, err := os.Stat(cleanPath)
				if os.IsNotExist(err) {
					return fmt.Errorf("❌ Path does not exist! Please enter a valid path")
				}
				if err != nil {
					return fmt.Errorf("❌ Error reading path: %v", err)
				}
				if !info.IsDir() {
					return fmt.Errorf("❌ Target path is a file, not a folder")
				}
				return nil
			},
		}
		targetDir, err := pathPrompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed: %v\n", err)
			return
		}

		// 3. Service name prompt
		namePrompt := promptui.Prompt{
			Label: "Enter microservice name",
			Validate: func(input string) error {
				if strings.TrimSpace(input) == "" {
					return fmt.Errorf("❌ Service name cannot be blank")
				}
				return nil
			},
		}
		serviceName, err := namePrompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed: %v\n", err)
			return
		}
		// Sanitize to lowercase for Kubernetes compatibility
		serviceName = strings.ToLower(strings.ReplaceAll(serviceName, " ", "-"))

		// 4. Backend selection prompt
		langSelect := promptui.Select{
			Label: "Select programming language for backend",
			Items: []string{
				"Go (Golang)",
				"Python (Django)",
				"Node.js (Express)",
				"Rust (Actix-web)",
			},
		}
		_, backend, err := langSelect.Run()
		if err != nil {
			fmt.Printf("Prompt failed: %v\n", err)
			return
		}

		// 5. Frontend selection prompt
		frontendSelect := promptui.Select{
			Label: "Select frontend Framework",
			Items: []string{
				"React (Vite)",
				"Next.js (React)",
				"Vue.Js",
				"Svelte",
				"None (Pure Backend API)",
			},
		}
		_, frontend, err := frontendSelect.Run()
		if err != nil {
			fmt.Printf("Selection failed: %v\n", err)
			return
		}

		// Determine backend port and health path based on backend tier
		var backendPort int
		switch backend {
		case "Node.js (Express)":
			backendPort = 8080
		case "Rust (Actix-web)":
			backendPort = 8080
		case "Python (Django)":
			backendPort = 8080
		default: // Go (Golang)
			backendPort = 8080
		}

		backendHealthPath := "/api/health"

		// Determine frontend port, service port, and health path based on frontend tier
		var frontendPort int
		var frontendServicePort int
		var frontendHealthPath string

		if frontend != "None (Pure Backend API)" {
			if strings.Contains(frontend, "Next.js") {
				frontendPort = 3000
				frontendServicePort = 3000
				frontendHealthPath = "/"
			} else {
				frontendPort = 80
				frontendServicePort = 80
				frontendHealthPath = "/"
			}
		}

		meta := templates.ProjectMetadata{
			TargetDir:           targetDir,
			ServiceName:         serviceName,
			Backend:             backend,
			Frontend:            frontend,
			GitHubUser:          "patel-jay",
			BackendPort:         backendPort,
			BackendHealthPath:   backendHealthPath,
			FrontendPort:        frontendPort,
			FrontendServicePort: frontendServicePort,
			FrontendHealthPath:  frontendHealthPath,
		}

		err = templates.GenerateBoilerplate(meta)
		if err != nil {
			fmt.Printf("Generation Failed: %v\n", err)
			return
		}

		// Calculate the exact target workspace subdirectory mapping
		var finalK8sDir string
		if targetDir == "." {
			wd, _ := os.Getwd()
			finalK8sDir = filepath.Join(wd, serviceName)
		} else {
			finalK8sDir = filepath.Join(filepath.Clean(targetDir), serviceName)
		}

		// Sanitize inputs to ensure zero formatting issues in Kubernetes templates
		cpuRequest = strings.TrimSpace(cpuRequest)
		memoryRequest = strings.TrimSpace(memoryRequest)

		// Dynamically assign correct internal ports based on backend type
		var backendPort int
		switch backend {
		case "Rust (Actix-web)":
			backendPort = 8080
		case "Node.js (Express)":
			backendPort = 8080
		case "Python (Django)":
			backendPort = 8000
		default:
			backendPort = 8080
		}

		// 8. Generate Separate Kubernetes Manifest files for BOTH Backend and Frontend
		
		// A. Generate Backend Deployment & Service Manifests
		backendVars := k8s.K8sManifestVars{
			ServiceName:   fmt.Sprintf("%s-backend", serviceName),
			ImageName:     fmt.Sprintf("patel-jay/%s-backend", serviceName),
			ContainerPort: backendPort,
			ServicePort:   backendPort,
			ServiceType:   serviceType,
			Replicas:      replicas,
			CpuRequest:    cpuRequest,
			MemoryRequest: memoryRequest,
			HealthPath:    "/api/health",
		}

		err = k8s.GenerateK8sManifestes(finalK8sDir, backendVars)
		if err != nil {
			fmt.Printf("❌ Backend Kubernetes manifest creation failed: %v\n", err)
			return
		}

		// B. Generate Frontend Deployment & Service Manifests (If frontend exists)
		if frontend != "None (Pure Backend API)" {
			var frontendPort int
			// Using strings.Contains for resilient string matching against promptui outputs
			if strings.Contains(frontend, "Next.js") {
				frontendPort = 3000
			} else {
				frontendPort = 80 // Direct native Nginx fallback route for production bundles (React/Vue/Svelte)
			}

			frontendVars := k8s.K8sManifestVars{
				ServiceName:   fmt.Sprintf("%s-frontend", serviceName),
				ImageName:     fmt.Sprintf("patel-jay/%s-frontend", serviceName),
				ContainerPort: frontendPort,
				ServicePort:   80,
				ServiceType:   "LoadBalancer", // Accessible standard edge routing 
				Replicas:      replicas,
				CpuRequest:    "100m",
				MemoryRequest: "128Mi",
				HealthPath:    "/",
			}

			err = k8s.GenerateK8sManifestes(finalK8sDir, frontendVars)
			if err != nil {
				fmt.Printf("❌ Frontend Kubernetes manifest creation failed: %v\n", err)
				return
			}
		}

		fmt.Println("\nYour customizable, ready-made microservice architecture is ready to launch!")
	},
}