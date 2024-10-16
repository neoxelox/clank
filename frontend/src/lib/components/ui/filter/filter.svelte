<script lang="ts">
  import { Badge } from "$lib/components/ui/badge";
  import { Button } from "$lib/components/ui/button";
  import * as Command from "$lib/components/ui/command";
  import { Picture } from "$lib/components/ui/picture";
  import * as Popover from "$lib/components/ui/popover";
  import { Separator } from "$lib/components/ui/separator";
  import { cn } from "$lib/utils/ui.js";
  import type { SvelteComponent } from "svelte";
  import Check from "svelte-radix/Check.svelte";
  import PlusCircled from "svelte-radix/PlusCircled.svelte";

  export let title: string;
  export let options: { value: string; label: string; icon?: string | typeof SvelteComponent }[] = [];
  export let counts: { [index: string]: number } = {};
  export let addable: boolean = false;
  export let values: string[] = [];

  values?.forEach((value) => {
    if (!options.find((option) => option.value === value)) options.push({ value: value, label: value });
  });

  let open: boolean = false;
  let newValue: string;

  function handleSelect(currentValue: string) {
    if (Array.isArray(values) && values.includes(currentValue)) {
      values = values.filter((v) => v !== currentValue);
    } else {
      values = [...(Array.isArray(values) ? values : []), currentValue];
    }
  }
</script>

<Popover.Root bind:open>
  <Popover.Trigger asChild let:builder>
    <Button builders={[builder]} variant="outline" class="h-9 border-dashed focus-visible:ring-2">
      <PlusCircled class="mr-2 h-4 w-4" />
      {title}
      {#if values?.length > 0}
        <Separator orientation="vertical" class="mx-2 h-4" />
        <div class="flex space-x-1">
          {#if values.length > 2}
            <Badge variant="secondary" class="rounded-sm px-1 font-normal">
              {values.length} selected
            </Badge>
          {:else}
            {#each values as value}
              <Badge variant="secondary" class="rounded-sm px-1 font-normal">
                {options.find((option) => option.value === value)?.label || value}
              </Badge>
            {/each}
          {/if}
        </div>
      {/if}
    </Button>
  </Popover.Trigger>
  <Popover.Content class="w-[200px] p-0" align="start" side="bottom">
    <Command.Root>
      <Command.Input
        placeholder={addable
          ? `Search or add ${title.toLowerCase().slice(0, 3)}...`
          : `Search ${title.toLowerCase()}...`}
        bind:value={newValue}
      />
      <Command.List>
        <Command.Empty class={addable && newValue ? "pb-0" : ""}>
          No {title.toLowerCase()} found.
          {#if addable && newValue}
            <div class="-mx-1 mt-6 h-px bg-border"></div>
            <button
              class="w-full p-1"
              on:click={() => {
                options.push({ value: newValue, label: newValue });
                handleSelect(newValue);
              }}
            >
              <div
                class="flex cursor-pointer select-none items-center overflow-hidden text-nowrap text-foreground rounded-sm px-2 py-1.5 text-sm outline-none hover:bg-accent hover:text-accent-foreground"
              >
                <PlusCircled class="mr-2 h-5 w-5 shrink-0" />
                Add to search
              </div>
            </button>
          {/if}
        </Command.Empty>
        <Command.Group>
          {#each options as option}
            <Command.Item value={option.label} onSelect={() => handleSelect(option.value)}>
              <div
                class={cn(
                  "mr-2 flex h-4 w-4 items-center justify-center rounded-sm border border-primary",
                  values?.includes(option.value)
                    ? "bg-primary text-primary-foreground"
                    : "opacity-50 [&_svg]:invisible",
                )}
              >
                <Check className={cn("h-4 w-4")} />
              </div>
              {#if option.icon}
                {#if typeof option.icon === "string"}
                  <Picture size="xs" src={option.icon} fallback={option.label} class="mr-2 mt-[0.05rem]" />
                {:else}
                  <svelte:component this={option.icon} class="mr-2 h-3 w-3 text-muted-foreground shrink-0" />
                {/if}
              {/if}
              <span class="overflow-hidden text-ellipsis text-nowrap">
                {option.label}
              </span>
              {#if counts[option.value]}
                <span class="ml-auto flex h-4 w-4 items-center justify-center font-mono text-xs">
                  {counts[option.value]}
                </span>
              {/if}
            </Command.Item>
          {/each}
        </Command.Group>
        {#if values?.length > 0}
          <Command.Separator />
          <Command.Item
            class="justify-center text-center"
            onSelect={() => {
              values = [];
            }}
          >
            Clear filters
          </Command.Item>
        {/if}
      </Command.List>
    </Command.Root>
  </Popover.Content>
</Popover.Root>
