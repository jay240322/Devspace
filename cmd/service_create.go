package cmd

import (
	"devspace/internal/config"
	"devspace/internal/templates"
	"fmt"
	"os"

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
		fmt.Println("Service Name: %s\n", serviceName)
		fmt.Println("Backend : %s\n", backend)
		fmt.Println("Frontend : %s\n", frontend)

		meta := templates.ProjectMetadata{
			ServiceName : serviceName,
			Backend : backend,
			Frontend : frontend,
			GitHubUser :"patel-jay",
		}

		err = templates.GenerateBoilerplate(meta)
		if err != nil {
			fmt.Printf("Generation Failed: %v\n", err)
			return
		}
		fmt.Println("You'r customizable, ready-made microservice architeture is ready to launch!")
	},
}