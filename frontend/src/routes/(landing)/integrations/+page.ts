import type { PageLoad } from "./$types";

export const load: PageLoad = () => {
  return {
    meta: {
      title: "Integrations",
      description:
        "Capture feedback wherever your customers are. Dozens of integrations available so you don't miss any feedback.",
    },
  };
};
