import { defineConfig } from 'vite';
import { processAssetsProd } from 'pixel-tools';
import { assetsConfig } from "./assets.mjs";

export default defineConfig({
    base: './',
    logLevel: 'warning',
    define: {},
    resolve: {
        dedupe: ['phaser'],
    },
    build: {
        rollupOptions: {
        },
        minify: 'terser',
        terserOptions: {
            compress: {
                passes: 2
            },
            mangle: true,
            format: {
                comments: false
            }
        }
    },
    server: {
        port: 8080
    },
    plugins: [
        processAssetsProd(assetsConfig),
        {
            name: 'build-notifications',
            buildStart() { console.log('Building application...') },
            renderStart() { console.log('Minifying application...') }
        }
    ]
});
