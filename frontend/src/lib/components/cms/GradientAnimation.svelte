<script lang="ts">
  import { cn } from "$lib/utils/ui";
</script>

<div class={cn("absolute left-0 top-0 -z-10 h-full w-full bg-transparent", $$props.class)}>
  <div
    class={cn(
      `absolute bg-gradient-to-tr from-[#ff80b5] to-[#9089fc] blur-3xl`,
      `-bottom-60 -left-96 h-2/5 w-2/5 mix-blend-hard-light`,
      `origin-center transform-gpu`,
      `animate-gradient-first`,
      `opacity-50`,
    )}
  ></div>
  <div
    class={cn(
      `absolute bg-gradient-to-tr from-[#ff80b5] to-[#9089fc] blur-3xl`,
      `left-[30%] top-[30%] h-2/5 w-2/5 mix-blend-hard-light`,
      `transform-gpu [transform-origin:calc(50%+400px)]`,
      `animate-gradient-third`,
      `opacity-30`,
    )}
  ></div>
  <div
    class={cn(
      `absolute bg-gradient-to-tr from-[#ff80b5] to-[#9089fc] blur-3xl`,
      `-right-96 -top-40 h-2/5 w-2/5 mix-blend-hard-light`,
      `origin-center transform-gpu`,
      `animate-gradient-fourth`,
      `opacity-50`,
    )}
  ></div>
  <div
    class={cn(
      `absolute bg-gradient-to-tr from-[#ff80b5] to-[#9089fc] blur-3xl`,
      `left-[30%] top-[30%] h-2/5 w-2/5 mix-blend-hard-light`,
      `transform-gpu [transform-origin:calc(50%-800px)_calc(50%+800px)]`,
      `animate-gradient-fifth`,
      `opacity-30`,
    )}
  ></div>
</div>
