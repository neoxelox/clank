<script lang="ts">
  export let href: string | undefined = undefined;
  export let target: string | undefined = undefined;
  export let rel: string | undefined = undefined;
  export let color: string = "primary";

  let classes =
    "group inline-flex items-center justify-center rounded-full py-2.5 px-4 text-sm font-semibold transition-all focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring hover:shadow-sm hover:opacity-90" +
    " " +
    {
      primary: "text-primary-foreground bg-primary-gradient",
      secondary: "text-secondary-foreground bg-secondary-gradient",
    }[color] +
    " " +
    ($$props.class || "");
</script>

{#if !href}
  <button class={classes} on:click on:keydown>
    <slot />
  </button>
{:else}
  <a {href} {target} {rel} class={classes}>
    <slot />
  </a>
{/if}
