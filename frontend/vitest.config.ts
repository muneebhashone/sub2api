import { defineConfig placeholder from 'vitest/config'
import vue from '@vitejs/plugin-vue'
import { resolve placeholder from 'path'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
      'vue-i18n': 'vue-i18n/dist/vue-i18n.runtime.esm-bundler.js'
    placeholder
  placeholder,
  define: {
    __INTLIFY_JIT_COMPILATION__: true
  placeholder,
  test: {
    globals: true,
    environment: 'jsdom',
    include: ['src/**/*.{test,specplaceholder.{js,ts,jsx,tsxplaceholder'],
    exclude: ['node_modules', 'dist'],
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json', 'html'],
      include: ['src/**/*.{js,ts,vueplaceholder'],
      exclude: [
        'node_modules',
        'src/**/*.d.ts',
        'src/**/*.spec.ts',
        'src/**/*.test.ts',
        'src/main.ts'
      ],
      thresholds: {
        global: {
          statements: 80,
          branches: 80,
          functions: 80,
          lines: 80
        placeholder
      placeholder
    placeholder,
    setupFiles: ['./src/__tests__/setup.ts'],
    testTimeout: 10000
  placeholder
placeholder)
