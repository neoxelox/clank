<script lang="ts">
  import { LinkIcon } from "@babeard/svelte-heroicons/outline";
  import { getContext, setContext } from "svelte";
  import { copy } from "svelte-copy";

  type $$Props = {
    id?: string;
    ["id:auto"]?: boolean;
  };
  export let id: $$Props["id"] = undefined;
  $$restProps;

  const scope: string = getContext("content-scope");
  setContext("content-scope", "Title");
</script>

{#if scope === "Introduction"}
  {#if id}
    <a href={`#${id}`} class="group relative" use:copy={`${window.location.href.split(/[?#]/)[0]}#${id}`}>
      <LinkIcon
        class="absolute -left-2.5 bottom-1 h-5 w-5 shrink-0 -translate-x-full stroke-2 text-foreground opacity-0 transition-opacity group-hover:opacity-100"
      />
      <h1 {id} class="pt-2 font-cal text-3xl sm:text-4xl">
        <slot />
      </h1>
    </a>
  {:else}
    <h1 {id} class="pt-2 font-cal text-3xl sm:text-4xl">
      <slot />
    </h1>
  {/if}
{:else if scope === "Body" || scope === "Conclusion"}
  {#if id}
    <a
      href={`#${id}`}
      class="group relative [&+*]:!mt-6 [&+h3]:!mt-0 [&+h3]:!pt-6"
      use:copy={`${window.location.href.split(/[?#]/)[0]}#${id}`}
    >
      <LinkIcon
        class="absolute -left-2.5 bottom-1 h-5 w-5 shrink-0 -translate-x-full stroke-2 text-foreground opacity-0 transition-opacity group-hover:opacity-100"
      />
      <h2 {id} class="pt-12 font-cal text-2xl">
        <slot />
      </h2>
    </a>
  {:else}
    <h2 {id} class="pt-12 font-cal text-2xl [&+*]:!mt-6 [&+h3]:!mt-0 [&+h3]:!pt-6">
      <slot />
    </h2>
  {/if}
{/if}
