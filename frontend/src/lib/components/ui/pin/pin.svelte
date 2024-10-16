<script lang="ts">
  import { cn } from "$lib/utils/ui.js";
  import { PinInput } from "bits-ui";
  import { createEventDispatcher, tick } from "svelte";

  const dispatch = createEventDispatcher();

  let rawValue: string[] | undefined = undefined;

  let className: string | undefined = undefined;
	export let value: string | undefined = undefined;
  export let size: number = 4;
  export let type: "text" | "password" = "text";
  export let uppercased: boolean = true;
  export let required: boolean = false;
  export let disabled: boolean = false;
	export { className as class };

  const f = (v?: string[]) => value = v?.join("");
  const fi = (v?: string) => rawValue = v?.split("");
  $: fi(value);
  $: f(rawValue);

  let dispatched: boolean = false;
  let finish = async () => {
    if (dispatched) return;
    await tick()
    dispatch("finish");
    dispatched = true;
  };

  $: (value?.length || 0) === size ? finish() : (dispatched = false);
</script>

<PinInput.Root
  bind:value={rawValue}
  class={cn(
    "flex items-center gap-2",
    className
  )}
  type={type}
  placeholder=""
>
  {#each {length: size} as _}
    <PinInput.Input
      class="flex h-12 w-10 select-none rounded-md border border-input bg-transparent p-1 text-center text-lg shadow-sm transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
      on:input={(event) => uppercased && (event.detail.currentTarget.value = event.detail.currentTarget.value.toUpperCase())}
      required={required}
      disabled={disabled}
    />
  {/each}
  <PinInput.HiddenInput />
</PinInput.Root>
