"use strict";

const fs = require("node:fs");
const path = require("node:path");
const { spawnSync } = require("node:child_process");
const { version } = require("../package.json");
const root = path.resolve(__dirname, "..");
const binDirectory = path.join(root, "bin");

fs.mkdirSync(binDirectory, { recursive: true });

const targets = [
  {
    goos: "windows",
    goarch: "amd64",
    filename: "devspace-win32-x64.exe"
  },
  {
    goos: "darwin",
    goarch: "amd64",
    filename: "devspace-darwin-x64"
  },
  {
    goos: "darwin",
    goarch: "arm64",
    filename: "devspace-darwin-arm64"
  },
  {
    goos: "linux",
    goarch: "amd64",
    filename: "devspace-linux-x64"
  },
  {
    goos: "linux",
    goarch: "arm64",
    filename: "devspace-linux-arm64"
  }
];

for (const target of targets) {
  console.log(`Building ${target.filename}...`);

  const result = spawnSync(
    "go",
    [
      "build",
      "-trimpath",
      "-ldflags=-s -w -X devspace/cmd.Version=" + version,
      "-o",
      path.join(binDirectory, target.filename),
      "."
    ],
    {
      cwd: root,
      stdio: "inherit",
      shell: false,
      env: {
        ...process.env,
        CGO_ENABLED: "0",
        GOOS: target.goos,
        GOARCH: target.goarch
      }
    }
  );

  if (result.error || result.status !== 0) {
    console.error(
      result.error?.message || `Build failed: ${target.filename}`
    );
    process.exit(1);
  }
}

console.log("All five executables built successfully.");