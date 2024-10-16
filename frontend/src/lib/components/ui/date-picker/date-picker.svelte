<script lang="ts">
  import { Button } from "$lib/components/ui/button/index.js";
  import * as Popover from "$lib/components/ui/popover";
  import { RangeCalendar } from "$lib/components/ui/range-calendar";
  import dayjs from "$lib/utils/datetime";
  import { fromAbsolute, fromDate, getLocalTimeZone, startOfMonth, startOfWeek, today } from "@internationalized/date";
  import type { DateRange } from "bits-ui";
  import CalendarRange from "lucide-svelte/icons/calendar-range";
  import Cross2 from "svelte-radix/Cross2.svelte";

  const timezone = getLocalTimeZone();

  let rawValue: DateRange | undefined;
  export let startValue: Date | undefined = undefined;
  export let endValue: Date | undefined = undefined;
  export let maxValue: Date | undefined = undefined;
  export let placeholder: string = "Pick a date";

  const f = (value?: DateRange) => {
    let sValue = value?.start?.toDate(timezone);
    sValue?.setHours(0, 0, 0, 0);
    if ((sValue?.getTime() || 0) <= 0) sValue = undefined;
    startValue = sValue;
    let eValue = value?.end?.toDate(timezone);
    eValue?.setHours(23, 59, 59, 999);
    endValue = eValue;
  };
  const fi = (sValue?: Date, eValue?: Date) => {
    if (!sValue && !eValue) rawValue = undefined;
    else if (!sValue && eValue)
      rawValue = {
        start: fromAbsolute(0, timezone),
        end: fromDate(eValue, timezone),
      };
    else
      rawValue = {
        start: sValue ? fromDate(sValue, timezone) : undefined,
        end: eValue ? fromDate(eValue, timezone) : undefined,
      };
  };
  $: fi(startValue, endValue);
  $: f(rawValue);

  let popoverOpen = false;
</script>

<div class={"grid gap-2 " + ($$props.class || "")}>
  <Popover.Root openFocus bind:open={popoverOpen}>
    <Popover.Trigger asChild let:builder>
      <Button variant="outline" builders={[builder]}>
        <CalendarRange class="mr-2 h-4 w-4 shrink-0 stroke-[1.65]" />
        {#if startValue && endValue}
          {dayjs(startValue).format("ll")} - {dayjs(endValue).format("ll")}
        {:else if !startValue && endValue}
          Up to {dayjs(endValue).format("ll")}
        {:else}
          {placeholder}
        {/if}
      </Button>
    </Popover.Trigger>
    <Popover.Content class="relative !z-30 flex w-auto items-center justify-center gap-7 p-0" align="start">
      <div class="mb-auto flex w-auto flex-col gap-3.5 pl-4 pt-5">
        <span
          class="cursor-pointer select-none text-xs font-semibold text-primary transition-colors hover:text-primary/80"
          on:click={() => {
            rawValue = { start: today(timezone), end: today(timezone) };
            popoverOpen = false;
          }}
        >
          Today
        </span>
        <span
          class="cursor-pointer select-none text-xs font-semibold text-primary transition-colors hover:text-primary/80"
          on:click={() => {
            rawValue = { start: today(timezone).subtract({ days: 1 }), end: today(timezone).subtract({ days: 1 }) };
            popoverOpen = false;
          }}
        >
          Yesterday
        </span>
        <span
          class="cursor-pointer select-none text-xs font-semibold text-primary transition-colors hover:text-primary/80"
          on:click={() => {
            rawValue = { start: today(timezone).subtract({ days: 6 }), end: today(timezone) };
            popoverOpen = false;
          }}
        >
          Last 7 days
        </span>
        <span
          class="cursor-pointer select-none text-xs font-semibold text-primary transition-colors hover:text-primary/80"
          on:click={() => {
            rawValue = { start: today(timezone).subtract({ days: 29 }), end: today(timezone) };
            popoverOpen = false;
          }}
        >
          Last 30 days
        </span>
        <span
          class="cursor-pointer select-none text-xs font-semibold text-primary transition-colors hover:text-primary/80"
          on:click={() => {
            rawValue = { start: startOfWeek(today(timezone), "es-ES"), end: today(timezone) };
            popoverOpen = false;
          }}
        >
          Current week
        </span>
        <span
          class="cursor-pointer select-none text-xs font-semibold text-primary transition-colors hover:text-primary/80"
          on:click={() => {
            rawValue = { start: startOfMonth(today(timezone)), end: today(timezone) };
            popoverOpen = false;
          }}
        >
          Current month
        </span>
        <span
          class="cursor-pointer select-none text-xs font-semibold text-primary transition-colors hover:text-primary/80"
          on:click={() => {
            rawValue = { start: fromAbsolute(0, timezone), end: today(timezone) };
            popoverOpen = false;
          }}
        >
          Up to today
        </span>
        <span
          class="cursor-pointer select-none text-xs font-semibold text-primary transition-colors hover:text-primary/80"
          on:click={() => {
            rawValue = { start: fromAbsolute(0, timezone), end: today(timezone).subtract({ days: 1 }) };
            popoverOpen = false;
          }}
        >
          Up to yesterday
        </span>
      </div>
      <span
        class="absolute bottom-4 left-4 flex cursor-pointer select-none items-center text-xs font-semibold text-secondary-foreground/75 transition-colors hover:text-secondary-foreground"
        on:click={async () => {
          rawValue = undefined;
          popoverOpen = false;
        }}
      >
        <Cross2 class="mr-1 h-3 w-3 shrink-0" />
        Clear
      </span>
      <RangeCalendar
        bind:value={rawValue}
        initialFocus
        numberOfMonths={2}
        weekStartsOn={1}
        maxValue={maxValue ? fromDate(maxValue, timezone) : undefined}
        class="relative after:absolute after:bottom-[5%] after:left-[-1%] after:h-[89%] after:w-[1px] after:bg-border"
      />
    </Popover.Content>
  </Popover.Root>
</div>
