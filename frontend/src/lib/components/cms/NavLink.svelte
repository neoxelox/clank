<script lang="ts">
  import { PopoverButton } from "@rgossiaux/svelte-headlessui";

  export let href: string;
  export let mobile: boolean = false;
</script>

{#if !mobile}
  <a
    {href}
    class={"hover:bg-secondary-gradient hover:shadow-xs inline-block rounded-xl bg-transparent px-2.5 py-1.5 text-sm font-medium text-secondary-foreground transition-colors " +
      ($$props.class || "")}
  >
    <slot />
  </a>
{:else}
  <PopoverButton as="a" {href} class={"block w-full p-2.5 " + ($$props.class || "")}>
    <slot />
  </PopoverButton>
{/if}
