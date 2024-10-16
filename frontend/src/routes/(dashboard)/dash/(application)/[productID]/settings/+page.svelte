<script lang="ts">
  import { Badge } from "$lib/components/ui/badge";
  import { Button } from "$lib/components/ui/button";
  import * as Card from "$lib/components/ui/card";
  import * as Command from "$lib/components/ui/command";
  import * as Dialog from "$lib/components/ui/dialog";
  import { Input } from "$lib/components/ui/input";
  import { Label } from "$lib/components/ui/label";
  import { Picture } from "$lib/components/ui/picture";
  import * as Popover from "$lib/components/ui/popover";
  import * as Select from "$lib/components/ui/select";
  import { Separator } from "$lib/components/ui/separator";
  import { Skeleton } from "$lib/components/ui/skeleton";
  import * as Tabs from "$lib/components/ui/tabs";
  import { Textarea } from "$lib/components/ui/textarea";
  import * as entities from "$lib/entities";
  import * as payloads from "$lib/payloads";
  import { api } from "$lib/services/api";
  import { organization, product, products, user } from "$lib/stores";
  import { clone } from "$lib/utils/object";
  import { capitalize, random } from "$lib/utils/string";
  import ExternalLink from "lucide-svelte/icons/external-link";
  import Plus from "lucide-svelte/icons/plus";
  import X from "lucide-svelte/icons/x";
  import { onDestroy, tick } from "svelte";
  import CaretSort from "svelte-radix/CaretSort.svelte";
  import { toast } from "svelte-sonner";

  let name: string;
  let picture: string;
  let language: string;
  let context: string;
  let newCategory: string;
  let categories: string[];
  let release: string;
  const unsubscribe1 = product.subscribe((product) => {
    if (!product) return;

    name = product.name;
    picture = product.picture;
    language = product.language;
    context = product.context;
    categories = clone(product.categories).map((category) => category.replaceAll("_", " "));
    release = product.release !== entities.NO_RELEASE ? product.release : "";
  });

  $: canUpdate = $user.role === entities.UserRole.ADMIN;

  let updateProduct = async () => {
    const newRelease = release !== "" ? release : entities.NO_RELEASE;
    const newCategories = categories.map((category) => category.toUpperCase().replaceAll(" ", "_"));

    const updateName = name !== $product.name;
    const updatePicture = picture !== $product.picture;
    const updateLanguage = language !== $product.language;
    const updateContext = context !== $product.context;
    const updateCategories = newCategories.toString() !== $product.categories.toString();
    const updateRelease = newRelease !== $product.release;

    if (!updateName && !updatePicture && !updateLanguage && !updateContext && !updateCategories && !updateRelease)
      return;

    const response = await api.put<payloads.PutProductRequest, payloads.PutProductResponse>(
      `/products/${$product.id}`,
      {
        name: updateName ? name : undefined,
        picture: updatePicture ? picture : undefined,
        context: updateContext ? context : undefined,
        categories: updateCategories ? newCategories : undefined,
        release: updateRelease ? newRelease : undefined,
      },
    );

    $product = payloads.toProduct(response);
    $products = $products.map((p) => (p.id === $product.id ? $product : p));

    toast.success("Product updated successfully!");
  };

  let openDeleteProductDialog: boolean = false;
  let deleteProduct = async () => {
    await api.delete<payloads.DeleteProductResponse>(`/products/${$product.id}`);
    $products = $products.filter((p) => p.id !== $product.id);
    $product = undefined;
    toast.success("Product deleted successfully!");
    openDeleteProductDialog = false;
  };

  type Integration = {
    kind: "collector" | "exporter";
    mode: "create" | "update";
    checked: boolean;
  } & entities.Collector &
    entities.Exporter;
  let integrations: Integration[] = [];
  let integration: Integration;
  let loadIntegrations: Promise<void>;

  const unsubscribe2 = product.subscribe((product) => {
    if (!product) return;

    loadIntegrations = (async () => {
      const collectorsResponse = await api.get<payloads.ListCollectorsResponse>(`/products/${product.id}/collectors`);
      integrations = collectorsResponse.collectors.map(
        (collector) => <Integration>{ kind: "collector", mode: "update", ...payloads.toCollector(collector) },
      );
    })();
  });

  let openAddIntegrationPopover: boolean = false;
  let initIntegration = (kind: string, type: string): Integration => {
    if (kind === "collector") {
      switch (type) {
        case entities.CollectorType.TRUSTPILOT:
          return <Integration>{
            kind: "collector",
            mode: "create",
            id: random(20),
            productID: $product.id,
            type: type,
            settings: <entities.TrustpilotCollectorSettings>{
              domain: "",
            },
          };
        case entities.CollectorType.PLAY_STORE:
          return <Integration>{
            kind: "collector",
            mode: "create",
            id: random(20),
            productID: $product.id,
            type: type,
            settings: <entities.PlayStoreCollectorSettings>{
              appID: "",
            },
          };
        case entities.CollectorType.APP_STORE:
          return <Integration>{
            kind: "collector",
            mode: "create",
            id: random(20),
            productID: $product.id,
            type: type,
            settings: <entities.AppStoreCollectorSettings>{
              appID: "",
            },
          };
        case entities.CollectorType.AMAZON:
          return <Integration>{
            kind: "collector",
            mode: "create",
            id: random(20),
            productID: $product.id,
            type: type,
            settings: <entities.AmazonCollectorSettings>{
              asin: "",
            },
          };
        case entities.CollectorType.IAGORA:
          return <Integration>{
            kind: "collector",
            mode: "create",
            id: random(20),
            productID: $product.id,
            type: type,
            settings: <entities.IAgoraCollectorSettings>{
              institution: "",
            },
          };
      }
    }

    throw new Error(`Unknown collector or exporter of type ${type}`);
  };

  let checkIntegration = (): string | undefined => {
    if (integration.kind === "collector") {
      let base = entities.CollectorDetails[integration.type].base;
      let identifier: string;
      switch (integration.type) {
        case entities.CollectorType.TRUSTPILOT:
          identifier = integration.settings.domain;
          break;
        case entities.CollectorType.PLAY_STORE:
          identifier = integration.settings.appID;
          break;
        case entities.CollectorType.APP_STORE:
          identifier = integration.settings.appID;
          break;
        case entities.CollectorType.AMAZON:
          identifier = integration.settings.asin;
          break;
        case entities.CollectorType.IAGORA:
          identifier = integration.settings.institution;
          break;
        default:
          throw new Error(`Unknown collector of type ${integration.type}`);
      }

      if (base && identifier) return base + identifier;
      return;
    }

    throw new Error(`Unknown integration of kind ${integration.kind}`);
  };

  let openAddIntegrationDialog: boolean = false;
  let addIntegration = async () => {
    let newIntegration: Integration;

    if (integration.kind === "collector") {
      let request;
      switch (integration.type) {
        case entities.CollectorType.TRUSTPILOT:
          if (!integration.settings.domain) return;
          request = <payloads.PostTrustpilotCollectorRequest>{
            type: integration.type,
            domain: integration.settings.domain,
          };
          break;
        case entities.CollectorType.PLAY_STORE:
          if (!integration.settings.appID) return;
          request = <payloads.PostPlayStoreCollectorRequest>{
            type: integration.type,
            app_id: integration.settings.appID,
          };
          break;
        case entities.CollectorType.APP_STORE:
          if (!integration.settings.appID) return;
          request = <payloads.PostAppStoreCollectorRequest>{
            type: integration.type,
            app_id: integration.settings.appID,
          };
          break;
        case entities.CollectorType.AMAZON:
          if (!integration.settings.asin) return;
          request = <payloads.PostAmazonCollectorRequest>{
            type: integration.type,
            asin: integration.settings.asin,
          };
          break;
        case entities.CollectorType.IAGORA:
          if (!integration.settings.institution) return;
          request = <payloads.PostIAgoraCollectorRequest>{
            type: integration.type,
            institution: integration.settings.institution,
          };
          break;
        default:
          throw new Error(`Unknown collector of type ${integration.type}`);
      }

      const response = await api.post<typeof request, payloads.PostCollectorResponse>(
        `/products/${$product.id}/collectors`,
        request,
      );

      newIntegration = <Integration>{ kind: "collector", mode: "update", ...payloads.toCollector(response) };
    }

    integrations = integrations.map((i) => (i.id === integration.id ? newIntegration : i));
    toast.success("Integration added successfully!");
    openAddIntegrationDialog = false;
  };

  let openRemoveIntegrationDialog: boolean = false;
  let removeIntegration = async () => {
    if (integration.kind === "collector")
      await api.delete<payloads.DeleteCollectorResponse>(`/products/${$product.id}/collectors/${integration.id}`);

    integrations = integrations.filter((i) => i.id !== integration.id);
    toast.success("Integration removed successfully!");
    openRemoveIntegrationDialog = false;
  };

  onDestroy(() => {
    unsubscribe1();
    unsubscribe2();
  });
