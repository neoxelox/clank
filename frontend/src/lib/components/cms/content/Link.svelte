<script lang="ts">
  import { CLANK_FRONTEND_PUBLIC_BASE_URL } from "$env/static/public";
  import { setContext } from "svelte";

  export let href: string;

  setContext("content-scope", "Link");
</script>

{#if href.startsWith(CLANK_FRONTEND_PUBLIC_BASE_URL)}
  <a {href} target="_blank" class="text-primary">
    <slot />
  </a>
{:else}
  <a {href} target="_blank" rel="noopener noreferrer" class="text-primary">
    <slot />
  </a>
{/if}
