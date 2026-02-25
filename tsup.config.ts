import { defineConfig } from 'tsup';

export default defineConfig({
  entry: ['plugins/vite-process-assets.ts', 'plugins/config.ts'],
  format: ['cjs'],
  dts: true,
  sourcemap: true,
  clean: true,
  splitting: false,
});
