<script lang="ts">
  import { page } from "$app/stores";
  import FeedbackCarousel from "$lib/components/cmp/FeedbackCarousel.svelte";
  import { Badge } from "$lib/components/ui/badge";
  import { BarChart } from "$lib/components/ui/bar-chart";
  import { Button } from "$lib/components/ui/button";
  import * as Dialog from "$lib/components/ui/dialog";
  import { Label } from "$lib/components/ui/label";
  import { Picture } from "$lib/components/ui/picture";
  import * as Select from "$lib/components/ui/select";
  import { Skeleton } from "$lib/components/ui/skeleton";
  import * as Tooltip from "$lib/components/ui/tooltip";
  import * as entities from "$lib/entities";
  import * as payloads from "$lib/payloads";
  import { api } from "$lib/services/api";
  import { product, users } from "$lib/stores";
  import dayjs from "$lib/utils/datetime";
  import { simplify, titlelize } from "$lib/utils/string";
  import CircleUserRound from "lucide-svelte/icons/circle-user-round";
  import { copy } from "svelte-copy";
  import QuestionMarkCircled from "svelte-radix/QuestionMarkCircled.svelte";
  import { toast } from "svelte-sonner";

  let load: Promise<void>;
  let suggestion: entities.Suggestion;

  $: assignee = $users.find((user) => user.id === suggestion?.assigneeID);
  $: categories = Object.entries(suggestion?.categories || {}).map(([category, value]) => ({
    label: titlelize(category.replaceAll("_", " ")),
    value: value,
  }));
  $: importances = Object.entries(suggestion?.importances || {}).map(([importance, value]) => ({
    label: entities.SuggestionImportanceDetails[importance].title,
    value: value,
  }));
  $: releases = Object.entries(suggestion?.releases || {}).map(([release, value]) => ({
    label: release === entities.NO_RELEASE ? "Unknown" : release,
    value: value,
  }));
  $: sources = Object.entries(suggestion?.sources || {}).map(([source, value]) => ({
    label: entities.FeedbackSourceDetails[source].title,
    value: value,
  }));

  const loadSuggestion = async () => {
    load = (async () => {
      const response = await api.get<payloads.GetSuggestionResponse>(
        `/products/${$product.id}/suggestions/${$page.params.suggestionID}`,
      );
      suggestion = payloads.toSuggestion(response);
    })();
  };
  loadSuggestion();

  let openArchiveSuggestionDialog = false;
  let archiveSuggestion = async () => {
    await api.put<payloads.PutSuggestionArchivedRequest, payloads.PutSuggestionArchivedResponse>(
      `/products/${$product.id}/suggestions/${suggestion.id}/archived`,
      {
        archived: true,
      },
    );

    loadSuggestion();

    toast.success("Suggestion archived successfully!");
    openArchiveSuggestionDialog = false;
  };

  let openUnarchiveSuggestionDialog = false;
  let unarchiveSuggestion = async () => {
    await api.put<payloads.PutSuggestionArchivedRequest, payloads.PutSuggestionArchivedResponse>(
      `/products/${$product.id}/suggestions/${suggestion.id}/archived`,
      {
        archived: false,
      },
    );

    loadSuggestion();

    toast.success("Suggestion unarchived successfully!");
    openUnarchiveSuggestionDialog = false;
  };

  let assignSuggestions = async (assignee: entities.User) => {
    await api.put<payloads.PutSuggestionAssigneeRequest, payloads.PutSuggestionAssigneeResponse>(
      `/products/${$product.id}/suggestions/${suggestion.id}/assignee`,
      {
        assignee_id: assignee.id,
      },
    );

    loadSuggestion();

    toast.success(`Suggestion assigned to ${assignee.name} successfully!`);
  };

  let unassignSuggestions = async () => {
    await api.put<payloads.PutSuggestionAssigneeRequest, payloads.PutSuggestionAssigneeResponse>(
      `/products/${$product.id}/suggestions/${suggestion.id}/assignee`,
      {
        assignee_id: undefined,
      },
    );

    loadSuggestion();

    toast.success("Suggestion unassigned successfully!");
  };
</script>

