<script lang="ts">
  import { Badge } from "$lib/components/ui/badge";
  import { Button } from "$lib/components/ui/button";
  import * as Dialog from "$lib/components/ui/dialog";
  import * as DropdownMenu from "$lib/components/ui/dropdown-menu";
  import { Input } from "$lib/components/ui/input";
  import { Label } from "$lib/components/ui/label";
  import { Picture } from "$lib/components/ui/picture";
  import * as Select from "$lib/components/ui/select";
  import * as Table from "$lib/components/ui/table";
  import * as entities from "$lib/entities";
  import * as payloads from "$lib/payloads";
  import { api } from "$lib/services/api";
  import { invitations, user, users } from "$lib/stores";
  import dayjs from "$lib/utils/datetime";
  import { clone } from "$lib/utils/object";
  import { capitalize } from "$lib/utils/string";
  import Plus from "lucide-svelte/icons/plus";
  import Search from "lucide-svelte/icons/search";
  import { onDestroy } from "svelte";
  import { Render, Subscribe, createTable } from "svelte-headless-table";
  import { addColumnFilters, addSortBy } from "svelte-headless-table/plugins";
  import ArrowDown from "svelte-radix/ArrowDown.svelte";
  import ArrowUp from "svelte-radix/ArrowUp.svelte";
  import CaretSort from "svelte-radix/CaretSort.svelte";
  import DotsHorizontal from "svelte-radix/DotsHorizontal.svelte";
  import { toast } from "svelte-sonner";
  import { writable } from "svelte/store";

  $: canUpdate = $user.role === entities.UserRole.ADMIN;

  type Member = {
    type: "user" | "invitation";
  } & entities.User &
    entities.Invitation;
  let members = writable<Member[]>([]);
  // A hacky and not very beautiful way to sync users and invitations status
  const unsubscribe = members.subscribe((members) => {
    if (!members.length) return;

    members = clone(members);

    $users = $users.map((user) => {
      const member = members.find((member) => user.id === member.id);
      if (member) return member;
      user.leftAt = new Date();
      return user;
    });

    $invitations = members.filter((member) => member.type === "invitation");
  });

  $members = clone($users)
    .filter((user) => !user.leftAt)
    .map((user) => <Member>{ type: "user", ...user })
    .concat(clone($invitations).map((invitation) => <Member>{ type: "invitation", ...invitation }));

  const table = createTable(members, {
    sort: addSortBy({ toggleOrder: ["asc", "desc", undefined] }),
    filter: addColumnFilters(),
  });

  const columns = table.createColumns([
    table.column({
      id: "member",
      accessor: (member) => {
        if (member.type === "user") return member.name + ":" + member.email + ":" + member.id;
        return member.email + ":" + member.id;
      },
      header: "Member",
      plugins: {
        filter: {
          fn: ({ filterValue, value }) => value.toLowerCase().includes(filterValue.toLowerCase()),
        },
      },
    }),
    table.column({
      id: "role",
      accessor: "role",
      header: "Role",
    }),
    table.column({
      id: "actions",
      accessor: "id",
      header: "",
    }),
  ]);

  const { headerRows, pageRows, tableAttrs, tableBodyAttrs, pluginStates } = table.createViewModel(columns);
  const { filterValues } = pluginStates.filter;

  let updateUser = async (user: entities.User) => {
    const response = await api.put<payloads.PutUserRequest, payloads.PutUserResponse>(`/users/${user.id}`, {
      role: user.role,
    });

    $members = $members.map((member) =>
      member.id === user.id ? <Member>{ type: "user", ...payloads.toUser(response) } : member,
    );

    toast.success("Member updated successfully!");
  };

  let openRemoveMemberDialog: boolean = false;
  let removeMemberDialog: Member;
  let removeMember = async () => {
    if (removeMemberDialog.type === "user")
      await api.delete<payloads.DeleteUserResponse>(`/users/${removeMemberDialog.id}`);
    else await api.delete<payloads.DeleteInvitationResponse>(`/invitations/${removeMemberDialog.id}`);
    $members = $members.filter((member) => member.id !== removeMemberDialog.id);
    toast.success("Member removed successfully!");
    openRemoveMemberDialog = false;
  };

  let openInviteUserDialog: boolean = false;
  let invitedEmail: string;
  let selInvitedRole: { value: string };
  $: invitedRole = selInvitedRole?.value;
  let inviteUser = async () => {
    if (!invitedEmail || !invitedRole) return;

    const response = await api.post<payloads.PostInvitationRequest, payloads.PostInvitationResponse>("/invitations", {
      email: invitedEmail,
      role: invitedRole,
    });

    $members = $members.concat([<Member>{ type: "invitation", ...payloads.toInvitation(response) }]);
    toast.success("Member invited successfully!");
    openInviteUserDialog = false;
  };
  $: openInviteUserDialog,
    (() => {
      invitedEmail = undefined;
      selInvitedRole = undefined;
    })();

  onDestroy(() => {
    unsubscribe();
  });
