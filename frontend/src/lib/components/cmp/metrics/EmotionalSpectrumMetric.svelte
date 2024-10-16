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
  import { scaleBand } from "d3-scale";
  import { curveCatmullRomClosed } from "d3-shape";
  import { Axis, Chart, Group, Spline, Svg, Text } from "layerchart";
  import Activity from "lucide-svelte/icons/activity";
  import QuestionMarkCircled from "svelte-radix/QuestionMarkCircled.svelte";

  let previousParams: entities.ReviewEmotionsParams;
  export let params: entities.ReviewEmotionsParams;
  let load: Promise<void>;
  let previousMetric: entities.ReviewEmotionsMetric | undefined;
  let metric: entities.ReviewEmotionsMetric;

  let loadMetric = debounce((params: entities.ReviewEmotionsParams) => {
    load = (async () => {
      let query = payloads
        .toGetReviewEmotionsMetricQuery({
          period_start_at: params.periodStartAt,
          period_end_at: params.periodEndAt,
        })
        .toString();
      let response = await api.get<payloads.GetReviewEmotionsMetricResponse>(
        `/products/${params.productID}/metrics/review-emotions?${query}`,
      );
      metric = payloads.toReviewEmotionsMetric(response);

      if (!params.periodEndAt) {
        previousMetric = undefined;
        return;
      }

      if (params.periodStartAt)
        query = payloads
          .toGetReviewEmotionsMetricQuery({
            period_start_at: dayjs(params.periodStartAt)
              .subtract(dayjs(params.periodEndAt).diff(params.periodStartAt, "ms") + 1, "ms")
              .toDate(),
            period_end_at: dayjs(params.periodStartAt).subtract(1, "ms").toDate(),
          })
          .toString();
      else
        query = payloads
          .toGetReviewEmotionsMetricQuery({
            period_start_at: undefined,
            period_end_at: dayjs(params.periodEndAt).subtract(1, "day").toDate(),
          })
          .toString();

      response = await api.get<payloads.GetReviewEmotionsMetricResponse>(
        `/products/${params.productID}/metrics/review-emotions?${query}`,
      );
      previousMetric = payloads.toReviewEmotionsMetric(response);
    })();
  }, 50);

  $: (() => {
    if (compare(params, previousParams)) return;
    previousParams = clone(params);
    loadMetric(params);
  })();

  let data = (emotions: Record<string, number>) =>
    Object.values(entities.ReviewEmotion).map((emotion) => {
      const details = entities.ReviewEmotionDetails[emotion];
      return {
        label: details.title,
        value: emotions[emotion] || 0,
      };
    });

  let isFlat = (emotions: Record<string, number>) => {
    const augmented = data(emotions);

    let last = augmented[0].value;
    for (let emotion = 0; emotion < augmented.length; emotion++) {
      if (last !== augmented[emotion].value) return false;
      last = augmented[emotion].value;
    }

    return true;
  };

  let isEqual = (current: Record<string, number>, previous: Record<string, number>) => {
    for (const emotion in entities.ReviewEmotion) {
      if ((current[emotion] || 0) !== (previous[emotion] || 0)) return false;
    }

    return true;
  };
</script>

