import { execSync } from 'child_process';
import path from 'path';
import { fileURLToPath } from 'url';
import fs from 'fs';

function assetPipeline() {
    atlaspack('sprites.yaml', 'sprites')
    tilepack('tilemaps')
    fontpack('fonts.yaml')
}

function fontpack(src, dst) { run('fontpack', src, dst) }
function atlaspack(src, dst) { run('atlaspack', src, dst) }
function tilepack(src, dst) { run('tilepack', src, dst) }

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
        assetPipeline()
    } catch (error) {
        console.error('Asset pipeline failed:', error.message);
        throw error;
    }
}
