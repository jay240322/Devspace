package templates

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func GenerateFrontendFramework(frontendPath string, _ string, meta ProjectMetadata) {
	fmt.Printf("📦 Initializing official ecosystem CLI setup for %s...\n", meta.Frontend)

	parentDir := filepath.Dir(frontendPath) 
	folderName := "frontend"                

	// Fix Windows Executable Extensions (.cmd)
	npmCmd := "npm"
	npxCmd := "npx"
	if runtime.GOOS == "windows" {
		npmCmd = "npm.cmd"
		npxCmd = "npx.cmd"
	}

	var cmd *exec.Cmd

	switch meta.Frontend {
	case "Next.js (React)":
		cmd = exec.Command(npxCmd, "create-next-app@latest", folderName, "--ts=false", "--src-dir=true", "--eslint=true", "--tailwind=false", "--app=true", "--import-alias", "@/*", "--use-npm")

	case "Vue.Js":
		cmd = exec.Command(npmCmd, "create", "vite@latest", folderName, "--", "--template", "vue")

	case "Svelte":
		cmd = exec.Command(npxCmd, "sv", "create", folderName, "--template", "minimal", "--no-types", "--no-add-ons", "--no-install")
		
	default: // React (Vite)
		cmd = exec.Command(npmCmd, "create", "vite@latest", folderName, "--", "--template", "react")
	}

	cmd.Dir = parentDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("❌ Failed to scaffold frontend via CLI execution: %v\n", err)
		return
	}

	fmt.Printf("✅ Official %s project skeleton generated perfectly inside /frontend!\n", meta.Frontend)
}