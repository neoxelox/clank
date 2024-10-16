<script lang="ts">
  import Container from "$lib/components/cms/Container.svelte";
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
  import Webhook from "$lib/components/icons/Webhook.svelte";
  import Widget from "$lib/components/icons/Widget.svelte";
  import Yelp from "$lib/components/icons/Yelp.svelte";
  import Zapier from "$lib/components/icons/Zapier.svelte";
  import Zendesk from "$lib/components/icons/Zendesk.svelte";
  import { ArrowRightCircleIcon } from "@babeard/svelte-heroicons/mini";

  let integrations = [
    Trustpilot,
    PlayStore,
    AppStore,
    Amazon,
    IAgora,
    Widget,
    Webhook,
    Typeform,
    Zapier,
    Gmail,
    Intercom,
    Zendesk,
    G2,
    GoogleMyBusiness,
    Tripadvisor,
    Yelp,
  ];
</script>

<section id="integrations" class="pb-16 pt-20 text-center lg:pt-32">
  <Container>
    <p class="text-pretty text-base font-medium text-foreground">Capture feedback wherever your customers are.</p>
    <div
      class="relative mt-8 flex overflow-x-hidden"
      style="mask-image: linear-gradient(to right, rgba(0, 0, 0, 0) 0%, rgb(0, 0, 0) 25%, rgb(0, 0, 0) 75%, rgba(0, 0, 0, 0) 100%);"
    >
      <div class="flex animate-marquee-lead">
        {#each integrations as integration}
          <svelte:component this={integration} class="mx-4 h-12 !w-auto flex-shrink-0 text-foreground" />
        {/each}
      </div>
      <div class="absolute top-0 flex animate-marquee-rear">
        {#each integrations as integration}
          <svelte:component this={integration} class="mx-4 h-12 !w-auto flex-shrink-0 text-foreground" />
        {/each}
      </div>
    </div>
    <div class="mt-12 flex justify-center">
      <p
        class="bg-secondary-gradient relative rounded-full p-8 px-4 py-1.5 text-sm leading-6 text-secondary-foreground transition-opacity hover:opacity-90 hover:shadow-sm"
      >
        <span class="hidden md:inline">Dozens of integrations available so you don't miss any feedback.</span>
        <a
          href="/integrations"
          class="inline-flex items-center justify-center font-semibold text-secondary-foreground md:text-primary"
        >
          <span class="absolute inset-0" />
          Explore all integrations
          <ArrowRightCircleIcon class="ml-1.5 h-4 w-4 flex-shrink-0 fill-secondary-foreground md:fill-primary" />
        </a>
      </p>
    </div>
  </Container>
</section>
