"use strict";

const fs = require("node:fs");
const path = require("node:path");

const target = process.platform + "-" + process.arch;
const extension = process.platform === "win32" ? ".exe" : "";

const executable = path.resolve(
  __dirname,
  "..",
  "bin",
  "devspace-" + target + extension
);

try {
  if (!fs.existsSync(executable)) {
    throw new Error("Executable not found: " + executable);
  }

  if (process.platform !== "win32") {
    fs.chmodSync(executable, 0o755);
  }

  console.log("Devspace is ready for " + target);
} catch (error) {
  console.error("Devspace installation failed: " + error.message);
  process.exit(1);
}