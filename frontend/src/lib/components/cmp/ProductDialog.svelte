<script lang="ts">
  import { goto } from "$app/navigation";
  import { Button } from "$lib/components/ui/button";
  import * as Dialog from "$lib/components/ui/dialog";
  import { Input } from "$lib/components/ui/input";
  import { Label } from "$lib/components/ui/label";
  import * as Select from "$lib/components/ui/select";
  import * as entities from "$lib/entities";
  import * as payloads from "$lib/payloads";
  import { api } from "$lib/services/api";
  import { organization, product, products, user } from "$lib/stores";
  import { capitalize } from "$lib/utils/string";
  import { toast } from "svelte-sonner";

  export let open: boolean = false;

  // TODO: This should be backend-guided
  $: maxProducts =
    {
      [entities.OrganizationPlan.STARTER]: 5,
      [entities.OrganizationPlan.TRIAL]: 1,
      [entities.OrganizationPlan.DEMO]: 1,
    }[$organization.plan] || Number.MAX_SAFE_INTEGER;

  $: canUpdate = $user.role === entities.UserRole.ADMIN && $products?.length < maxProducts;

  let name: string;
  let selLanguage: { value: string };
  $: language = selLanguage?.value;

  let addProduct = async () => {
    if (!name || !language) return;

    const response = await api.post<payloads.PostProductRequest, payloads.PostProductResponse>("/products", {
      name: name,
      language: language,
    });

    $product = payloads.toProduct(response);
    $products = $products.concat([$product]);

    toast.success(`Product added successfully!`);
    setTimeout(async () => await goto(`/dash/${$product.id}/settings`), 0);

    open = false;
  };

  $: open,
    (() => {
      name = undefined;
      selLanguage = undefined;
    })();
</script>

<Dialog.Root bind:open>
  <Dialog.Content>
    <Dialog.Header>
      <Dialog.Title>Add a product</Dialog.Title>
      <Dialog.Description>Add a new product to recollect and analyze feedback automatically.</Dialog.Description>
    </Dialog.Header>
    <div class="space-y-4">
      {#if $organization.plan === entities.OrganizationPlan.STARTER || $organization.plan === entities.OrganizationPlan.TRIAL || $organization.plan === entities.OrganizationPlan.DEMO}
        <div
          class="mb-4 flex w-full flex-row items-center justify-center gap-x-4 gap-y-2 rounded-lg border border-border bg-muted/40 px-3.5 py-2.5 shadow-sm"
        >
          <h3 class="text-wrap text-sm font-semibold leading-none tracking-tight text-muted-foreground">
            Unlock unlimited products per organization
          </h3>
          <Button href="/#pricing" variant="outline" size="sm" class="rounded-full px-4 py-2.5 text-xs font-semibold">
            Upgrade now &rarr;
          </Button>
        </div>
        {#if $products?.length >= maxProducts}
          <h3 class="text-justify text-sm font-medium leading-none tracking-tight text-destructive">
            You have added the maximum products available for the {entities.OrganizationPlanDetails[$organization.plan]
              .title}
            plan!
          </h3>
        {/if}
      {/if}
      <div class="space-y-2">
        <Label>Name</Label>
        <Input
          bind:value={name}
          type="text"
          maxlength={100}
          placeholder="Awesome Maps"
          required
          disabled={!canUpdate}
        />
      </div>
      <div class="space-y-2">
        <Label>Language</Label>
        <Select.Root bind:selected={selLanguage} required disabled={!canUpdate}>
          <Select.Trigger>
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
        <p class="text-[0.8rem] text-muted-foreground">Feedbacks will be translated to this language for coherence.</p>
      </div>
    </div>
    <Dialog.Footer>
      <Button variant="outline" on:click={() => (open = false)}>Cancel</Button>
      <Button on:click={() => addProduct()} disabled={!canUpdate}>Add</Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
