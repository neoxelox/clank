<script lang="ts">
	import { cn, flyAndScale } from "$lib/utils/ui.js";
	import { Tooltip as TooltipPrimitive } from "bits-ui";

	type $$Props = TooltipPrimitive.ContentProps;

	let className: $$Props["class"] = undefined;
	export let sideOffset: $$Props["sideOffset"] = 4;
	export let transition: $$Props["transition"] = flyAndScale;
	export let transitionConfig: $$Props["transitionConfig"] = {
		y: 8,
		duration: 150,
	};
	export { className as class };
</script>

<TooltipPrimitive.Content
	{transition}
	{transitionConfig}
	{sideOffset}
	class={cn(
		"z-50 rounded-md px-3 py-1.5 text-xs border border-input bg-background text-center text-foreground shadow-sm",
		className
	)}
	{...$$restProps}
>
  <div class="bg-background">
    <TooltipPrimitive.Arrow class="rounded-[2px] border-l border-t border-input" />
  </div>
	<slot />
</TooltipPrimitive.Content>
