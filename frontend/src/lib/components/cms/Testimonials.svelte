<script lang="ts">
  import FirstAvatar from "$assets/avatars/avatar-1.png?enhanced";
  import SecondAvatar from "$assets/avatars/avatar-2.png?enhanced";
  import ThirdAvatar from "$assets/avatars/avatar-3.png?enhanced";
  import FourthAvatar from "$assets/avatars/avatar-4.png?enhanced";
  import FifthAvatar from "$assets/avatars/avatar-5.png?enhanced";
  import Container from "$lib/components/cms/Container.svelte";
  import Quote from "$lib/components/icons/Quote.svelte";

  let testimonials_1 = [
    {
      content:
        "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam nec dui sed nisi porta fermentum a nec diam eros.",
      author: {
        name: "Sheryl Berge",
        role: "CEO at Lynch LLC",
        image: FirstAvatar,
      },
    },
    {
      content:
        "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Integer non purus ut massa imperdiet vehicula. Nunc lacinia ullamcorper lorem, eu blandit libero. Etiam nec egestas.",
      author: {
        name: "Amy Hahn",
        role: "Director at Velocity Industries",
        image: FourthAvatar,
      },
    },
  ];

  let testimonials_2 = [
    {
      content:
        "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Quisque tempus lectus sed sapien placerat, vel pharetra leo consectetur. Proin pellentesque suscipit tortor, eu hendrerit neque lobortis sit dui.",
      author: {
        name: "Leland Kiehn",
        role: "Founder of Kiehn and Sons",
        image: FifthAvatar,
      },
    },
    {
      content:
        "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Pellentesque iaculis tristique odio ac auctor. Sed eget odio nisi. Nulla sodales sed.",
      author: {
        name: "Erin Powlowski",
        role: "COO at Armstrong Inc",
        image: SecondAvatar,
      },
    },
  ];

  let testimonials_3 = [
    {
      content:
        "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Ut eu nisi diam. Ut vitae nibh aliquet, interdum odio vitae, tempus nunc.",
      author: {
        name: "Peter Renolds",
        role: "Founder of West Inc",
        image: ThirdAvatar,
      },
    },
    {
      content:
        "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam vulputate leo non massa posuere, placerat tincidunt magna rutrum. Ut non ipsum at dui turpis.",
      author: {
        name: "Amy Hahn",
        role: "Director at Velocity Industries",
        image: FourthAvatar,
      },
    },
  ];
</script>

<section id="customers" class="bg-background py-24 sm:py-32">
  <Container>
    <div class="mx-auto max-w-2xl text-center">
      <h2 class="font-cal text-3xl text-foreground sm:text-4xl md:text-5xl">
        Don't just take our word&nbsp;for&nbsp;it.
      </h2>
      <p class="mt-4 text-pretty text-lg leading-8 tracking-tight text-foreground">Join our happy ❤️ customers.</p>
    </div>
    <ul
      role="list"
      class="mx-auto mt-16 grid max-w-2xl grid-cols-1 gap-6 sm:gap-8 lg:mt-20 lg:max-w-none lg:grid-cols-3"
    >
      <li>
        <ul role="list" class="flex flex-col gap-y-6 sm:gap-y-8">
          {#each testimonials_1 as testimonial}
            <li>
              <figure
                class="bg-secondary-gradient relative rounded-2xl p-6 transition-opacity hover:opacity-80 hover:shadow-sm"
              >
                <Quote class="absolute left-6 top-6 fill-background" />
                <blockquote class="relative">
                  <p class="text-lg leading-loose tracking-tight text-secondary-foreground">
                    {testimonial.content}
                  </p>
                </blockquote>
                <figcaption class="relative mt-6 flex items-center justify-between border-t border-background pt-6">
                  <div>
                    <div class="text-base font-semibold text-secondary-foreground">
                      {testimonial.author.name}
                    </div>
                    <div class="mt-1 text-sm text-secondary-foreground">
                      {testimonial.author.role}
                    </div>
                  </div>
                  <div class="overflow-hidden rounded-full border-2 border-background bg-background">
                    <enhanced:img
                      src={testimonial.author.image}
                      fetchpriority="high"
                      loading="eager"
                      decoding="async"
                      class="h-14 w-14 object-cover"
                    />
                  </div>
                </figcaption>
              </figure>
            </li>
          {/each}
        </ul>
      </li>
      <li>
        <ul role="list" class="flex flex-col gap-y-6 sm:gap-y-8">
          {#each testimonials_2 as testimonial}
            <li>
              <figure
                class="bg-secondary-gradient relative rounded-2xl p-6 transition-opacity hover:opacity-80 hover:shadow-sm"
              >
                <Quote class="absolute left-6 top-6 fill-background" />
                <blockquote class="relative">
                  <p class="text-lg leading-loose tracking-tight text-secondary-foreground">
                    {testimonial.content}
                  </p>
                </blockquote>
                <figcaption class="relative mt-6 flex items-center justify-between border-t border-background pt-6">
                  <div>
                    <div class="text-base font-semibold text-secondary-foreground">
                      {testimonial.author.name}
                    </div>
                    <div class="mt-1 text-sm text-secondary-foreground">
                      {testimonial.author.role}
                    </div>
                  </div>
                  <div class="overflow-hidden rounded-full border-2 border-background bg-background">
                    <enhanced:img
                      src={testimonial.author.image}
                      fetchpriority="high"
                      loading="eager"
                      decoding="async"
                      class="h-14 w-14 object-cover"
                    />
                  </div>
                </figcaption>
              </figure>
            </li>
          {/each}
        </ul>
      </li>
      <li>
        <ul role="list" class="flex flex-col gap-y-6 sm:gap-y-8">
          {#each testimonials_3 as testimonial}
            <li>
              <figure
                class="bg-secondary-gradient relative rounded-2xl p-6 transition-opacity hover:opacity-80 hover:shadow-sm"
              >
                <Quote class="absolute left-6 top-6 fill-background" />
                <blockquote class="relative">
                  <p class="text-lg leading-loose tracking-tight text-secondary-foreground">
                    {testimonial.content}
                  </p>
                </blockquote>
                <figcaption class="relative mt-6 flex items-center justify-between border-t border-background pt-6">
                  <div>
                    <div class="text-base font-semibold text-secondary-foreground">
                      {testimonial.author.name}
                    </div>
                    <div class="mt-1 text-sm text-secondary-foreground">
                      {testimonial.author.role}
                    </div>
                  </div>
                  <div class="overflow-hidden rounded-full border-2 border-background bg-background">
                    <enhanced:img
                      src={testimonial.author.image}
                      fetchpriority="high"
                      loading="eager"
                      decoding="async"
                      class="h-14 w-14 object-cover"
                    />
                  </div>
                </figcaption>
              </figure>
            </li>
          {/each}
        </ul>
      </li>
    </ul>
  </Container>
</section>
