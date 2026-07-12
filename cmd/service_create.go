package cmd

import (
	"devspace/internal/config"
	"devspace/internal/templates"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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

		// 6. Advanced Kubernetes Interactive Prompts
		
		// A. Replicas Prompt
		replicaPrompt := promptui.Prompt{
			Label:   "Enter number of Kubernetes Replicas",
			Default: "1",
			Validate: func(input string) error {
				val, err := strconv.Atoi(input)
				if err != nil || val <= 0 {
					return fmt.Errorf("❌ Replicas must be a positive integer greater than 0")
				}
				return nil
			},
		}
		replicaStr, err := replicaPrompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed: %v\n", err)
			return
		}
		replicas, _ := strconv.Atoi(replicaStr)

		// B. CPU Allocation Prompt
		cpuPrompt := promptui.Prompt{
			Label:   "Enter CPU Request limit (e.g., 100m, 250m, 500m)",
			Default: "100m",
			Validate: func(input string) error {
				if strings.TrimSpace(input) == "" {
					return fmt.Errorf("❌ CPU allocation cannot be empty")
				}
				return nil
			},
		}
		cpuRequest, err := cpuPrompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed: %v\n", err)
			return
		}
		cpuRequest = strings.TrimSpace(cpuRequest)

		// C. Memory Allocation Prompt
		memPrompt := promptui.Prompt{
			Label:   "Enter Memory Request limit (e.g., 128Mi, 256Mi, 512Mi)",
			Default: "128Mi",
			Validate: func(input string) error {
				if strings.TrimSpace(input) == "" {
					return fmt.Errorf("❌ Memory allocation cannot be empty")
				}
				return nil
			},
		}
		memoryRequest, err := memPrompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed: %v\n", err)
			return
		}
		memoryRequest = strings.TrimSpace(memoryRequest)

		// D. Backend Service Edge Protocol Selection Prompt
		serviceTypeSelect := promptui.Select{
			Label: "Select Kubernetes Service Routing Type",
			Items: []string{"ClusterIP", "NodePort", "LoadBalancer"},
		}
		_, serviceType, err := serviceTypeSelect.Run()
		if err != nil {
			fmt.Printf("Selection failed: %v\n", err)
			return
		}

		// Calculate the exact target workspace subdirectory mapping safely
		var finalK8sDir string
		if targetDir == "." {
			wd, _ := os.Getwd()
			finalK8sDir = filepath.Join(wd, serviceName)
		} else {
			finalK8sDir = filepath.Join(filepath.Clean(targetDir), serviceName)
		}

		// Ship data parameters downstream to structural engine using a pointer
		meta := templates.ProjectMetadata{
			TargetDir:      targetDir,
			ServiceName:    serviceName,
			Backend:        backend,
			Frontend:       frontend,
			GitHubUser:     "patel-jay",
			K8sReplicas:    replicas,
			K8sServiceType: serviceType,
			K8sCpuRequest:  cpuRequest,
			K8sMemRequest:  strings.ReplaceAll(strings.ToUpper(memoryRequest), "MI", "Mi"),
		}

		// 🟢 Call workspace generation seamlessly. Your generator.go pointer switch block 
		// now correctly defines and applies ports and health paths out of the box!
		err = templates.GenerateServiceWorkspace(finalK8sDir, &meta)
		if err != nil {
			fmt.Printf("Generation Failed: %v\n", err)
			return
		}

		fmt.Println("\nYour customizable, ready-made microservice architecture is ready to launch!")
	},
}