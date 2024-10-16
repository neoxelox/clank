<script lang="ts">
  import { page } from "$app/stores";
  import { CLANK_FRONTEND_PUBLIC_BASE_URL } from "$env/static/public";
  import { Button } from "$lib/components/ui/button";
  import { Label } from "$lib/components/ui/label";
  import { Pin } from "$lib/components/ui/pin";
  import { Progress } from "$lib/components/ui/progress";
  import * as entities from "$lib/entities";
  import * as payloads from "$lib/payloads";
  import { api } from "$lib/services/api";
  import { signInCodeID, signInCodeState } from "$lib/stores";
  import LoaderCircle from "lucide-svelte/icons/loader-circle";
  import Mail from "lucide-svelte/icons/mail";
  import { toast } from "svelte-sonner";

  let loading = false;
  let provider = $page.params.provider;

  let id = $page.url.searchParams.get("id") || $signInCodeID || undefined;
  let code = $page.url.searchParams.get("code") || undefined;
  let state = $page.url.searchParams.get("state") || $signInCodeState || undefined;

  let signInEnd = (redirectTo?: string) => {
    if (redirectTo && new URL(redirectTo).origin == new URL(CLANK_FRONTEND_PUBLIC_BASE_URL).origin) {
      window.location.href = redirectTo;
    } else {
      window.location.href = CLANK_FRONTEND_PUBLIC_BASE_URL + "/dash";
    }
  };

  let signInEmailEnd = async () => {
    if (!id || !code || !state) {
      const error = entities.ApiErrorMessage[entities.ApiErrorCode.ERR_UNAUTHORIZED];
      toast.error("Something went wrong", { description: error });
      throw new entities.ApiError(entities.ApiErrorCode.ERR_UNAUTHORIZED, error, 401);
    }

    const response = await api.post<payloads.PostSignInEmailEndRequest, payloads.PostSignInEmailEndResponse>(
      "/signin/email/end",
      {
        state: state,
        sign_in_code_id: id,
        sign_in_code_code: code,
      },
    );

    $signInCodeState = "";
    $signInCodeID = "";

    signInEnd(response.redirect_to);
  };

  let signInOAuthEnd = async (provider: string) => {
    if (provider === "apple") {
      toast.warning("Sign in with Apple not available.");
      return;
    }

    if (!state) {
      const error = entities.ApiErrorMessage[entities.ApiErrorCode.ERR_UNAUTHORIZED];
      toast.error("Something went wrong", { description: error });
      throw new entities.ApiError(entities.ApiErrorCode.ERR_UNAUTHORIZED, error, 401);
    }

    const response = await api.post<payloads.PostSignInOAuthEndRequest, payloads.PostSignInOAuthEndResponse>(
      `/signin/${provider}/end`,
      {
        state: state,
        auth_result: $page.url.searchParams.toString(),
      },
    );

    signInEnd(response.redirect_to);
  };

  let signInSAMLEnd = async () => {
    toast.warning("Sign in with SAML not available for organization.");
  };

  if (provider === "saml") {
    signInSAMLEnd();
  } else if (provider !== "email") {
    signInOAuthEnd(provider);
  }
</script>

{#if provider === "email"}
  <div class="mx-auto grid w-96 gap-12">
    <h1 class="text-center font-cal text-3xl text-foreground sm:text-4xl md:text-5xl">Welcome back</h1>
    <div class="grid gap-4">
      <div class="grid gap-2">
        <Label>Sign in Code</Label>
        <Pin
          bind:value={code}
          type="text"
          size={6}
          required
          on:finish={() => (loading = true) && signInEmailEnd().finally(() => (loading = false))}
          disabled={loading}
        />
        <p class="text-sm text-muted-foreground">
          Check your email for the sign in code or access through the magic link we have sent.
        </p>
      </div>
      <Button
        class="w-full"
        on:click={() => (loading = true) && signInEmailEnd().finally(() => (loading = false))}
        disabled={loading}
      >
        {#if loading}
          <LoaderCircle class="mr-2 h-4 w-4 animate-spin" />
        {:else}
          <Mail class="mr-2 h-4 w-4" />
        {/if}
        Sign in with Email
      </Button>
    </div>
  </div>
{:else if provider === "saml"}
  <div class="mx-auto grid w-96 gap-12">
    <h1 class="text-center font-cal text-3xl text-foreground sm:text-4xl md:text-5xl">Welcome back</h1>
    <Progress value={66} max={100} />
  </div>
{:else}
  <div class="mx-auto grid w-96 gap-12">
    <h1 class="text-center font-cal text-3xl text-foreground sm:text-4xl md:text-5xl">Welcome back</h1>
    <Progress value={66} max={100} />
  </div>
{/if}
