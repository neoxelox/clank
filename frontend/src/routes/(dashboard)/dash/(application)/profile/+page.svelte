<script lang="ts">
  import { CLANK_FRONTEND_PUBLIC_BASE_URL } from "$env/static/public";
  import { Button } from "$lib/components/ui/button";
  import * as Dialog from "$lib/components/ui/dialog";
  import { Input } from "$lib/components/ui/input";
  import { Label } from "$lib/components/ui/label";
  import { Picture } from "$lib/components/ui/picture";
  import * as Tabs from "$lib/components/ui/tabs";
  import * as payloads from "$lib/payloads";
  import { api } from "$lib/services/api";
  import { user } from "$lib/stores";
  import { onDestroy } from "svelte";
  import { toast } from "svelte-sonner";

  let name: string;
  let picture: string;
  let email: string;
  const unsubscribe = user.subscribe((user) => {
    name = user.name;
    picture = user.picture;
    email = user.email;
  });

  let updateAccount = async () => {
    const updateName = name !== $user.name;
    const updatePicture = picture !== $user.picture;
    const updateEmail = email !== $user.email;

    if (!updateName && !updatePicture && !updateEmail) return;

    const response = await api.put<payloads.PutMeRequest, payloads.PutMeResponse>("/user", {
      name: updateName ? name : undefined,
      picture: updatePicture ? picture : undefined,
    });

    $user = payloads.toUser(response);

    toast.success("User updated successfully!");
  };

  let openDeleteAccountDialog: boolean = false;
  let deleteAccount = async () => {
    await api.delete<payloads.DeleteMeResponse>("/user");
    window.location.href = CLANK_FRONTEND_PUBLIC_BASE_URL + "/dash/signin";
    openDeleteAccountDialog = false;
  };

  onDestroy(() => {
    unsubscribe();
  });
</script>

<Tabs.Root class="w-full">
  <Tabs.List>
    <Tabs.Trigger value="general">General</Tabs.Trigger>
    <Tabs.Trigger value="security">Security</Tabs.Trigger>
  </Tabs.List>
  <Tabs.Content value="general">
    <div class="grid grid-cols-2 gap-12 pt-2">
      <div class="space-y-6">
        <div class="space-y-2">
          <Label>Name</Label>
          <Input bind:value={name} type="text" maxlength={100} placeholder="Alex Rodriguez" />
        </div>
        <div class="space-y-2">
          <Label>Email</Label>
          <Input bind:value={email} type="email" maxlength={100} placeholder="alex@clank.so" disabled />
          <p class="text-[0.8rem] text-muted-foreground">
            Updating the email is currently not supported. Contact <a
              href="mailto:support@clank.so"
              class="text-primary">our team</a
            > to change it.
          </p>
        </div>
        <Button on:click={() => updateAccount()}>Save</Button>
      </div>
      <div class="space-y-6">
        <div class="space-y-2">
          <Label>Picture</Label>
          <Picture size="4xl" src={picture} fallback={$user.name} />
        </div>
      </div>
    </div>
  </Tabs.Content>
  <Tabs.Content value="security">
    <div class="grid grid-cols-2 gap-12 pt-2">
      <div class="space-y-6">
        <div class="space-y-2">
          <Button variant="destructive" on:click={() => (openDeleteAccountDialog = true)}>Delete account</Button>
          <p class="text-[0.8rem] text-muted-foreground">
            Your organization will be deleted if there are no more members left. This action is irreversible.
          </p>
        </div>
      </div>
    </div>
  </Tabs.Content>
</Tabs.Root>
<Dialog.Root bind:open={openDeleteAccountDialog}>
  <Dialog.Content>
    <Dialog.Header>
      <Dialog.Title>Are you absolutely sure?</Dialog.Title>
      <p class="pt-2 text-sm text-muted-foreground">
        This action cannot be undone. This will permanently delete your account and your organization too if there are
        no more members left.
      </p>
    </Dialog.Header>
    <Dialog.Footer>
      <Button variant="outline" on:click={() => (openDeleteAccountDialog = false)}>Cancel</Button>
      <Button variant="destructive" on:click={() => deleteAccount()}>Delete</Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
