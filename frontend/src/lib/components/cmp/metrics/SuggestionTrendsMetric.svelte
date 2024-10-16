<script lang="ts">
  import * as Card from "$lib/components/ui/card";
  import * as Carousel from "$lib/components/ui/carousel";
  import { type CarouselAPI } from "$lib/components/ui/carousel/context";
  import { Skeleton } from "$lib/components/ui/skeleton";
  import * as Tooltip from "$lib/components/ui/tooltip";
  import * as entities from "$lib/entities";
  import * as payloads from "$lib/payloads";
  import { api } from "$lib/services/api";
  import dayjs from "$lib/utils/datetime";
  import { debounce } from "$lib/utils/delay";
  import { clone, compare } from "$lib/utils/object";
  import { titlelize } from "$lib/utils/string";
  import Autoplay from "embla-carousel-autoplay";
  import TrendingUp from "lucide-svelte/icons/trending-up";
  import { afterUpdate } from "svelte";
  import ArrowDown from "svelte-radix/ArrowDown.svelte";
  import Minus from "svelte-radix/Minus.svelte";
  import QuestionMarkCircled from "svelte-radix/QuestionMarkCircled.svelte";
  import Width from "svelte-radix/Width.svelte";

  let previousParams: entities.SuggestionTrendsParams;
  export let params: entities.SuggestionTrendsParams;
  let load: Promise<void>;
  let previousMetric: entities.SuggestionTrendsMetric | undefined;
  let metric: entities.SuggestionTrendsMetric;

  let loadSourceTrend = async (params: entities.SuggestionSourcesParams): Promise<[string?, string?]> => {
    let query = payloads
      .toGetSuggestionSourcesMetricQuery({
        period_start_at: params.periodStartAt,
        period_end_at: params.periodEndAt,
      })
      .toString();

    let response = await api.get<payloads.GetSuggestionSourcesMetricResponse>(
      `/products/${params.productID}/metrics/suggestion-sources?${query}`,
    );
    let trend = rank(payloads.toSuggestionSourcesMetric(response).sources)[0]?.label;
    if (trend) trend = entities.FeedbackSourceDetails[trend].title;

    if (!params.periodEndAt) return [trend, undefined];

    if (params.periodStartAt)
      query = payloads
        .toGetSuggestionSourcesMetricQuery({
          period_start_at: dayjs(params.periodStartAt)
            .subtract(dayjs(params.periodEndAt).diff(params.periodStartAt, "ms") + 1, "ms")
            .toDate(),
          period_end_at: dayjs(params.periodStartAt).subtract(1, "ms").toDate(),
        })
        .toString();
    else
      query = payloads
        .toGetSuggestionSourcesMetricQuery({
          period_start_at: undefined,
          period_end_at: dayjs(params.periodEndAt).subtract(1, "day").toDate(),
        })
        .toString();

    response = await api.get<payloads.GetSuggestionSourcesMetricResponse>(
      `/products/${params.productID}/metrics/suggestion-sources?${query}`,
    );
    let previousTrend = rank(payloads.toSuggestionSourcesMetric(response).sources)[0]?.label;
    if (previousTrend) previousTrend = entities.FeedbackSourceDetails[previousTrend].title;

    return [trend, previousTrend];
  };

  let loadImportanceTrend = async (params: entities.SuggestionImportancesParams): Promise<[string?, string?]> => {
    let query = payloads
      .toGetSuggestionImportancesMetricQuery({
        period_start_at: params.periodStartAt,
        period_end_at: params.periodEndAt,
      })
      .toString();

    let response = await api.get<payloads.GetSuggestionImportancesMetricResponse>(
      `/products/${params.productID}/metrics/suggestion-importances?${query}`,
    );
    let trend = rank(payloads.toSuggestionImportancesMetric(response).importances)[0]?.label;
    if (trend) trend = entities.SuggestionImportanceDetails[trend].title;

    if (!params.periodEndAt) return [trend, undefined];

    if (params.periodStartAt)
      query = payloads
        .toGetSuggestionImportancesMetricQuery({
          period_start_at: dayjs(params.periodStartAt)
            .subtract(dayjs(params.periodEndAt).diff(params.periodStartAt, "ms") + 1, "ms")
            .toDate(),
          period_end_at: dayjs(params.periodStartAt).subtract(1, "ms").toDate(),
        })
        .toString();
    else
      query = payloads
        .toGetSuggestionImportancesMetricQuery({
          period_start_at: undefined,
          period_end_at: dayjs(params.periodEndAt).subtract(1, "day").toDate(),
        })
        .toString();

    response = await api.get<payloads.GetSuggestionImportancesMetricResponse>(
      `/products/${params.productID}/metrics/suggestion-importances?${query}`,
    );
    let previousTrend = rank(payloads.toSuggestionImportancesMetric(response).importances)[0]?.label;
    if (previousTrend) previousTrend = entities.SuggestionImportanceDetails[previousTrend].title;

    return [trend, previousTrend];
  };

  let loadCategoryTrend = async (params: entities.SuggestionCategoriesParams): Promise<[string?, string?]> => {
    let query = payloads
      .toGetSuggestionCategoriesMetricQuery({
        period_start_at: params.periodStartAt,
        period_end_at: params.periodEndAt,
      })
      .toString();

    let response = await api.get<payloads.GetSuggestionCategoriesMetricResponse>(
      `/products/${params.productID}/metrics/suggestion-categories?${query}`,
    );
    let trend = rank(payloads.toSuggestionCategoriesMetric(response).categories)[0]?.label;
    if (trend) trend = titlelize(trend.replaceAll("_", " "));

    if (!params.periodEndAt) return [trend, undefined];

    if (params.periodStartAt)
      query = payloads
        .toGetSuggestionCategoriesMetricQuery({
          period_start_at: dayjs(params.periodStartAt)
            .subtract(dayjs(params.periodEndAt).diff(params.periodStartAt, "ms") + 1, "ms")
            .toDate(),
          period_end_at: dayjs(params.periodStartAt).subtract(1, "ms").toDate(),
        })
        .toString();
    else
      query = payloads
        .toGetSuggestionCategoriesMetricQuery({
          period_start_at: undefined,
          period_end_at: dayjs(params.periodEndAt).subtract(1, "day").toDate(),
        })
        .toString();

    response = await api.get<payloads.GetSuggestionCategoriesMetricResponse>(
      `/products/${params.productID}/metrics/suggestion-categories?${query}`,
    );
    let previousTrend = rank(payloads.toSuggestionCategoriesMetric(response).categories)[0]?.label;
    if (previousTrend) previousTrend = titlelize(previousTrend.replaceAll("_", " "));

    return [trend, previousTrend];
  };

  let loadReleaseTrend = async (params: entities.SuggestionReleasesParams): Promise<[string?, string?]> => {
    let query = payloads
      .toGetSuggestionReleasesMetricQuery({
        period_start_at: params.periodStartAt,
        period_end_at: params.periodEndAt,
      })
      .toString();

    let response = await api.get<payloads.GetSuggestionReleasesMetricResponse>(
      `/products/${params.productID}/metrics/suggestion-releases?${query}`,
    );
    let trend = rank(payloads.toSuggestionReleasesMetric(response).releases)[0]?.label;
    if (trend) trend = trend === entities.NO_RELEASE ? "Unknown" : trend;

    if (!params.periodEndAt) return [trend, undefined];

    if (params.periodStartAt)
      query = payloads
        .toGetSuggestionReleasesMetricQuery({
          period_start_at: dayjs(params.periodStartAt)
            .subtract(dayjs(params.periodEndAt).diff(params.periodStartAt, "ms") + 1, "ms")
            .toDate(),
          period_end_at: dayjs(params.periodStartAt).subtract(1, "ms").toDate(),
        })
        .toString();
    else
      query = payloads
        .toGetSuggestionReleasesMetricQuery({
          period_start_at: undefined,
          period_end_at: dayjs(params.periodEndAt).subtract(1, "day").toDate(),
        })
        .toString();

    response = await api.get<payloads.GetSuggestionReleasesMetricResponse>(
      `/products/${params.productID}/metrics/suggestion-releases?${query}`,
    );
    let previousTrend = rank(payloads.toSuggestionReleasesMetric(response).releases)[0]?.label;
    if (previousTrend) previousTrend = previousTrend === entities.NO_RELEASE ? "Unknown" : previousTrend;

    return [trend, previousTrend];
  };

  let loadMetric = debounce((params: entities.SuggestionTrendsParams) => {
    load = (async () => {
      const [
        [sourceTrend, previousSourceTrend],
        [importanceTrend, previousImportanceTrend],
        [categoryTrend, previousCategoryTrend],
        [releaseTrend, previousReleaseTrend],
      ] = await Promise.all([
        loadSourceTrend(params),
        loadImportanceTrend(params),
        loadCategoryTrend(params),
        loadReleaseTrend(params),
      ]);

      metric = {
        source: sourceTrend,
        importance: importanceTrend,
        category: categoryTrend,
        release: releaseTrend,
      };

      if (params.periodEndAt)
        previousMetric = {
          source: previousSourceTrend,
          importance: previousImportanceTrend,
          category: previousCategoryTrend,
          release: previousReleaseTrend,
        };
      else previousMetric = undefined;
    })();
  }, 50);

  $: (() => {
    if (compare(params, previousParams)) return;
    previousParams = clone(params);
    loadMetric(params);
  })();

  let rank = (options: Record<string, number>) =>
    Object.entries(options)
      .toSorted((a, b) => b[1] - a[1])
      .map(([label], index) => ({ label: label, rank: index + 1 }));

  const trends: string[] = ["category", "release", "importance", "source"];
  let trend: string = trends[0];
  let carousel: CarouselAPI;
  $: if (carousel) carousel.on("select", () => (trend = trends[carousel.selectedScrollSnap()]));

  afterUpdate(() => {
    if (carousel) trend = trends[carousel.selectedScrollSnap()];
  });
