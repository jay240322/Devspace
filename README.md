<h1 align="center"><b>Devspace</b></h1>

<p align="center">
  <em>An interactive full-stack scaffolding engine built in Go.</em>
</p>

---

Devspace empowers developers to instantly spin up complete application structures—handling frontend and backend framework setups, generating optimized Docker configurations, and building cloud-ready Kubernetes manifests with a single command-line interface.

---

## 📦 Fast Installation

Instead of cloning the source code and compiling dependencies manually, you can download the standalone executable immediately:

1. Head over to the **[Latest Releases](https://github.com/jay240322/Devspace/releases/latest)** section on GitHub.
2. Download `devspace.exe`.
3. Drop it into your target workspace directory.

---

## Step-by-Step Guide: Creating Your First Service

Follow this complete walkthrough to initialize your workspace, containerize your app, and deploy it to a cluster.


## 1. Initialize Devspace

Run the following command:

```bash
.\devspace.exe create
```

If this is your first time using Devspace, GitHub authentication is required.

Set your GitHub Personal Access Token:

```bash
set GITHUB_TOKEN=YOUR_GITHUB_TOKEN
```

After setting the token, run the initialization command again:

```bash
.\devspace.exe create
```

---

## 2. Configure Your Project

Devspace will launch an interactive setup wizard.

You'll be prompted to provide the following information:

### Project Directory

Specify where the project should be generated.

```text
C:\Users\YourName\Desktop
```

###  Microservice Name

Provide a name for your service.

```text
my-service
```

###  Backend Framework

Choose one of the available backend frameworks.

```text
Go (Golang)
Python (Django)
Node.js (Express)
Rust (Actix-web)
```

###  Frontend Framework

Choose your preferred frontend framework.

```text
React (Vite)
Next.js
Vue.js
Svelte
None (Backend API Only)
```

###  Kubernetes Replicas

Enter the number of replicas to deploy.

```text
3
```

###  CPU Resource Profile

Select the CPU request profile.

```text
100m  (Lightweight)
250m  (Medium)
500m  (High Performance)
```

###  Memory Resource Profile

Select the memory request profile.

```text
128Mi
256Mi
512Mi
1024Mi (1Gi)
```

### Kubernetes Service Type

Choose how your application will be exposed.

```text
ClusterIP
NodePort
LoadBalancer
```

After confirming the configuration, Devspace automatically generates your complete project.

---

# 📂 Generated Project Structure

Your workspace will contain a production-ready project similar to the following:

```bash
my-service/
├── backend/
├── frontend/
├── k8s/
   ├──myservice-backend-deployment.yaml
   ├──myservice-backend-service.yaml
   ├──myservice-frontend-deployment.yaml
   └──myservice-frontend-deployment.yaml
├── Dockerfile
├── .dockerignorefile
```

> **Note:** The generated files may vary depending on the frameworks you selected.

---

# Develop Your Application

After project generation, implement your business logic by modifying the generated files.

Typical files you'll edit include:

- Backend source code
- Frontend application
- Dockerfile
- Kubernetes manifests
- Environment variables

---

# 🐳 Build Docker Images

Devspace uses a multi-stage Dockerfile, allowing frontend and backend images to be built independently.

Build the frontend image:

```bash
docker build --target frontend -t my-frontend .
```

Build the backend image:

```bash
docker build --target backend -t my-backend .
```

---

# ☸ Deploy to Kubernetes

Deploy the generated Kubernetes manifests.

```bash
kubectl apply -f kubernetes/
```

Verify that the deployment was successful.

Check Pods:

```bash
kubectl get pods
```

Check Services:

```bash
kubectl get svc
```

Check Deployments:

```bash
kubectl get deployments
```

If all resources are running successfully, your application is now deployed on Kubernetes.

---

#  Complete Workflow

```text
Download Devspace
        │
        ▼
Run:
.\devspace.exe create
        │
        ▼
Configure Project
        │
        ▼
Project Generated
        │
        ▼
Develop Your Application
        │
        ▼
Build Docker Images
        │
        ▼
Deploy to Kubernetes
        │
        ▼
Application Running 🚀
```
