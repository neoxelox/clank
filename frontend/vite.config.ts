import { sentryVitePlugin } from "@sentry/vite-plugin";
import { enhancedImages } from "@sveltejs/enhanced-img";
import { sveltekit } from "@sveltejs/kit/vite";
import { svelteInspector } from "@sveltejs/vite-plugin-svelte-inspector";
import { FontaineTransform } from "fontaine";
import { defineConfig } from "vite";
import enhancedHeadings from "./src/lib/plugins/enhanced_headings";

export default defineConfig({
  build: {
    sourcemap: true,
  },
  ssr: {
    noExternal: ["@jill64/sentry-sveltekit-cloudflare"],
  },
  plugins: [
    enhancedImages(),
    enhancedHeadings(),
    sveltekit(),
    FontaineTransform.vite({
      fallbacks: [
        "Avenir",
        "Helvetica",
        "Arial",
        "Georgia",
        "Cambria",
        "Times",
        "BlinkMacSystemFont",
        "Segoe UI",
        "Helvetica Neue",
        "Noto Sans",
      ],
      resolvePath: (id) => `file://static/fonts/${id}`,
    }),
    svelteInspector({
      showToggleButton: "always",
      toggleButtonPos: "bottom-right",
    }),
    sentryVitePlugin({
      org: process.env.CLANK_FRONTEND_SENTRY_ORG || "",
      project: "frontend",
      release: {
        name: process.env.CLANK_RELEASE || "wip",
      },
      authToken: process.env.CLANK_FRONTEND_SENTRY_AUTH_TOKEN || "",
      disable: !process.env.CLANK_FRONTEND_SENTRY_AUTH_TOKEN,
    }),
  ],
  server: {
    host: "0.0.0.0",
    port: Number(process.env.CLANK_FRONTEND_PORT) || 3333,
    strictPort: true,
  },
  preview: {
    host: "0.0.0.0",
    port: Number(process.env.CLANK_FRONTEND_PORT) || 3333,
    strictPort: true,
  },
  clearScreen: false,
});
