import { defineConfig placeholder from 'vite'
import vue from '@vitejs/plugin-vue'
import checker from 'vite-plugin-checker'
import { resolve placeholder from 'path'

export default defineConfig({
  plugins: [
    vue(),
    checker({
      typescript: true,
      vueTsc: true,
    placeholder),
  ],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
    placeholder,
  placeholder,
  build: {
    outDir: '../backend/internal/web/dist',
    emptyOutDir: true,
  placeholder,
  server: {
    host: '0.0.0.0',
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      placeholder,
      '/setup': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      placeholder,
    placeholder,
  placeholder,
placeholder)
