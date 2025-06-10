import { defineConfig, loadEnv } from 'vite'
import react from '@vitejs/plugin-react'
import fs from 'fs'

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  process.env = {...process.env, ...loadEnv(mode, process.cwd())}

  if (!process.env.VITE_SSL_KEY_PATH || !process.env.VITE_SSL_CERT_PATH)
    throw new Error("No SSL set up. You must obtain SSL certificate and set up VITE_SSL_KEY_PATH and VITE_SSL_CERT_PATH environment variables")

  const SSL_KEY_PATH = fs.realpathSync(process.env.VITE_SSL_KEY_PATH)
  const SSL_CERT_PATH = fs.realpathSync(process.env.VITE_SSL_CERT_PATH)

  return defineConfig({
    plugins: [react()],
    server: {
      https: {
        key: fs.readFileSync(SSL_KEY_PATH),
        cert: fs.readFileSync(SSL_CERT_PATH)
      }
    }
  })
})
