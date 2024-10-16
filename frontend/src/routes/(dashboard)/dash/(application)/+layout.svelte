<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";
  import { CLANK_FRONTEND_PUBLIC_BASE_URL } from "$env/static/public";
  import ProductSwitcher from "$lib/components/cmp/ProductSwitcher.svelte";
  import Slack from "$lib/components/icons/Slack.svelte";
  import { Badge } from "$lib/components/ui/badge";
  import * as Breadcrumb from "$lib/components/ui/breadcrumb";
  import { Button } from "$lib/components/ui/button";
  import * as Card from "$lib/components/ui/card";
  import * as DropdownMenu from "$lib/components/ui/dropdown-menu";
  import { Input } from "$lib/components/ui/input";
  import { Picture } from "$lib/components/ui/picture";
  import { Skeleton } from "$lib/components/ui/skeleton";
  import * as Tooltip from "$lib/components/ui/tooltip";
  import * as entities from "$lib/entities";
  import * as payloads from "$lib/payloads";
  import { api } from "$lib/services/api";
  import { organization, product, products, user } from "$lib/stores";
  import dayjs from "$lib/utils/datetime";
  import { debounce } from "$lib/utils/delay";
  import { titlelize, trim } from "$lib/utils/string";
  import Bell from "lucide-svelte/icons/bell";
  import LayoutDashboard from "lucide-svelte/icons/layout-dashboard";
  import Lightbulb from "lucide-svelte/icons/lightbulb";
  import MessageSquare from "lucide-svelte/icons/message-square";
  import MessageSquareText from "lucide-svelte/icons/message-square-text";
  import Settings from "lucide-svelte/icons/settings";
  import ShieldAlert from "lucide-svelte/icons/shield-alert";
  import Sparkles from "lucide-svelte/icons/sparkles";
  import UserRound from "lucide-svelte/icons/user-round";
  import { onDestroy } from "svelte";

  const loadProducts = (async () => {
    const response = await api.get<payloads.ListProductsResponse>("/products");
    $products = response.products.map((product) => payloads.toProduct(product));
    if ($products) $product = $products.find((product) => product.id === $page.params.productID) || $products[0];
  })();

  let newIssues = 0;
  let newSuggestions = 0;
  let loadNewBadges = debounce(async (product: entities.Product) => {
    let query = payloads
      .toGetIssueCountMetricQuery({
        period_start_at: dayjs()
          .subtract(entities.ISSUE_NEW_MAX_DAYS - 1, "days")
          .toDate(),
        period_end_at: dayjs().toDate(),
      })
      .toString();
    const issuesResponse = await api.get<payloads.GetIssueCountMetricResponse>(
      `/products/${product.id}/metrics/issue-count?${query}`,
    );
    newIssues = payloads.toIssueCountMetric(issuesResponse).newIssues;

    query = payloads
      .toGetSuggestionCountMetricQuery({
        period_start_at: dayjs()
          .subtract(entities.SUGGESTION_NEW_MAX_DAYS - 1, "days")
          .toDate(),
        period_end_at: dayjs().toDate(),
      })
      .toString();
    const suggestionsResponse = await api.get<payloads.GetSuggestionCountMetricResponse>(
      `/products/${product.id}/metrics/suggestion-count?${query}`,
    );
    newSuggestions = payloads.toSuggestionCountMetric(suggestionsResponse).newSuggestions;
  }, 50);

  let getNewIssuesURL = (): string => {
    const params = new URLSearchParams();

    params.set("first_seen_end_at", dayjs().toISOString());
    params.set(
      "last_seen_start_at",
      dayjs()
        .subtract(entities.ISSUE_NEW_MAX_DAYS - 1, "days")
        .toISOString(),
    );
    params.set("_refresh", "");

    return `/dash/${$product.id}/issues?${params.toString()}`;
  };

  let getNewSuggestionsURL = (): string => {
    const params = new URLSearchParams();

    params.set("first_seen_end_at", dayjs().toISOString());
    params.set(
      "last_seen_start_at",
      dayjs()
        .subtract(entities.SUGGESTION_NEW_MAX_DAYS - 1, "days")
        .toISOString(),
    );
    params.set("_refresh", "");

    return `/dash/${$product.id}/suggestions?${params.toString()}`;
  };

  const unsubscribe = product.subscribe(async (product) => {
    if (!product) return;

    newIssues = 0;
    newSuggestions = 0;
    loadNewBadges(product);
  });

  let signOut = async () => {
    await api.post<payloads.PostSignOutRequest, payloads.PostSignOutResponse>("/signout");
    window.location.href = CLANK_FRONTEND_PUBLIC_BASE_URL + "/dash/signin";
  };

  $: (async () => {
    await loadProducts;
    if ($product) {
      if ($page.url.pathname === "/dash" || ($page.params.productID && $product.id !== $page.params.productID))
        await goto(`/dash/${$product.id}/overview`);
    } else {
      if ($page.params.productID && $page.url.pathname !== "/dash") await goto("/dash");
    }
  })();

  const classNavlinkInactive =
    "flex items-center gap-3 rounded-lg px-3 py-2 text-muted-foreground transition-all hover:text-primary h-10";
  const classNavlinkActive =
    "flex items-center gap-3 rounded-lg bg-muted px-3 py-2 text-primary transition-all hover:text-primary pointer-events-none h-10";

  let breadcrumb = [{ page: "Dashboard" }];
  $: (() => {
    const path = trim($page.url.pathname, "/").split("/");
    breadcrumb = [];
    for (let index = 0; index < path.length; index++) {
      if (path[index] === "dash" && path.length != 1) continue;
      if (path[index] === $page.params.productID) continue;

      let page = titlelize(path[index]);
      if (path[index] === "dash") page = "Dashboard";
      if (Object.values($page.params).includes(path[index])) page = path[index];

      let link = "/" + path.slice(0, index + 1).join("/");

      breadcrumb.push({ page: page, link: link });
    }
    breadcrumb[breadcrumb.length - 1].link = undefined;
  })();

  onDestroy(() => {
    unsubscribe();
  });
