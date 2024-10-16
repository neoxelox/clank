<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";
  import ChurnRateMetric from "$lib/components/cmp/metrics/ChurnRateMetric.svelte";
  import CustomerSatisfactionScoreMetric from "$lib/components/cmp/metrics/CustomerSatisfactionScoreMetric.svelte";
  import EmotionalSpectrumMetric from "$lib/components/cmp/metrics/EmotionalSpectrumMetric.svelte";
  import IssueRateMetric from "$lib/components/cmp/metrics/IssueRateMetric.svelte";
  import IssueTrendsMetric from "$lib/components/cmp/metrics/IssueTrendsMetric.svelte";
  import NegativePerceptionMetric from "$lib/components/cmp/metrics/NegativePerceptionMetric.svelte";
  import NetPromoterScoreMetric from "$lib/components/cmp/metrics/NetPromoterScoreMetric.svelte";
  import NeutralPerceptionMetric from "$lib/components/cmp/metrics/NeutralPerceptionMetric.svelte";
  import PositivePerceptionMetric from "$lib/components/cmp/metrics/PositivePerceptionMetric.svelte";
  import RetentionRateMetric from "$lib/components/cmp/metrics/RetentionRateMetric.svelte";
  import SuggestionRateMetric from "$lib/components/cmp/metrics/SuggestionRateMetric.svelte";
  import SuggestionTrendsMetric from "$lib/components/cmp/metrics/SuggestionTrendsMetric.svelte";
  import { Badge } from "$lib/components/ui/badge";
  import { Button } from "$lib/components/ui/button";
  import { DatePicker } from "$lib/components/ui/date-picker";
  import * as Tabs from "$lib/components/ui/tabs";
  import * as entities from "$lib/entities";
  import { organization, product } from "$lib/stores";
  import dayjs from "$lib/utils/datetime";
  import Download from "lucide-svelte/icons/download";
  import { onDestroy } from "svelte";

  let periodStartAt: Date | undefined;
  let periodEndAt: Date | undefined;

  if ($page.url.searchParams.has("period_start_at"))
    periodStartAt = dayjs($page.url.searchParams.get("period_start_at")).toDate();
  if ($page.url.searchParams.has("period_end_at"))
    periodEndAt = dayjs($page.url.searchParams.get("period_end_at")).toDate();

  const unsubscribe = product.subscribe(async (product) => {
    if (!product) return;

    setTimeout(async () => {
      periodStartAt = periodStartAt || undefined;
      periodEndAt = periodEndAt || new Date();
    }, 0);
  });

  $: (async () => {
    const params = new URLSearchParams();

    if (periodStartAt) params.set("period_start_at", dayjs(periodStartAt).toISOString());
    if (periodEndAt) params.set("period_end_at", dayjs(periodEndAt).toISOString());

    await goto(`?${params}`, { noScroll: true, keepFocus: true });
  })();

  onDestroy(() => {
    unsubscribe();
  });
</script>

<div class="h-full w-full">
  {#if $organization.plan === entities.OrganizationPlan.TRIAL || $organization.plan === entities.OrganizationPlan.DEMO}
    <div
      class="mb-4 flex w-full flex-row items-center justify-center gap-x-4 gap-y-2 rounded-lg border border-border bg-muted/40 px-3.5 py-2.5 shadow-sm"
    >
      {#if $organization.plan === entities.OrganizationPlan.TRIAL}
        <h3 class="text-wrap font-semibold leading-none tracking-tight text-muted-foreground">
          Unlock advanced UX & CX research qualitative metrics
        </h3>
      {:else if $organization.plan === entities.OrganizationPlan.DEMO}
        <h3 class="text-wrap font-semibold leading-none tracking-tight text-muted-foreground">
          While in demo mode new feedbacks won't be analyzed
        </h3>
      {/if}
      <Button href="/#pricing" variant="outline" size="sm" class="rounded-full px-4 py-2.5 text-sm font-semibold">
        Upgrade now &rarr;
      </Button>
    </div>
  {/if}
  <Tabs.Root class="w-full">
    <div class="mb-4 flex flex-wrap items-center justify-start gap-4">
      <Tabs.List class="w-auto">
        <Tabs.Trigger value="vitals">Vitals</Tabs.Trigger>
        <Tabs.Trigger value="insights" disabled>
          Insights
          <Badge variant="secondary" class="ml-1.5 mt-px px-0 opacity-50">Soon</Badge>
        </Tabs.Trigger>
      </Tabs.List>
      <div class="ml-auto flex items-center gap-4">
        <DatePicker
          placeholder="Entire history"
          maxValue={new Date()}
          bind:startValue={periodStartAt}
          bind:endValue={periodEndAt}
        />
        <Button size="sm" class="h-[2.125rem] font-semibold" disabled>
          <Download class="mr-2 h-4 w-4 shrink-0" />
          Export
        </Button>
      </div>
    </div>
    <Tabs.Content value="vitals">
      <div class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-4">
        <NetPromoterScoreMetric
          params={{ productID: $product.id, periodStartAt: periodStartAt, periodEndAt: periodEndAt }}
          class="h-[134px]"
        />
        <CustomerSatisfactionScoreMetric
          params={{ productID: $product.id, periodStartAt: periodStartAt, periodEndAt: periodEndAt }}
          class="h-[134px]"
        />
        <RetentionRateMetric
          params={{ productID: $product.id, periodStartAt: periodStartAt, periodEndAt: periodEndAt }}
          class="h-[134px]"
        />
        <ChurnRateMetric
          params={{ productID: $product.id, periodStartAt: periodStartAt, periodEndAt: periodEndAt }}
          class="h-[134px]"
        />
        <div class="col-span-full grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
          <IssueRateMetric
            params={{ productID: $product.id, periodStartAt: periodStartAt, periodEndAt: periodEndAt }}
            class="h-[134px]"
          />
          <SuggestionRateMetric
            params={{ productID: $product.id, periodStartAt: periodStartAt, periodEndAt: periodEndAt }}
            class="h-[134px]"
          />
          <EmotionalSpectrumMetric
            params={{ productID: $product.id, periodStartAt: periodStartAt, periodEndAt: periodEndAt }}
            class="row-span-2 h-[284px]"
          />
          <IssueTrendsMetric
            params={{ productID: $product.id, periodStartAt: periodStartAt, periodEndAt: periodEndAt }}
            class="h-[134px]"
          />
          <SuggestionTrendsMetric
            params={{ productID: $product.id, periodStartAt: periodStartAt, periodEndAt: periodEndAt }}
            class="h-[134px]"
          />
          <PositivePerceptionMetric
            params={{ productID: $product.id, periodStartAt: periodStartAt, periodEndAt: periodEndAt }}
            class="row-span-2 h-[284px]"
          />
          <NeutralPerceptionMetric
            params={{ productID: $product.id, periodStartAt: periodStartAt, periodEndAt: periodEndAt }}
            class="row-span-2 h-[284px]"
          />
          <NegativePerceptionMetric
            params={{ productID: $product.id, periodStartAt: periodStartAt, periodEndAt: periodEndAt }}
            class="row-span-2 h-[284px]"
          />
        </div>
      </div>
    </Tabs.Content>
  </Tabs.Root>
</div>
