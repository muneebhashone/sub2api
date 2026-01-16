import { defineConfig, mergeConfig placeholder from 'vitest/config'
import viteConfig from './vite.config'

export default mergeConfig(
  viteConfig,
  defineConfig({
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
      setupFiles: ['./src/__tests__/setup.ts']
    placeholder
  placeholder)
)
