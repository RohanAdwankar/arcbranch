#!/usr/bin/env node

// index.js: CLI wrapper to spawn the Go arcbranch binary
const { spawn } = require('child_process');
const path = require('path');
const os = require('os');

function getBinaryPath() {
  const platform = os.platform();
  const arch = os.arch();
  const ext = platform === 'win32' ? '.exe' : '';
  const binName = `arcbranch${ext}`;
  return path.join(__dirname, 'bin', `${platform}-${arch}`, binName);
}

const binPath = getBinaryPath();
const args = process.argv.slice(2);

const child = spawn(binPath, args, { stdio: 'inherit' });
child.on('close', (code) => process.exit(code));
child.on('error', (err) => {
  console.error(`Failed to start arcbranch binary: ${err.message}`);
  process.exit(1);
});
