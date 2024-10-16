<script lang="ts">
  import Step1Drawing from "$assets/drawings/step-1.png?enhanced";
  import Step2Drawing from "$assets/drawings/step-2.png?enhanced";
  import Step3Drawing from "$assets/drawings/step-3.png?enhanced";
  import Step4Drawing from "$assets/drawings/step-4.png?enhanced";
  import Container from "$lib/components/cms/Container.svelte";
  import { ArrowRightCircleIcon } from "@babeard/svelte-heroicons/mini";
</script>

<section id="how-it-works" class="bg-background py-24 sm:py-32">
  <Container class="relative">
    <div class="mx-auto max-w-7xl px-6 lg:px-8">
      <div class="mx-auto max-w-2xl lg:max-w-none">
        <h2 class="text-center font-cal text-3xl text-foreground sm:text-4xl md:text-5xl">How it works.</h2>
        <p class="mx-auto mt-6 max-w-2xl text-pretty text-center text-lg leading-8 tracking-tight text-foreground">
          Focus on shipping while we take care of your customer feedback.
        </p>
        <u class="mt-20 flex list-none flex-col gap-20 no-underline sm:gap-24 md:gap-36">
          <li class="flex flex-col items-center gap-12 md:flex-row md:self-start">
            <enhanced:img
              src={Step1Drawing}
              fetchpriority="high"
              loading="eager"
              decoding="async"
              class="h-auto w-full shrink-0 rounded-3xl bg-muted object-cover shadow-xl md:h-[336px] md:w-[500px]"
            />
            <span class="max-w-xl text-center md:text-left">
              <h3 class="text-pretty font-cal text-xl text-foreground md:text-2xl">
                1. Bring in customer feedback, from various sources.
              </h3>
              <p class="mt-3 text-pretty text-base text-muted-foreground md:text-lg">
                Clank integrates with platforms customers already use, meaning you can get started right now. Clank's
                infrastructure is scalable by default and custom solutions can be set-up effortlessly.
              </p>
              <p
                class="bg-secondary-gradient relative mt-6 inline-block rounded-full p-8 px-4 py-1.5 text-sm leading-6 text-secondary-foreground transition-opacity hover:opacity-90 hover:shadow-sm"
              >
                <a
                  href="/integrations"
                  class="inline-flex items-center justify-center font-semibold text-secondary-foreground md:text-primary"
                >
                  <span class="absolute inset-0" />
                  Explore all integrations
                  <ArrowRightCircleIcon
                    class="ml-1.5 h-4 w-4 flex-shrink-0 fill-secondary-foreground md:fill-primary"
                  />
                </a>
              </p>
            </span>
          </li>
          <li class="flex flex-col items-center gap-12 md:flex-row-reverse md:self-end">
            <enhanced:img
              src={Step2Drawing}
              fetchpriority="high"
              loading="eager"
              decoding="async"
              class="h-auto w-full shrink-0 rounded-3xl bg-muted object-cover shadow-xl md:h-[336px] md:w-[500px]"
            />
            <span class="max-w-xl text-center md:text-left">
              <h3 class="text-pretty font-cal text-xl text-foreground md:text-2xl">
                2. Leverage our state-of-the-art AI engine and save days spent on manual analysis.
              </h3>
              <p class="mt-3 text-pretty text-base text-muted-foreground md:text-lg">
                Clank is a tool that works for you, literally. No more wasting time reading thousands of reviews. Forget
                about never-ending spreadsheets. Our AI-powered engine is able to identify common issues and
                suggestions, synthesize interactive customer personas, predict churn and retention and measure vital
                KPIs.
              </p>
            </span>
          </li>
          <li class="flex flex-col items-center gap-12 md:flex-row md:self-start">
            <enhanced:img
              src={Step3Drawing}
              fetchpriority="high"
              loading="eager"
              decoding="async"
              class="h-auto w-full shrink-0 rounded-3xl bg-muted object-cover shadow-xl md:h-[336px] md:w-[500px]"
            />
            <span class="max-w-xl text-center md:text-left">
              <h3 class="text-pretty font-cal text-xl text-foreground md:text-2xl">
                3. Uncover valuable insights from trending issues, suggestions and customer personas.
              </h3>
              <p class="mt-3 text-pretty text-base text-muted-foreground md:text-lg">
                Clank's analytics transforms raw feedback into clear, actionable insights. From KPIs to customer
                personas, Clank makes it easy to know your customers better than ever before. Pinpoint critical customer
                pains and ship requested features. Interact with your customer personas and discover valuable insights
                in real time.
              </p>
            </span>
          </li>
          <li class="flex flex-col items-center gap-12 md:flex-row-reverse md:self-end">
            <enhanced:img
              src={Step4Drawing}
              fetchpriority="high"
              loading="eager"
              decoding="async"
              class="h-auto w-full shrink-0 rounded-3xl bg-muted object-cover shadow-xl md:h-[336px] md:w-[500px]"
            />
            <span class="max-w-xl text-center md:text-left">
              <h3 class="text-pretty font-cal text-xl text-foreground md:text-2xl">
                4. Make informed product decisions and delight customers.
              </h3>
              <p class="mt-3 text-pretty text-base text-muted-foreground md:text-lg">
                Clank helps your entire organization stay connected to customer feedback. With in-depth insights at your
                fingertips, you can now build what customers want, not what you think they need. Understand how your new
                feature is seen by your customers before it flops in the market. Discover emerging trends and stay ahead
                of your competitors.
              </p>
            </span>
          </li>
        </u>
      </div>
    </div>
  </Container>
</section>
