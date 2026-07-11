package cmd

import (
	"devspace/internal/config"
	"devspace/internal/templates"
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
	Use: "create",
	Short: "Create and provision a new  microsetvice instantly",
	Run: func(cmd *cobra.Command, args []string){
		fmt.Println("Initializing Devspace Service Provisioner..")
        //    security configuration
		_, err := config.LoadConfig()
		if err != nil {
			fmt.Println("Security validation failed: %v\n", err)
			fmt.Println("Run 'set GITHUB_TOKEN=your_token' in terminal before running this tool")
			os.Exit(1)
		}
		//input for custom folder path destination
		pathPrompt := promptui.Prompt{
			Label: "Enter full path destination where microservice should be built(use . for curent directory)",
		    Default: ".",
			Validate: func(input string)error {
				if input == "." {
					return nil
				}
				cleanPath := filepath.Clean(input)
				
				info, err := os.Stat(cleanPath)
				if os.IsNotExist(err){
					return fmt.Errorf("❌ Path does not exit! please enter valid path")
				}
				if err != nil {
					return fmt.Errorf("❌ Error readng path: %v",err)
				}
				if !info.IsDir(){
					return fmt.Errorf("❌ Target path is a file, not a folder")
				}
				return nil
			},
		}
		targetDir, err := pathPrompt.Run()
		if err != nil {
			fmt.Println("Prompt failed: %v\n", err)
			return
		}
        // input for service name
		namePrompt := promptui.Prompt{
			Label: "Enter microservice name",
		}
		serviceName, err := namePrompt.Run()
		if err != nil {
			fmt.Println("Prompt failed: %v\n", err)
			return
		}
		// input for selecting backend language
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
		// input for selecting frontend framework
		frontendSelect := promptui.Select{
			Label: "Select frontend Framework",
			Items: []string{
				"React (Vite)",
				"Next.js (Ract)",
				"Vue.Js",
				"Svelte",
				"None (Pure Backend API)",
			},
		}
		_, frontend, err := frontendSelect.Run()
		if err !=nil{
			fmt.Printf("Selection failed: %v\n", err)
			return
		}
        // plan Summay Output
		fmt.Println("\n -- Full-stack plan Created Successfully!\n")
		fmt.Println("Path: %s\n", targetDir)
		fmt.Println("Service Name: %s\n", serviceName)
		fmt.Println("Backend : %s\n", backend)
		fmt.Println("Frontend : %s\n", frontend)

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
		fmt.Println("You'r customizable, ready-made microservice architeture is ready to launch!")
	},
}