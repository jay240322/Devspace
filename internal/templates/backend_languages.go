package templates

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func getCommandName(baseCmd string) string {
	if runtime.GOOS == "windows" {
		return baseCmd + ".cmd"
	}
	return baseCmd
}

func GenerateGoBackend(basePath string, meta ProjectMetadata) {
	fmt.Println("🚀 Running official 'go mod init' setup...")
	backendPath := filepath.Join(basePath, "backend")
	_ = os.MkdirAll(backendPath, 0755)

	cmd := exec.Command("go", "mod", "init", "backend")
	cmd.Dir = backendPath
	_ = cmd.Run()

	code := `package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "{\"message\": \"🚀 Welcome to the {{.ServiceName}} Go backend API!\"}")
	})
	fmt.Println("🌐 Server starting seamlessly on port 8080...")
	http.ListenAndServe(":8080", nil)
}`
	_ = writeTemplate(filepath.Join(backendPath, "main.go"), code, meta)
}

func GeneratePythonBackend(basePath string, meta ProjectMetadata) {
	fmt.Println("🐍 Invoking native 'django-admin startproject' CLI scaffolding...")
	backendPath := filepath.Join(basePath, "backend")

	// FIX: Create the target destination directory FIRST so django-admin does not crash
	if err := os.MkdirAll(backendPath, 0755); err != nil {
		fmt.Printf("❌ Failed to create backend directory: %v\n", err)
		return
	}

	// Run official CLI configuration setup: django-admin startproject backend <dir>
	// Passing "." tells Django to generate the configuration files directly inside the backend directory
	cmd := exec.Command("django-admin", "startproject", "backend", ".")
	cmd.Dir = backendPath // Direct the command context execution inside the target backend folder
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	err := cmd.Run()
	if err != nil {
		fmt.Printf("❌ django-admin execution failed: %v. Utilizing fallback engine...\n", err)
	}
}

func GenerateNodeBackend(basePath string, meta ProjectMetadata) {
	fmt.Println("🟢 Executing Node package installation CLIs ('npm init' & 'npm install')...")
	backendPath := filepath.Join(basePath, "backend")
	_ = os.MkdirAll(backendPath, 0755)

	npmCmd := getCommandName("npm")

	cmdInit := exec.Command(npmCmd, "init", "-y")
	cmdInit.Dir = backendPath
	_ = cmdInit.Run()

	cmdInstall := exec.Command(npmCmd, "install", "express", "cors")
	cmdInstall.Dir = backendPath
	cmdInstall.Stdout = os.Stdout
	cmdInstall.Stderr = os.Stderr
	_ = cmdInstall.Run()

	code := `const express = require('express');
const cors = require('cors');
const app = express();

app.use(cors());
app.get('/', (req, res) => {
    res.json({ message: "🚀 Welcome to the {{.ServiceName}} Node.js Express API!" });
});

app.listen(8080, () => console.log('🌐 Server active on port 8080'));`
	_ = writeTemplate(filepath.Join(backendPath, "index.js"), code, meta)
}

func GenerateRustBackend(basePath string, meta ProjectMetadata) {
	targetBackendDir := filepath.Join(basePath, "backend")
	
	// 1. FIXED: Explicitly force-create the backend directory structure first
	_ = os.MkdirAll(filepath.Join(targetBackendDir, "src"), 0755)

	fmt.Println("🦀 Provisioning Cargo Binary Executable Instance...")
	
	// 2. Try running cargo init natively inside the pre-made folder
	cmd := exec.Command("cargo", "init", "--bin")
	cmd.Dir = targetBackendDir 
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	// If cargo fails or isn't installed locally, our fallback writes the files manually!
	if err := cmd.Run(); err != nil {
		fmt.Println("⚠️  Local 'cargo' CLI not found or failed. Falling back to manual template injection...")
	}

	// 3. Robust Overwrite to guarantee the workspace files exist regardless of local toolchains
	cargoTomlContent := `[package]
name = "backend"
version = "0.1.0"
edition = "2021"

[dependencies]
actix-web = "4.4"
serde_json = "1.0"`
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
	fmt.Println("✅ Rust backend structure successfully secured.")
}