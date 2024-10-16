<script lang="ts">
  import FeaturesBackground from "$assets/backgrounds/features.jpg";
  import IssuesScreenshot from "$assets/screenshots/issues.png?enhanced";
  import MetricsScreenshot from "$assets/screenshots/metrics.png?enhanced";
  import PersonasScreenshot from "$assets/screenshots/personas.png?enhanced";
  import SuggestionsScreenshot from "$assets/screenshots/suggestions.png?enhanced";
  import Container from "$lib/components/cms/Container.svelte";
  import { Tab, TabGroup, TabList, TabPanel, TabPanels } from "@rgossiaux/svelte-headlessui";

  let features = [
    {
      title: "Metrics & Insights",
      description:
        "React to trends as they are happening, predict churn and retention, track NPS, CSAT and other CX KPIs, and measure vital brand metrics.",
      image: MetricsScreenshot,
    },
    {
      title: "Issues",
      description:
        "Identify frequent problems, pinpoint critical customer pains and improve your product based on customer insights.",
      image: IssuesScreenshot,
    },
    {
      title: "Suggestions",
      description: "Address common requests, ship what customers really want and reach product-market fit.",
      image: SuggestionsScreenshot,
    },
    {
      title: "Customer Personas",
      description:
        "Synthesize interactive personas from the knowledge of all your customers. Discover what they think of your product and uncover valuable insights in real time.",
      image: PersonasScreenshot,
    },
  ];
</script>

<section id="features" class="relative overflow-hidden bg-primary pb-28 pt-20 sm:py-32">
  <img
    src={FeaturesBackground}
    width={2245}
    height={1636}
    fetchpriority="auto"
    loading="lazy"
    decoding="async"
    class="absolute left-1/2 top-1/2 max-w-none translate-x-[-44%] translate-y-[-42%]"
  />
  <Container class="relative">
    <div class="max-w-2xl text-center md:mx-auto xl:max-w-none">
      <h2 class="font-cal text-3xl text-primary-foreground sm:text-4xl md:text-5xl">
        Build better products with customer feedback.
      </h2>
      <p class="mt-6 text-pretty text-lg leading-8 tracking-tight text-primary-foreground">
        Centralize customer feedback to uncover insights and make informed product decisions.
      </p>
    </div>
    <TabGroup
      class="mt-16 grid grid-cols-1 items-center gap-y-2 pt-10 sm:gap-y-6 md:mt-20 lg:grid-cols-12 lg:pt-0"
      let:selectedIndex
    >
      <div class="no-scrollbar -mx-4 flex overflow-x-auto pb-4 sm:mx-0 sm:overflow-visible sm:pb-0 lg:col-span-5">
        <TabList
          class="relative z-10 flex gap-x-4 whitespace-nowrap px-4 sm:mx-auto sm:px-0 lg:mx-0 lg:block lg:gap-x-0 lg:gap-y-1 lg:whitespace-normal"
        >
          {#each features as feature, featureIndex}
            <div
              class={"group relative rounded-full px-4 py-1.5 lg:rounded-l-xl lg:rounded-r-none lg:p-6 " +
                (selectedIndex === featureIndex
                  ? "bg-primary-foreground lg:bg-primary-foreground/10 lg:ring-1 lg:ring-inset lg:ring-primary-foreground/10"
                  : "transition-colors hover:bg-primary-foreground/10 lg:hover:bg-primary-foreground/5")}
            >
              <h3>
                <Tab
                  class={"font-cal text-lg " +
                    (selectedIndex === featureIndex
                      ? "text-primary lg:text-primary-foreground"
                      : "text-primary-foreground/80 transition-colors hover:text-primary-foreground")}
                >
                  <span class="absolute inset-0 rounded-full lg:rounded-l-xl lg:rounded-r-none" />
                  {feature.title}
                </Tab>
              </h3>
              <p
                class={"mt-2 hidden text-pretty text-base lg:block " +
                  (selectedIndex === featureIndex
                    ? "text-primary-foreground"
                    : "text-primary-foreground/80 transition-colors group-hover:text-primary-foreground")}
              >
                {feature.description}
              </p>
              <!-- Trick to preload images behind a button -->
              <enhanced:img
                src={feature.image}
                sizes="(min-width: 1024px) 67.8125rem, (min-width: 640px) 100vw, 45rem"
                fetchpriority="high"
                loading="eager"
                decoding="async"
                style="display: none;"
              />
            </div>
          {/each}
        </TabList>
      </div>
      <TabPanels class="lg:col-span-7">
        {#each features as feature (feature.title)}
          <TabPanel>
            <div class="relative sm:px-6 lg:hidden">
              <div
                class="absolute -inset-x-4 bottom-[-4.25rem] top-[-6.5rem] bg-primary-foreground/10 ring-1 ring-inset ring-primary-foreground/10 sm:inset-x-0 sm:rounded-t-xl"
              />
              <p class="relative mx-auto max-w-2xl text-pretty text-base text-primary-foreground sm:text-center">
                {feature.description}
              </p>
            </div>
            <div
              class="mt-10 w-[45rem] overflow-hidden rounded-xl bg-primary-foreground shadow-xl shadow-primary/20 sm:w-auto lg:mt-0 lg:w-[67.8125rem]"
            >
              <enhanced:img
                src={feature.image}
                sizes="(min-width: 1024px) 67.8125rem, (min-width: 640px) 100vw, 45rem"
                fetchpriority="high"
                loading="eager"
                decoding="sync"
                class="w-full"
              />
            </div>
          </TabPanel>
        {/each}
      </TabPanels>
    </TabGroup>
  </Container>
</section>
