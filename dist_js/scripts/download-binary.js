const https = require('https');
const fs = require('fs');
const path = require('path');
const os = require('os');
const { execSync } = require('child_process');

// download-binary.js: downloads the correct arcbranch binary into dist_js/bin/<platform-arch>

// map for matching GitHub assets
const PLATFORM_MAP = { darwin: 'Darwin', linux: 'Linux', win32: 'Windows' };
const ARCH_MAP = { x64: 'x86_64', ia32: 'i386', arm64: 'arm64' };

function getDownloadUrl() {
  const pkg = require(path.join(__dirname, '..', 'package.json'));
  const version = pkg.version;
  const rawPlat = os.platform();
  const rawArch = os.arch();
  const plat = PLATFORM_MAP[rawPlat] || rawPlat;
  const arch = ARCH_MAP[rawArch] || rawArch;
  const ext = plat === 'Windows' ? '.zip' : '.tar.gz';
  return `https://github.com/RohanAdwankar/arcbranch/releases/download/v${version}/arcbranch_${plat}_${arch}${ext}`;
}

function download(url, dest, cb) {
  const file = fs.createWriteStream(dest);
  https.get(url, (res) => {
    if (res.statusCode !== 200) {
      return cb(new Error(`Download failed: ${res.statusCode}`));
    }
    res.pipe(file);
    file.on('finish', () => file.close(cb));
  }).on('error', (err) => {
    fs.unlink(dest, () => cb(err));
  });
}

function extract(archive, dest, cb) {
  if (archive.endsWith('.zip')) {
    execSync(`unzip -o ${archive} -d ${dest}`);
  } else {
    execSync(`tar -xzf ${archive} -C ${dest}`);
  }
  cb();
}

function ensureDir(dir) {
  if (!fs.existsSync(dir)) fs.mkdirSync(dir, { recursive: true });
}

function main() {
  const url = getDownloadUrl();
  const rawPlat = os.platform();
  const rawArch = os.arch();
  const plat = PLATFORM_MAP[rawPlat] || rawPlat;
  const arch = ARCH_MAP[rawArch] || rawArch;
  const ext = plat === 'Windows' ? '.zip' : '.tar.gz';
  const binDir = path.join(__dirname, '..', 'bin', `${rawPlat}-${rawArch}`);
  ensureDir(binDir);
  const archive = path.join(binDir, `arcbranch_${plat}_${arch}${ext}`);

  console.log(`Downloading arcbranch binary from ${url}...`);
  download(url, archive, (err) => {
    if (err) {
      console.error(err);
      process.exit(1);
    }
    console.log('Extracting...');
    extract(archive, binDir, (err2) => {
      if (err2) {
        console.error(err2);
        process.exit(1);
      }
      console.log('Binary downloaded and extracted to', binDir);
      fs.chmodSync(path.join(binDir, `arcbranch${plat === 'Windows' ? '.exe' : ''}`), 0o755);
    });
  });
}

main();
