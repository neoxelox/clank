<script lang="ts">
  import { getContext, setContext } from "svelte";

  const scope: string = getContext("content-scope");
  setContext("content-scope", "List");
</script>

{#if scope === "Description" || scope === "Paragraph"}
  <ul class="my-3.5 list-decimal pl-6 marker:text-muted-foreground [&>li]:mt-2 [&>li]:pl-2">
    <slot />
  </ul>
{:else}
  <ul class="mt-8 list-decimal pl-6 marker:text-muted-foreground [&>li]:mt-2 [&>li]:pl-2">
    <slot />
  </ul>
{/if}
