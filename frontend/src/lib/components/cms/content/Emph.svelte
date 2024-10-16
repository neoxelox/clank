<script lang="ts">
  import { getContext, setContext } from "svelte";

  const scope: string = getContext("content-scope");
  setContext("content-scope", "Emph");
</script>

{#if scope === "Description"}
  <code class="relative rounded bg-muted px-[0.3rem] py-[0.2rem] font-mono text-lg font-semibold">
    <slot />
  </code>
{:else}
  <code class="relative rounded bg-muted px-[0.3rem] py-[0.2rem] font-mono text-sm font-semibold">
    <slot />
  </code>
{/if}
