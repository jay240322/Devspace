"use strict";

const fs = require("node:fs");
const path = require("node:path");
const { spawnSync } = require("node:child_process");

const supported = new Set([
  "win32-x64",
  "darwin-x64",
  "darwin-arm64",
  "linux-x64",
  "linux-arm64",
]);

const target = process.platform + "-" + process.arch;
const extension = process.platform === "win32" ? ".exe" : "";

const executable = path.resolve(
  __dirname,
  "..",
  "bin",
  "devspace-" + target + extension
);

try {
  if (!supported.has(target)) {
    throw new Error("Unsupported platform: " + target);
  }

  if (!fs.existsSync(executable)) {
    throw new Error("Executable not found: " + executable);
  }

  if (process.platform !== "win32") {
    fs.chmodSync(executable, 0o755);
  }

  const result = spawnSync(executable, ["welcome"], {
    stdio: "inherit",
    shell: false,
  });

  if (result.error) {
    throw result.error;
  }

  if (result.status !== 0) {
    throw new Error(
      "Welcome command failed: " +
        (result.signal || "exit code " + result.status)
    );
  }
} catch (error) {
  console.error("Devspace installation failed: " + error.message);
  process.exit(1);
}