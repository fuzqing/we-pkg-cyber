import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  base: './',
  build: {
    // 确保 public 目录下的文件保持原样复制，不经过资源处理
    copyPublicDir: true,
    // 配置资源处理，排除 WASM 文件
    rollupOptions: {
      external: [],
    },
  },
  // 配置静态资源处理
  assetsInclude: [],
})
