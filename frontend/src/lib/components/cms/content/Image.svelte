<script lang="ts">
  import { imgPlaceholder } from "$lib/utils/image";
  import { setContext } from "svelte";

  export let src: string;

  setContext("content-scope", "Image");
</script>

<figure class="mt-8">
  <img
    {src}
    onerror={imgPlaceholder}
    fetchpriority="auto"
    loading="lazy"
    decoding="async"
    class="aspect-video max-h-[500px] w-full rounded-3xl bg-muted object-cover"
  />
  {#if $$slots.default}
    <figcaption class="mt-4 text-center text-sm leading-5 text-muted-foreground">
      <slot />
    </figcaption>
  {/if}
</figure>
