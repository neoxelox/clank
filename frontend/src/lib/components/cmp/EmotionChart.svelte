<script lang="ts">
  import * as entities from "$lib/entities";
  import { scaleBand } from "d3-scale";
  import { curveCatmullRomClosed } from "d3-shape";
  import { Axis, Chart, Group, Spline, Svg, Text } from "layerchart";

  export let emotions: Record<string, number>;

  let data = Object.values(entities.ReviewEmotion).map((emotion) => {
    const details = entities.ReviewEmotionDetails[emotion];
    return {
      label: details.title,
      value: 0,
    };
  });

  Object.entries(emotions).forEach(([emotion, value]) => {
    const details = entities.ReviewEmotionDetails[emotion];
    data[data.findIndex(({ label }) => label === details.title)].value += value;
  });
</script>

<div class={"aspect-video flex h-full w-full items-center justify-center " + ($$props.class || "")}>
  <Chart
    {data}
    x="label"
    xScale={scaleBand()}
    xDomain={data.map((v) => v.label)}
    xRange={[0, 2 * Math.PI]}
    y="value"
    yRange={({ height }) => [20, height / 2.5]}
    yPadding={[0, 10]}
    padding={{ top: 0, bottom: 0, left: 0, right: 0 }}
  >
    <Svg>
      <Group center>
        <Axis placement="radius" grid={{ class: "stroke-border fill-background/20" }} ticks={3} format={() => ""} />
        <Axis
          placement="angle"
          grid={{ class: "stroke-border" }}
          format={(label) => {
            if (label.toUpperCase() == entities.ReviewEmotion.ANNOYANCE) return "Anyce.";
            if (label.toUpperCase() == entities.ReviewEmotion.APPREHENSION) return "Aprhns.";
            return label;
          }}
        >
          <svelte:fragment slot="tickLabel" let:labelProps>
            <Text {...labelProps} scaleToFit class="cursor-default select-none fill-muted-foreground text-xs" />
          </svelte:fragment>
        </Axis>
        <Spline radial curve={curveCatmullRomClosed} class="fill-primary/20 stroke-primary" />
      </Group>
    </Svg>
  </Chart>
</div>
