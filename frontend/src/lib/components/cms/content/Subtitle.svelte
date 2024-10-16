<script lang="ts">
  import { getContext, setContext } from "svelte";

  const scope: string = getContext("content-scope");
  setContext("content-scope", "Subtitle");
</script>

{#if scope === "Introduction"}
  <h2 class="pt-8 font-cal text-xl first-of-type:pt-10 [&+*]:!mt-4">
    <slot />
  </h2>
{:else if scope === "Body" || scope === "Conclusion"}
  <h3 class="pt-8 font-cal text-lg [&+*]:!mt-2">
    <slot />
  </h3>
{/if}
