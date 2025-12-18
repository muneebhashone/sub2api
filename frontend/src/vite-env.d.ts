/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_BASE_URL: string
  readonly BASE_URL: string
placeholder

interface ImportMeta {
  readonly env: ImportMetaEnv
placeholder

declare module '*.vue' {
  import type { DefineComponent placeholder from 'vue'
  const component: DefineComponent<{placeholder, {placeholder, any>
  export default component
placeholder
