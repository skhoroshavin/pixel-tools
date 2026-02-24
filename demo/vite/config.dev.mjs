import { defineConfig } from 'vite';
import path from 'path';
import { runAssetPipeline } from './asset-pipeline.mjs';

export default defineConfig({
    base: './',
    define: {},
    server: {
        port: 8080,
        open: true,
    },
    plugins: [
        {
            name: 'asset-pipeline-watcher',
            configureServer(server) {
                const assetDir = path.resolve(server.config.root, 'assets')
                server.watcher.add(assetDir);

                const update = (filePath) => {
                    if (filePath.includes('/public/packed_assets/')) return
                    console.log(`Asset changed: ${filePath}`)
                    runAssetPipeline()
                    server.ws.send({ type: 'full-reload' })
                }

                server.watcher.on('change', update)
                server.watcher.on('add', update)
                server.watcher.on('unlink', update)
            }
        }
    ]
});
