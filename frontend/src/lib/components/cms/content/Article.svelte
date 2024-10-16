<script lang="ts">
  import Container from "$lib/components/cms/Container.svelte";
  import GradientAnimation from "$lib/components/cms/GradientAnimation.svelte";
  import dayjs from "$lib/utils/datetime";
  import { titlelize } from "$lib/utils/string";
  import { setContext } from "svelte";

  export let id: string;
  export let group: string;
  export let date: Date | string | undefined = undefined;
  if (date) date = new Date(date);

  setContext("content-scope", "Article");
</script>

<section {id} class="relative bg-background pt-10 sm:pt-14">
  <GradientAnimation class="z-0 h-96" />
  <Container class="relative z-10 px-8 sm:px-8 lg:px-8">
    <div class="mx-auto max-w-3xl text-pretty text-base leading-7 text-foreground [&_b]:font-semibold">
      <p class="font-cal text-base leading-7 text-primary">
        {titlelize(group)}
        {#if date}
          <span class="text-xs text-primary/50 transition-colors hover:text-primary">
            / {dayjs(date).format("MMMM D, YYYY")}
          </span>
        {/if}
      </p>
      <slot />
    </div>
  </Container>
</section>
