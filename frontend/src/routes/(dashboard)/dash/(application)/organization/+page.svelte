<script lang="ts">
  import { CLANK_FRONTEND_PUBLIC_BASE_URL } from "$env/static/public";
  import MemberTable from "$lib/components/cmp/MemberTable.svelte";
  import { Button } from "$lib/components/ui/button";
  import * as Dialog from "$lib/components/ui/dialog";
  import { Input } from "$lib/components/ui/input";
  import { Label } from "$lib/components/ui/label";
  import { Picture } from "$lib/components/ui/picture";
  import { Switch } from "$lib/components/ui/switch";
  import * as Tabs from "$lib/components/ui/tabs";
  import * as entities from "$lib/entities";
  import * as payloads from "$lib/payloads";
  import { api } from "$lib/services/api";
  import { organization, user } from "$lib/stores";
  import { clone } from "$lib/utils/object";
  import { onDestroy } from "svelte";
  import { toast } from "svelte-sonner";

  let name: string;
  let picture: string;
  let domain: string;
  let settings: entities.OrganizationSettings;
  const unsubscribe = organization.subscribe((organization) => {
    name = organization.name;
    picture = organization.picture;
    domain = organization.domain;
    settings = clone(organization.settings);
  });

  $: canUpdate = $user.role === entities.UserRole.ADMIN;

  let updateOrganizationProfile = async (): Promise<entities.Organization | undefined> => {
    const updateName = name !== $organization.name;
    const updatePicture = picture !== $organization.picture;
    const updateDomain = domain !== $organization.domain;

    if (!updateName && !updatePicture && !updateDomain) return;

    const response = await api.put<payloads.PutOrganizationRequest, payloads.PutOrganizationResponse>("/organization", {
      name: updateName ? name : undefined,
      picture: updatePicture ? picture : undefined,
    });

    return payloads.toOrganization(response);
  };

  let updateOrganizationSettings = async (): Promise<entities.OrganizationSettings | undefined> => {
    const updateDomainSignIn = settings.domainSignIn !== $organization.settings.domainSignIn;

    if (!updateDomainSignIn) return;

    const response = await api.put<payloads.PutOrganizationSettingsRequest, payloads.PutOrganizationSettingsResponse>(
      "/organization/settings",
      {
        domain_sign_in: updateDomainSignIn ? settings.domainSignIn : undefined,
      },
    );

    return payloads.toOrganizationSettings(response);
  };

  let updateOrganization = async () => {
    const updatedOrganization = await updateOrganizationProfile();
    const updatedOrganizationSettings = await updateOrganizationSettings();

    if (updatedOrganization) $organization = updatedOrganization;
    if (updatedOrganizationSettings) $organization.settings = updatedOrganizationSettings;
    if (updatedOrganization || updatedOrganizationSettings) toast.success("Organization updated successfully!");
  };

  let openDeleteOrganizationDialog: boolean = false;
  let deleteOrganization = async () => {
    await api.delete<payloads.DeleteOrganizationResponse>("/organization");
    window.location.href = CLANK_FRONTEND_PUBLIC_BASE_URL + "/dash/signin";
    openDeleteOrganizationDialog = false;
  };

  onDestroy(() => {
    unsubscribe();
  });
</script>

<Tabs.Root class="w-full">
  <Tabs.List>
    <Tabs.Trigger value="general">General</Tabs.Trigger>
    <Tabs.Trigger value="security">Security</Tabs.Trigger>
    <Tabs.Trigger value="members">Members</Tabs.Trigger>
  </Tabs.List>
  <Tabs.Content value="general">
    <div class="grid grid-cols-2 gap-12 pt-2">
      <div class="space-y-6">
        <div class="space-y-2">
          <Label>Name</Label>
          <Input bind:value={name} type="text" maxlength={100} placeholder="Awesome Maps" disabled={!canUpdate} />
        </div>
        <div class="space-y-2">
          <Label>Domain</Label>
          <Input bind:value={domain} type="text" maxlength={100} placeholder="awesome.maps" disabled />
          <p class="text-[0.8rem] text-muted-foreground">
            Updating the domain is currently not supported. Contact <a
              href="mailto:support@clank.so"
              class="text-primary">our team</a
            > to change it.
          </p>
        </div>
        <div class="space-y-2">
          <Label>Domain sign in</Label>
          <Switch
            bind:checked={settings.domainSignIn}
            disabled={!canUpdate || !settings.isDomainSignInSupported}
            class="block"
          />
          <p class="text-[0.8rem] text-muted-foreground">
            Everyone with a <b>@{domain}</b> email will be able to join your organization automatically.
          </p>
        </div>
        <Button on:click={() => updateOrganization()} disabled={!canUpdate}>Save</Button>
      </div>
      <div class="space-y-6">
        <div class="space-y-2">
          <Label>Picture</Label>
          <Picture size="4xl" src={picture} fallback={$organization.name} />
        </div>
      </div>
    </div>
  </Tabs.Content>
  <Tabs.Content value="security">
    <div class="grid grid-cols-2 gap-12 pt-2">
      <div class="space-y-6">
        <div class="space-y-2">
          <Button variant="destructive" on:click={() => (openDeleteOrganizationDialog = true)} disabled={!canUpdate}>
            Delete organization
          </Button>
          <p class="text-[0.8rem] text-muted-foreground">
            Your organization and all of its members will be deleted. This action is irreversible.
          </p>
        </div>
      </div>
    </div>
  </Tabs.Content>
  <Tabs.Content value="members">
    {#if $organization.plan === entities.OrganizationPlan.TRIAL || $organization.plan === entities.OrganizationPlan.DEMO}
      <div
        class="mt-4 flex w-full flex-row items-center justify-center gap-x-4 gap-y-2 rounded-lg border border-border bg-muted/40 px-3.5 py-2.5 shadow-sm"
      >
        <h3 class="text-wrap font-semibold leading-none tracking-tight text-muted-foreground">
          Unlock unlimited members per organization
        </h3>
        <Button href="/#pricing" variant="outline" size="sm" class="rounded-full px-4 py-2.5 text-sm font-semibold">
          Upgrade now &rarr;
        </Button>
      </div>
    {/if}
    <MemberTable class="mt-4" />
    <p class="mt-2 text-[0.8rem] text-muted-foreground">
      Members can browse and manage issues, suggestions, feedbacks and metrics. They cannot create nor manage products,
      the organization and other users, or its billing.
    </p>
  </Tabs.Content>
</Tabs.Root>
<Dialog.Root bind:open={openDeleteOrganizationDialog}>
  <Dialog.Content>
    <Dialog.Header>
      <Dialog.Title>Are you absolutely sure?</Dialog.Title>
      <p class="pt-2 text-sm text-muted-foreground">
        This action cannot be undone. This will permanently delete your organization and all of its members.
      </p>
    </Dialog.Header>
    <Dialog.Footer>
      <Button variant="outline" on:click={() => (openDeleteOrganizationDialog = false)}>Cancel</Button>
      <Button variant="destructive" on:click={() => deleteOrganization()} disabled={!canUpdate}>Delete</Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
