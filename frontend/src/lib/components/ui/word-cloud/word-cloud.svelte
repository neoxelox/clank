<script lang="ts">
  import { cn } from "$lib/utils/ui.js";
  import type d3 from "d3";
  import { extent } from "d3-array";
  import d3Cloud from "d3-cloud";
  import { scaleLinear, scaleLog } from "d3-scale";
  import { onDestroy } from "svelte";

  export let width: number = 420;
  export let height: number = 300;
  export let padding: number = 4;
  export let fontFamily: string = "Helvetica";
  export let fontSize: number = 10;
  export let values: { label: string; value: number }[];
  let className: string | undefined = undefined;
  export { className as class };

  let data: { text: string; size: number }[];
  let words: d3.layout.cloud.Word[];
  let cloud: d3.layout.Cloud<d3.layout.cloud.Word>;
  let sizeScale: d3.scale.Log<number, number>;
  let opacityScale: d3.scale.Linear<number, number>;
  let minSize: number;
  let maxSize: number;

  $: (() => {
    if (cloud) cloud.stop();

    data = values.toSorted((a, b) => b.value - a.value).map(({ label, value }) => ({ text: label, size: value }));
    words = [];

    minSize = Math.min(...data.map(d => d.size));
    maxSize = Math.max(...data.map(d => d.size));

    sizeScale = scaleLog()
      .domain([minSize, maxSize])
      .range([fontSize * 0.8, fontSize * 3]);

    opacityScale = scaleLinear()
      .domain([fontSize * 0.8, fontSize * 3])
      .range([0.3, 1])
      .clamp(true);

    cloud = d3Cloud()
      .size([width, height])
      .words(data)
      .padding(padding)
      .rotate(0)
      .font(fontFamily)
      .fontSize(d => sizeScale(Math.max(d.size || 1, 1)))
      .random(() => 0)
      .on("word", ({ size, x, y, rotate, text }) => (words = words.concat([{ size, x, y, rotate, text }])))
      .on("end", () => (maxSize = extent(words, (d) => d.size)[1] || 0));

    cloud.start();
  })();

  onDestroy(() => {
    if (cloud) cloud.stop();
  });
</script>

<svg
  {width}
  {height}
  viewBox={`0 0 ${width} ${height}`}
  text-anchor="middle"
  font-family={fontFamily}
  class={cn("", className)}
>
  <g>
    {#each words as word}
      <text
        font-size={word.size}
        transform={`translate(${word.x}, ${word.y}) rotate(${word.rotate})`}
        opacity={opacityScale(word.size || 0)}
        class="cursor-default select-none fill-inherit font-semibold leading-none tracking-tight [word-spacing:-0.025em]"
      >
        {word.text}
      </text>
    {/each}
  </g>
</svg>
