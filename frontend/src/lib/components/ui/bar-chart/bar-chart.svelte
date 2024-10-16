<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import { major } from "$lib/utils/object";
  import { cn } from "$lib/utils/ui.js";
  import { Accordion } from "bits-ui";
  import CaretSort from "svelte-radix/CaretSort.svelte";
  import { slide } from "svelte/transition";

  let className: string | undefined = undefined;
  let items: { label: string; value: number; percentage: number; step: number }[] = [];
  export let values: { label: string; value: number }[];
  export let open: boolean = false;
  export { className as class };

  let total = values.reduce((acc, { value }) => acc + value, 0);
  let acc = 0;
  let temp = values.reduce((accu, { label, value }) => ({ ...accu, [label]: value }), {});
  while (Object.keys(temp).length) {
    const key = major(temp);
    const item = values.find(({ label }) => label === key);
    const percentage = Math.round((item.value / total) * 100);
    acc += Object.keys(temp).length > 1 ? percentage : 0;
    items.push({
      ...item,
      percentage: Object.keys(temp).length > 1 ? percentage : 100 - acc,
      step: Object.keys(temp).length > 1 ? acc : 100,
    });
    delete temp[key];
  }
</script>

<Accordion.Root value={open ? "item" : undefined} class={cn("w-full", className)}>
  <Accordion.Item value="item">
    <Accordion.Trigger asChild let:builder>
      <Button variant="outline" builders={[builder]} class="group relative w-full">
        <span class="z-10 overflow-hidden text-ellipsis text-nowrap text-primary-foreground">{items[0].label}</span>
        <CaretSort
          class={"z-10 ml-auto h-4 w-4 shrink-0 " + (items.length > 4 ? "text-foreground/50" : "text-primary-foreground")}
        />
        {#each items.toReversed() as item, index}
          <div
            class={"absolute left-0 top-0 box-border flex h-full items-center justify-end rounded-md bg-primary text-xs text-primary-foreground transition-colors group-hover:bg-primary/90 " +
              (item.percentage >= 90 ? "pr-10" : "pr-4")}
            style={`width: ${item.step}%; opacity: ${1 - (items.length - index - 1) * 0.25};`}
          >
            {#if index === items.length - 1 && item.percentage >= 40}
              {item.percentage}% ({item.value})
            {/if}
          </div>
        {/each}
      </Button>
    </Accordion.Trigger>
    <Accordion.Content
      transition={slide}
      transitionConfig={{ duration: 200 }}
      class="flex flex-col items-center justify-start gap-2 px-1 pb-0 pt-2 text-sm"
    >
      {#each items as item, index}
        <div class="flex w-full items-center justify-between">
          <span class="flex items-center justify-center gap-2">
            <span
              class={"h-3 w-3 rounded-full " +
                (index < 4 ? "bg-primary" : "border border-primary/25 bg-background")}
              style={index < 4 ? `opacity: ${1 - index * 0.25};` : ""}
            ></span>
            <span>{item.label}</span>
          </span>
          <span>
            {items.length > 1 && index === items.length - 1 ? `~${item.percentage}` : item.percentage}% ({item.value})
          </span>
        </div>
      {/each}
    </Accordion.Content>
  </Accordion.Item>
</Accordion.Root>
