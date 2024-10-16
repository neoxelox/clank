import { CLANK_FRONTEND_PUBLIC_BASE_URL } from "$env/static/public";
import type { RequestHandler } from "@sveltejs/kit";
import * as sitemap from "super-sitemap";

export const GET: RequestHandler = async () => {
  return await sitemap.response({
    origin: CLANK_FRONTEND_PUBLIC_BASE_URL,
    changefreq: "daily",
    priority: 1.0,
    sort: "alpha",
    headers: {
      "Cache-Control": "max-age=0, s-maxage=3600",
      "Content-Type": "application/xml",
    },
    excludeRoutePatterns: [".*\\(dashboard\\).*"],
  });
};