</script>

<Tabs.Root class="w-full">
  <Tabs.List>
    <Tabs.Trigger value="general">General</Tabs.Trigger>
    <Tabs.Trigger value="security">Security</Tabs.Trigger>
    <Tabs.Trigger value="integrations">Integrations</Tabs.Trigger>
  </Tabs.List>
  <Tabs.Content value="general">
    <div class="grid grid-cols-2 gap-12 pt-2">
      <div class="space-y-6">
        <div class="space-y-2">
          <Label>Name</Label>
          <Input bind:value={name} type="text" maxlength={100} placeholder="Awesome Maps" disabled={!canUpdate} />
        </div>
        <div class="space-y-2">
          <Label>Language</Label>
          <Select.Root
            selected={{
              value: language,
              label: capitalize(language),
            }}
            disabled
          >
            <Select.Trigger class="w-full">
              <Select.Value placeholder="Select a language" />
            </Select.Trigger>
            <Select.Content class="!mt-0">
              {#each Object.values(entities.ProductLanguage) as language}
                <Select.Item value={language}>
                  {capitalize(language)}
                </Select.Item>
              {/each}
            </Select.Content>
          </Select.Root>
          <p class="text-[0.8rem] text-muted-foreground">
            Feedbacks will be translated to this language for coherence. It cannot be updated once the product is
            created.
          </p>
        </div>
        <div class="space-y-2">
          <Label>Context</Label>
          <Textarea
            bind:value={context}
            maxlength={2500}
            indicator
            placeholder={`The product is called ${name}...`}
            class="min-h-28"
            disabled={!canUpdate}
          />
          <p class="text-[0.8rem] text-muted-foreground">
            Providing information about your product helps us identify more meaningful issues and suggestions,
            synthesize coherent customer personas and accurately measure vital metrics.
          </p>
        </div>
        <div class="space-y-2">
          <Label>Categories</Label>
          <div
            class={"flex h-auto max-h-28 w-full cursor-text flex-wrap gap-2 overflow-y-auto rounded-md border border-input bg-transparent p-3 shadow-sm " +
              (!canUpdate ? "!cursor-not-allowed !border-input/50" : "")}
          >
            {#each categories as category}
              <Badge
                variant="secondary"
                class={"flex h-6 shrink-0 select-text items-center justify-center font-medium " +
                  (!canUpdate ? "opacity-50" : "")}
              >
                {category}
                <Separator orientation="vertical" class="mx-2 h-4" />
                <Button
                  variant="ghost"
                  size="icon"
                  class="h-3 w-3"
                  on:click={() => {
                    categories = categories.filter((c) => c !== category);
                  }}
                  disabled={!canUpdate}
                >
                  <X class="h-3 w-3 shrink-0 opacity-50" />
                </Button>
              </Badge>
            {/each}
            {#if categories.length < 25}
              <Popover.Root open={!!newCategory} disableFocusTrap={true} openFocus={null}>
                <Popover.Trigger
                  class="focus-visible:outline-none"
                  on:click={(event) => (!newCategory ? event.preventDefault() : undefined)}
                  on:keydown={(event) => event.preventDefault()}
                >
                  <div class="relative">
                    <Plus
                      class={"absolute left-1 top-1 h-4 w-4 cursor-text stroke-2 text-muted-foreground " +
                        (!canUpdate ? "opacity-50" : "")}
                    />
                    <Input
                      bind:value={newCategory}
                      type="text"
                      maxlength={50}
                      placeholder="Add category"
                      class="h-6 w-32 shrink-0 rounded-md border border-none px-2.5 py-0 pl-6 align-middle text-xs font-medium leading-loose shadow-none focus-visible:ring-2"
                      disabled={!canUpdate}
                    />
                  </div>
                </Popover.Trigger>
                <Popover.Content class="w-auto p-0">
                  <Button
                    variant="ghost"
                    on:click={() => {
                      categories = categories.concat(newCategory.toUpperCase().replaceAll("_", " "));
                      newCategory = "";
                    }}
                    disabled={!canUpdate}
                  >
                    <Plus class="mr-2 h-4 w-4 stroke-[2.5]" />
                    Add "{newCategory}" as a category
                  </Button>
                </Popover.Content>
              </Popover.Root>
            {/if}
          </div>
          <p class="text-[0.8rem] text-muted-foreground">
            Feedbacks will be smartly classified into these categories to enhance organization. Already analyzed
            feedbacks won't be recategorized.
          </p>
        </div>
        <div class="space-y-2">
          <Label>Release</Label>
          <Input bind:value={release} type="text" maxlength={50} placeholder="1.31.7" disabled={!canUpdate} />
          <p class="text-[0.8rem] text-muted-foreground">
            Knowing the latest version of your product (if any) allows us to prioritize and address your recent issues
            and suggestions more effectively.
          </p>
        </div>
        <Button on:click={() => updateProduct()} disabled={!canUpdate}>Save</Button>
      </div>
      <div class="space-y-6">
        <div class="space-y-2">
          <Label>Picture</Label>
          <Picture size="4xl" src={picture} fallback={$product.name} />
        </div>
      </div>
    </div>
  </Tabs.Content>
  <Tabs.Content value="security">
    <div class="grid grid-cols-2 gap-12 pt-2">
      <div class="space-y-6">
        <div class="space-y-2">
          <Button variant="destructive" on:click={() => (openDeleteProductDialog = true)} disabled={!canUpdate}>
            Delete product
          </Button>
          <p class="text-[0.8rem] text-muted-foreground">
            This product and its issues, suggestions and feedbacks will be deleted. This action is irreversible.
          </p>
        </div>
      </div>
    </div>
  </Tabs.Content>
  <Tabs.Content value="integrations">
    <div class="space-y-4 pt-2">
      {#if $organization.plan === entities.OrganizationPlan.TRIAL || $organization.plan === entities.OrganizationPlan.DEMO}
        <div
          class="flex w-full flex-row items-center justify-center gap-x-4 gap-y-2 rounded-lg border border-border bg-muted/40 px-3.5 py-2.5 shadow-sm"
        >
          <h3 class="text-wrap font-semibold leading-none tracking-tight text-muted-foreground">
            Unlock unlimited integrations per organization
          </h3>
          <Button href="/#pricing" variant="outline" size="sm" class="rounded-full px-4 py-2.5 text-sm font-semibold">
            Upgrade now &rarr;
          </Button>
        </div>
      {/if}
      {#await loadIntegrations}
        <Card.Root class="rounded-lg border border-dashed shadow-sm">
          <Card.Content class="flex h-24 flex-1 items-center justify-center p-6">
            <div class="flex flex-1 items-center justify-center space-x-4">
              <Skeleton class="h-12 w-12 rounded-full" />
              <div class="space-y-2">
                <Skeleton class="h-4 w-[250px]" />
                <Skeleton class="h-4 w-[200px]" />
              </div>
            </div>
          </Card.Content>
        </Card.Root>
      {:then}
        {#if !loadIntegrations}
          <Card.Root class="rounded-lg border border-dashed shadow-sm">
            <Card.Content class="flex h-24 flex-1 items-center justify-center p-6">
              <div class="flex flex-1 items-center justify-center space-x-4">
                <Skeleton class="h-12 w-12 rounded-full" />
                <div class="space-y-2">
                  <Skeleton class="h-4 w-[250px]" />
                  <Skeleton class="h-4 w-[200px]" />
                </div>
              </div>
            </Card.Content>
          </Card.Root>
        {:else}
          <div class="grid gap-4">
            {#each integrations as i}
              {@const d = i.kind === "collector" ? entities.CollectorDetails[i.type] : entities.ExporterDetails[i.type]}
              <Card.Root class="rounded-lg shadow-sm">
                <Card.Header>
                  <Card.Title class="flex">
                    <svelte:component this={d.icon} class="mr-2 h-4 w-4" />
                    {d.title}
                  </Card.Title>
                  <Card.Description>{d.description}</Card.Description>
                </Card.Header>
                <Card.Content>
                  <div class="space-y-6">
                    {#if i.kind === "collector"}
                      {#if i.type === entities.CollectorType.TRUSTPILOT}
                        <div class="space-y-2">
                          <Label>Domain</Label>
                          <Input
                            bind:value={i.settings.domain}
                            type="text"
                            maxlength={100}
                            placeholder="awesomemaps.com"
                            disabled={!canUpdate || i.mode === "update"}
                          />
                          <p class="text-[0.8rem] text-muted-foreground">
                            The domain of your company's website as seen in the Trustpilot website. It cannot be updated
                            once the integration is added.
                          </p>
                        </div>
                      {:else if i.type === entities.CollectorType.PLAY_STORE}
                        <div class="space-y-2">
                          <Label>App ID</Label>
                          <Input
                            bind:value={i.settings.appID}
                            type="text"
                            maxlength={100}
                            placeholder="com.awesome.maps"
                            disabled={!canUpdate || i.mode === "update"}
                          />
                          <p class="text-[0.8rem] text-muted-foreground">
                            The ID of your Android App as seen in the Play Store. It cannot be updated once the
                            integration is added.
                          </p>
                        </div>
                      {:else if i.type === entities.CollectorType.APP_STORE}
                        <div class="space-y-2">
                          <Label>App ID</Label>
                          <Input
                            bind:value={i.settings.appID}
                            type="text"
                            maxlength={100}
                            placeholder="1248762770"
                            disabled={!canUpdate || i.mode === "update"}
                          />
                          <p class="text-[0.8rem] text-muted-foreground">
                            The ID of your iOS App as seen in the App Store (not the Bundle ID). It cannot be updated
                            once the integration is added.
                          </p>
                        </div>
                      {:else if i.type === entities.CollectorType.AMAZON}
                        <div class="space-y-2">
                          <Label>ASIN</Label>
                          <Input
                            bind:value={i.settings.asin}
                            type="text"
                            maxlength={100}
                            placeholder="B007HZO5WK"
                            disabled={!canUpdate || i.mode === "update"}
                          />
                          <p class="text-[0.8rem] text-muted-foreground">
                            The ASIN of your product as seen in the Amazon website. Size or color variations can have
                            their own ASINs, but it's best to use the ASIN of the main product to collect all reviews.
                            It cannot be updated once the integration is added.
                          </p>
                        </div>
                      {:else if i.type === entities.CollectorType.IAGORA}
                        <div class="space-y-2">
                          <Label>Institution</Label>
                          <Input
                            bind:value={i.settings.institution}
                            type="text"
                            maxlength={100}
                            placeholder="Massachusetts_Institute_of_Technology"
                            disabled={!canUpdate || i.mode === "update"}
                          />
                          <p class="text-[0.8rem] text-muted-foreground">
                            The name of your institution as seen in the iAgora website. It cannot be updated once the
                            integration is added.
                          </p>
                        </div>
                      {/if}
                    {/if}
                  </div>
                </Card.Content>
                <Card.Footer class="gap-2 border-t px-6 py-4">
                  {#if i.mode === "create"}
                    <Button
                      variant="outline"
                      on:click={() => (integrations = integrations.filter((ii) => ii.id !== i.id))}
                      disabled={!canUpdate}
                    >
                      Remove
                    </Button>
                    <Button
                      on:click={() => {
                        integration = i;
                        integration.checked = false;
                        openAddIntegrationDialog = true;
                      }}
                      disabled={!canUpdate}
                    >
                      Add
                    </Button>
                  {:else}
                    <Button
                      variant="outline"
                      on:click={() => {
                        integration = i;
                        openRemoveIntegrationDialog = true;
                      }}
                      disabled={!canUpdate}
                    >
                      Remove
                    </Button>
                    <!-- Comment save button in edit mode because confuses users after adding an integration -->
                    <!-- <Button disabled={!canUpdate}>Save</Button> -->
                  {/if}
                </Card.Footer>
              </Card.Root>
            {/each}
            <Card.Root class="rounded-lg border border-dashed shadow-sm">
              <Card.Content class="flex flex-1 items-center justify-center p-6">
                <Popover.Root bind:open={openAddIntegrationPopover} let:ids>
                  <Popover.Trigger asChild let:builder>
                    <Button
                      builders={[builder]}
                      variant="outline"
                      role="combobox"
                      class="w-[calc(280px-3rem)] justify-between"
                      disabled={!canUpdate}
                    >
                      <Plus class="mr-2 h-4 w-4 stroke-2 text-muted-foreground" />
                      <span class="text-muted-foreground">Add a new integration</span>
                      <CaretSort class="ml-auto h-4 w-4 shrink-0 opacity-50" />
                    </Button>
                  </Popover.Trigger>
                  <Popover.Content class="w-[calc(280px-3rem-1px)] p-0">
                    <Command.Root>
                      <Command.Input placeholder="Search integration..." />
                      <Command.List>
                        <Command.Empty>No integration found.</Command.Empty>
                        <Command.Group heading="Collectors">
                          {#each Object.values(entities.CollectorType) as type}
                            {@const details = entities.CollectorDetails[type]}
                            <Command.Item
                              onSelect={() => {
                                integrations = integrations.concat([initIntegration("collector", type)]);
                                openAddIntegrationPopover = false;
                                tick().then(() => document.getElementById(ids.trigger)?.focus());
                              }}
                              value={details.title}
                              class="text-sm"
                            >
                              <svelte:component this={details.icon} class="mr-2 h-3 w-3" />
                              {details.title}
                            </Command.Item>
                          {/each}
                        </Command.Group>
                        <Command.Group heading="Exporters" class="relative">
                          <Badge
                            variant="secondary"
                            class="absolute right-3 top-1/2 flex h-6 shrink-0 -translate-y-1/2 items-center justify-center"
                          >
                            Soon
                          </Badge>
                        </Command.Group>
                      </Command.List>
                    </Command.Root>
                  </Popover.Content>
                </Popover.Root>
              </Card.Content>
            </Card.Root>
          </div>
        {/if}
      {:catch}
        <Card.Root class="rounded-lg border border-dashed shadow-sm">
          <Card.Content class="flex h-24 flex-1 items-center justify-center p-6 text-center">
            <p>Something went wrong!</p>
          </Card.Content>
        </Card.Root>
      {/await}
    </div>
  </Tabs.Content>
</Tabs.Root>
<Dialog.Root bind:open={openDeleteProductDialog}>
  <Dialog.Content>
    <Dialog.Header>
      <Dialog.Title>Are you absolutely sure?</Dialog.Title>
      <p class="pt-2 text-sm text-muted-foreground">
        This action cannot be undone. This will permanently delete this product and its issues, suggestions and
        feedbacks.
      </p>
    </Dialog.Header>
    <Dialog.Footer>
      <Button variant="outline" on:click={() => (openDeleteProductDialog = false)}>Cancel</Button>
      <Button variant="destructive" on:click={() => deleteProduct()} disabled={!canUpdate}>Delete</Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
<Dialog.Root bind:open={openAddIntegrationDialog}>
  <Dialog.Content>
    <Dialog.Header>
      <Dialog.Title>Are you absolutely sure?</Dialog.Title>
      <p class="pt-2 text-sm text-muted-foreground">
        New feedbacks from this integration will start being collected or exported impacting your monthly usage. Monitor
        the organization's usage in the <a href="/dash/billing" target="_blank" class="pointer-events-none text-primary"
          >billing</a
        >
        dashboard.
        <br />
        <br />
        If applies, we will collect some existing feedbacks to kick things off (this may take some minutes).
      </p>
    </Dialog.Header>
    <Dialog.Footer>
      <Button variant="outline" on:click={() => (openAddIntegrationDialog = false)}>Cancel</Button>
      {#if !checkIntegration() || integration.checked}
        <Button on:click={() => addIntegration()} disabled={!canUpdate}>Add</Button>
      {:else}
        <Button
          href={checkIntegration()}
          target="_blank"
          rel="noopener noreferrer"
          on:click={() => (integration.checked = true)}
          disabled={!canUpdate}
        >
          Check <ExternalLink class="ml-1.5 h-4 w-4 shrink-0 stroke-[2.25]" />
        </Button>
      {/if}
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
<Dialog.Root bind:open={openRemoveIntegrationDialog}>
  <Dialog.Content>
    <Dialog.Header>
      <Dialog.Title>Are you absolutely sure?</Dialog.Title>
      <p class="pt-2 text-sm text-muted-foreground">
        This will remove this integration and new feedbacks from there will stop being collected or exported. You can
        always add it again.
      </p>
    </Dialog.Header>
    <Dialog.Footer>
      <Button variant="outline" on:click={() => (openRemoveIntegrationDialog = false)}>Cancel</Button>
      <Button variant="destructive" on:click={() => removeIntegration()} disabled={!canUpdate}>Remove</Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