</script>

<Card.Root class={"overflow-hidden " + ($$props.class || "")}>
  <Card.Header class="flex flex-row items-center justify-between space-y-0 pb-2">
    <Card.Title class="text-sm font-medium">
      <Tooltip.Root openDelay={333}>
        <Tooltip.Trigger class="flex cursor-default items-center">
          Suggestion Trends
          <span class="ml-1 text-muted-foreground">/ {titlelize(trend)}</span>
          <QuestionMarkCircled class="ml-1 h-3.5 w-3.5 shrink-0 text-foreground/50" />
        </Tooltip.Trigger>
        <Tooltip.Content>
          <p class="max-w-[390px] text-wrap">
            Compilation of the trending category, release, importance and source of suggestions identified in feedbacks
            posted in the selected period.
          </p>
        </Tooltip.Content>
      </Tooltip.Root>
    </Card.Title>
    <TrendingUp class="h-4 w-4 shrink-0 text-muted-foreground" />
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
        <Carousel.Root
          opts={{ watchDrag: false, duration: 20, loop: true }}
          plugins={[Autoplay({ stopOnInteraction: false, stopOnMouseEnter: true, delay: 5000 })]}
          bind:api={carousel}
          class="w-full"
        >
          <Carousel.Next
            variant="outline"
            class="absolute bottom-0.5 right-0 top-[unset] z-10 flex h-5 w-7 translate-x-0 translate-y-0 rotate-0 items-center justify-center rounded-md p-0"
          />
          <Carousel.Content>
            {#each trends as trend}
              {@const metricTrend = metric[trend]}
              <Carousel.Item>
                <div class="overflow-hidden text-ellipsis text-nowrap text-2xl font-bold">
                  {metricTrend || "No trend"}
                </div>
                {#if previousMetric}
                  {@const previousMetricTrend = previousMetric[trend]}
                  <p class="inline-flex w-full items-center text-sm text-muted-foreground">
                    {#if previousMetricTrend !== metricTrend}
                      <span>from</span>
                      <ArrowDown class="ml-px mr-0.5 h-3.5 w-3.5 shrink-0 text-destructive" />
                      <span class="w-full overflow-hidden text-ellipsis text-nowrap">
                        <span class="text-destructive">
                          {previousMetricTrend || "No trend"}
                        </span>
                        in last period
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
              </Carousel.Item>
            {/each}
          </Carousel.Content>
        </Carousel.Root>
      {/if}
    {:catch}
      <p>Something went wrong!</p>
    {/await}
  </Card.Content>
</Card.Root>
