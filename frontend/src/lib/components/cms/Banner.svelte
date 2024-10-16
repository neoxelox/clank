<script lang="ts">
  import Button from "$lib/components/cms/Button.svelte";
  import { ArrowRightCircleIcon } from "@babeard/svelte-heroicons/mini";
  import { XMarkIcon } from "@babeard/svelte-heroicons/solid";

  export let title: string;
  export let description: string;

  let open: boolean = true;
</script>

{#if open}
  <div
    class="relative isolate z-50 flex items-center gap-x-6 overflow-hidden bg-[#f9fafb] px-6 py-2.5 sm:px-3.5 sm:before:flex-1"
  >
    <div
      class="absolute left-[max(-7rem,calc(50%-52rem))] top-1/2 -z-10 -translate-y-1/2 transform-gpu blur-2xl"
      aria-hidden="true"
    >
      <div
        class="aspect-[577/310] w-[36.0625rem] bg-gradient-to-r from-[#ff80b5] to-[#9089fc] opacity-30"
        style="clipPath: polygon(74.8% 41.9%, 97.2% 73.2%, 100% 34.9%, 92.5% 0.4%, 87.5% 0%, 75% 28.6%, 58.5% 54.6%, 50.1% 56.8%, 46.9% 44%, 48.3% 17.4%, 24.7% 53.9%, 0% 27.9%, 11.9% 74.2%, 24.9% 54.1%, 68.6% 100%, 74.8% 41.9%)"
      />
    </div>
    <div
      class="absolute left-[max(45rem,calc(50%+8rem))] top-1/2 -z-10 -translate-y-1/2 transform-gpu blur-2xl"
      aria-hidden="true"
    >
      <div
        class="aspect-[577/310] w-[36.0625rem] bg-gradient-to-r from-[#ff80b5] to-[#9089fc] opacity-30"
        style="clipPath: polygon(74.8% 41.9%, 97.2% 73.2%, 100% 34.9%, 92.5% 0.4%, 87.5% 0%, 75% 28.6%, 58.5% 54.6%, 50.1% 56.8%, 46.9% 44%, 48.3% 17.4%, 24.7% 53.9%, 0% 27.9%, 11.9% 74.2%, 24.9% 54.1%, 68.6% 100%, 74.8% 41.9%)"
      />
    </div>
    <div class="flex flex-row items-center gap-x-4 gap-y-2 md:flex-wrap">
      <p class="text-wrap text-sm leading-6 text-foreground">
        <strong class="font-semibold">{title}</strong>
        <span class="mx-1 text-base font-bold">·</span>
        {description}
      </p>
      <Button href="/dash/signin" color="secondary" class="text-nowrap !px-3.5 !py-1 shadow-sm">
        Register now
        <ArrowRightCircleIcon class="ml-1.5 h-4 w-4 flex-shrink-0 fill-secondary-foreground" />
      </Button>
    </div>
    <div class="flex flex-1 justify-end">
      <button type="button" class="-m-3 p-3 focus-visible:outline-offset-[-4px]" on:click={() => (open = false)}>
        <XMarkIcon class="h-5 w-5 text-foreground" />
      </button>
    </div>
  </div>
{/if}