{#await load}
  <div class="flex flex-1 items-center justify-center space-x-4 rounded-lg border border-dashed shadow-sm">
    <Skeleton class="h-12 w-12 rounded-full" />
    <div class="space-y-2">
      <Skeleton class="h-4 w-[250px]" />
      <Skeleton class="h-4 w-[200px]" />
    </div>
  </div>
{:then}
  <div class="flex flex-1 flex-col">
    <div class="flex items-start justify-between border-b pb-4">
      <div class="w-3/4 text-wrap">
        <h1 class="text-xl font-semibold">
          <span class="mr-1">{suggestion.title}</span>
          <h3
            class="inline-block cursor-default font-mono text-sm font-normal text-muted-foreground transition-colors hover:text-foreground"
            use:copy={suggestion.id}
          >
            #{suggestion.id}
          </h3>
        </h1>
        <div class="mt-1 space-x-2">
          {#if suggestion.archivedAt && dayjs(suggestion.archivedAt).isBefore(suggestion.lastSeenAt)}
            <Badge variant="destructive" class="shrink-0">REGRESSION</Badge>
          {/if}
          {#if dayjs(suggestion.firstSeenAt).isAfter(dayjs().subtract(entities.SUGGESTION_NEW_MAX_DAYS - 1, "days"))}
            <Badge class="shrink-0">NEW</Badge>
          {/if}
          {#if suggestion.archivedAt}
            <Badge variant="outline" class="shrink-0">COMPLETED</Badge>
          {:else}
            <Badge variant="outline" class="shrink-0">ONGOING</Badge>
          {/if}
        </div>
      </div>
      <div class="ml-auto flex items-center justify-start space-x-2">
        <Button
          variant="outline"
          on:click={() =>
            suggestion.archivedAt ? (openUnarchiveSuggestionDialog = true) : (openArchiveSuggestionDialog = true)}
        >
          {suggestion.archivedAt ? "Unarchive" : "Archive"}
        </Button>
      </div>
    </div>
    <div class="grid h-full grid-cols-3 gap-8">
      <div class="col-span-2 space-y-6 pt-6">
        <div class="space-y-1">
          <Label>Description</Label>
          <div class="whitespace-pre-line text-justify leading-relaxed">{suggestion.description}</div>
          {#if suggestion.reason.length > 0}
            <div class="!mt-1">
              <span
                class="rounded-md bg-secondary px-4 py-1 text-justify text-sm leading-loose text-secondary-foreground"
              >
                <Label class="mr-2 font-normal text-muted-foreground">Reason</Label>
                {suggestion.reason}
              </span>
            </div>
          {/if}
        </div>
        <div class="space-y-3">
          <Tooltip.Root openDelay={333}>
            <Tooltip.Trigger class="mb-4 flex cursor-default items-center">
              <Label>Feedbacks</Label>
              <QuestionMarkCircled class="ml-1 h-4 w-4 shrink-0 text-foreground/50" />
            </Tooltip.Trigger>
            <Tooltip.Content>
              <p class="max-w-[325px] text-wrap">
                This suggestion has been inferred from these feedbacks. The same feedback can describe multiple
                suggestions.
              </p>
            </Tooltip.Content>
          </Tooltip.Root>
          <FeedbackCarousel {suggestion} />
        </div>
      </div>
      <div class="space-y-6 border-l pl-6 pt-6">
        <div class="space-y-2">
          <Label>Customers</Label>
          <div class="text-xl">{suggestion.customers}</div>
        </div>
        <div class="space-y-2">
          <Label>Assigned to</Label>
          <Select.Root selected={{ value: assignee?.id }} required>
            <Select.Trigger>
              <div class="flex items-center justify-center">
                {#if assignee}
                  <Picture size="md" src={assignee.picture} fallback={assignee.name} class="mr-2" />
                  <span class="overflow-hidden text-ellipsis text-nowrap">
                    {assignee.leftAt ? `${assignee.name} (removed)` : assignee.name}
                  </span>
                {:else}
                  <CircleUserRound class="mr-2 h-5 w-5 stroke-[1.5]" />
                  <span class="overflow-hidden text-ellipsis text-nowrap">No one</span>
                {/if}
              </div>
            </Select.Trigger>
            <Select.Content class="!mt-0">
              <Select.Group>
                <Select.Label class="text-xs">Assign to</Select.Label>
                {#each $users as user (user.id)}
                  {#if !user.leftAt}
                    <Select.Item value={user.id} on:click={() => assignSuggestions(user)}>
                      <Picture size="sm" src={user.picture} fallback={user.name} class="mr-2" />
                      <span class="overflow-hidden text-ellipsis text-nowrap">
                        {user.name}
                      </span>
                    </Select.Item>
                  {/if}
                {/each}
                <Select.Item value={undefined} on:click={() => unassignSuggestions()}>
                  <CircleUserRound class="mr-2 h-4 w-4 text-muted-foreground" />
                  <span>No one</span>
                </Select.Item>
              </Select.Group>
            </Select.Content>
          </Select.Root>
        </div>
        <div class="space-y-2">
          <Label>Last seen</Label>
          <div>{dayjs(suggestion.lastSeenAt).toString()} ({simplify(dayjs(suggestion.lastSeenAt).fromNow())})</div>
        </div>
        <div class="space-y-2">
          <Label>First seen</Label>
          <div>
            {dayjs(suggestion.firstSeenAt).toString()} ({simplify(dayjs(suggestion.firstSeenAt).fromNow(true))} old)
          </div>
        </div>
        <div class="space-y-2">
          <Tooltip.Root openDelay={333}>
            <Tooltip.Trigger class="mb-3 flex cursor-default items-center">
              <Label>Categories</Label>
              <QuestionMarkCircled class="ml-1 h-4 w-4 shrink-0 text-foreground/50" />
            </Tooltip.Trigger>
            <Tooltip.Content>
              <p class="max-w-[300px] text-wrap">
                Feedbacks are smartly classified into categories to enhance organization.
                <a href={`/dash/${$product.id}/settings`} class="text-primary">Modify product categories.</a>
              </p>
            </Tooltip.Content>
          </Tooltip.Root>
          <BarChart open={true} values={categories} />
        </div>
        <div class="space-y-2">
          <Tooltip.Root openDelay={333}>
            <Tooltip.Trigger class="mb-3 flex cursor-default items-center">
              <Label>Importances</Label>
              <QuestionMarkCircled class="ml-1 h-4 w-4 shrink-0 text-foreground/50" />
            </Tooltip.Trigger>
            <Tooltip.Content>
              <p class="max-w-[300px] text-wrap">
                Interest is smartly inferred from the context of your product and how customers describe the suggestion.
              </p>
            </Tooltip.Content>
          </Tooltip.Root>
          <BarChart values={importances} />
        </div>
        <div class="space-y-2">
          <Label>Releases</Label>
          <BarChart values={releases} />
        </div>
        <div class="space-y-2">
          <Label>Sources</Label>
          <BarChart values={sources} />
        </div>
      </div>
    </div>
  </div>
{:catch}
  <div class="flex flex-1 items-center justify-center space-x-4 rounded-lg border border-dashed shadow-sm">
    Something went wrong!
  </div>
{/await}
<Dialog.Root bind:open={openArchiveSuggestionDialog}>
  <Dialog.Content>
    <Dialog.Header>
      <Dialog.Title>Are you absolutely sure?</Dialog.Title>
      <p class="pt-2 text-sm text-muted-foreground">
        This will mark the suggestion as completed and hide it. We will notify you of any regression. You can always
        unarchive it again.
      </p>
    </Dialog.Header>
    <Dialog.Footer>
      <Button variant="outline" on:click={() => (openArchiveSuggestionDialog = false)}>Cancel</Button>
      <Button on:click={() => archiveSuggestion()}>Archive</Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
<Dialog.Root bind:open={openUnarchiveSuggestionDialog}>
  <Dialog.Content>
    <Dialog.Header>
      <Dialog.Title>Are you absolutely sure?</Dialog.Title>
      <p class="pt-2 text-sm text-muted-foreground">
        This will mark the suggestion as ongoing. You can always archive it again.
      </p>
    </Dialog.Header>
    <Dialog.Footer>
      <Button variant="outline" on:click={() => (openUnarchiveSuggestionDialog = false)}>Cancel</Button>
      <Button on:click={() => unarchiveSuggestion()}>Unarchive</Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