<Card.Root class={"overflow-hidden " + ($$props.class || "")}>
  <Card.Header class="flex flex-row items-center justify-between space-y-0 pb-2">
    <Card.Title class="text-sm font-medium">
      <Tooltip.Root openDelay={333}>
        <Tooltip.Trigger class="relative z-10 flex cursor-default items-center">
          Emotional Spectrum
          <QuestionMarkCircled class="ml-1 h-3.5 w-3.5 shrink-0 text-foreground/50" />
        </Tooltip.Trigger>
        <Tooltip.Content>
          <p class="max-w-[355px] text-wrap">
            The range of emotions that your product produces in customers, based on inferred feelings of feedbacks
            posted in the selected period. Analysis is based on the <a
              href="https://en.wikipedia.org/wiki/Emotion_classification#Plutchik's_wheel_of_emotions"
              target="_blank"
              rel="noopener noreferrer"
              class="text-primary">Plutchik wheel of emotions</a
            >.
          </p>
        </Tooltip.Content>
      </Tooltip.Root>
    </Card.Title>
    <Activity class="h-4 w-4 shrink-0 text-muted-foreground" />
  </Card.Header>
  <Card.Content class="-mt-7 flex h-full w-full flex-col items-center justify-center">
    {#await load}
      <div class="aspect-video relative -ml-3 mt-3 h-full w-full">
        <Skeleton class="absolute left-1/2 top-1/2 z-[15] h-10 w-10 -translate-x-1/2 -translate-y-1/2 rounded-full" />
        <div
          class="absolute left-1/2 top-1/2 z-[14] h-[4.5rem] w-[4.5rem] -translate-x-1/2 -translate-y-1/2 rounded-full bg-card"
        ></div>
        <Skeleton class="absolute left-1/2 top-1/2 z-[13] h-24 w-24 -translate-x-1/2 -translate-y-1/2 rounded-full" />
        <div
          class="absolute left-1/2 top-1/2 z-[12] h-32 w-32 -translate-x-1/2 -translate-y-1/2 rounded-full bg-card"
        ></div>
        <Skeleton class="absolute left-1/2 top-1/2 z-[11] h-40 w-40 -translate-x-1/2 -translate-y-1/2 rounded-full" />
      </div>
    {:then}
      {#if !metric}
        <div class="aspect-video relative -ml-3 mt-3 h-full w-full">
          <Skeleton class="absolute left-1/2 top-1/2 z-[15] h-10 w-10 -translate-x-1/2 -translate-y-1/2 rounded-full" />
          <div
            class="absolute left-1/2 top-1/2 z-[14] h-[4.5rem] w-[4.5rem] -translate-x-1/2 -translate-y-1/2 rounded-full bg-card"
          ></div>
          <Skeleton class="absolute left-1/2 top-1/2 z-[13] h-24 w-24 -translate-x-1/2 -translate-y-1/2 rounded-full" />
          <div
            class="absolute left-1/2 top-1/2 z-[12] h-32 w-32 -translate-x-1/2 -translate-y-1/2 rounded-full bg-card"
          ></div>
          <Skeleton class="absolute left-1/2 top-1/2 z-[11] h-40 w-40 -translate-x-1/2 -translate-y-1/2 rounded-full" />
        </div>
      {:else}
        {@const metricData = data(metric.emotions)}
        <div class="aspect-video relative -ml-3 mt-3 flex h-full w-full items-center justify-center">
          <Chart
            data={!isFlat(metric.emotions)
              ? metricData
              : metricData.map((point) => {
                  // Data needs to be at least 1n different for the X axis to draw
                  if (point.label.toUpperCase() == entities.ReviewEmotion.SERENITY) point.value += 1;
                  return point;
                })}
            x="label"
            xScale={scaleBand()}
            xDomain={metricData.map((v) => v.label)}
            xRange={[0, 2 * Math.PI]}
            y="value"
            yRange={({ height }) => [20, height / 2.6]}
            yPadding={[0, 16]}
            padding={{ top: 0, bottom: 0, left: 0, right: 0 }}
          >
            <Svg>
              <Group center>
                <Axis
                  placement="radius"
                  grid={{ class: "stroke-border fill-background/20" }}
                  ticks={(scale) => {
                    const domain = scale.domain();
                    return [domain[0] * 0.8, (domain[1] / 2) * 0.8, domain[1] * 0.8];
                  }}
                  format={() => ""}
                />
                <Axis placement="angle" grid={{ class: "stroke-border" }}>
                  <svelte:fragment slot="tickLabel" let:labelProps>
                    <Text {...labelProps} scaleToFit class="cursor-default select-none fill-muted-foreground text-xs" />
                  </svelte:fragment>
                </Axis>
                {#if !isFlat(metric.emotions)}
                  <Spline radial curve={curveCatmullRomClosed} class="fill-primary/20 stroke-primary" />
                {/if}
              </Group>
            </Svg>
          </Chart>
          {#if isFlat(metric.emotions)}
            <div
              class="absolute left-1/2 top-1/2 z-10 h-10 w-10 -translate-x-1/2 -translate-y-1/2 rounded-full border border-primary bg-primary/20"
            ></div>
          {/if}
          {#if previousMetric && !isEqual(metric.emotions, previousMetric.emotions)}
            {@const previousMetricData = data(previousMetric.emotions)}
            <div
              class="aspect-video absolute left-1/2 top-1/2 flex h-full w-full -translate-x-1/2 -translate-y-1/2 items-center justify-center opacity-25"
            >
              <Chart
                data={!isFlat(previousMetric.emotions)
                  ? previousMetricData
                  : previousMetricData.map((point) => {
                      // Data needs to be at least 1n different for the X axis to draw
                      if (point.label.toUpperCase() == entities.ReviewEmotion.SERENITY) point.value += 1;
                      return point;
                    })}
                x="label"
                xScale={scaleBand()}
                xDomain={previousMetricData.map((v) => v.label)}
                xRange={[0, 2 * Math.PI]}
                y="value"
                yRange={({ height }) => [20, height / 2.6]}
                yPadding={[0, 16]}
                padding={{ top: 0, bottom: 0, left: 0, right: 0 }}
              >
                <Svg>
                  <Group center>
                    {#if !isFlat(previousMetric.emotions)}
                      <Spline radial curve={curveCatmullRomClosed} class="fill-destructive/20 stroke-destructive" />
                    {/if}
                  </Group>
                </Svg>
              </Chart>
              {#if isFlat(previousMetric.emotions)}
                <div
                  class="absolute left-1/2 top-1/2 z-10 h-10 w-10 -translate-x-1/2 -translate-y-1/2 rounded-full border border-destructive bg-destructive/20"
                ></div>
              {/if}
            </div>
          {/if}
        </div>
      {/if}
    {:catch}
      <p>Something went wrong!</p>
    {/await}
  </Card.Content>
</Card.Root>
