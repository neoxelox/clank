<script lang="ts">
  import * as Card from "$lib/components/ui/card";
  import { Skeleton } from "$lib/components/ui/skeleton";
  import * as Tooltip from "$lib/components/ui/tooltip";
  import { WordCloud } from "$lib/components/ui/word-cloud";
  import * as entities from "$lib/entities";
  import * as payloads from "$lib/payloads";
  import { api } from "$lib/services/api";
  import dayjs from "$lib/utils/datetime";
  import { debounce } from "$lib/utils/delay";
  import { clone, compare } from "$lib/utils/object";
  import UserRoundX from "lucide-svelte/icons/user-round-x";
  import QuestionMarkCircled from "svelte-radix/QuestionMarkCircled.svelte";

  let previousParams: entities.ReviewKeywordsParams;
  export let params: entities.ReviewKeywordsParams;
  let load: Promise<void>;
  let previousMetric: entities.ReviewKeywordsMetric | undefined;
  let metric: entities.ReviewKeywordsMetric;

  let loadMetric = debounce((params: entities.ReviewKeywordsParams) => {
    load = (async () => {
      let query = payloads
        .toGetReviewKeywordsMetricQuery({
          period_start_at: params.periodStartAt,
          period_end_at: params.periodEndAt,
        })
        .toString();
      let response = await api.get<payloads.GetReviewKeywordsMetricResponse>(
        `/products/${params.productID}/metrics/review-keywords?${query}`,
      );
      metric = payloads.toReviewKeywordsMetric(response);

      if (!params.periodEndAt) {
        previousMetric = undefined;
        return;
      }

      if (params.periodStartAt)
        query = payloads
          .toGetReviewKeywordsMetricQuery({
            period_start_at: dayjs(params.periodStartAt)
              .subtract(dayjs(params.periodEndAt).diff(params.periodStartAt, "ms") + 1, "ms")
              .toDate(),
            period_end_at: dayjs(params.periodStartAt).subtract(1, "ms").toDate(),
          })
          .toString();
      else
        query = payloads
          .toGetReviewKeywordsMetricQuery({
            period_start_at: undefined,
            period_end_at: dayjs(params.periodEndAt).subtract(1, "day").toDate(),
          })
          .toString();

      response = await api.get<payloads.GetReviewKeywordsMetricResponse>(
        `/products/${params.productID}/metrics/review-keywords?${query}`,
      );
      previousMetric = payloads.toReviewKeywordsMetric(response);
    })();
  }, 50);

  $: (() => {
    if (compare(params, previousParams)) return;
    previousParams = clone(params);
    loadMetric(params);
  })();

  let data = (keywords: Record<string, number>) =>
    Object.entries(keywords).map(([keyword, value]) => ({
      label: keyword,
      value: value,
    }));

  let size = (keywords: Record<string, number>) => {
    const words = data(keywords).length;

    if (words <= 5) return 24;
    if (words <= 10) return 18;
    if (words <= 20) return 14;

    return 10;
  };

  let isEqual = (current: Record<string, number>, previous: Record<string, number>) => compare(current, previous);
</script>

