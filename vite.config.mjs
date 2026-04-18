import path from "node:path";
import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import tailwindcss from "@tailwindcss/vite";

export default defineConfig(({ mode }) => ({
  base: "/static/dist/",
  plugins: [vue(), tailwindcss()],
  publicDir: false,
  build: {
    outDir: "public/dist",
    emptyOutDir: true,
    sourcemap: mode === "development",
    minify: mode !== "development",
    rollupOptions: {
      input: path.resolve("resources/js/main.js"),
      output: {
        entryFileNames: "app.js",
        chunkFileNames: "chunks/[name].js",
        assetFileNames: (assetInfo) => {
          if (assetInfo.name && assetInfo.name.endsWith(".css")) return "app.css";
          return "assets/[name][extname]";
        },
        manualChunks: undefined
      }
    }
  }
}));
