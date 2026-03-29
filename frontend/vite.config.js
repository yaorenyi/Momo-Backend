import { defineConfig } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';
import tailwindcss from '@tailwindcss/vite'
import path from 'path';
import cssInjectedByJsPlugin from 'vite-plugin-css-injected-by-js';

export default defineConfig({
  plugins: [
    tailwindcss(),
    svelte({
    compilerOptions: {
      css: 'injected' 
    }
  }), cssInjectedByJsPlugin()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  build: {
    lib: {
      entry: path.resolve(__dirname, 'src/main.ts'),
      name: 'momo',
      fileName: (format) => `momo-comment.min.js`,
      formats: ['iife']
    },
    rollupOptions: {
    }
  }
});