<script lang="ts">
  import Container from "$lib/components/cms/Container.svelte";
  import Typeform from "$lib/components/icons/Typeform.svelte";
  import { ArrowRightCircleIcon } from "@babeard/svelte-heroicons/mini";

  let companies_1 = [Typeform, Typeform, Typeform];
  let companies_2 = [Typeform, Typeform, Typeform];
</script>

<Container class="pb-16 pt-20 text-center lg:pt-32">
  <p class="text-base font-medium text-foreground">Trusted by these six companies so far</p>
  <ul
    role="list"
    class="mt-8 flex items-center justify-center gap-x-8 sm:flex-col sm:gap-x-0 sm:gap-y-10 xl:flex-row xl:gap-x-12 xl:gap-y-0"
  >
    <li>
      <ul role="list" class="flex flex-col items-center gap-y-8 sm:flex-row sm:gap-x-12 sm:gap-y-0">
        {#each companies_1 as company}
          <li class="flex">
            <svelte:component this={company} class="h-12 w-auto text-foreground" />
          </li>
        {/each}
      </ul>
    </li>
    <li>
      <ul role="list" class="flex flex-col items-center gap-y-8 sm:flex-row sm:gap-x-12 sm:gap-y-0">
        {#each companies_2 as company}
          <li class="flex">
            <svelte:component this={company} class="h-12 w-auto text-foreground" />
          </li>
        {/each}
      </ul>
    </li>
  </ul>
  <div class="mt-12 flex justify-center">
    <p
      class="bg-secondary-gradient relative rounded-full p-8 px-4 py-1.5 text-sm leading-6 text-secondary-foreground transition-opacity hover:opacity-90 hover:shadow-sm"
    >
      <span class="hidden md:inline">Over 2500 companies use our tools to better their business.</span>
      <a
        href="/#features"
        class="inline-flex items-center justify-center font-semibold text-secondary-foreground md:text-primary"
      >
        <span class="absolute inset-0" />
        Read our customer stories
        <ArrowRightCircleIcon class="ml-1.5 h-4 w-4 flex-shrink-0 fill-secondary-foreground md:fill-primary" />
      </a>
    </p>
  </div>
</Container>
