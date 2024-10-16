<script lang="ts">
  import { cn } from "$lib/utils/ui.js";

  let className: string | undefined = undefined;
  export let size: "sm" | "default" = "default"
  export let steps: string[];
  export { className as class };
</script>

<ul
  class={cn(
    "relative grid gap-8 pl-6 after:absolute after:inset-y-0 after:left-0 after:w-[3px] after:bg-border",
    className,
  )}
>
  {#each steps as step, index}
    <li class="relative grid gap-2">
      <div
        class={cn("absolute left-0 top-[0.125rem] z-10 flex h-5 w-5 translate-x-[-2rem] cursor-default items-center justify-center rounded-full border bg-background text-center align-middle text-xs font-semibold text-foreground outline outline-[5px] outline-background",
          size === "sm" && "top-0"
        )}
      >
        {index + 1}
      </div>
      <div class={cn(size === "sm" && "text-sm")}>{step}</div>
    </li>
  {/each}
</ul>
