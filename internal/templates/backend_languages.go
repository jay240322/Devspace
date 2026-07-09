package templates

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func GenerateGoBackend(basePath string, meta ProjectMetadata) {
	backendPath := filepath.Join(basePath, "backend")
	_ = os.MkdirAll(backendPath, 0755)

	fmt.Println("🚀 Scaffolding Native Go Workspace Module...")
	cmd := exec.Command("go", "mod", "init", fmt.Sprintf("github.com/%s/%s/backend", meta.GitHubUser, meta.ServiceName))
	cmd.Dir = backendPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()

	mainCode := `package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, "{\"status\": \"healthy\", \"engine\": \"Go Alpine Core\"}")
	})

	fmt.Println("🚀 Go Microservice running smoothly on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}`
	_ = os.WriteFile(filepath.Join(backendPath, "main.go"), []byte(mainCode), 0644)
}

func GeneratePythonBackend(basePath string, meta ProjectMetadata) {
	backendPath := filepath.Join(basePath, "backend")
	_ = os.MkdirAll(backendPath, 0755)

	fmt.Println("🐍 Initializing Django Admin Project Core Layout...")
	
	pythonCmd := "python"
	if runtime.GOOS == "windows" {
		pythonCmd = "python"
	}

	cmd := exec.Command(pythonCmd, "-m", "django", "startproject", "backend", ".")
	cmd.Dir = backendPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		// Fallback directly to regular django-admin global CLI call if module path behaves oddly
		fallbackCmd := exec.Command("django-admin", "startproject", "backend", ".")
		fallbackCmd.Dir = backendPath
		_ = fallbackCmd.Run()
	}
}

func GenerateNodeBackend(basePath string, meta ProjectMetadata) {
	backendPath := filepath.Join(basePath, "backend")
	_ = os.MkdirAll(backendPath, 0755)

	fmt.Println("🟢 Initializing Express Node.js Workspace Context...")
	
	npmCmd := "npm"
	if runtime.GOOS == "windows" {
		npmCmd = "npm.cmd"
	}

	initCmd := exec.Command(npmCmd, "init", "-y")
	initCmd.Dir = backendPath
	_ = initCmd.Run()

	// FIXED: Concatenating string representation to prevent JavaScript backticks from breaking Go literals
	indexCode := `const express = require('express');
const app = express();
const PORT = 8080;

app.get('/api/health', (req, res) => {
    res.json({ status: 'healthy', engine: 'Node.js Express Runtime' });
});

app.listen(PORT, () => {
    console.log(` + "`" + `🚀 Express Backend server tracking actively on port ${PORT}` + "`" + `);
});`
	_ = os.WriteFile(filepath.Join(backendPath, "index.js"), []byte(indexCode), 0644)
}

func GenerateRustBackend(basePath string, meta ProjectMetadata) {
	targetBackendDir := filepath.Join(basePath, "backend")
	_ = os.MkdirAll(basePath, 0755)

	fmt.Println("🦀 Provisioning Cargo Binary Executable Instance...")
	
	cmd := exec.Command("cargo", "new", "backend", "--bin")
	cmd.Dir = basePath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()

	cargoTomlContent := `[package]
name = "backend"
version = "0.1.0"
edition = "2021"

[dependencies]
actix-web = "4.4"`
	_ = os.WriteFile(filepath.Join(targetBackendDir, "Cargo.toml"), []byte(cargoTomlContent), 0644)

	mainRustCode := `use actix_web::{get, HttpResponse, App, HttpServer, Responder};

#[get("/api/health")]
async fn health_check() -> impl Responder {
    HttpResponse::Ok().json(serde_json::json!({"status": "healthy", "engine": "Rust Actix-Web Execution Core"}))
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    println!("🚀 Rust Actix Web Server executing natively on port 8080...");
    HttpServer::new(|| {
        App::new().service(health_check)
    })
    .bind(("0.0.0.0", 8080))?
    .run()
    .await
}`
	_ = os.WriteFile(filepath.Join(targetBackendDir, "src", "main.rs"), []byte(mainRustCode), 0644)
}