#!/usr/bin/env node

const { spawn } = require('child_process');
const path = require('path');
const fs = require('fs');

function getBinaryPath(name) {
  const { platform, arch } = process;
  const binDir = path.resolve(__dirname, '../bin');
  const exe = platform === 'win32' ? '.exe' : '';

  const binaryPath = path.join(binDir, `${name}-${platform}-${arch}${exe}`);
  if (!fs.existsSync(binaryPath)) {
    console.error(`Binary ${binaryPath} not found.`);
    process.exit(1);
  }

  return binaryPath;
}

function run() {
  const binaryName = path.basename(process.argv[1]);
  const args = process.argv.slice(2);

  spawn(getBinaryPath(binaryName), args, { stdio: 'inherit' })
    .on('error', (err) => {
      console.error(`Failed to start: ${err.message}`);
      process.exit(1);
    })
    .on('close', (code) => process.exit(code || 0));
}

run();