<Card.Root class={"overflow-hidden " + ($$props.class || "")}>
  <Card.Header class="flex flex-row items-center justify-between space-y-0 pb-2">
    <Card.Title class="text-sm font-medium">
      <Tooltip.Root openDelay={333}>
        <Tooltip.Trigger class="flex cursor-default items-center">
          Negative Perception
          <QuestionMarkCircled class="ml-1 h-3.5 w-3.5 shrink-0 text-foreground/50" />
        </Tooltip.Trigger>
        <Tooltip.Content>
          <p class="max-w-[325px] text-wrap">
            Key terms from the point of view of unsatisfied customers, highlighting how they perceive your product,
            based on negative feedbacks posted in the selected period.
          </p>
        </Tooltip.Content>
      </Tooltip.Root>
    </Card.Title>
    <UserRoundX class="h-4 w-4 shrink-0 text-muted-foreground" />
  </Card.Header>
  <Card.Content class="flex h-[232px] w-full flex-col items-center justify-center px-6 pb-4 pt-1">
    {#await load}
      <div class="mt-1 flex h-full w-full flex-wrap items-start justify-center gap-4">
        <Skeleton class="h-10 w-1/2" />
        <Skeleton class="h-8 w-1/5" />
        <div class="flex w-2/12 flex-col items-end gap-3">
          <Skeleton class="h-5 w-full" />
          <Skeleton class="h-2 w-2/3" />
        </div>
        <Skeleton class="h-9 w-1/3" />
        <Skeleton class="h-7 w-3/5" />
        <Skeleton class="h-6 w-1/4" />
        <Skeleton class="h-6 w-1/4" />
        <Skeleton class="h-5 w-2/5" />
        <div class="flex w-2/4 flex-col gap-3">
          <Skeleton class="h-3 w-full" />
          <div class="flex w-full gap-2">
            <Skeleton class="h-2 w-1/3" />
            <Skeleton class="h-2 w-1/3" />
            <Skeleton class="h-2 w-1/3" />
          </div>
        </div>
        <Skeleton class="h-8 w-2/5" />
      </div>
    {:then}
      {#if !metric}
        <div class="mt-1 flex h-full w-full flex-wrap items-start justify-center gap-4">
          <Skeleton class="h-10 w-1/2" />
          <Skeleton class="h-8 w-1/5" />
          <div class="flex w-2/12 flex-col items-end gap-3">
            <Skeleton class="h-5 w-full" />
            <Skeleton class="h-2 w-2/3" />
          </div>
          <Skeleton class="h-9 w-1/3" />
          <Skeleton class="h-7 w-3/5" />
          <Skeleton class="h-6 w-1/4" />
          <Skeleton class="h-6 w-1/4" />
          <Skeleton class="h-5 w-2/5" />
          <div class="flex w-2/4 flex-col gap-3">
            <Skeleton class="h-3 w-full" />
            <div class="flex w-full gap-2">
              <Skeleton class="h-2 w-1/3" />
              <Skeleton class="h-2 w-1/3" />
              <Skeleton class="h-2 w-1/3" />
            </div>
          </div>
          <Skeleton class="h-8 w-2/5" />
        </div>
      {:else}
        {@const metricData = data(metric.negative)}
        <div class="group relative flex h-full w-full items-center justify-center">
          {#if metricData.length}
            <WordCloud
              values={metricData}
              width={329}
              height={212}
              fontFamily={'"Geist Sans", "Geist Sans fallback", ui-sans-serif, system-ui, -apple-system, "system-ui", "Helvetica"'}
              fontSize={size(metric.negative)}
              class={"shrink-0 fill-foreground opacity-100 transition-opacity " +
                (previousMetric && !isEqual(metric.negative, previousMetric.negative) ? "group-hover:opacity-25" : "")}
            />
          {:else}
            <div class="h-full w-full pb-2 pt-1">
              <div
                class="flex h-full w-full flex-1 items-center justify-center space-x-4 rounded-lg border border-dashed"
              >
                <p class="text-muted-foreground">No keywords</p>
              </div>
            </div>
          {/if}
          {#if previousMetric && !isEqual(metric.negative, previousMetric.negative)}
            {@const previousMetricData = data(previousMetric.negative)}
            <WordCloud
              values={previousMetricData}
              width={329}
              height={212}
              fontFamily={'"Geist Sans", "Geist Sans fallback", ui-sans-serif, system-ui, -apple-system, "system-ui", "Helvetica"'}
              fontSize={size(previousMetric.negative)}
              class="absolute left-1/2 top-1/2 shrink-0 -translate-x-1/2 -translate-y-1/2 fill-destructive opacity-0 transition-opacity group-hover:opacity-100"
            />
          {/if}
        </div>
      {/if}
    {:catch}
      <p class="-mt-3">Something went wrong!</p>
    {/await}
  </Card.Content>
</Card.Root>
