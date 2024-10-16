<script lang="ts">
  import PricingBackground from "$assets/backgrounds/pricing.jpg";
  import Button from "$lib/components/cms/Button.svelte";
  import Container from "$lib/components/cms/Container.svelte";
  import { CheckCircleIcon } from "@babeard/svelte-heroicons/outline";
  import { RadioGroup, RadioGroupLabel, RadioGroupOption } from "@rgossiaux/svelte-headlessui";

  let plans = [
    {
      name: "Starter",
      link: "/",
      price: { month: "$49", year: "$39" },
      description: "The essentials to guide product decisions and delight customers.",
      cta: "Start 14-day free trial",
      features: [
        "10k feedbacks/month included",
        "Track up to 5 products",
        "Connect unlimited integrations",
        "Unlimited users",
        "Metrics & Insights",
        "Issues & Suggestions",
        "Basic Customer Personas",
        "Standard support",
      ],
      emphasized: false,
    },
    {
      name: "Business",
      link: "/",
      price: { month: "$349", year: "$279" },
      description: "A plan that scales with your rapidly growing business.",
      cta: "Start 14-day free trial",
      features: [
        "100k feedbacks/month included",
        "Track unlimited products",
        "Connect unlimited integrations",
        "Unlimited users",
        "Metrics & Insights",
        "Issues & Suggestions",
        "Interactive Customer Personas",
        "Standard support",
      ],
      emphasized: true,
    },
    {
      name: "Enterprise",
      link: "https://meet.brevo.com/clank-alex/demo",
      external: true,
      price: { month: "Custom", year: "Custom" },
      description: "Dedicated support and infrastructure for your company.",
      cta: "Book a demo",
      features: [
        "Custom feedbacks/month included",
        "Track unlimited products",
        "Connect unlimited integrations",
        "Unlimited users",
        "Metrics & Insights",
        "Issues & Suggestions",
        "Interactive Customer Personas",
        "Priority support",
        "Export data",
        "Dedicated infrastructure",
        "SAML-based SSO",
        "Custom solutions",
      ],
      emphasized: false,
    },
  ];

  let frequency = "year";
</script>

<section id="pricing" class="relative overflow-hidden bg-[#f8fafc] py-24 sm:py-32">
  <img
    src={PricingBackground}
    width={1558}
    height={946}
    fetchpriority="auto"
    loading="lazy"
    decoding="async"
    class="absolute left-1/2 top-0 max-w-none -translate-y-1/4 translate-x-[-30%]"
  />
  <Container class="relative">
    <div class="mx-auto max-w-7xl px-6 lg:px-8">
      <h2 class="text-pretty text-center font-cal text-3xl text-foreground sm:text-4xl md:text-5xl">
        Clank costs less than building the wrong features.
      </h2>
      <p class="mx-auto mt-6 max-w-2xl text-pretty text-center text-lg leading-8 tracking-tight text-foreground">
        Whether your customer base is large or small, we're here to help you grow. Start making informed product
        decisions with a free 14-day trial, no credit card is needed.
      </p>
      <div class="mt-16 flex justify-center">
        <fieldset>
          <RadioGroup
            bind:value={frequency}
            class="bg-secondary-gradient grid grid-cols-2 gap-x-1 rounded-full p-1 text-center text-xs font-semibold leading-5 shadow-sm"
          >
            <RadioGroupOption as="span" value="month" let:checked>
              <RadioGroupLabel
                class={"inline-block w-full cursor-pointer rounded-full px-2.5 py-1 " +
                  (checked
                    ? "bg-primary-gradient text-primary-foreground"
                    : "bg-secondary-gradient text-secondary-foreground transition-colors hover:text-secondary-foreground/60")}
              >
                Monthly
              </RadioGroupLabel>
            </RadioGroupOption>
            <RadioGroupOption as="span" value="year" let:checked>
              <RadioGroupLabel
                class={"inline-block w-full cursor-pointer rounded-full px-2.5 py-1 " +
                  (checked
                    ? "bg-primary-gradient text-primary-foreground"
                    : "bg-secondary-gradient text-secondary-foreground transition-colors hover:text-secondary-foreground/60")}
              >
                Annually <span class="font-medium">(20% off!)</span>
              </RadioGroupLabel>
            </RadioGroupOption>
          </RadioGroup>
        </fieldset>
      </div>
      <div class="isolate mx-auto mt-10 grid max-w-md grid-cols-1 gap-8 lg:mx-0 lg:max-w-none lg:grid-cols-3">
        {#each plans as plan}
          <div
            class={"rounded-3xl p-8 shadow-sm xl:p-10 " +
              (plan.emphasized
                ? "bg-primary-gradient text-primary-foreground"
                : "bg-secondary-gradient text-secondary-foreground")}
          >
            <div class="flex items-center justify-between gap-x-4">
              <h3 class="font-cal text-2xl leading-8">
                {plan.name}
              </h3>
              {#if plan.emphasized}
                <p class="rounded-full bg-primary-foreground px-2.5 py-1 text-xs font-bold leading-5 text-primary">
                  Most popular
                </p>
              {/if}
            </div>
            <p class="mt-4 text-base leading-6">
              {plan.description}
            </p>
            <p class="mt-6 flex items-baseline gap-x-1">
              <span class="font-cal text-4xl">{plan.price[frequency]}</span>
              <span class="text-sm font-semibold leading-6">/month</span>
            </p>
            <Button
              href={plan.link}
              target={plan.external ? "_blank" : undefined}
              rel={plan.external ? "noopener noreferrer" : undefined}
              color={plan.emphasized ? "secondary" : "primary"}
              class={"mt-6 w-full " + (plan.emphasized ? "!bg-primary-foreground-gradient" : "")}
            >
              {plan.cta}
            </Button>
            <ul role="list" class="mt-8 space-y-3 text-sm leading-6 xl:mt-10">
              {#each plan.features as feature}
                <li class="flex gap-x-3">
                  <CheckCircleIcon
                    class={"h-6 w-6 flex-none stroke-2 " +
                      (plan.emphasized ? "text-primary-foreground" : "text-primary")}
                  />
                  {feature}
                </li>
              {/each}
            </ul>
          </div>
        {/each}
      </div>
    </div>
  </Container>
</section>
