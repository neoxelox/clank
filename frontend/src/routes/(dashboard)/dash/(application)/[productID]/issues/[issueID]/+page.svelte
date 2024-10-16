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
  import { Timeline } from "$lib/components/ui/timeline";
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
  let issue: entities.Issue;

  $: assignee = $users.find((user) => user.id === issue?.assigneeID);
  $: categories = Object.entries(issue?.categories || {}).map(([category, value]) => ({
    label: titlelize(category.replaceAll("_", " ")),
    value: value,
  }));
  $: severities = Object.entries(issue?.severities || {}).map(([severity, value]) => ({
    label: entities.IssueSeverityDetails[severity].title,
    value: value,
  }));
  $: releases = Object.entries(issue?.releases || {}).map(([release, value]) => ({
    label: release === entities.NO_RELEASE ? "Unknown" : release,
    value: value,
  }));
  $: sources = Object.entries(issue?.sources || {}).map(([source, value]) => ({
    label: entities.FeedbackSourceDetails[source].title,
    value: value,
  }));

  const loadIssue = async () => {
    load = (async () => {
      const response = await api.get<payloads.GetIssueResponse>(
        `/products/${$product.id}/issues/${$page.params.issueID}`,
      );
      issue = payloads.toIssue(response);
    })();
  };
  loadIssue();

  let openArchiveIssueDialog = false;
  let archiveIssue = async () => {
    await api.put<payloads.PutIssueArchivedRequest, payloads.PutIssueArchivedResponse>(
      `/products/${$product.id}/issues/${issue.id}/archived`,
      {
        archived: true,
      },
    );

    loadIssue();

    toast.success("Issue archived successfully!");
    openArchiveIssueDialog = false;
  };

  let openUnarchiveIssueDialog = false;
  let unarchiveIssue = async () => {
    await api.put<payloads.PutIssueArchivedRequest, payloads.PutIssueArchivedResponse>(
      `/products/${$product.id}/issues/${issue.id}/archived`,
      {
        archived: false,
      },
    );

    loadIssue();

    toast.success("Issue unarchived successfully!");
    openUnarchiveIssueDialog = false;
  };

  let assignIssues = async (assignee: entities.User) => {
    await api.put<payloads.PutIssueAssigneeRequest, payloads.PutIssueAssigneeResponse>(
      `/products/${$product.id}/issues/${issue.id}/assignee`,
      {
        assignee_id: assignee.id,
      },
    );

    loadIssue();

    toast.success(`Issue assigned to ${assignee.name} successfully!`);
  };

  let unassignIssues = async () => {
    await api.put<payloads.PutIssueAssigneeRequest, payloads.PutIssueAssigneeResponse>(
      `/products/${$product.id}/issues/${issue.id}/assignee`,
      {
        assignee_id: undefined,
      },
    );

    loadIssue();

    toast.success("Issue unassigned successfully!");
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
          <span class="mr-1">{issue.title}</span>
          <h3
            class="inline-block cursor-default font-mono text-sm font-normal text-muted-foreground transition-colors hover:text-foreground"
            use:copy={issue.id}
          >
            #{issue.id}
          </h3>
        </h1>
        <div class="mt-1 space-x-2">
          {#if issue.archivedAt && dayjs(issue.archivedAt).isBefore(issue.lastSeenAt)}
            <Badge variant="destructive" class="shrink-0">REGRESSION</Badge>
          {/if}
          {#if dayjs(issue.firstSeenAt).isAfter(dayjs().subtract(entities.ISSUE_NEW_MAX_DAYS - 1, "days"))}
            <Badge class="shrink-0">NEW</Badge>
          {/if}
          {#if issue.archivedAt}
            <Badge variant="outline" class="shrink-0">COMPLETED</Badge>
          {:else}
            <Badge variant="outline" class="shrink-0">ONGOING</Badge>
          {/if}
        </div>
      </div>
      <div class="ml-auto flex items-center justify-start space-x-2">
        <Button
          variant="outline"
          on:click={() => (issue.archivedAt ? (openUnarchiveIssueDialog = true) : (openArchiveIssueDialog = true))}
        >
          {issue.archivedAt ? "Unarchive" : "Archive"}
        </Button>
      </div>
    </div>
    <div class="grid h-full grid-cols-3 gap-8">
      <div class="col-span-2 space-y-6 pt-6">
        <div class="space-y-1">
          <Label>Description</Label>
          <div class="whitespace-pre-line text-justify leading-relaxed">{issue.description}</div>
        </div>
        {#if issue.steps.length > 0}
          <div class="!mb-7 !mt-5 space-y-3">
            <Label>Steps to reproduce</Label>
            <Timeline size="sm" steps={issue.steps} class="ml-1 gap-6" />
          </div>
        {/if}
        <div class="space-y-3">
          <Tooltip.Root openDelay={333}>
            <Tooltip.Trigger class="mb-4 flex cursor-default items-center">
              <Label>Feedbacks</Label>
              <QuestionMarkCircled class="ml-1 h-4 w-4 shrink-0 text-foreground/50" />
            </Tooltip.Trigger>
            <Tooltip.Content>
              <p class="max-w-[300px] text-wrap">
                This issue has been inferred from these feedbacks. The same feedback can describe multiple issues.
              </p>
            </Tooltip.Content>
          </Tooltip.Root>
          <FeedbackCarousel {issue} />
        </div>
      </div>
      <div class="space-y-6 border-l pl-6 pt-6">
        <div class="space-y-2">
          <Label>Customers</Label>
          <div class="text-xl">{issue.customers}</div>
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
                    <Select.Item value={user.id} on:click={() => assignIssues(user)}>
                      <Picture size="sm" src={user.picture} fallback={user.name} class="mr-2" />
                      <span class="overflow-hidden text-ellipsis text-nowrap">
                        {user.name}
                      </span>
                    </Select.Item>
                  {/if}
                {/each}
                <Select.Item value={undefined} on:click={() => unassignIssues()}>
                  <CircleUserRound class="mr-2 h-4 w-4 text-muted-foreground" />
                  <span>No one</span>
                </Select.Item>
              </Select.Group>
            </Select.Content>
          </Select.Root>
        </div>
        <div class="space-y-2">
          <Label>Last seen</Label>
          <div>{dayjs(issue.lastSeenAt).toString()} ({simplify(dayjs(issue.lastSeenAt).fromNow())})</div>
        </div>
        <div class="space-y-2">
          <Label>First seen</Label>
          <div>{dayjs(issue.firstSeenAt).toString()} ({simplify(dayjs(issue.firstSeenAt).fromNow(true))} old)</div>
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
              <Label>Severities</Label>
              <QuestionMarkCircled class="ml-1 h-4 w-4 shrink-0 text-foreground/50" />
            </Tooltip.Trigger>
            <Tooltip.Content>
              <p class="max-w-[300px] text-wrap">
                Severeness is smartly inferred from the context of your product and how customers describe the issue.
              </p>
            </Tooltip.Content>
          </Tooltip.Root>
          <BarChart values={severities} />
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
<Dialog.Root bind:open={openArchiveIssueDialog}>
  <Dialog.Content>
    <Dialog.Header>
      <Dialog.Title>Are you absolutely sure?</Dialog.Title>
      <p class="pt-2 text-sm text-muted-foreground">
        This will mark the issue as completed and hide it. We will notify you of any regression. You can always
        unarchive it again.
      </p>
    </Dialog.Header>
    <Dialog.Footer>
      <Button variant="outline" on:click={() => (openArchiveIssueDialog = false)}>Cancel</Button>
      <Button on:click={() => archiveIssue()}>Archive</Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
<Dialog.Root bind:open={openUnarchiveIssueDialog}>
  <Dialog.Content>
    <Dialog.Header>
      <Dialog.Title>Are you absolutely sure?</Dialog.Title>
      <p class="pt-2 text-sm text-muted-foreground">
        This will mark the issue as ongoing. You can always archive it again.
      </p>
    </Dialog.Header>
    <Dialog.Footer>
      <Button variant="outline" on:click={() => (openUnarchiveIssueDialog = false)}>Cancel</Button>
      <Button on:click={() => unarchiveIssue()}>Unarchive</Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
