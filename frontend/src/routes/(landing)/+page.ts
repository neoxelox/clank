import type { PageLoad } from "./$types";

export const load: PageLoad = () => {
  return {
    meta: {
      title: "Customer feedback analysis on autopilot",
      description:
        "Collect and analyze customer feedback, identify issues and suggestions, predict churn and retention, track CX KPIs and synthesize interactive customer personas.",
    },
  };
};
