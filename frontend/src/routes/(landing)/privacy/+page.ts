import type { PageLoad } from "./$types";

export const load: PageLoad = () => {
  return {
    meta: {
      title: "Privacy Policy",
      description:
        "Discover how Clank processes your personal information, your privacy rights, and our data protection measures. Stay informed and confident using our services.",
    },
  };
};
