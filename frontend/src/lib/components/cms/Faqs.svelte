<script lang="ts">
  import Container from "$lib/components/cms/Container.svelte";
  import { MinusSmallIcon, PlusSmallIcon } from "@babeard/svelte-heroicons/outline";
  import { Disclosure, DisclosureButton, DisclosurePanel } from "@rgossiaux/svelte-headlessui";

  let faqs = [
    {
      question: "How is feedback counted?",
      answer:
        "Each customer feedback that is ingested through our API or collected through our different integrations is counted towards your organization's usage. Further processing of the feedback is included at no cost.",
    },
    {
      question: "How is feedback collected?",
      answer:
        'We know that customers can leave feedback on a lot of different platforms. We believe that centralizing feedback in one place is key to making informed product decisions. For your convenience, we aim to provide various integrations to automatically collect feedback from these sources every single day. <a href="/integrations" class="text-primary">Explore all integrations</a>.',
    },
    {
      question: "How is feedback analyzed?",
      answer:
        "Our in-house advanced analysis engine, leverages AI to identify issues and suggestions from each customer feedback, and aggregates them so you can focus on the most recurring problems and requests. It is also capable of predicting churn and retention, along with other CX KPIs such as NPS or CSAT. Our latest advancements in feedback analysis, have allowed us to synthesize interactive personas from the knowledge of all your customers.",
    },
    {
      question: "What sets Clank apart from other feedback tools?",
      answer:
        "Clank goes beyond traditional feedback analysis tools. Incorporating advanced AI processing, allows us to uncover more detailed and valuable insights at a scale (and a price) never seen before. Clank can identify and aggregate the most common issues and suggestions, synthesize interactive personas from the knowledge of all your customers, predict churn, retention, NPS, CSAT and much, much more.",
    },
    {
      question: "How do I start using Clank?",
      answer: `You can start analyzing customer feedback free in no time. Sign in to the <a href="/dash/signin" class="text-primary">Dashboard</a>, add your product, integrate Play Store, App Store… or other sources, and… that's it, really!`,
    },
    {
      question: "What if I decide to cancel?",
      answer:
        "If you no longer wish to use Clank, you may cancel at any time. Your subscription will stop at the end of the current month or year depending on your plan.",
    },
  ];
</script>

<section id="faqs" class="bg-background pb-8 pt-24 sm:pb-16 sm:pt-32">
  <Container class="relative">
    <div class="mx-auto max-w-7xl px-6 lg:px-8">
      <h2 class="text-center font-cal text-3xl text-foreground sm:text-4xl md:text-5xl">Frequently asked questions.</h2>
      <p class="mt-4 text-pretty text-center text-lg leading-8 tracking-tight text-foreground">
        Can't find the answer you're looking for? Reach out to <a href="mailto:support@clank.so" class="text-primary"
          >our team</a
        >.
      </p>
      <dl class="mt-10 space-y-6 divide-y divide-border">
        {#each faqs as faq}
          <Disclosure as="div" class="pt-6" let:open>
            <dt>
              <DisclosureButton class="flex w-full items-start justify-between text-left text-foreground">
                <span class="text-pretty text-base font-semibold leading-7">{faq.question}</span>
                <span class="ml-6 flex h-7 items-center">
                  {#if open}
                    <MinusSmallIcon class="h-6 w-6 text-foreground" />
                  {:else}
                    <PlusSmallIcon class="h-6 w-6 text-foreground" />
                  {/if}
                </span>
              </DisclosureButton>
            </dt>
            <DisclosurePanel as="dd" class="mt-2 pr-12">
              <p class="text-pretty text-base leading-7 text-muted-foreground">
                <!-- eslint-disable-next-line svelte/no-at-html-tags -->
                {@html faq.answer}
              </p>
            </DisclosurePanel>
          </Disclosure>
        {/each}
      </dl>
    </div>
  </Container>
</section>
