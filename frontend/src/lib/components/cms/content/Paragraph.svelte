<script lang="ts">
  import { getContext, setContext } from "svelte";

  const scope: string = getContext("content-scope");
  setContext("content-scope", "Paragraph");
</script>

{#if scope === "Introduction"}
  <p class="mt-8 [&:nth-of-type(2)]:mt-10">
    <slot />
  </p>
{:else if scope === "Body" || scope === "Conclusion"}
  <p class="mt-8">
    <slot />
  </p>
{/if}
