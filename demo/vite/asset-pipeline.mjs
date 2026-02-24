import { execSync } from 'child_process';
import path from 'path';
import { fileURLToPath } from 'url';
import fs from 'fs';

function assetPipeline() {
    // TODO: Put your commands for packing assets here
}

function fontpack(src, dst) { run('fontpack', src, dst) }
function atlaspack(src, dst) { run('atlaspack', src, dst) }

function copy(src, dst) {
    const srcPath = path.join(rootDir, sourceAssetDir, src)
    const dstPath = path.join(rootDir, packedAssetDir, dst ?? '', src)
    console.log(`Copying ${src} to ${dstShort}...`)
    fs.copyFileSync(srcPath, dstPath)
}

function run(cmd, src, dst) {
    const srcPath = path.join(sourceAssetDir, src)
    const dstPath = path.join(packedAssetDir, dst ?? '')
    execSync(`${cmd} ${srcPath} ${dstPath}`, { cwd: rootDir, stdio: 'inherit' })
}

const __dirname = path.dirname(fileURLToPath(import.meta.url))
const rootDir = path.resolve(__dirname, '..')

const sourceAssetDir = 'assets'
const packedAssetDir = 'public/packed_assets'

export function runAssetPipeline() {
    console.log('Running asset pipeline...')
    try {
        fs.rmSync(path.join(rootDir, packedAssetDir), { recursive: true, force: true })
        assetPipeline()
    } catch (error) {
        console.error('Asset pipeline failed:', error.message);
        throw error;
    }
}
