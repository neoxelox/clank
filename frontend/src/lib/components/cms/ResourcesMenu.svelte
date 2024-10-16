<script lang="ts">
  import Slack from "$lib/components/icons/Slack.svelte";
  import { BookOpenIcon, ChatBubbleLeftRightIcon, ChevronDownIcon } from "@babeard/svelte-heroicons/outline";
  import { CalendarIcon } from "@babeard/svelte-heroicons/solid";
  import { Popover, PopoverButton, PopoverPanel, Transition, TransitionChild } from "@rgossiaux/svelte-headlessui";

  export let mobile: boolean = false;
</script>

<Popover class="relative" let:open>
  {#if !mobile}
    <PopoverButton
      class={"hover:bg-secondary-gradient hover:shadow-xs inline-block rounded-xl bg-transparent px-2.5 py-1.5 text-sm font-semibold text-secondary-foreground transition-colors " +
        ($$props.class || "")}
    >
      <span class="inline-flex items-center gap-x-1">
        Resources
        <ChevronDownIcon class="mt-px h-3.5 w-3.5 shrink-0 stroke-[2.5]" />
      </span>
    </PopoverButton>
  {:else}
    <PopoverButton class={"block w-full p-2.5 text-start " + ($$props.class || "")}>
      <span class="inline-flex items-center gap-x-1">
        Resources
        <ChevronDownIcon class="mt-px h-3.5 w-3.5 shrink-0 stroke-[2.5]" />
      </span>
    </PopoverButton>
  {/if}
  <Transition show={open}>
    <TransitionChild
      enter="transition ease-out duration-200"
      enterFrom="opacity-0 translate-y-1"
      enterTo="opacity-100 translate-y-0"
      leave="transition ease-in duration-150"
      leaveFrom="opacity-100 translate-y-0"
      leaveTo="opacity-0 translate-y-1"
    >
      <PopoverPanel class="absolute left-1/2 z-20 mt-2.5 flex w-screen max-w-max -translate-x-1/2 px-4">
        <div
          class="w-screen max-w-md flex-auto overflow-hidden rounded-3xl bg-background text-sm leading-6 shadow-lg ring-1 ring-foreground/5"
        >
          <div class="p-4">
            <div class="hover:bg-secondary-gradient group relative flex gap-x-6 rounded-lg p-4">
              <div
                class="bg-secondary-gradient mt-1 flex h-11 w-11 flex-none items-center justify-center rounded-lg group-hover:bg-background group-hover:!bg-none"
              >
                <ChatBubbleLeftRightIcon
                  class="h-6 w-6 shrink-0 text-secondary-foreground/70 group-hover:text-primary"
                />
              </div>
              <div>
                <a href="/#features" class="font-semibold text-foreground">
                  Blog
                  <span class="absolute inset-0" />
                </a>
                <p class="mt-1 text-pretty text-muted-foreground">
                  Read the latest UX & CX tips to thrive in the upcoming AI era.
                </p>
              </div>
            </div>
            <div class="hover:bg-secondary-gradient group relative flex gap-x-6 rounded-lg p-4">
              <div
                class="bg-secondary-gradient mt-1 flex h-11 w-11 flex-none items-center justify-center rounded-lg group-hover:bg-background group-hover:!bg-none"
              >
                <BookOpenIcon class="h-6 w-6 shrink-0 text-secondary-foreground/70 group-hover:text-primary" />
              </div>
              <div>
                <a href="/#features" class="font-semibold text-foreground">
                  Documentation
                  <span class="absolute inset-0" />
                </a>
                <p class="mt-1 text-pretty text-muted-foreground">Learn all about the tools that Clank offers you.</p>
              </div>
            </div>
          </div>
          <div class="grid grid-cols-2 divide-x divide-secondary-foreground/5 bg-background">
            <a
              href="https://meet.brevo.com/clank-alex/demo"
              target="_blank"
              rel="noopener noreferrer"
              class="bg-secondary-gradient flex items-center justify-center gap-x-2.5 p-3 font-semibold text-secondary-foreground transition-opacity hover:opacity-80 hover:shadow-sm"
            >
              <CalendarIcon class="h-5 w-5 flex-none shrink-0 text-secondary-foreground/70" />
              Book demo
            </a>
            <a
              href="https://join.slack.com/t/clank-so-community/shared_invite/zt-2lja8kgir-H2k99~E1tcahgiYAsxN~HA"
              target="_blank"
              rel="noopener noreferrer"
              class="bg-secondary-gradient flex items-center justify-center gap-x-2.5 p-3 font-semibold text-secondary-foreground transition-opacity hover:opacity-80 hover:shadow-sm"
            >
              <Slack class="h-4 !w-4 flex-none shrink-0 text-secondary-foreground/70" />
              Join community
            </a>
          </div>
        </div>
      </PopoverPanel>
    </TransitionChild>
  </Transition>
</Popover>
