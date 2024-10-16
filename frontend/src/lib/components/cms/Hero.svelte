<script lang="ts">
  import Button from "$lib/components/cms/Button.svelte";
  import Container from "$lib/components/cms/Container.svelte";
  import GradientAnimation from "$lib/components/cms/GradientAnimation.svelte";
  import Underline from "$lib/components/icons/Underline.svelte";
  import { PlayIcon } from "@babeard/svelte-heroicons/mini";
  import { CheckCircleIcon, CreditCardIcon } from "@babeard/svelte-heroicons/outline";
</script>

<Container class="relative z-40 pb-16 pt-20 text-center lg:pt-32">
  <GradientAnimation />
  <h1
    class="md:!leading-xl mx-auto max-w-4xl font-cal text-5xl font-bold leading-[100%] tracking-[-0.002em] text-foreground sm:text-7xl"
  >
    Customer feedback
    <span class="relative whitespace-nowrap">
      <Underline />
      <span class="text-primary-gradient relative">analysis</span>
    </span>
    on autopilot.
  </h1>
  <p class="mx-auto mt-12 max-w-2xl text-pretty text-lg tracking-tight text-foreground md:max-w-4xl">
    Leverage Clank's advanced AI engine to automatically collect and analyze customer feedback at scale.
  </p>
  <div class="mx-auto mt-7 flex max-w-lg justify-center gap-x-3 md:max-w-3xl md:gap-x-6">
    <ul
      role="list"
      class="space-y-3 text-pretty text-left text-base leading-6 tracking-tight text-foreground md:text-lg"
    >
      <li class="flex gap-x-3 md:items-center">
        <CheckCircleIcon class="h-6 w-6 flex-none stroke-2 text-primary" />
        Identify common issues and suggestions
      </li>
      <li class="flex gap-x-3 md:items-center">
        <CheckCircleIcon class="h-6 w-6 flex-none stroke-2 text-primary" />
        Predict churn and retention
      </li>
    </ul>
    <ul
      role="list"
      class="space-y-3 text-pretty text-left text-base leading-6 tracking-tight text-foreground md:text-lg"
    >
      <li class="flex gap-x-3 md:items-center">
        <CheckCircleIcon class="h-6 w-6 flex-none stroke-2 text-primary" />
        Track NPS, CSAT and other CX KPIs
      </li>
      <li class="flex gap-x-3 md:items-center">
        <CheckCircleIcon class="h-6 w-6 flex-none stroke-2 text-primary" />
        Interact with synthetic customer personas
      </li>
    </ul>
  </div>
  <div class="mt-10 flex justify-center gap-x-6">
    <Button href="/dash/signin" color="primary">Try Clank for free</Button>
    <Button href="https://www.youtube.com/@clank-so" target="_blank" rel="noopener noreferrer" color="secondary">
      <PlayIcon class="h-3 w-3 flex-none fill-secondary-foreground" />
      <span class="ml-2">Watch demo</span>
    </Button>
  </div>
  <p
    class="mx-auto mt-6 flex items-center justify-center text-sm font-normal leading-none tracking-tight text-muted-foreground"
  >
    <CreditCardIcon class="mr-2 h-5 w-5 shrink-0 stroke-[1.25]" />
    <span>No credit card required</span>
  </p>
</Container>
