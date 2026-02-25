import { defineConfig } from 'vite';
import { processAssetsDev } from 'pixel-tools';
import { assetsConfig } from "./assets.mjs";

export default defineConfig({
    base: './',
    define: {},
    server: {
        port: 8080,
        open: true,
    },
    plugins: [
        processAssetsDev(assetsConfig)
    ]
});
