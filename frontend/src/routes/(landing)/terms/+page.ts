import type { PageLoad } from "./$types";

export const load: PageLoad = () => {
  return {
    meta: {
      title: "Terms and Conditions",
      description:
        "Read the Terms and Conditions for Clank's SaaS services, including usage restrictions, payment details, confidentiality, and liability limitations.",
    },
  };
};
