import MagicString from "magic-string";
import type { PreprocessorGroup } from "svelte/compiler";
import { parse, walk } from "svelte/compiler";
import type { TemplateNode } from "svelte/types/compiler/interfaces";
import type { Plugin } from "vite";

export default function enhancedHeadings(): Plugin {
  return {
    name: "vite-plugin-enhanced-headings",
    api: {
      sveltePreprocess: <PreprocessorGroup>{
        markup: ({ content, filename }) => {
          if (!content.includes("id:auto")) return;

          const ast = parse(content, { filename });
          const str = new MagicString(content);

          // @ts-expect-error: switch to zimmerframe with Svelte 5
          walk(ast.html, {
            enter(node: TemplateNode) {
              if (node.type !== "Element" && node.type !== "InlineComponent") return;
              if (!node.attributes || node.attributes.length === 0) return;
              if (node.attributes.find(({ name }) => name === "id")) return;
              if (!node.attributes.find(({ name }) => name === "id:auto")) return;
              if (!node.children || node.children.length === 0) return;

              let id = node.children[0].data as string;
              id = id
                .trim()
                .toLowerCase()
                .replace(/\s+/g, "-")
                .replace(/[^a-z0-9-]/gi, "");
              id = `id="${id}"`;

              str.appendLeft(node.attributes[0].start, ` ${id}`);
            },
          });

          return {
            code: str.toString(),
            map: str.generateMap({
              source: filename,
              includeContent: true,
              hires: true,
            }),
          };
        },
      },
    },
  };
}
