<script lang="ts">
  import { page } from "$app/stores";
  import Button from "$lib/components/cms/Button.svelte";
  import Container from "$lib/components/cms/Container.svelte";
  import NavLink from "$lib/components/cms/NavLink.svelte";
  import Logo from "$lib/components/icons/Logo.svelte";
  import MobileMenu from "$lib/components/icons/MobileMenu.svelte";
  import { ArrowRightCircleIcon } from "@babeard/svelte-heroicons/mini";
  import {
    Popover,
    PopoverButton,
    PopoverOverlay,
    PopoverPanel,
    Transition,
    TransitionChild,
  } from "@rgossiaux/svelte-headlessui";
  import ResourcesMenu from "./ResourcesMenu.svelte";
</script>

<header class="px-4 py-10 lg:px-0">
  <Container>
    <nav class="relative z-50 flex justify-between">
      <div class="flex items-center md:gap-x-12">
        <a href="/">
          <Logo class="h-10 w-auto shrink-0 text-primary" gradient="bg-primary-gradient" />
        </a>
        <div class="hidden md:flex md:gap-x-6">
          <NavLink
            href="/#features"
            class={"font-semibold " +
              ($page.url.pathname.startsWith("/features") ? "bg-secondary-gradient shadow-xs" : "")}
          >
            Features
          </NavLink>
          <NavLink
            href="/integrations"
            class={"font-semibold " +
              ($page.url.pathname.startsWith("/integrations") ? "bg-secondary-gradient shadow-xs" : "")}
          >
            Integrations
          </NavLink>
          <NavLink
            href="/#features"
            class={"font-semibold " +
              ($page.url.pathname.startsWith("/customers") ? "bg-secondary-gradient shadow-xs" : "")}
          >
            Customers
          </NavLink>
          <NavLink
            href="/#pricing"
            class={"font-semibold " +
              ($page.url.pathname.startsWith("/pricing") ? "bg-secondary-gradient shadow-xs" : "")}
          >
            Pricing
          </NavLink>
          <ResourcesMenu />
        </div>
      </div>
      <div class="flex items-center gap-x-5 md:gap-x-8">
        <div class="hidden md:block">
          <NavLink href="/dash/signin" class="flex items-center justify-center font-semibold">
            Sign in
            <ArrowRightCircleIcon class="ml-1.5 h-4 w-4 flex-shrink-0 fill-foreground" />
          </NavLink>
        </div>
        <Button href="/dash/signin" color="primary">Try Clank for free</Button>
        <div class="-mr-1 md:hidden">
          <Popover let:open>
            <PopoverButton class="relative z-10 flex h-8 w-8 items-center justify-center">
              <MobileMenu {open} />
            </PopoverButton>
            <Transition show={open}>
              <TransitionChild
                enter="duration-150 ease-out"
                enterFrom="opacity-0"
                enterTo="opacity-100"
                leave="duration-150 ease-in"
                leaveFrom="opacity-100"
                leaveTo="opacity-0"
                class="fixed inset-0 bg-background/70 backdrop-blur-sm"
              >
                <PopoverOverlay />
              </TransitionChild>
              <TransitionChild
                enter="duration-150 ease-out"
                enterFrom="opacity-0 scale-95"
                enterTo="opacity-100 scale-100"
                leave="duration-100 ease-in"
                leaveFrom="opacity-100 scale-100"
                leaveTo="opacity-0 scale-95"
                class="absolute inset-x-0 top-full mt-4 flex origin-top flex-col rounded-2xl border border-border bg-background p-4 text-lg tracking-tight text-foreground shadow-lg"
              >
                <PopoverPanel>
                  <NavLink href="/#features" mobile>Features</NavLink>
                  <NavLink href="/integrations" mobile>Integrations</NavLink>
                  <NavLink href="/#features" mobile>Customers</NavLink>
                  <NavLink href="/#pricing" mobile>Pricing</NavLink>
                  <ResourcesMenu mobile />
                  <hr class="m-2 border-border" />
                  <NavLink href="/dash/signin" mobile class="flex items-center justify-start font-semibold">
                    Sign in
                    <ArrowRightCircleIcon class="ml-1.5 h-5 w-5 flex-shrink-0 fill-foreground" />
                  </NavLink>
                </PopoverPanel>
              </TransitionChild>
            </Transition>
          </Popover>
        </div>
      </div>
    </nav>
  </Container>
</header>
