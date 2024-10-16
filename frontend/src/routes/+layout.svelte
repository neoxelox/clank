<script lang="ts">
  import { page } from "$app/stores";
  import { CLANK_FRONTEND_PUBLIC_BASE_URL } from "$env/static/public";
  import { trimSuffix } from "$lib/utils/string";
  import "../app.scss";

  $: meta = <App.MetaData>{
    title: ($page.data.meta?.title || "Customer feedback analysis on autopilot") + " - Clank",
    description:
      $page.data.meta?.description ||
      "Collect and analyze customer feedback, identify issues and suggestions, predict churn and retention, track CX KPIs and synthesize interactive customer personas.",
    image: $page.data.meta?.image || CLANK_FRONTEND_PUBLIC_BASE_URL + "/images/link-preview.png",
  };
</script>

<svelte:head>
  <title>{meta.title}</title>
  <meta name="description" content={meta.description} />
  <meta name="twitter:card" content="summary" />
  <meta name="twitter:site" content="@clank_so" />
  <meta name="twitter:creator" content="@clank_so" />
  <meta name="twitter:title" content={meta.title} />
  <meta name="twitter:description" content={meta.description} />
  <meta name="twitter:image" content={meta.image} />
  <meta property="og:type" content="website" />
  <meta property="og:site_name" content={meta.title} />
  <meta property="og:title" content={meta.title} />
  <meta property="og:description" content={meta.description} />
  <meta property="og:url" content={trimSuffix(CLANK_FRONTEND_PUBLIC_BASE_URL + $page.url.pathname, "/")} />
  <meta property="og:image" content={meta.image} />
</svelte:head>

<slot />
