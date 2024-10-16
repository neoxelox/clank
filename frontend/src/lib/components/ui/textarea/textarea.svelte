<script lang="ts">
	import { cn } from "$lib/utils/ui.js";
	import type { HTMLTextareaAttributes } from "svelte/elements";
	import type { TextareaEvents } from "./index.js";

	type $$Props = HTMLTextareaAttributes & { indicator?: boolean };
	type $$Events = TextareaEvents;

	let className: $$Props["class"] = undefined;
	export let value: $$Props["value"] = undefined;
	export { className as class };

	// Workaround for https://github.com/sveltejs/svelte/issues/9305
	// Fixed in Svelte 5, but not backported to 4.x.
	export let readonly: $$Props["readonly"] = undefined;

  export let indicator: boolean = false;
  let textarea: HTMLTextAreaElement;
  let scrollbar: boolean = false;
  $: value, (scrollbar = textarea && (textarea.clientHeight < textarea.scrollHeight))
</script>

{#if indicator && typeof value === "string"}
  <div class="relative">
    <textarea
      bind:this={textarea}
      class={cn(
        "flex min-h-[60px] w-full rounded-md border border-input bg-transparent px-3 py-2 text-sm shadow-sm placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50",
        className
      )}
      bind:value
      {readonly}
      on:blur
      on:change
      on:click
      on:focus
      on:keydown
      on:keypress
      on:keyup
      on:mouseover
      on:mouseenter
      on:mouseleave
      on:paste
      on:input
      {...$$restProps}
    ></textarea>
    <span class={cn(
      "absolute right-2 text-right bottom-1 text-[0.8rem] text-muted-foreground select-none",
      scrollbar && "right-[1.375rem]",
      ($$props.maxlength && ((value?.length || 0) / $$props.maxlength) > 0.95) && "text-amber-600",
      ($$props.maxlength && (value?.length || 0) > $$props.maxlength) && "text-red-600",
    )}>
      <span>{value?.length || 0}</span>
      {#if $$props.maxlength}
        <span>/ {$$props.maxlength}</span>
      {/if}
    </span>
  </div>
{:else}
  <textarea
    class={cn(
      "flex min-h-[60px] w-full rounded-md border border-input bg-transparent px-3 py-2 text-sm shadow-sm placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50",
      className
    )}
    bind:value
    {readonly}
    on:blur
    on:change
    on:click
    on:focus
    on:keydown
    on:keypress
    on:keyup
    on:mouseover
    on:mouseenter
    on:mouseleave
    on:paste
    on:input
    {...$$restProps}
  ></textarea>
{/if}