</script>

<div class="grid min-h-screen w-full grid-cols-[220px_1fr] lg:grid-cols-[280px_1fr]">
  <div class="relative block border-r bg-muted/40">
    <div class="sticky top-0 flex h-full max-h-screen flex-col gap-2 backdrop-blur-md">
      <div class="flex h-14 shrink-0 items-center border-b px-4 lg:h-[60px] lg:px-6">
        <ProductSwitcher />
      </div>
      <div class="flex-1">
        <nav class="grid items-start px-2 text-sm font-medium lg:px-4">
          <a
            href={$product ? `/dash/${$product.id}/overview` : "/dash"}
            class={$page.url.pathname.startsWith(`/dash/${$product?.id}/overview`)
              ? classNavlinkActive
              : classNavlinkInactive}
          >
            <LayoutDashboard class="h-4 w-4" />
            Overview
          </a>
          <a
            href={$product ? `/dash/${$product.id}/issues` : "/dash"}
            class={$page.url.pathname.startsWith(`/dash/${$product?.id}/issues`)
              ? classNavlinkActive
              : classNavlinkInactive}
          >
            <ShieldAlert class="h-4 w-4" />
            Issues
            {#if newIssues > 0}
              <Tooltip.Root openDelay={333}>
                <Tooltip.Trigger
                  class="pointer-events-auto ml-auto animate-[banner-appear_250ms_ease-in] cursor-pointer"
                >
                  <a href={$product ? getNewIssuesURL() : "/dash"}>
                    <Badge class="flex h-6 w-6 shrink-0 items-center justify-center rounded-full">
                      {Math.min(newIssues, 99)}
                    </Badge>
                  </a>
                </Tooltip.Trigger>
                <Tooltip.Content>
                  <p class="max-w-[300px] text-wrap">
                    {newIssues} new issues
                  </p>
                </Tooltip.Content>
              </Tooltip.Root>
            {/if}
          </a>
          <a
            href={$product ? `/dash/${$product.id}/suggestions` : "/dash"}
            class={$page.url.pathname.startsWith(`/dash/${$product?.id}/suggestions`)
              ? classNavlinkActive
              : classNavlinkInactive}
          >
            <Lightbulb class="h-4 w-4" />
            Suggestions
            {#if newSuggestions > 0}
              <Tooltip.Root openDelay={333}>
                <Tooltip.Trigger
                  class="pointer-events-auto ml-auto animate-[banner-appear_250ms_ease-in] cursor-pointer"
                >
                  <a href={$product ? getNewSuggestionsURL() : "/dash"}>
                    <Badge class="flex h-6 w-6 shrink-0 items-center justify-center rounded-full">
                      {Math.min(newSuggestions, 99)}
                    </Badge>
                  </a>
                </Tooltip.Trigger>
                <Tooltip.Content>
                  <p class="max-w-[300px] text-wrap">
                    {newSuggestions} new suggestions
                  </p>
                </Tooltip.Content>
              </Tooltip.Root>
            {/if}
          </a>
          <a
            href={$product ? `/dash/${$product.id}/feedbacks` : "/dash"}
            class={$page.url.pathname.startsWith(`/dash/${$product?.id}/feedbacks`)
              ? classNavlinkActive
              : classNavlinkInactive}
          >
            <MessageSquareText class="h-4 w-4" />
            Feedbacks
          </a>
          <a
            href={$product ? `/dash/${$product.id}/personas/coonm8ldnifcucu7qsu0` : "/dash"}
            class={$page.url.pathname.startsWith(`/dash/${$product?.id}/personas`)
              ? classNavlinkActive
              : classNavlinkInactive}
          >
            <UserRound class="h-4 w-4" />
            Personas
            <Badge variant="secondary" class="ml-auto flex h-6 shrink-0 items-center justify-center">Soon</Badge>
          </a>
          <a
            href={$product ? `/dash/${$product.id}/settings` : "/dash"}
            class={$page.url.pathname.startsWith(`/dash/${$product?.id}/settings`)
              ? classNavlinkActive
              : classNavlinkInactive}
          >
            <Settings class="h-4 w-4" />
            Settings
          </a>
        </nav>
      </div>
      {#if $organization.plan === entities.OrganizationPlan.STARTER}
        <div class="mt-auto p-4 pb-3">
          <Card.Root>
            <Card.Header class="p-4">
              <Card.Title>Upgrade to Business</Card.Title>
              <Card.Description>
                Unlock all features, add unlimited products and analyze up to 100k feedbacks a month.
              </Card.Description>
            </Card.Header>
            <Card.Content class="p-4 pt-0">
              <Button href="/#pricing" size="sm" class="w-full">Upgrade</Button>
            </Card.Content>
          </Card.Root>
        </div>
      {:else if $organization.plan === entities.OrganizationPlan.TRIAL}
        <div class="mt-auto p-4 pb-3">
          <Card.Root>
            <Card.Header class="p-4">
              <Card.Title>Free trial ends {dayjs($organization.trialEndsAt).fromNow()}</Card.Title>
              <Card.Description>Upgrade to Business plan and analyze up to 100k feedbacks a month.</Card.Description>
            </Card.Header>
            <Card.Content class="p-4 pt-0">
              <Button href="/#pricing" size="sm" class="w-full">Upgrade</Button>
            </Card.Content>
          </Card.Root>
        </div>
      {:else if $organization.plan === entities.OrganizationPlan.DEMO}
        <div class="mt-auto p-4 pb-3">
          <Card.Root>
            <Card.Header class="p-4">
              <Card.Title>Demo mode</Card.Title>
              <Card.Description>
                You can browse your existing products, but new feedbacks won't be analyzed.
              </Card.Description>
            </Card.Header>
            <Card.Content class="p-4 pt-0">
              <Button href="/#pricing" size="sm" class="w-full">Unlock</Button>
            </Card.Content>
          </Card.Root>
        </div>
      {/if}
      <div class="mt-auto border-t py-2">
        <nav class="grid items-start px-2 text-sm font-medium lg:px-4">
          <a
            href="https://join.slack.com/t/clank-so-community/shared_invite/zt-2lja8kgir-H2k99~E1tcahgiYAsxN~HA"
            target="_blank"
            rel="noopener noreferrer"
            class={classNavlinkInactive}
          >
            <Slack class="h-4 !w-4" />
            Community
          </a>
          <a href="mailto:support@clank.so" class={classNavlinkInactive}>
            <MessageSquare class="h-4 w-4" />
            Give Feedback
          </a>
          <a href="/blog" class={classNavlinkInactive}>
            <Bell class="h-4 w-4" />
            What's New
            <Badge variant="secondary" class="ml-auto flex h-6 w-6 shrink-0 items-center justify-center rounded-full"
              >1</Badge
            >
          </a>
        </nav>
      </div>
    </div>
  </div>
  <div class="relative flex flex-col">
    <header
      class="sticky top-0 z-40 flex h-14 shrink-0 items-center gap-6 border-b bg-muted/40 px-4 backdrop-blur-md lg:h-[60px] lg:px-6"
    >
      <div class="w-full flex-1">
        <Breadcrumb.Root>
          <Breadcrumb.List>
            {#each breadcrumb as crumb, index}
              {#if index > 0}
                <Breadcrumb.Separator />
              {/if}
              <Breadcrumb.Item>
                {#if crumb.link}
                  <Breadcrumb.Link href={crumb.link}>{crumb.page}</Breadcrumb.Link>
                {:else}
                  <Breadcrumb.Page>{crumb.page}</Breadcrumb.Page>
                {/if}
              </Breadcrumb.Item>
            {/each}
          </Breadcrumb.List>
        </Breadcrumb.Root>
      </div>
      <div class="relative w-[calc(330px+0.5rem)]">
        <Sparkles class="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
        <Badge variant="secondary" class="absolute right-2 top-1/2 h-6 shrink-0 -translate-y-1/2">Soon</Badge>
        <Input
          type="search"
          placeholder="Ask anything about your customers..."
          class="w-full !cursor-default appearance-none bg-background pl-8 shadow-none"
          disabled
        />
      </div>
      <DropdownMenu.Root>
        <DropdownMenu.Trigger asChild let:builder>
          <Button builders={[builder]} variant="outline" size="icon" class="rounded-full">
            <Picture size="2xl" src={$user.picture} fallback={$user.name} class="h-full w-full" />
          </Button>
        </DropdownMenu.Trigger>
        <DropdownMenu.Content align="end">
          <DropdownMenu.Label class="font-normal">
            <div class="flex flex-col space-y-1">
              <p class="text-sm font-medium leading-none">{$user.name}</p>
              <p class="flex items-center justify-start gap-1 text-xs leading-none text-muted-foreground">
                <Picture size="xs" src={$organization.picture} fallback={$organization.name} />
                {$organization.name}
              </p>
            </div>
          </DropdownMenu.Label>
          <DropdownMenu.Separator />
          <DropdownMenu.Item href="/dash/profile">Profile</DropdownMenu.Item>
          <DropdownMenu.Item href="/dash/organization">Organization</DropdownMenu.Item>
          <DropdownMenu.Item href="/dash/billing" class="pointer-events-none flex items-center justify-between gap-9">
            Billing
            <Badge variant="secondary" class="flex h-6 shrink-0 items-center justify-center">Soon</Badge>
          </DropdownMenu.Item>
          <DropdownMenu.Separator />
          <DropdownMenu.Item on:click={() => signOut()}>Sign out</DropdownMenu.Item>
        </DropdownMenu.Content>
      </DropdownMenu.Root>
    </header>
    <main class="flex flex-1 flex-col bg-background p-6">
      {#await loadProducts}
        <div class="flex flex-1 items-center justify-center space-x-4 rounded-lg border border-dashed shadow-sm">
          <Skeleton class="h-12 w-12 rounded-full" />
          <div class="space-y-2">
            <Skeleton class="h-4 w-[250px]" />
            <Skeleton class="h-4 w-[200px]" />
          </div>
        </div>
      {:then}
        <slot />
      {:catch}
        <div class="flex flex-1 items-center justify-center space-x-4 rounded-lg border border-dashed shadow-sm">
          <p>Something went wrong!</p>
        </div>
      {/await}
    </main>
  </div>
</div>