</script>

<div class={"w-full " + ($$props.class || "")}>
  <div class="mb-4 flex items-center justify-between gap-4">
    <div class="relative">
      <Search class="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
      <Input class="w-[300px] pl-8" placeholder="Search member..." type="text" bind:value={$filterValues.member} />
    </div>
    <Button on:click={() => (openInviteUserDialog = true)} disabled={!canUpdate}>
      <Plus class="mr-2 h-4 w-4 stroke-[2.5]" />
      Invite
    </Button>
  </div>
  <div class="rounded-md border shadow-sm">
    <Table.Root {...$tableAttrs}>
      <Table.Header>
        {#each $headerRows as headerRow}
          <Subscribe rowAttrs={headerRow.attrs()}>
            <Table.Row class="bg-muted/50">
              {#each headerRow.cells as cell (cell.id)}
                <Subscribe attrs={cell.attrs()} let:attrs props={cell.props()} let:props>
                  <Table.Head {...attrs}>
                    {#if cell.id === "member"}
                      <Button variant="ghost" on:click={props.sort.toggle}>
                        <Render of={cell.render()} />
                        {#if props.sort.order === "asc"}
                          <ArrowUp class="ml-2 h-4 w-4" />
                        {:else if props.sort.order === "desc"}
                          <ArrowDown class="ml-2 h-4 w-4" />
                        {:else}
                          <CaretSort class="ml-2 h-4 w-4" />
                        {/if}
                      </Button>
                    {:else}
                      <Render of={cell.render()} />
                    {/if}
                  </Table.Head>
                </Subscribe>
              {/each}
            </Table.Row>
          </Subscribe>
        {/each}
      </Table.Header>
      <Table.Body {...$tableBodyAttrs}>
        {#if $pageRows.length}
          {#each $pageRows as row (row.id)}
            <Subscribe rowAttrs={row.attrs()} let:rowAttrs>
              <Table.Row {...rowAttrs}>
                {#each row.cells as cell (cell.id)}
                  <Subscribe attrs={cell.attrs()} let:attrs>
                    <Table.Cell {...attrs}>
                      {#if cell.id === "member"}
                        {#if cell.row.original.type === "user"}
                          <div class="ml-4 flex items-center gap-3">
                            <Picture size="xl" src={cell.row.original.picture} fallback={cell.row.original.name} />
                            <div class="flex flex-col items-start">
                              <span class="flex items-center gap-2">
                                {cell.row.original.name}
                                {#if cell.row.original.id === $user.id}
                                  <Badge variant="outline" class="flex h-6 shrink-0 items-center justify-center">
                                    You
                                  </Badge>
                                {/if}
                              </span>
                              <span class="text-xs">{cell.row.original.email}</span>
                            </div>
                          </div>
                        {:else}
                          <div class="ml-4 flex flex-col items-start">
                            <span class="flex items-center gap-3">
                              {cell.row.original.email}
                            </span>
                            <span class="text-xs">
                              {dayjs().isBefore(cell.row.original.expiresAt) ? "Invite pending" : "Invite expired"}
                            </span>
                          </div>
                        {/if}
                      {:else if cell.id === "role"}
                        <Select.Root
                          selected={{
                            value: cell.row.original.role,
                            label: capitalize(cell.row.original.role),
                          }}
                          onSelectedChange={(role) => {
                            updateUser({ ...cell.row.original, role: role?.value });
                          }}
                          disabled={!canUpdate || cell.row.original.type !== "user"}
                        >
                          <Select.Trigger class="w-40">
                            <Select.Value placeholder="Select a role" />
                          </Select.Trigger>
                          <Select.Content class="!mt-0">
                            {#each Object.values(entities.UserRole) as role}
                              <Select.Item value={role}>
                                {capitalize(role)}
                              </Select.Item>
                            {/each}
                          </Select.Content>
                        </Select.Root>
                      {:else if cell.id === "actions"}
                        <DropdownMenu.Root>
                          <DropdownMenu.Trigger asChild let:builder>
                            <Button
                              variant="ghost"
                              builders={[builder]}
                              class="ml-auto mr-4 flex h-8 w-8 p-0 data-[state=open]:bg-muted"
                            >
                              <DotsHorizontal class="h-4 w-4" />
                            </Button>
                          </DropdownMenu.Trigger>
                          <DropdownMenu.Content class="w-[160px]" align="end">
                            <DropdownMenu.Item
                              on:click={() => {
                                removeMemberDialog = cell.row.original;
                                openRemoveMemberDialog = true;
                              }}
                              disabled={!canUpdate}
                            >
                              Remove
                            </DropdownMenu.Item>
                          </DropdownMenu.Content>
                        </DropdownMenu.Root>
                      {:else}
                        <Render of={cell.render()} />
                      {/if}
                    </Table.Cell>
                  </Subscribe>
                {/each}
              </Table.Row>
            </Subscribe>
          {/each}
        {:else}
          <Table.Row>
            <Table.Cell colspan={columns.length} class="h-24 text-center">No members found.</Table.Cell>
          </Table.Row>
        {/if}
      </Table.Body>
    </Table.Root>
  </div>
</div>
<Dialog.Root bind:open={openInviteUserDialog}>
  <Dialog.Content>
    <Dialog.Header>
      <Dialog.Title>Invite a user</Dialog.Title>
      <Dialog.Description>Invite your teammates to easily collaborate with them.</Dialog.Description>
    </Dialog.Header>
    <div class="space-y-4">
      <div class="space-y-2">
        <Label>Email</Label>
        <Input bind:value={invitedEmail} type="email" maxlength={100} placeholder="alex@clank.so" required />
      </div>
      <div class="space-y-2">
        <Label>Role</Label>
        <Select.Root bind:selected={selInvitedRole} required>
          <Select.Trigger>
            <Select.Value placeholder="Select a role" />
          </Select.Trigger>
          <Select.Content class="!mt-0">
            {#each Object.values(entities.UserRole) as role}
              <Select.Item value={role}>
                {capitalize(role)}
              </Select.Item>
            {/each}
          </Select.Content>
        </Select.Root>
        <p class="text-[0.8rem] text-muted-foreground">
          Only administrators can manage the organization and its billing information.
        </p>
      </div>
    </div>
    <Dialog.Footer>
      <Button variant="outline" on:click={() => (openInviteUserDialog = false)}>Cancel</Button>
      <Button on:click={() => inviteUser()} disabled={!canUpdate}>Invite</Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
<Dialog.Root bind:open={openRemoveMemberDialog}>
  <Dialog.Content>
    <Dialog.Header>
      <Dialog.Title>Are you absolutely sure?</Dialog.Title>
      <p class="pt-2 text-sm text-muted-foreground">
        This will remove the member from your organization. You can always invite them again.
      </p>
    </Dialog.Header>
    <Dialog.Footer>
      <Button variant="outline" on:click={() => (openRemoveMemberDialog = false)}>Cancel</Button>
      <Button variant="destructive" on:click={() => removeMember()} disabled={!canUpdate}>Remove</Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
