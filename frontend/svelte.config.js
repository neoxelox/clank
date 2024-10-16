import adapter from "@sveltejs/adapter-cloudflare";
import { vitePreprocess } from "@sveltejs/vite-plugin-svelte";

/** @type {import('@sveltejs/kit').Config} */
export default {
  preprocess: vitePreprocess(),
  kit: {
    adapter: adapter(),
    version: {
      name: process.env.CLANK_RELEASE || "wip",
    },
    alias: {
      $assets: "src/assets",
      $lib: "src/lib",
      $routes: "src/routes",
    },
    env: {
      publicPrefix: "CLANK_FRONTEND_PUBLIC",
      privatePrefix: "",
    },
    csp: {
      mode: "auto",
    },
    csrf: {
      checkOrigin: true,
    },
  },
  onwarn: (warning, handler) => {
    if (warning.code.startsWith("a11y")) return;
    if (warning.code === "illegal-attribute-character") return;
    handler(warning);
  },
};
