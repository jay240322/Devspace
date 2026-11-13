<div align="center">

# Devspace

**From an empty folder to an application workspace.**

Scaffold backends, frontends, Dockerfiles, and Kubernetes manifests through an interactive CLI.

[![npm](https://img.shields.io/npm/v/devspace-kit?style=for-the-badge&color=FACC15&labelColor=18181B)](https://www.npmjs.com/package/devspace-kit)
[![Go](https://img.shields.io/badge/Built_with-Go-FACC15?style=for-the-badge&labelColor=18181B)](https://go.dev/)
[![License](https://img.shields.io/badge/License-Apache_2.0-FACC15?style=for-the-badge&labelColor=18181B)](LICENSE)

**Windows · macOS · Linux**

[Installation](#installation) · [Quick Start](#quick-start) · [Commands](#commands) · [Contributing](#contributing)

</div>

---

## Overview

Devspace guides you through creating an application workspace from your terminal.

Choose your project location, backend framework, optional frontend, and Kubernetes settings. Devspace scaffolds application files and generates container and deployment configuration for you to customize.

The npm package is **`devspace-kit`**. The installed terminal command is **`devspace`**.

## Features

- **Interactive setup** — configure your project through terminal prompts.
- **Backend choices** — Go, Django, Express, and Actix-web.
- **Optional frontend** — React, Next.js, Vue, Svelte, or backend only.
- **Flexible location** — generate projects inside your chosen directory.
- **Input validation** — validate paths and normalize service names.
- **Docker configuration** — generate a Dockerfile with separate runtime targets.
- **Kubernetes manifests** — generate Deployments and Services with replicas, resource settings, and health probes.
- **Terminal interface** — yellow branding, white sparkle animation, and cyan help headings.

## Installation

Install **Node.js 22 or newer**, including npm.

Check that both are available:

```bash
node --version
npm --version
```

Install Devspace globally:

```bash
npm install -g devspace-kit
```

Verify the installation:

```bash
devspace --version
```

### Show the installation welcome screen

```bash
npm install -g devspace-kit --foreground-scripts
```

npm normally hides successful installation-script output. The standard installation command still installs Devspace.

### Add npm’s global command directory to PATH, if needed

If `devspace --version` works, skip this step.

If your terminal cannot find the command after installation, check npm’s global installation directory:

```bash
npm prefix -g
```

**Windows**

Globally installed npm commands are placed directly in the returned directory. Add that directory to your **user PATH**:

1. Search Windows for **Edit environment variables for your account**.
2. Under **User variables**, select **Path → Edit → New**.
3. Paste the directory returned by `npm prefix -g`.
4. Save and fully restart your terminal application.

A typical location is:

```text
C:\Users\YourName\AppData\Roaming\npm
```

Use the actual path returned on your computer.

Verify in a new CMD window:

```cmd
where devspace
devspace --version
```

Adding this directory to your user PATH makes npm’s global commands available from any folder in your account. A system-wide PATH change is not required.

**macOS / Linux**

Globally installed npm commands are placed in the `bin` directory under the returned prefix. Ensure that directory is on your shell’s `PATH`, then restart your terminal.

Verify:

```bash
command -v devspace
devspace --version
```

### Update

```bash
npm install -g devspace-kit@latest
```

### Uninstall

```bash
npm uninstall -g devspace-kit
```

Uninstalling Devspace does not delete projects you generated.

## Quick Start

Create your first project:

```bash
devspace create
```

Choose your settings in the wizard:

| Setting | Example |
| --- | --- |
| Project location | `.` |
| Service name | `inventory-api` |
| Backend | Node.js (Express) |
| Frontend | React (Vite) |
| Replicas | `1` |
| Backend CPU request | `100m` |
| Backend memory request | `128Mi` |
| Backend Service type | `ClusterIP` |

The project is created inside:

```text
<selected-directory>/inventory-api
```

**No GitHub token is required for local project generation.**

To open the animated welcome screen:

```bash
devspace
```

Press **Enter** or **Ctrl+C** to return to your terminal.

## Commands

| Command | Description |
| --- | --- |
| `devspace` | Open the animated welcome screen |
| `devspace create` | Start the interactive project wizard |
| `devspace --help` | Show commands, features, and requirements |
| `devspace create --help` | Explain the wizard settings |
| `devspace --version` | Show the installed version |
| `devspace -v` | Show the installed version |
| `devspace completion --help` | Explore shell completion options |

Project settings are entered through interactive prompts.

## Supported Platforms

| Operating system | Packaged architectures |
| --- | --- |
| Windows | x64 |
| macOS | x64 and ARM64 / Apple Silicon |
| Linux, including Ubuntu | x64 and ARM64 |

A native Windows ARM64 executable is not currently included.

## Requirements

To install and launch Devspace, you need Node.js 22 or newer with npm.

Install the additional tools required by your selected stack and make them available on `PATH`.

| Selection | Required tools |
| --- | --- |
| Go backend | Go |
| Django backend | Python and Django, including `django-admin` |
| Express backend | Node.js and npm |
| Actix-web backend | Rust and Cargo for development and builds |
| Frontend frameworks | Node.js, npm, and npx |
| Container builds | Docker |
| Kubernetes deployment | kubectl and access to a configured cluster |

The npm package includes compiled executables. Go is not required just to launch Devspace.

Framework generators may download dependencies and require newer runtime versions than the CLI launcher’s minimum.

## Supported Stacks

### Backends

| Language | Framework or starter |
| --- | --- |
| Go | Standard-library HTTP server |
| Python | Django |
| JavaScript | Express |
| Rust | Actix-web |

### Frontends

| Framework | Generator |
| --- | --- |
| React | Vite |
| Next.js | Create Next App |
| Vue | Vite |
| Svelte | Svelte CLI |
| None | Backend-only workspace |

Generated frontend containers serve static files through Nginx. Next.js is configured for static export, and Svelte targets static output. Server-rendered workloads require changes to the generated setup.

## Configure Your Project

### Project location

Choose an existing parent directory.

| Input | Meaning |
| --- | --- |
| Enter or `.` | Current terminal directory |
| `~` | Your home directory |
| `"C:\Users\YourName\Desktop"` | Windows directory |
| `/home/yourname/projects` | Linux directory |
| `/Users/yourname/projects` | macOS directory |

Quoted paths are accepted. Create the parent directory first if it does not exist.

Use a new project name to avoid overwriting files in an existing workspace.

### Service name

Service names:

- Contain between 1 and 46 characters after normalization.
- Start with a letter.
- Use letters, numbers, and hyphens.
- End with a letter or number.
- Are converted to lowercase.
- Have spaces converted to hyphens.
- Cannot use Windows reserved names such as `con` or `nul`.

For example:

```text
Inventory API → inventory-api
```

### Kubernetes settings

| Setting | Available values |
| --- | --- |
| Replicas | Positive integer; default `1` |
| Backend CPU request | `100m`, `250m`, `500m`, `1000m` |
| Backend memory request | `128Mi`, `256Mi`, `512Mi`, `1024Mi` |
| Backend Service type | `ClusterIP`, `NodePort`, `LoadBalancer` |

Resource limits are generated from the selected requests.

When a frontend is included, its current template uses a `LoadBalancer` Service, a `100m` CPU request, a `128Mi` memory request, and the selected replica count.

Edit the generated manifests to adjust these settings.

### Wizard controls

| Key | Action |
| --- | --- |
| Arrow keys | Navigate choices |
| Enter | Confirm an answer |
| Ctrl+C | Cancel the current prompt |

## Generated Workspace

For a full-stack project named `inventory-api`:

| Path | Purpose |
| --- | --- |
| `backend/` | Backend application and framework files |
| `frontend/` | Frontend application and framework files |
| `Dockerfile` | Container build stages and runtime targets |
| `.dockerignore` | Docker build-context exclusions |
| `k8s/inventory-api-backend-deployment.yaml` | Backend Deployment |
| `k8s/inventory-api-backend-service.yaml` | Backend Service |
| `k8s/inventory-api-frontend-deployment.yaml` | Frontend Deployment |
| `k8s/inventory-api-frontend-service.yaml` | Frontend Service |

All paths are relative to the generated project directory.

Backend-only projects omit the frontend directory and frontend manifests. Framework-specific files depend on the selected stack and generator version.

## Run Your Application Locally

For an Express backend:

```bash
cd inventory-api/backend
node index.js
```

The generated Express starter listens on port `8080`.

For a React/Vite frontend, open another terminal:

```bash
cd inventory-api/frontend
npm install
npm run dev
```

Use the development commands appropriate to your selected framework. Configure API URLs, environment variables, and application routes for your project.

## Containerize Your Application

From the generated project root:

```bash
cd inventory-api
```

Build the backend image:

```bash
docker build --target backend-runtime -t inventory-api-backend:local .
```

If you included a frontend:

```bash
docker build --target frontend-runtime -t inventory-api-frontend:local .
```

Review base images, dependency versions, and build commands before building. Framework-specific adjustments may be necessary.

## Deploy to Kubernetes

Before applying the generated manifests:

1. Build your application images.
2. Push them to an accessible registry or load them into your local cluster.
3. Update Deployment `image:` fields to match your image names and tags.
4. Review ports, health probes, resources, Service types, and application configuration.
5. Confirm the target Kubernetes context.

The current generated image references use the `patel-jay/` namespace. Replace them with your own registry paths.

Check your context:

```bash
kubectl config current-context
```

Apply the reviewed manifests from the project root:

```bash
kubectl apply -f k8s/
```

Inspect the resources:

```bash
kubectl get deployments
kubectl get pods
kubectl get services
```

A `LoadBalancer` Service requires support from your cluster environment. Local clusters may need additional configuration.

## Terminal Appearance

Devspace displays yellow branding, a file icon, white sparkles, and cyan help headings.

For the best experience:

- Use a Unicode-capable terminal with a monospace font.
- Keep the window at least 90 columns wide and 24 lines tall.
- Avoid resizing while the animation runs.
- Press Enter or Ctrl+C before entering your next shell command.

The installation welcome screen and help output are static. Redirected output uses plain text.

### Disable colors and animation

**Windows CMD:**

```cmd
set "NO_COLOR=1"
devspace
```

Restore colors:

```cmd
set "NO_COLOR="
```

**macOS / Linux:**

```bash
NO_COLOR=1 devspace
```

## Troubleshooting

### Installation shows no banner

Successful installer output is normally hidden by npm. Use:

```bash
npm install -g devspace-kit --foreground-scripts
```

Or run `devspace` after installation.

### The terminal cannot find `devspace`

Follow the optional PATH setup under [Installation](#installation). Ensure you add the directory returned by npm for your computer, then restart your terminal application.

### A framework command is missing

Install the tools required by your selected stack and make them available on `PATH`.

Depending on the selection, Devspace may invoke `go`, `django-admin`, `cargo`, `npm`, or `npx`.

### Generation leaves missing files

Inspect the terminal output and generated workspace. Some framework subprocess failures currently print an error without stopping the entire generation workflow.

Resolve the reported toolchain or dependency issue before continuing.

## Current Scope

Devspace provides scaffolding and editable container and deployment templates.

- Cloud infrastructure and Kubernetes clusters are not provisioned automatically.
- Images are not automatically built or pushed.
- Generated manifests are not automatically applied.
- Applications require configuration and review before production use.
- The Django template currently uses Django’s development server.
- Frontend-to-backend integration remains application-specific.
- Framework and toolchain changes may require template adjustments.
- Project generation is interactive rather than a fully unattended CI workflow.

## Contributing

Contributions to documentation, framework templates, terminal behavior, validation, and platform compatibility are welcome.

### Suggest an improvement or report a bug

Open a [GitHub issue](https://github.com/jay240322/Devspace/issues).

For bug reports, include:

- Your operating system and architecture.
- Your Devspace version.
- The selected backend and frontend.
- Steps to reproduce the problem.
- Relevant error output, with credentials removed.

### Submit a change

1. Fork the repository.
2. Create a branch for your change.
3. Keep the change focused and update relevant documentation.
4. Verify the affected behavior.
5. Open a pull request describing what changed, why, and how it was tested.

For framework changes, generate a sample project and verify its relevant development, container, or deployment workflow.

For terminal changes, check the welcome screen, help output, and behavior with colors disabled.

## License

Licensed under the [Apache License 2.0](LICENSE).

---

<div align="center">

[View on npm](https://www.npmjs.com/package/devspace-kit) · [Report an issue](https://github.com/jay240322/Devspace/issues) · [Contribute](https://github.com/jay240322/Devspace/pulls)

</div>
