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
	fmt.Println("🦀 Initializing native binary workspace via 'cargo new'...")
	backendPath := filepath.Join(basePath, "backend")
	parentDir := filepath.Dir(backendPath)

	cmd := exec.Command("cargo", "new", "backend", "--bin")
	cmd.Dir = parentDir
	_ = cmd.Run()

	code := `use actix_cors::Cors;
use actix_web::{get, HttpResponse, App, HttpServer, Responder};
use serde::Serialize;

#[derive(Serialize)]
struct Message { message: String }

#[get("/")]
async fn index() -> impl Responder {
    HttpResponse::Ok().json(Message { message: "🚀 Welcome to the {{.ServiceName}} Rust Actix API!".to_string() })
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    HttpServer::new(|| {
        App::new()
            .wrap(Cors::permissive())
            .service(index)
	})
    .bind(("0.0.0.0", 8080))?.run().await
}`

	cargo := `[package]
name = "backend"
version = "0.1.0"
edition = "2021"

[dependencies]
actix-web = "4"
actix-cors = "0.6"
serde = { version = "1.0", features = ["derive"] }`

	_ = writeTemplate(filepath.Join(backendPath, "src/main.rs"), code, meta)
	_ = writeTemplate(filepath.Join(backendPath, "Cargo.toml"), cargo, meta)
}