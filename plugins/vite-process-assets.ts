import { spawnSync } from 'child_process';
import path from 'path';
import fs from 'fs';
import type { Plugin } from 'vite';
import type {AssetsConfig} from './config';

export function processAssetsProd(config: AssetsConfig): Plugin {
    return {
        name: 'production-asset-pipeline',
        apply: 'build',
        buildStart() {
            try {
                runPipeline(config, process.cwd());
            } catch (error: any) {
                console.error('Asset pipeline failed:', error.message);
                process.exit(1);
            }
        }
    };
}

export function processAssetsDev(config: AssetsConfig): Plugin {
    return {
        name: 'dev-assets',
        apply: 'serve',
        configResolved(resolvedConfig) {
            try {
                runPipeline(config, resolvedConfig.root);
            } catch (error: any) {
                console.error('Initial asset generation failed:', error.message);
            }
        },
        configureServer(server) {
            const { root } = server.config;
            server.watcher.add(path.resolve(root, config.source_path));

            const update = (filePath: string) => {
                if (filePath.startsWith(path.resolve(root, config.destination_path))) return;

                console.log(`Asset changed: ${filePath}`)
                try {
                    runPipeline(config, root)
                    server.ws.send({ type: 'full-reload' })
                } catch (error: any) {
                    console.error('Asset pipeline update failed:', error.message)
                }
            }

            server.watcher.on('change', update)
            server.watcher.on('add', update)
            server.watcher.on('unlink', update)
        }
    }
}

function runPipeline(config: AssetsConfig, rootDir: string): void {
    const fullSourcePath = path.resolve(rootDir, config.source_path);
    const fullDestPath = path.resolve(rootDir, config.destination_path);

    if (!fs.existsSync(fullDestPath)) {
        fs.mkdirSync(fullDestPath, { recursive: true })
    }

    const run = (cmd: string, src: string, dst: string) => {
        const resolvedSrc = path.join(fullSourcePath, src)
        const resolvedDst = path.join(fullDestPath, dst)
        const args: string[] = []
        if (config.padding) {
            args.push('-padding', String(config.padding))
        }
        args.push(resolvedSrc, resolvedDst)
        const result = spawnSync(cmd, args, { cwd: rootDir, stdio: 'inherit' })
        if (result.error) throw result.error
        if (result.status !== 0) throw new Error(`${cmd} exited with code ${result.status}`)
    }

    for (const entry of config.copy || []) {
        const src = path.join(fullSourcePath, entry.source)
        const dst = path.join(fullDestPath, entry.target ?? entry.source)
        if (fs.statSync(src).isDirectory()) {
            fs.cpSync(src, dst, { recursive: true })
        } else {
            const dstDir = path.dirname(dst)
            if (!fs.existsSync(dstDir)) {
                fs.mkdirSync(dstDir, { recursive: true })
            }
            fs.copyFileSync(src, dst)
        }
    }

    for (const font of config.fonts || []) {
        const target = font.target ?? path.dirname(font.source)
        run('fontpack', font.source, target)
    }

    for (const atlas of config.atlases || []) {
        const parsedSource = path.parse(atlas.source)
        const target = atlas.target ?? path.join(parsedSource.dir, parsedSource.name)
        run('atlaspack', atlas.source, target)
    }

    for (const tilemap of config.tilemaps || []) {
        const target = tilemap.target ?? tilemap.source
        run('tilepack', tilemap.source, target)
    }
}
