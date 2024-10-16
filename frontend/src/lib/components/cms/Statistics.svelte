<script lang="ts">
  import Container from "$lib/components/cms/Container.svelte";

  let stats = [
    {
      title: "Feedbacks analyzed",
      value: "100K+",
    },
    {
      title: "Customer satisfaction increased",
      value: "30%",
    },
    {
      title: "Issues identified",
      value: "15K",
    },
    {
      title: "Time saved",
      value: "70%",
    },
  ];
</script>

<section id="statistics" class="bg-background py-24 sm:py-32">
  <Container class="relative">
    <div class="mx-auto max-w-7xl px-6 lg:px-8">
      <div class="mx-auto max-w-2xl lg:max-w-none">
        <h2 class="text-center font-cal text-3xl text-foreground sm:text-4xl md:text-5xl">
          Analyze customer feedback at scale.
        </h2>
        <dl
          class="bg-secondary-gradient mt-16 grid grid-cols-1 overflow-hidden rounded-2xl text-center sm:grid-cols-2 lg:grid-cols-4"
        >
          {#each stats as stat}
            <div
              class="flex flex-col bg-transparent p-8 outline outline-1 outline-background transition-opacity hover:opacity-80 hover:shadow-sm"
            >
              <dt class="font-cal text-sm leading-6 text-secondary-foreground/80">{stat.title}</dt>
              <dd class="order-first font-cal text-3xl text-secondary-foreground">
                {stat.value}
              </dd>
            </div>
          {/each}
        </dl>
      </div>
    </div>
  </Container>
</section>
