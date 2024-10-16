<script lang="ts">
  import Container from "$lib/components/cms/Container.svelte";
  import GradientAnimation from "$lib/components/cms/GradientAnimation.svelte";
  import Amazon from "$lib/components/icons/Amazon.svelte";
  import AppStore from "$lib/components/icons/AppStore.svelte";
  import G2 from "$lib/components/icons/G2.svelte";
  import Gmail from "$lib/components/icons/Gmail.svelte";
  import GoogleMyBusiness from "$lib/components/icons/GoogleMyBusiness.svelte";
  import IAgora from "$lib/components/icons/IAgora.svelte";
  import Intercom from "$lib/components/icons/Intercom.svelte";
  import PlayStore from "$lib/components/icons/PlayStore.svelte";
  import Tripadvisor from "$lib/components/icons/Tripadvisor.svelte";
  import Trustpilot from "$lib/components/icons/Trustpilot.svelte";
  import Typeform from "$lib/components/icons/Typeform.svelte";
  import Underline from "$lib/components/icons/Underline.svelte";
  import Webhook from "$lib/components/icons/Webhook.svelte";
  import Widget from "$lib/components/icons/Widget.svelte";
  import Yelp from "$lib/components/icons/Yelp.svelte";
  import Zapier from "$lib/components/icons/Zapier.svelte";
  import Zendesk from "$lib/components/icons/Zendesk.svelte";
  import { ArrowRightCircleIcon } from "@babeard/svelte-heroicons/mini";
  import { SquaresPlusIcon } from "@babeard/svelte-heroicons/outline";

  const integrations = [
    {
      name: "Trustpilot",
      description: "Collect reviews from your company page on Trustpilot.",
      category: "Websites",
      icon: Trustpilot,
      available: true,
    },
    {
      name: "Play Store",
      description: "Collect reviews from your application page on Play Store (all storefronts).",
      category: "Apps",
      icon: PlayStore,
      available: true,
    },
    {
      name: "App Store",
      description: "Collect reviews from your application page on App Store (all storefronts).",
      category: "Apps",
      icon: AppStore,
      available: true,
    },
    {
      name: "Amazon",
      description: "Collect reviews from your product page on Amazon (all storefronts).",
      category: "Retail",
      icon: Amazon,
      available: true,
    },
    {
      name: "iAgora",
      description: "Collect reviews from your institution page on iAgora.",
      category: "Institutions",
      icon: IAgora,
      available: true,
    },
    {
      name: "Widget",
      description: "Collect feedbacks directly from your website through a Widget.",
      category: "Websites",
      icon: Widget,
      available: false,
    },
    {
      name: "Webhook",
      description: "Collect feedbacks directly from your custom solution through Webhooks.",
      category: "Custom",
      icon: Webhook,
      available: false,
    },
    {
      name: "Typeform",
      description: "Collect feedbacks from your surveys on Typeform.",
      category: "Surveys",
      icon: Typeform,
      available: false,
    },
    {
      name: "Zapier",
      description: "Collect feedbacks from your workflows on Zapier.",
      category: "Custom",
      icon: Zapier,
      available: false,
    },
    {
      name: "Email",
      description: "Collect feedbacks from support inquiries on your Email accounts.",
      category: "Support",
      icon: Gmail,
      available: false,
    },
    {
      name: "Intercom",
      description: "Collect feedbacks from your customer conversations on Intercom.",
      category: "Support",
      icon: Intercom,
      available: false,
    },
    {
      name: "Zendesk",
      description: "Collect feedbacks from your customer support tickets on Zendesk.",
      category: "Support",
      icon: Zendesk,
      available: false,
    },
    {
      name: "G2",
      description: "Collect reviews from your company page on G2.",
      category: "Websites",
      icon: G2,
      available: false,
    },
    {
      name: "Google My Business",
      description: "Collect reviews from your business profile on Google Search & Maps.",
      category: "Businesses",
      icon: GoogleMyBusiness,
      available: false,
    },
    {
      name: "Tripadvisor",
      description: "Collect reviews from your hotel or airbnb page on Tripadvisor.",
      category: "Accommodations",
      icon: Tripadvisor,
      available: false,
    },
    {
      name: "Yelp",
      description: "Collect reviews from your restaurant page on Yelp.",
      category: "Restaurants",
      icon: Yelp,
      available: false,
    },
  ];
</script>

<GradientAnimation class="z-0 h-[700px]" />
<Container class="relative z-10 pb-24 pt-20 text-center sm:pb-32 lg:pt-32">
  <h1
    class="md:!leading-xl mx-auto max-w-4xl font-cal text-5xl font-bold leading-[100%] tracking-[-0.002em] text-foreground sm:text-7xl"
  >
    Capture feedback
    <span class="relative whitespace-nowrap">
      <Underline />
      <span class="text-primary-gradient relative">wherever</span>
    </span>
    your customers are.
  </h1>
  <p class="mx-auto mt-8 max-w-2xl text-lg tracking-tight text-foreground">
    Dozens of integrations available so you don't miss any feedback.
  </p>
  <ul
    class="mt-20 grid h-auto w-full grid-cols-1 place-items-center gap-6 sm:mt-40 sm:gap-8 md:grid-cols-2 lg:grid-cols-3"
  >
    {#each integrations as integration}
      <li
        class="bg-secondary-gradient relative h-56 w-96 space-y-5 overflow-hidden rounded-2xl p-8 text-left text-secondary-foreground shadow-sm transition-opacity hover:opacity-80 data-[disabled]:pointer-events-none data-[disabled]:opacity-50"
        data-disabled={!integration.available || undefined}
      >
        <div class="flex w-full items-center justify-start gap-4">
          <svelte:component this={integration.icon} class="h-16 w-16 shrink-0" />
          <div class="flex w-full flex-col items-start justify-center">
            <h2 class="font-cal text-xl">{integration.name}</h2>
            <span class="text-base tracking-tight">{integration.category}</span>
          </div>
        </div>
        <p class="text-pretty px-px text-base tracking-tight">{integration.description}</p>
        <a
          href="/dash/signin"
          class="float-right !mt-3.5 inline-flex items-center justify-center text-sm font-semibold text-primary"
        >
          <span class="absolute inset-0" />
          Integrate {integration.name}
          <ArrowRightCircleIcon class="ml-1.5 h-4 w-4 flex-shrink-0 fill-primary" />
        </a>
        {#if !integration.available}
          <span class="text-md absolute -top-1 right-4 rounded-md px-2.5 py-0.5 font-cal text-muted-foreground">
            Soon
          </span>
        {/if}
      </li>
    {/each}
    <li
      class="relative h-56 w-96 space-y-5 overflow-hidden rounded-2xl border-4 border-dashed border-secondary-darker p-8 text-left text-secondary-foreground/80 transition-opacity hover:opacity-80"
    >
      <div class="flex w-full items-center justify-start gap-4">
        <SquaresPlusIcon class="h-16 w-16 shrink-0" />
        <div class="flex w-full flex-col items-start justify-center">
          <h2 class="font-cal text-xl">Integration</h2>
          <span class="text-base tracking-tight">Category</span>
        </div>
      </div>
      <p class="text-pretty px-px text-base tracking-tight">Don't see what you need? Tell us what you're missing.</p>
      <a
        href="mailto:support@clank.so"
        class="float-right !mt-3.5 inline-flex items-center justify-center text-sm font-semibold text-primary/80"
      >
        <span class="absolute inset-0" />
        Suggest an integration
        <ArrowRightCircleIcon class="ml-1.5 h-4 w-4 flex-shrink-0 fill-primary/80" />
      </a>
    </li>
  </ul>
</Container>
