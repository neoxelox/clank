module.exports = {
  overrides: [
    { files: "*.svelte", options: { parser: "svelte" } },
    { files: "src/worker.d.ts", options: { requirePragma: true } },
    { files: "src/lib/components/ui/**/*", options: { requirePragma: true } },
    { files: "src/lib/utils/ui.ts", options: { requirePragma: true } },
    { files: "src/_ui.scss", options: { requirePragma: true } },
  ],
  plugins: ["prettier-plugin-svelte", "prettier-plugin-tailwindcss"],
  tailwindConfig: ".tailwindrc.cjs",
  semi: true,
  singleQuote: false,
  trailingComma: "all",
  printWidth: 120,
  useTabs: false,
  tabWidth: 2,
  endOfLine: "auto",
};
