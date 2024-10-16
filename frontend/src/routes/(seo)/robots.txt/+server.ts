import { CLANK_FRONTEND_PUBLIC_BASE_URL } from "$env/static/public";
import type { RequestHandler } from "@sveltejs/kit";
import { text } from "@sveltejs/kit";

export const GET: RequestHandler = async () => {
  return text(
    `
User-agent: *
Disallow: /dash/

Sitemap: ${CLANK_FRONTEND_PUBLIC_BASE_URL}/sitemap.xml
    `.trim(),
    {
      headers: {
        "Cache-Control": "max-age=0, s-maxage=3600",
        "Content-Type": "text/plain",
      },
    },
  );
};
