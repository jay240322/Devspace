package cmd

import (
	"devspace/internal/templates"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create and provision a new microservice instantly",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Initializing Devspace Service Provisioner..")

		// 2. Project location
		pathPrompt := promptui.Prompt{
			Label:       "Project location (Enter = current folder)",
			HideEntered: true,
			Validate: func(input string) error {
				_, err := resolveProjectDirectory(input)
				return err
			},
		}

		enteredPath, err := pathPrompt.Run()
		if err != nil {
			fmt.Printf("Prompt cancelled: %v\n", err)
			return
		}

		targetDir, err := resolveProjectDirectory(enteredPath)
		if err != nil {
			fmt.Printf("Invalid project location: %v\n", err)
			return
		}

		fmt.Println("Project directory:", targetDir)

		// 3. Service name
		namePrompt := promptui.Prompt{
			Label:       "Service name (example: trial-1)",
			HideEntered: true,
			Validate:    validateServiceName,
		}

		enteredName, err := namePrompt.Run()
		if err != nil {
			fmt.Printf("Prompt cancelled: %v\n", err)
			return
		}

		serviceName := normalizeServiceName(enteredName)
		fmt.Println("Service name:", serviceName)

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

		// B. CPU Allocation selection
		cpuSelect := promptui.Select{
			Label: "Select CPU Request Limit Profile",
			Items: []string{
				"100m (Lightweight - Recommended for Node/Go)",
				"250m (Medium - Good for Python/Django)",
				"500m (High Performance)",
				"1000m (1 Full Core - Heavy Load)",
			},
		}
		_, selectedCpu, err := cpuSelect.Run()
		if err != nil {
			fmt.Printf("Selection failed: %v\n", err)
			return
		}
		// Extract just the core Kubernetes value (e.g., "100m") from the selection string
		cpuRequest := strings.Split(selectedCpu, " ")[0]

		// C. Memory Allocation selection
		memSelect := promptui.Select{
			Label: "Select Memory Request Limit Profile",
			Items: []string{
				"128Mi (Micro - Light API instances)",
				"256Mi (Standard - Recommended)",
				"512Mi (Enhanced - Heavy caching/dependencies)",
				"1024Mi (1GiB - High volume production)",
			},
		}
		_, selectedMem, err := memSelect.Run()
		if err != nil {
			fmt.Printf("Selection failed: %v\n", err)
			return
		}
		// Extract just the core Kubernetes profile value (e.g., "128Mi") safely
		memoryRequest := strings.Split(selectedMem, " ")[0]

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
