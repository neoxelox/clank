<script lang="ts">
  import CallToActionBackground from "$assets/backgrounds/call-to-action.jpg";
  import Button from "$lib/components/cms/Button.svelte";
  import Container from "$lib/components/cms/Container.svelte";
  import { CreditCardIcon } from "@babeard/svelte-heroicons/outline";
</script>

<section id="call-to-action" class="relative overflow-hidden bg-primary py-24 sm:py-32">
  <img
    src={CallToActionBackground}
    width={2347}
    height={1244}
    fetchpriority="auto"
    loading="lazy"
    decoding="async"
    class="absolute left-1/2 top-1/2 max-w-none -translate-x-1/2 -translate-y-1/2"
  />
  <Container class="relative">
    <div class="mx-auto max-w-lg text-center">
      <h2 class="font-cal text-3xl text-primary-foreground sm:text-4xl">Up and running in 5 minutes.</h2>
      <p class="mt-4 text-pretty text-lg tracking-tight text-primary-foreground">
        Don't waste time reading thousands of reviews. Delight customers, reduce churn, and accelerate growth.
      </p>
      <Button href="/dash/signin" color="secondary" class="!bg-primary-foreground-gradient mt-10">
        Start 14-day free trial
      </Button>
      <p
        class="mx-auto mt-6 flex items-center justify-center text-sm font-normal leading-none tracking-tight text-primary-foreground"
      >
        <CreditCardIcon class="mr-2 h-5 w-5 shrink-0 stroke-[1.25]" />
        <span>No credit card required</span>
      </p>
    </div>
  </Container>
</section>
