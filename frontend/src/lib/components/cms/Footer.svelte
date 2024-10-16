<script lang="ts">
  import Button from "$lib/components/cms/Button.svelte";
  import Container from "$lib/components/cms/Container.svelte";
  import NavLink from "$lib/components/cms/NavLink.svelte";
  import LinkedIn from "$lib/components/icons/LinkedIn.svelte";
  import Slack from "$lib/components/icons/Slack.svelte";
  import Twitter from "$lib/components/icons/Twitter.svelte";
  import YouTube from "$lib/components/icons/YouTube.svelte";
  import { PlayIcon } from "@babeard/svelte-heroicons/mini";

  let socials = [
    {
      link: "https://www.linkedin.com/company/clank-so/",
      logo: LinkedIn,
    },
    {
      link: "https://join.slack.com/t/clank-so-community/shared_invite/zt-2lja8kgir-H2k99~E1tcahgiYAsxN~HA",
      logo: Slack,
    },
    {
      link: "https://x.com/clank_so",
      logo: Twitter,
    },
    {
      link: "https://www.youtube.com/@clank-so",
      logo: YouTube,
    },
  ];
</script>

<footer>
  <Container>
    <nav class="pb-8 pt-20 sm:pt-24 lg:pt-32">
      <div class="xl:grid xl:grid-cols-3 xl:gap-8">
        <div class="grid grid-cols-2 gap-8 text-center md:text-left xl:col-span-2">
          <div class="md:grid md:grid-cols-2 md:gap-8">
            <div>
              <h3 class="text-sm font-semibold leading-6 text-foreground">Product</h3>
              <ul role="list" class="mt-6 space-y-4">
                <li>
                  <NavLink href="/#features" class="-ml-2.5 font-normal">Features</NavLink>
                </li>
                <li>
                  <NavLink href="/#how-it-works" class="-ml-2.5 font-normal">How it works</NavLink>
                </li>
                <li>
                  <NavLink href="/integrations" class="-ml-2.5 font-normal">Integrations</NavLink>
                </li>
                <li>
                  <NavLink href="/#pricing" class="-ml-2.5 font-normal">Pricing</NavLink>
                </li>
              </ul>
            </div>
            <div class="mt-10 md:mt-0">
              <h3 class="text-sm font-semibold leading-6 text-foreground">Support</h3>
              <ul role="list" class="mt-6 space-y-4">
                <li>
                  <NavLink href="mailto:support@clank.so" class="-ml-2.5 font-normal">Contact us</NavLink>
                </li>
                <li>
                  <NavLink
                    href="https://join.slack.com/t/clank-so-community/shared_invite/zt-2lja8kgir-H2k99~E1tcahgiYAsxN~HA"
                    target="_blank"
                    rel="noopener noreferrer"
                    class="-ml-2.5 font-normal"
                  >
                    Slack community
                  </NavLink>
                </li>
                <li>
                  <NavLink href="mailto:support@clank.so" class="-ml-2.5 font-normal">Give feedback</NavLink>
                </li>
              </ul>
            </div>
          </div>
          <div class="md:grid md:grid-cols-2 md:gap-8">
            <div>
              <h3 class="text-sm font-semibold leading-6 text-foreground">Resources</h3>
              <ul role="list" class="mt-6 space-y-4">
                <li>
                  <NavLink href="/#features" class="-ml-2.5 font-normal">Blog</NavLink>
                </li>
                <li>
                  <NavLink href="/#features" class="-ml-2.5 font-normal">Documentation</NavLink>
                </li>
                <li>
                  <NavLink href="/#features" class="-ml-2.5 font-normal">Security</NavLink>
                </li>
                <li>
                  <NavLink
                    href="https://meet.brevo.com/clank-alex/demo"
                    target="_blank"
                    rel="noopener noreferrer"
                    class="-ml-2.5 font-normal"
                  >
                    Book a demo
                  </NavLink>
                </li>
              </ul>
            </div>
            <div class="mt-10 md:mt-0">
              <h3 class="text-sm font-semibold leading-6 text-foreground">Legal</h3>
              <ul role="list" class="mt-6 space-y-4">
                <li>
                  <NavLink href="/terms" class="-ml-2.5 font-normal">Terms</NavLink>
                </li>
                <li>
                  <NavLink href="/privacy" class="-ml-2.5 font-normal">Privacy</NavLink>
                </li>
                <li>
                  <NavLink href="/privacy#8-cookies-and-tracking-technologies" class="-ml-2.5 font-normal">
                    Cookies
                  </NavLink>
                </li>
              </ul>
            </div>
          </div>
        </div>
        <div class="mt-10 text-center md:text-left xl:mt-0">
          <h3 class="text-sm font-semibold leading-6 text-foreground">Customer feedback analysis on autopilot.</h3>
          <p class="mt-2 text-sm font-normal leading-6 text-muted-foreground">
            Collect and analyze customer feedback at scale with AI. Identify common issues and suggestions, predict
            churn and retention, track CX KPIs and synthesize interactive customer personas.
          </p>
          <div class="mt-6 flex justify-center gap-x-6 md:justify-start">
            <Button href="/dash/signin" color="primary">Try Clank for free</Button>
            <Button
              href="https://www.youtube.com/@clank-so"
              target="_blank"
              rel="noopener noreferrer"
              color="secondary"
            >
              <PlayIcon class="h-3 w-3 flex-none fill-secondary-foreground" />
              <span class="ml-2">Watch demo</span>
            </Button>
          </div>
        </div>
      </div>
      <div
        class="mt-16 flex flex-col items-center justify-between border-t border-border pt-8 sm:mt-20 md:flex-row lg:mt-24"
      >
        <div class="flex space-x-6 md:order-2">
          {#each socials as social}
            <a href={social.link} target="_blank" rel="noopener noreferrer">
              <svelte:component
                this={social.logo}
                class="h-6 w-6 text-muted-foreground/70 transition-colors hover:text-muted-foreground"
              />
            </a>
          {/each}
        </div>
        <p class="mt-8 text-xs leading-5 text-muted-foreground md:order-1 md:mt-0">
          Copyright &copy; {new Date().getFullYear()} Clank. All rights reserved.
        </p>
      </div>
    </nav>
  </Container>
</footer>
