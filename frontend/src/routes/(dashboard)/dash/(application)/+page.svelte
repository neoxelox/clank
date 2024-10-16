<script lang="ts">
  import ProductDialog from "$lib/components/cmp/ProductDialog.svelte";
  import ProductSwitcher from "$lib/components/cmp/ProductSwitcher.svelte";
  import { Button } from "$lib/components/ui/button";
  import * as entities from "$lib/entities";
  import { products, user } from "$lib/stores";

  $: canUpdate = $user.role === entities.UserRole.ADMIN;

  let openProductDialog = false;
</script>

<div class="flex flex-1 items-center justify-center rounded-lg border border-dashed shadow-sm">
  <div class="flex flex-col items-center gap-1 text-center">
    {#if $products?.length}
      <h3 class="text-2xl font-bold tracking-tight">Select a product</h3>
      <p class="text-sm text-muted-foreground">You can start analyzing feedback as soon as you select a product.</p>
      <ProductSwitcher class="mt-4" />
    {:else}
      <h3 class="text-2xl font-bold tracking-tight">You have no products</h3>
      <p class="text-sm text-muted-foreground">You can start analyzing feedback as soon as you add a product.</p>
      <Button
        class="mt-4"
        on:click={() => {
          openProductDialog = true;
        }}
        disabled={!canUpdate}
      >
        Add product
      </Button>
    {/if}
  </div>
</div>
<ProductDialog bind:open={openProductDialog} />
