<script lang="ts">
  import { browser } from "$app/environment";
  import Button from "$lib/components/cms/Button.svelte";
  import { cookies } from "$lib/stores";
</script>

{#if browser && !$cookies}
  <div class="pointer-events-none fixed inset-x-0 bottom-0 z-50 px-6 pb-6">
    <div
      class="animate-banner-appear pointer-events-auto max-w-lg rounded-xl bg-background p-6 shadow-lg ring-1 ring-border"
    >
      <p class="text-sm leading-6 text-foreground">
        We use cookies in this website to give you the best experience on our site. By continuing, you acknowledge that
        you understand and agree to our
        <a href="/privacy" target="_blank" class="font-semibold text-primary">Privacy Policy</a>
        and
        <a href="/privacy#8-cookies-and-tracking-technologies" target="_blank" class="font-semibold text-primary"
          >Cookies Policy</a
        >.
      </p>
      <div class="mt-4 flex items-center gap-x-5">
        <Button color="primary" class="!px-3.5 !py-1" on:click={() => ($cookies = "accept")}>Accept all</Button>
        <Button color="secondary" class="!px-3.5 !py-1" on:click={() => ($cookies = "reject")}>Reject all</Button>
      </div>
    </div>
  </div>
{/if}
