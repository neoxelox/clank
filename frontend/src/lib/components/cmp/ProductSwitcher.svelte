<script lang="ts">
  import ProductDialog from "$lib/components/cmp/ProductDialog.svelte";
  import { Button } from "$lib/components/ui/button";
  import * as Command from "$lib/components/ui/command";
  import { Picture } from "$lib/components/ui/picture";
  import * as Popover from "$lib/components/ui/popover";
  import * as entities from "$lib/entities";
  import { product, products, user } from "$lib/stores";
  import { compare } from "$lib/utils/object";
  import { tick } from "svelte";
  import CaretSort from "svelte-radix/CaretSort.svelte";
  import Check from "svelte-radix/Check.svelte";
  import PlusCircled from "svelte-radix/PlusCircled.svelte";

  $: canUpdate = $user.role === entities.UserRole.ADMIN;

  let open = false;
  let openDialog = false;
</script>

<Popover.Root bind:open let:ids>
  <Popover.Trigger asChild let:builder>
    <Button
      builders={[builder]}
      variant="outline"
      role="combobox"
      class={"w-[calc(220px-2rem)] justify-between lg:w-[calc(280px-3rem)] " + ($$props.class || "")}
    >
      {#if $product}
        <Picture size="md" src={$product.picture} fallback={$product.name} class="mr-2" />
        <span class="overflow-hidden text-ellipsis text-nowrap">{$product.name}</span>
      {:else}
        <span class="text-muted-foreground">Select a product</span>
      {/if}
      <CaretSort class="ml-auto h-4 w-4 shrink-0 opacity-50" />
    </Button>
  </Popover.Trigger>
  <Popover.Content class="w-[calc(220px-2rem-1px)] p-0 lg:w-[calc(280px-3rem-1px)]">
    <Command.Root>
      {#if $products?.length > 0}
        <Command.Input placeholder="Search product..." />
        <Command.List>
          <Command.Empty>No product found.</Command.Empty>
          <Command.Group heading="Starred">
            {#each $products.slice(0, 1) as p}
              <Command.Item
                onSelect={() => {
                  if (!compare(p, $product)) $product = p;
                  open = false;
                  tick().then(() => document.getElementById(ids.trigger)?.focus());
                }}
                value={p.name + ":" + p.id}
                class="text-sm"
              >
                <Picture size="md" src={p.picture} fallback={p.name} class="mr-2" />
                {p.name}
                <Check class={"ml-auto h-4 w-4 " + (p.id === $product?.id ? "text-primary" : "text-transparent")} />
              </Command.Item>
            {/each}
          </Command.Group>
          {#if $products?.length > 1}
            <Command.Group heading="Other">
              {#each $products.slice(1) as p}
                <Command.Item
                  onSelect={() => {
                    if (!compare(p, $product)) $product = p;
                    open = false;
                    tick().then(() => document.getElementById(ids.trigger)?.focus());
                  }}
                  value={p.name + ":" + p.id}
                  class="text-sm"
                >
                  <Picture size="md" src={p.picture} fallback={p.name} class="mr-2" />
                  {p.name}
                  <Check class={"ml-auto h-4 w-4 " + (p.id === $product?.id ? "text-primary" : "text-transparent")} />
                </Command.Item>
              {/each}
            </Command.Group>
          {/if}
        </Command.List>
        <Command.Separator />
      {/if}
      <Command.List>
        <Command.Group>
          <Command.Item
            onSelect={() => {
              open = false;
              openDialog = true;
            }}
            disabled={!canUpdate}
          >
            <PlusCircled class="mr-2 h-5 w-5" />
            Add product
          </Command.Item>
        </Command.Group>
      </Command.List>
    </Command.Root>
  </Popover.Content>
</Popover.Root>
<ProductDialog bind:open={openDialog} />
