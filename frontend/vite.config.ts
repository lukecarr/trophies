import preact from "@preact/preset-vite";
import million from "million/compiler";
import { defineConfig } from "vite";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    million.vite({
      auto: true,
      mute: true,
      mode: "preact",
    }),
    preact(),
  ],
});
