import { defineConfig } from '@hey-api/openapi-ts'

export default defineConfig({
  input: '../backend/docs/swagger.yaml',
  output: 'src/api',
  plugins: ['@hey-api/client-fetch', '@tanstack/react-query'],
})
