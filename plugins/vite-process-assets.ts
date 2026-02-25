import { execSync } from 'child_process';
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

                console.log(`Asset changed: ${filePath}`);
                try {
                    runPipeline(config, root);
                    server.ws.send({ type: 'full-reload' });
                } catch (error: any) {
                    console.error('Asset pipeline update failed:', error.message);
                }
            };

            server.watcher.on('change', update);
            server.watcher.on('add', update);
            server.watcher.on('unlink', update);
        }
    };
}

function runPipeline(config: AssetsConfig, rootDir: string): void {
    const fullSourcePath = path.resolve(rootDir, config.source_path);
    const fullDestPath = path.resolve(rootDir, config.destination_path);

    if (!fs.existsSync(fullDestPath)) {
        fs.mkdirSync(fullDestPath, { recursive: true });
    }

    const paddingArg = config.padding ? ` -padding ${config.padding}` : '';
    const run = (cmd: string, src: string, dst: string) => {
        src = path.join(fullSourcePath, src);
        dst = path.join(fullDestPath, dst)
        execSync(`${cmd}${paddingArg} ${src} ${dst}`, { cwd: rootDir, stdio: 'inherit' });
    }

    for (const font of config.fonts || []) {
        font.target ??= path.dirname(font.source)
        run('fontpack', font.source, font.target);
    }

    for (const atlas of config.atlases || []) {
        const parsed = path.parse(atlas.source);
        atlas.target ??= path.join(parsed.dir, parsed.name);
        run('atlaspack', atlas.source, atlas.target);
    }

    for (const tilemap of config.tilemaps || []) {
        tilemap.target ??= tilemap.source
        run('tilepack', tilemap.source, tilemap.target);
    }
}
