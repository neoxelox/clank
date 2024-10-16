<script lang="ts">
  import * as Card from "$lib/components/ui/card";
  import { Skeleton } from "$lib/components/ui/skeleton";
  import * as Tooltip from "$lib/components/ui/tooltip";
  import * as entities from "$lib/entities";
  import * as payloads from "$lib/payloads";
  import { api } from "$lib/services/api";
  import dayjs from "$lib/utils/datetime";
  import { debounce } from "$lib/utils/delay";
  import { clone, compare } from "$lib/utils/object";
  import Waypoints from "lucide-svelte/icons/waypoints";
  import ArrowDown from "svelte-radix/ArrowDown.svelte";
  import ArrowUp from "svelte-radix/ArrowUp.svelte";
  import Minus from "svelte-radix/Minus.svelte";
  import QuestionMarkCircled from "svelte-radix/QuestionMarkCircled.svelte";
  import Width from "svelte-radix/Width.svelte";

  let previousParams: entities.NetPromoterScoreParams;
  export let params: entities.NetPromoterScoreParams;
  let load: Promise<void>;
  let previousMetric: entities.NetPromoterScoreMetric | undefined;
  let metric: entities.NetPromoterScoreMetric;

  let loadMetric = debounce((params: entities.NetPromoterScoreParams) => {
    load = (async () => {
      let query = payloads
        .toGetNetPromoterScoreMetricQuery({
          period_start_at: params.periodStartAt,
          period_end_at: params.periodEndAt,
        })
        .toString();
      let response = await api.get<payloads.GetNetPromoterScoreMetricResponse>(
        `/products/${params.productID}/metrics/nps?${query}`,
      );
      metric = payloads.toNetPromoterScoreMetric(response);

      if (!params.periodEndAt) {
        previousMetric = undefined;
        return;
      }

      if (params.periodStartAt)
        query = payloads
          .toGetNetPromoterScoreMetricQuery({
            period_start_at: dayjs(params.periodStartAt)
              .subtract(dayjs(params.periodEndAt).diff(params.periodStartAt, "ms") + 1, "ms")
              .toDate(),
            period_end_at: dayjs(params.periodStartAt).subtract(1, "ms").toDate(),
          })
          .toString();
      else
        query = payloads
          .toGetNetPromoterScoreMetricQuery({
            period_start_at: undefined,
            period_end_at: dayjs(params.periodEndAt).subtract(1, "day").toDate(),
          })
          .toString();

      response = await api.get<payloads.GetNetPromoterScoreMetricResponse>(
        `/products/${params.productID}/metrics/nps?${query}`,
      );
      previousMetric = payloads.toNetPromoterScoreMetric(response);
    })();
  }, 50);

  $: (() => {
    if (compare(params, previousParams)) return;
    previousParams = clone(params);
    loadMetric(params);
  })();
</script>

<Card.Root class={"overflow-hidden " + ($$props.class || "")}>
  <Card.Header class="flex flex-row items-center justify-between space-y-0 pb-2">
    <Card.Title class="text-sm font-medium">
      <Tooltip.Root openDelay={333}>
        <Tooltip.Trigger class="flex cursor-default items-center">
          Net Promoter Score
          <QuestionMarkCircled class="ml-1 h-3.5 w-3.5 shrink-0 text-foreground/50" />
        </Tooltip.Trigger>
        <Tooltip.Content>
          <p class="max-w-[355px] text-wrap">
            How likely customers are to recommend your product, based on inferred customer intention of feedbacks posted
            in the selected period. It is normally presented as a score up to ±100 points.
          </p>
        </Tooltip.Content>
      </Tooltip.Root>
    </Card.Title>
    <Waypoints class="h-4 w-4 shrink-0 text-muted-foreground" />
  </Card.Header>
  <Card.Content>
    {#await load}
      <div class="mt-1 space-y-2">
        <Skeleton class="h-6 w-3/4" />
        <Skeleton class="h-5 w-1/2" />
      </div>
    {:then}
      {#if !metric}
        <div class="mt-1 space-y-2">
          <Skeleton class="h-6 w-3/4" />
          <Skeleton class="h-5 w-1/2" />
        </div>
      {:else}
        <div class="text-2xl font-bold">
          {metric.score}
        </div>
        {#if previousMetric}
          <p class="inline-flex items-center text-sm text-muted-foreground">
            {#if previousMetric.score > metric.score}
              <ArrowDown class="mr-0.5 h-3.5 w-3.5 shrink-0 text-destructive" />
              <span>
                <span class="text-destructive">
                  {parseFloat((previousMetric.score - metric.score).toFixed(2))}
                </span>
                from last period
              </span>
            {:else if previousMetric.score < metric.score}
              <ArrowUp class="text-constructive mr-0.5 h-3.5 w-3.5 shrink-0" />
              <span>
                <span class="text-constructive">
                  {parseFloat((metric.score - previousMetric.score).toFixed(2))}
                </span>
                from last period
              </span>
            {:else}
              <Minus class="mr-1 h-3.5 w-3.5 shrink-0" />
              <span>equal as last period</span>
            {/if}
          </p>
        {:else}
          <p class="inline-flex items-center text-sm text-muted-foreground">
            <Width class="mr-1 h-3.5 w-3.5 shrink-0" />
            <span>all time</span>
          </p>
        {/if}
      {/if}
    {:catch}
      <p>Something went wrong!</p>
    {/await}
  </Card.Content>
</Card.Root>
