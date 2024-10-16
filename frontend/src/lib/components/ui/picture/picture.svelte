<script lang="ts">
	import * as Avatar from "$lib/components/ui/avatar";
	import { initials } from "$lib/utils/string";
	import { cn } from "$lib/utils/ui.js";

  export let size: "xs" | "sm" | "md" | "lg" | "xl" | "2xl" | "3xl" | "4xl" = "md";
  export let src: string | undefined = undefined;
  export let fallback: string | undefined = undefined;
  let className: string | undefined = undefined;
  export { className as class };

  // TODO: This doesn't make sense, backend should make pictures optional instead
  $: if (src?.startsWith("https://clank.so/images/pictures")) src = undefined;

  const rootClasses = cn(
    "shrink-0 select-none",
    size === "xs" && "h-3 w-3",
    size === "sm" && "h-4 w-4",
    size === "md" && "h-5 w-5",
    size === "lg" && "h-6 w-6",
    size === "xl" && "h-8 w-8",
    size === "2xl" && "h-9 w-9",
    size === "3xl" && "h-12 w-12",
    size === "4xl" && "h-40 w-40",
    className
  );

  const imageClasses = cn(
    "object-cover",
    size === "xs" && "",
    size === "sm" && "",
    size === "md" && "",
    size === "lg" && "",
    size === "xl" && "",
    size === "2xl" && "",
    size === "3xl" && "",
    size === "4xl" && "",
  );

  const fallbackClasses = cn(
    "bg-primary font-bold text-primary-foreground",
    size === "xs" && "text-[0.3rem]",
    size === "sm" && "text-[0.4rem]",
    size === "md" && "text-[0.55rem]",
    size === "lg" && "text-[0.625rem]",
    size === "xl" && "text-[0.75rem]",
    size === "2xl" && "text-[0.8rem]",
    size === "3xl" && "text-[1.25rem]",
    size === "4xl" && "text-[4.5rem]",
  );
</script>

<Avatar.Root class={rootClasses}>
  {#if src}
    <Avatar.Image src={src} class={imageClasses} />
  {/if}
  {#if fallback}
    <Avatar.Fallback class={fallbackClasses}>
      {initials(fallback)}
    </Avatar.Fallback>
  {/if}
</Avatar.Root>
