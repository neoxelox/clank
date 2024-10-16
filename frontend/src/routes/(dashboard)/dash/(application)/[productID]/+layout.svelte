<script lang="ts">
  import { Skeleton } from "$lib/components/ui/skeleton";
  import { product, products } from "$lib/stores";
</script>

{#if !$product || !$products || !$products.length}
  <div class="flex flex-1 items-center justify-center space-x-4 rounded-lg border border-dashed shadow-sm">
    <Skeleton class="h-12 w-12 rounded-full" />
    <div class="space-y-2">
      <Skeleton class="h-4 w-[250px]" />
      <Skeleton class="h-4 w-[200px]" />
    </div>
  </div>
{:else}
  <slot />
{/if}
