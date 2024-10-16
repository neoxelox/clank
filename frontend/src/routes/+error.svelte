<script lang="ts">
  import { page } from "$app/stores";
  import GradientAnimation from "$lib/components/cms/GradientAnimation.svelte";
  import Logo from "$lib/components/icons/Logo.svelte";
  import { ArrowLeftCircleIcon, ChevronRightIcon } from "@babeard/svelte-heroicons/mini";
  import { BookOpenIcon, ChatBubbleLeftRightIcon, SquaresPlusIcon } from "@babeard/svelte-heroicons/outline";

  $: error = {
    404: {
      status: 404,
      title: "This page does not exist",
      description: "Sorry, we couldn't find the page you're looking for.",
    },
  }[$page.status] || {
    status: 500,
    title: "Something went wrong",
    description: `Having issues? Reach out to our team at <a href="mailto:support@clank.so" class="text-primary">support@clank.so</a>`,
  };

  let pages = [
    {
      link: "/#features",
      title: "Blog",
      description: "Read the latest UX & CX tips to thrive in the upcoming AI era.",
      icon: ChatBubbleLeftRightIcon,
    },
    {
      link: "/#features",
      title: "Documentation",
      description: "Learn all about the tools that Clank offers you.",
      icon: BookOpenIcon,
    },
    {
      link: "/dash",
      title: "Dashboard",
      description: "Discover what your customers think about your product.",
      icon: SquaresPlusIcon,
    },
  ];
</script>

<GradientAnimation class="z-0 overflow-hidden" />
<main
  class="no-scrollbar relative z-10 mx-auto flex h-full w-full max-w-7xl flex-col overflow-y-auto px-6 pt-6 sm:pt-4 md:pt-8 lg:px-8"
>
  <Logo class="mt-auto h-12 w-auto shrink-0 text-primary" gradient="bg-primary-gradient" />
  <div class="mx-auto mt-12 max-w-2xl text-center sm:mt-16">
    <p class="text-base font-bold leading-8 text-primary">{error.status}</p>
    <h1 class="mt-4 font-cal text-3xl text-foreground sm:text-4xl md:text-5xl">{error.title}</h1>
    <p class="mt-4 text-pretty text-lg leading-8 tracking-tight text-foreground sm:mt-6">
      <!-- eslint-disable-next-line svelte/no-at-html-tags -->
      {@html error.description}
    </p>
  </div>
  <div class="mx-auto mb-auto mt-12 flow-root max-w-lg sm:mt-14">
    <ul role="list" class="-mt-6 divide-y divide-border border-b border-border">
      {#each pages as page}
        <li class="relative flex gap-x-6 py-6">
          <div class="bg-secondary-gradient flex h-10 w-10 flex-none items-center justify-center rounded-lg shadow-sm">
            <svelte:component this={page.icon} class="h-6 w-6 shrink-0 text-secondary-foreground/70" />
          </div>
          <div class="flex-auto">
            <h3 class="text-sm font-semibold leading-6 text-foreground">
              <a href={page.link}>
                <span class="absolute inset-0" aria-hidden="true" />
                {page.title}
              </a>
            </h3>
            <p class="mt-2 text-pretty text-sm leading-6 text-muted-foreground">{page.description}</p>
          </div>
          <div class="flex-none self-center">
            <ChevronRightIcon class="h-5 w-5 text-muted-foreground" />
          </div>
        </li>
      {/each}
    </ul>
    <div class="mt-10 flex justify-center">
      <a
        href="/"
        class="flex items-center justify-center text-sm font-semibold leading-6 text-primary transition-opacity hover:opacity-80"
      >
        <ArrowLeftCircleIcon class="mr-1.5 h-4 w-4 flex-shrink-0 fill-primary" />
        Back to home
      </a>
    </div>
  </div>
</main>
