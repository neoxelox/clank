<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";
  import { CLANK_FRONTEND_PUBLIC_BASE_URL } from "$env/static/public";
  import Amazon from "$lib/components/icons/Amazon.svelte";
  import Apple from "$lib/components/icons/Apple.svelte";
  import Google from "$lib/components/icons/Google.svelte";
  import SAML from "$lib/components/icons/SAML.svelte";
  import { Button } from "$lib/components/ui/button";
  import { Input } from "$lib/components/ui/input";
  import { Label } from "$lib/components/ui/label";
  import * as payloads from "$lib/payloads";
  import { api } from "$lib/services/api";
  import { signInCodeID, signInCodeState } from "$lib/stores";
  import LoaderCircle from "lucide-svelte/icons/loader-circle";
  import { toast } from "svelte-sonner";

  let loading = false;
  let redirectTo = $page.url.searchParams.get("redirect_to") || undefined;
  if (redirectTo) {
    redirectTo = CLANK_FRONTEND_PUBLIC_BASE_URL + redirectTo;
  }

  let email: string;

  let signInEmailStart = async () => {
    if (!email) return;

    const response = await api.post<payloads.PostSignInEmailStartRequest, payloads.PostSignInEmailStartResponse>(
      "/signin/email/start",
      { redirect_to: redirectTo, email: email },
    );

    $signInCodeState = response.sign_in_code_state;
    $signInCodeID = response.sign_in_code_id;

    await goto("/dash/signin/email/end");
  };

  let signInOAuthStart = async (provider: string) => {
    if (provider === "apple") {
      toast.warning("Sign in with Apple not available.");
      return;
    }

    const response = await api.post<payloads.PostSignInOAuthStartRequest, payloads.PostSignInOAuthStartResponse>(
      `/signin/${provider}/start`,
      { redirect_to: redirectTo },
    );

    window.location.href = response.auth_url;
  };

  let signInSAMLStart = async () => {
    toast.warning("Sign in with SAML not available for organization.");
  };

  // Try to sign in directly
  api.get<payloads.GetMeResponse>("/user").then(() => {
    if (redirectTo && new URL(redirectTo).origin == new URL(CLANK_FRONTEND_PUBLIC_BASE_URL).origin) {
      window.location.href = redirectTo;
    } else {
      window.location.href = CLANK_FRONTEND_PUBLIC_BASE_URL + "/dash";
    }
  });
</script>

<div class="mx-auto grid w-96 gap-12">
  <h1 class="text-center font-cal text-3xl text-foreground sm:text-4xl md:text-5xl">Welcome back</h1>
  <div class="grid gap-4">
    <div class="grid gap-2">
      <Label for="email">Email</Label>
      <Input id="email" bind:value={email} type="email" placeholder="email@example.com" required disabled={loading} />
      <p class="text-sm text-muted-foreground">Use an organization email to easily collaborate with your teammates.</p>
    </div>
    <Button
      class="w-full"
      on:click={() => (loading = true) && signInEmailStart().finally(() => (loading = false))}
      disabled={loading}
    >
      {#if loading}
        <LoaderCircle class="mr-2 h-4 w-4 animate-spin" />
      {/if}
      Continue
    </Button>
    <hr class="m-2 border-border" />
    <Button
      variant="outline"
      class="w-full"
      on:click={() => (loading = true) && signInOAuthStart("google").finally(() => (loading = false))}
      disabled={loading}
    >
      <Google class="mr-2 h-4 w-4" />
      Sign in with Google
    </Button>
    <Button
      variant="outline"
      class="w-full"
      on:click={() => (loading = true) && signInOAuthStart("apple").finally(() => (loading = false))}
      disabled={loading}
    >
      <Apple class="mr-2 h-4 w-4" />
      Sign in with Apple
    </Button>
    <Button
      variant="outline"
      class="w-full"
      on:click={() => (loading = true) && signInOAuthStart("amazon").finally(() => (loading = false))}
      disabled={loading}
    >
      <Amazon class="mr-2 h-4 w-4" />
      Sign in with Amazon
    </Button>
    <hr class="m-2 border-border" />
    <Button
      variant="outline"
      class="w-full"
      on:click={() => (loading = true) && signInSAMLStart().finally(() => (loading = false))}
      disabled={loading}
    >
      <SAML class="mr-2 h-4 w-4 stroke-foreground" />
      Sign in with SAML
    </Button>
    <p class="text-center text-sm text-muted-foreground">We will create an account if you don't have one.</p>
    <p class="text-center text-sm text-muted-foreground">
      By continuing, you acknowledge that you understand and agree to our <a
        href="/terms"
        target="_blank"
        class="text-primary">Terms and Conditions</a
      >
      and <a href="/privacy" target="_blank" class="text-primary">Privacy Policy</a>.
    </p>
  </div>
</div>
