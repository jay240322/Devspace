package templates

import "path/filepath"

func GenerateGoBackend(path string, meta ProjectMetadata) {
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

	docker := `FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]`

	_ = writeTemplate(filepath.Join(path, "main.go"), code, meta)
	_ = writeTemplate(filepath.Join(path, "Dockerfile"), docker, meta)
}

func GeneratePythonBackend(path string, meta ProjectMetadata) {
	code := `import os
import sys
from django.conf import settings
from django.core.wsgi import get_wsgi_application
from django.http import JsonResponse
from django.urls import path

if not settings.configured:
    settings.configure(
        DEBUG=True,
        SECRET_KEY="devspace-super-secret-key-cluster",
        ROOT_URLCONF=__name__,
        ALLOWED_HOSTS=["*"],
        MIDDLEWARE=[
            "corsheaders.middleware.CorsMiddleware",
            "django.middleware.common.CommonMiddleware",
        ],
        CORS_ALLOW_ALL_ORIGINS=True,
        INSTALLED_APPS=[
            "corsheaders",
        ],
    )

def home(request):
    return JsonResponse({"message": "🚀 Welcome to the {{.ServiceName}} Python Django API cluster!"})

urlpatterns = [
    path("", home),
]

application = get_wsgi_application()

if __name__ == "__main__":
    from django.core.management import execute_from_command_line
    print("🌐 Python Django Server starting seamlessly on port 8080...")
    execute_from_command_line([sys.argv[0], "runserver", "0.0.0.0:8080"])`

	docker := `FROM python:3.11-slim
WORKDIR /app
RUN pip install django django-cors-headers
COPY . .
EXPOSE 8080
CMD ["python", "app.py"]`

	_ = writeTemplate(filepath.Join(path, "app.py"), code, meta)
	_ = writeTemplate(filepath.Join(path, "Dockerfile"), docker, meta)
}

func GenerateNodeBackend(path string, meta ProjectMetadata) {
	code := `const express = require('express');
const cors = require('cors');
const app = express();

app.use(cors());
app.get('/', (req, res) => {
    res.json({ message: "🚀 Welcome to the {{.ServiceName}} Node.js Express API!" });
});

app.listen(8080, () => console.log('🌐 Server active on port 8080'));`

	docker := `FROM node:20-alpine
WORKDIR /app
RUN npm install express cors
COPY . .
EXPOSE 8080
CMD ["node", "index.js"]`

	_ = writeTemplate(filepath.Join(path, "index.js"), code, meta)
	_ = writeTemplate(filepath.Join(path, "Dockerfile"), docker, meta)
}

func GenerateRustBackend(path string, meta ProjectMetadata) {
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
name = "{{.ServiceName}}-backend"
version = "0.1.0"
edition = "2021"

[dependencies]
actix-web = "4"
actix-cors = "0.6"
serde = { version = "1.0", features = ["derive"] }`

	docker := `FROM rust:1.75 as builder
WORKDIR /app
COPY . .
RUN cargo build --release

FROM debian:bookworm-slim
WORKDIR /root/
COPY --from=builder /app/target/release/{{.ServiceName}}-backend ./main
EXPOSE 8080
CMD ["./main"]`

	_ = writeTemplate(filepath.Join(path, "src/main.rs"), code, meta)
	_ = writeTemplate(filepath.Join(path, "Cargo.toml"), cargo, meta)
	_ = writeTemplate(filepath.Join(path, "Dockerfile"), docker, meta)
}