#!/usr/bin/env node
"use strict";

const path = require("node:path");
const { spawnSync } = require("node:child_process");

const supported = new Set([
  "win32-x64",
  "darwin-x64",
  "darwin-arm64",
  "linux-x64",
  "linux-arm64"
]);

const target = process.platform + "-" + process.arch;

if (!supported.has(target)) {
  console.error("Devspace does not support " + target + " yet.");
  process.exit(1);
}

const extension = process.platform === "win32" ? ".exe" : "";

const executable = path.join(
  __dirname,
  "devspace-" + target + extension
);

const result = spawnSync(executable, process.argv.slice(2), {
  stdio: "inherit",
  cwd: process.cwd(),
  env: process.env,
  shell: false
});

if (result.error) {
  console.error("Unable to start Devspace: " + result.error.message);
  process.exit(1);
}

process.exit(
  result.status ?? (result.signal === "SIGINT" ? 130 : 1)
);