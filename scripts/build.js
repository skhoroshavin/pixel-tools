const { execSync } = require('child_process');
const path = require('path');
const fs = require('fs');

const platforms = [
  { os: 'darwin', arch: 'x64', goos: 'darwin', goarch: 'amd64' },
  { os: 'darwin', arch: 'arm64', goos: 'darwin', goarch: 'arm64' },
  { os: 'linux', arch: 'x64', goos: 'linux', goarch: 'amd64' },
  { os: 'linux', arch: 'arm64', goos: 'linux', goarch: 'arm64' },
  { os: 'win32', arch: 'x64', goos: 'windows', goarch: 'amd64', ext: '.exe' },
];

const binDir = path.join(__dirname, '..', 'bin');

if (!fs.existsSync(binDir)) {
  fs.mkdirSync(binDir);
}

for (const platform of platforms) {
  const suffix = `-${platform.os}-${platform.arch}${platform.ext || ''}`;
  console.log(`Building for ${platform.os}-${platform.arch}...`);
  try {
    // Build all binaries using Makefile with suffix
    execSync('make build', {
      stdio: 'inherit',
      env: {
        ...process.env,
        GOOS: platform.goos,
        GOARCH: platform.goarch,
        BINARY_SUFFIX: suffix,
      },
    });
  } catch (err) {
    console.warn(`Failed to build for ${platform.os}-${platform.arch}: ${err.message}`);
  }
}
