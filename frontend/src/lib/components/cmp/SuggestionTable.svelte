<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";
  import { Badge } from "$lib/components/ui/badge";
  import { Button } from "$lib/components/ui/button";
  import { Checkbox } from "$lib/components/ui/checkbox";
  import { DatePicker } from "$lib/components/ui/date-picker";
  import * as Dialog from "$lib/components/ui/dialog";
  import * as DropdownMenu from "$lib/components/ui/dropdown-menu";
  import { Filter } from "$lib/components/ui/filter";
  import { Input } from "$lib/components/ui/input";
  import { Picture } from "$lib/components/ui/picture";
  import * as Select from "$lib/components/ui/select";
  import { Skeleton } from "$lib/components/ui/skeleton";
  import * as Table from "$lib/components/ui/table";
  import * as Tabs from "$lib/components/ui/tabs";
  import * as Tooltip from "$lib/components/ui/tooltip";
  import * as entities from "$lib/entities";
  import * as payloads from "$lib/payloads";
  import { api } from "$lib/services/api";
  import { product, users } from "$lib/stores";
  import dayjs from "$lib/utils/datetime";
  import { debounce } from "$lib/utils/delay";
  import { clone, compare, major } from "$lib/utils/object";
  import { simplify, titlelize, unitize } from "$lib/utils/string";
  import CircleUserRound from "lucide-svelte/icons/circle-user-round";
  import Search from "lucide-svelte/icons/search";
  import { onDestroy } from "svelte";
  import { Render, Subscribe, createTable } from "svelte-headless-table";
  import ArrowDown from "svelte-radix/ArrowDown.svelte";
  import ArrowUp from "svelte-radix/ArrowUp.svelte";
  import CaretSort from "svelte-radix/CaretSort.svelte";
  import ChevronLeft from "svelte-radix/ChevronLeft.svelte";
  import ChevronRight from "svelte-radix/ChevronRight.svelte";
  import Cross2 from "svelte-radix/Cross2.svelte";
  import DoubleArrowLeft from "svelte-radix/DoubleArrowLeft.svelte";
  import DoubleArrowRight from "svelte-radix/DoubleArrowRight.svelte";
  import { toast } from "svelte-sonner";
  import { writable } from "svelte/store";

  let defaults = (request: payloads.ListSuggestionsRequest): payloads.ListSuggestionsRequest => {
    request.filters.status = request.filters.status || payloads.SuggestionFilterStatus.ACTIVE;
    request.filters.last_seen_start_at = request.filters.last_seen_start_at || undefined;
    request.filters.first_seen_end_at = request.filters.first_seen_end_at || new Date();
    request.orders.relevance = request.orders.relevance || "DESC";
    request.pagination.limit = request.pagination.limit || 50;

    return request;
  };

  let previousRequest: payloads.ListSuggestionsRequest;
  let request = writable<payloads.ListSuggestionsRequest>(
    defaults(payloads.fromListSuggestionsQuery($page.url.searchParams)),
  );
  let load: Promise<void>;
  let suggestions = writable<entities.Suggestion[]>([]);
  let pages: (string | undefined)[] = [undefined];
  let selected: string[] = [];

  const sourceOptions = Object.values(entities.FeedbackSource).map((source) => {
    const details = entities.FeedbackSourceDetails[source];
    return {
      value: source,
      label: details.title,
      icon: details.icon,
    };
  });
  const importanceOptions = Object.values(entities.SuggestionImportance).map((importance) => {
    const details = entities.SuggestionImportanceDetails[importance];
    return {
      value: importance,
      label: details.title,
      icon: details.icon,
    };
  });
  const releaseOptions =
    $product.release !== entities.NO_RELEASE
      ? [
          {
            value: $product.release,
            label: $product.release,
          },
        ]
      : [];
  releaseOptions.push({
    value: entities.NO_RELEASE,
    label: "Unknown",
  });
  const categoryOptions = $product.categories.map((category) => {
    return {
      value: category,
      label: titlelize(category.replaceAll("_", " ")),
    };
  });
  categoryOptions.push({
    value: entities.NO_CATEGORY,
    label: "Unknown",
  });
  const assigneeOptions = $users.map((user) => {
    return {
      value: user.id,
      label: user.leftAt ? `${user.name} (removed)` : user.name,
      icon: user.picture,
    };
  });
  assigneeOptions.push({
    value: payloads.SUGGESTION_FILTER_UNASSIGNED,
    label: "Unassigned",
    icon: CircleUserRound,
  });
  const orderOptions = ["ASC", "DESC", undefined];
  let currentPage = 0;

  let loadSuggestions = debounce((params: string) => {
    load = (async () => {
      const response = await api.get<payloads.ListSuggestionsResponse>(
        `/products/${$product.id}/suggestions?${params}`,
      );
      $suggestions = response.suggestions.map((suggestion) => payloads.toSuggestion(suggestion));
      const page = response.next ? decodeURIComponent(response.next) : undefined;
      if (!pages.includes(page)) pages = pages.concat([page]);
    })();
  }, 500);

  $: if ($page.url.searchParams.has("_refresh"))
    $request = defaults(payloads.fromListSuggestionsQuery($page.url.searchParams));

  let resetQuery = () => {
    previousRequest = undefined;
    $request = defaults(payloads.fromListSuggestionsQuery(new URLSearchParams()));
  };

  const unsubscribe = request.subscribe(async (request) => {
    if (compare(request, previousRequest)) return;

    selected = [];

    if (!previousRequest || request.pagination.from === previousRequest.pagination.from) {
      pages = [undefined];
      currentPage = 0;
      request.pagination.from = undefined;
    }

    previousRequest = clone(request);

    const params = payloads.toListSuggestionsQuery(request).toString();
    await goto(`?${params}`, { noScroll: true, keepFocus: true });
    loadSuggestions(params);
  });

  const table = createTable(suggestions);

  const columns = table.createColumns([
    table.column({
      id: "suggestion",
      accessor: "title",
      header: "",
    }),
    table.column({
      id: "importance",
      accessor: "importances",
      header: "Importance",
    }),
    table.column({
      id: "seen",
      accessor: "lastSeenAt",
      header: "Seen at",
    }),
    table.column({
      id: "customers",
      accessor: "customers",
      header: "Customers",
    }),
    table.column({
      id: "assignee",
      accessor: "assigneeID",
      header: "Assignee",
    }),
  ]);

  const { headerRows, pageRows, tableAttrs, tableBodyAttrs } = table.createViewModel(columns);

  let openArchiveSuggestionsDialog = false;
  let archiveSuggestions = async () => {
    for (const id of selected) {
      await api.put<payloads.PutSuggestionArchivedRequest, payloads.PutSuggestionArchivedResponse>(
        `/products/${$product.id}/suggestions/${id}/archived`,
        {
          archived: true,
        },
      );
    }

    $request.reload = !$request.reload;

    toast.success("Suggestions archived successfully!");
    openArchiveSuggestionsDialog = false;
  };

  let openUnarchiveSuggestionsDialog = false;
  let unarchiveSuggestions = async () => {
    for (const id of selected) {
      await api.put<payloads.PutSuggestionArchivedRequest, payloads.PutSuggestionArchivedResponse>(
        `/products/${$product.id}/suggestions/${id}/archived`,
        {
          archived: false,
        },
      );
    }

    $request.reload = !$request.reload;

    toast.success("Suggestions unarchived successfully!");
    openUnarchiveSuggestionsDialog = false;
  };

  let assignSuggestions = async (assignee: entities.User) => {
    for (const id of selected) {
      await api.put<payloads.PutSuggestionAssigneeRequest, payloads.PutSuggestionAssigneeResponse>(
        `/products/${$product.id}/suggestions/${id}/assignee`,
        {
          assignee_id: assignee.id,
        },
      );
    }

    $request.reload = !$request.reload;

    toast.success(`Suggestions assigned to ${assignee.name} successfully!`);
  };

  let unassignSuggestions = async () => {
    for (const id of selected) {
      await api.put<payloads.PutSuggestionAssigneeRequest, payloads.PutSuggestionAssigneeResponse>(
        `/products/${$product.id}/suggestions/${id}/assignee`,
        {
          assignee_id: undefined,
        },
      );
    }

    $request.reload = !$request.reload;

    toast.success("Suggestions unassigned successfully!");
  };

  onDestroy(() => {
    unsubscribe();
  });
</script>

<div class={"w-full " + ($$props.class || "")}>
  <div class="mb-4 flex flex-wrap items-center justify-start gap-4">
    <Tabs.Root class="w-auto" bind:value={$request.filters.status}>
      <Tabs.List>
        <Tabs.Trigger value={payloads.SuggestionFilterStatus.ACTIVE}>Active</Tabs.Trigger>
        <Tabs.Trigger value={payloads.SuggestionFilterStatus.REGRESSED}>Regressed</Tabs.Trigger>
        <Tabs.Trigger value={payloads.SuggestionFilterStatus.ARCHIVED}>Archived</Tabs.Trigger>
      </Tabs.List>
    </Tabs.Root>
    <div class="relative">
      <Search class="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
      <Input
        class="w-[300px] pl-8"
        placeholder="Search content..."
        type="text"
        maxlength={250}
        bind:value={$request.filters.content}
      />
    </div>
    {#if $page.url.searchParams.toString().length > 1}
      <Button variant="ghost" on:click={() => resetQuery()}>
        Reset
        <Cross2 class="ml-2 h-4 w-4" />
      </Button>
    {/if}
    <DatePicker
      class="ml-auto"
      placeholder="Seen at"
      maxValue={new Date()}
      bind:startValue={$request.filters.last_seen_start_at}
      bind:endValue={$request.filters.first_seen_end_at}
    />
  </div>
  <div class="mb-4 flex flex-wrap items-center justify-start gap-4">
    <Filter title="Sources" options={sourceOptions} bind:values={$request.filters.sources} />
    <Filter title="Importances" options={importanceOptions} bind:values={$request.filters.importances} />
    <Filter title="Releases" options={releaseOptions} addable bind:values={$request.filters.releases} />
    <Filter title="Categories" options={categoryOptions} addable bind:values={$request.filters.categories} />
    <Filter title="Assignees" options={assigneeOptions} bind:values={$request.filters.assignees} />
    <Button
      variant="outline"
      class="ml-auto"
      on:click={() =>
        ($request.orders.relevance = orderOptions[(orderOptions.indexOf($request.orders.relevance) + 1) % 3])}
    >
      Relevance
      {#if $request.orders.relevance === "ASC"}
        <ArrowUp class="ml-2 h-4 w-4" />
      {:else if $request.orders.relevance === "DESC"}
        <ArrowDown class="ml-2 h-4 w-4" />
      {:else}
        <CaretSort class="ml-2 h-4 w-4" />
      {/if}
    </Button>
  </div>
  <div class="rounded-lg border shadow-sm">
    <Table.Root {...$tableAttrs}>
      <Table.Header>
        {#each $headerRows as headerRow}
          <Subscribe rowAttrs={headerRow.attrs()}>
            <Table.Row class="bg-muted/50">
              {#each headerRow.cells as cell (cell.id)}
                <Subscribe attrs={cell.attrs()} let:attrs props={cell.props()} let:props>
                  <Table.Head {...attrs}>
                    {#if cell.id === "suggestion"}
                      {@const checked =
                        $suggestions.length > 0
                          ? selected.length === $suggestions.length
                            ? true
                            : selected.length > 0
                              ? "indeterminate"
                              : false
                          : false}
                      {@const unarchiving =
                        $request.filters.status === payloads.SuggestionFilterStatus.REGRESSED ||
                        $request.filters.status === payloads.SuggestionFilterStatus.ARCHIVED}
                      <div class="flex items-center justify-start space-x-2">
                        <Checkbox
                          class="ml-2 mr-1"
                          {checked}
                          onCheckedChange={(change) => {
                            if (change && checked !== "indeterminate")
                              selected = $suggestions.map((suggestion) => suggestion.id);
                            else selected = [];
                          }}
                        />
                        <Button
                          variant="outline"
                          class="h-6 px-2.5"
                          on:click={() =>
                            unarchiving
                              ? (openUnarchiveSuggestionsDialog = true)
                              : (openArchiveSuggestionsDialog = true)}
                          disabled={selected.length === 0}
                        >
                          {#if unarchiving}
                            Unarchive
                          {:else}
                            Archive
                          {/if}
                        </Button>
                        <DropdownMenu.Root>
                          <DropdownMenu.Trigger asChild let:builder>
                            <Button
                              variant="outline"
                              builders={[builder]}
                              class="h-6 px-2.5"
                              disabled={selected.length === 0}
                            >
                              Assign
                              <CaretSort class="ml-2 h-4 w-4" />
                            </Button>
                          </DropdownMenu.Trigger>
                          <DropdownMenu.Content class="w-[160px]" align="end">
                            <DropdownMenu.Label class="text-xs">Assign to</DropdownMenu.Label>
                            <DropdownMenu.Separator />
                            {#each $users as user (user.id)}
                              {#if !user.leftAt}
                                <DropdownMenu.Item on:click={() => assignSuggestions(user)}>
                                  <Picture size="sm" src={user.picture} fallback={user.name} class="mr-2" />
                                  <span class="overflow-hidden text-ellipsis text-nowrap">
                                    {user.name}
                                  </span>
                                </DropdownMenu.Item>
                              {/if}
                            {/each}
                            <DropdownMenu.Item on:click={() => unassignSuggestions()}>
                              <CircleUserRound class="mr-2 h-4 w-4 text-muted-foreground" />
                              <span>No one</span>
                            </DropdownMenu.Item>
                          </DropdownMenu.Content>
                        </DropdownMenu.Root>
                      </div>
                    {:else if cell.id === "importance"}
                      <Tooltip.Root openDelay={333}>
                        <Tooltip.Trigger class="cursor-text">
                          <Render of={cell.render()} />
                        </Tooltip.Trigger>
                        <Tooltip.Content>
                          <p class="max-w-[310px] text-wrap">
                            Interest is smartly inferred from the context of your product and how customers describe the
                            suggestion.
                          </p>
                        </Tooltip.Content>
                      </Tooltip.Root>
                    {:else if cell.id === "customers"}
                      <span class="flex justify-end">
                        <Render of={cell.render()} />
                      </span>
                    {:else if cell.id === "assignee"}
                      <span class="mr-2 flex justify-end">
                        <Render of={cell.render()} />
                      </span>
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
        {#await load}
          <Table.Row>
            <Table.Cell colspan={columns.length} class="h-24 text-center">
              <div class="flex flex-1 items-center justify-center space-x-4">
                <Skeleton class="h-12 w-12 rounded-full" />
                <div class="space-y-2">
                  <Skeleton class="h-4 w-[250px]" />
                  <Skeleton class="h-4 w-[200px]" />
                </div>
              </div>
            </Table.Cell>
          </Table.Row>
        {:then}
          {#if $pageRows.length}
            {#each $pageRows as row (row.id)}
              <Subscribe rowAttrs={row.attrs()} let:rowAttrs>
                <Table.Row {...rowAttrs}>
                  {#each row.cells as cell (cell.id)}
                    {@const suggestion = cell.row.original}
                    <Subscribe attrs={cell.attrs()} let:attrs>
                      <Table.Cell {...attrs}>
                        {#if cell.id === "suggestion"}
                          {@const included = selected.includes(suggestion.id)}
                          {@const category = major(suggestion.categories).replaceAll("_", " ")}
                          <div class="flex shrink-0 flex-nowrap items-center space-x-2">
                            <Checkbox
                              class="ml-2 mr-1"
                              checked={included}
                              onCheckedChange={(checked) => {
                                if (checked && !included) selected = [...selected, suggestion.id];
                                else if (!checked && included) selected = selected.filter((id) => id !== suggestion.id);
                              }}
                            />
                            <a
                              class="flex max-w-[650px] shrink-0 cursor-pointer flex-nowrap items-center space-x-2"
                              href={`/dash/${$product.id}/suggestions/${suggestion.id}`}
                            >
                              {#if suggestion.archivedAt && dayjs(suggestion.archivedAt).isBefore(suggestion.lastSeenAt)}
                                <Badge variant="destructive" class="shrink-0">REGRESSION</Badge>
                              {/if}
                              {#if dayjs(suggestion.firstSeenAt).isAfter(dayjs().subtract(entities.SUGGESTION_NEW_MAX_DAYS - 1, "days"))}
                                <Badge class="shrink-0">NEW</Badge>
                              {/if}
                              {#if category && category !== entities.NO_CATEGORY}
                                <Badge variant="outline" class="shrink-0">{category}</Badge>
                              {/if}
                              <Tooltip.Root openDelay={750}>
                                <Tooltip.Trigger class="truncate text-nowrap font-medium">
                                  {suggestion.title}
                                </Tooltip.Trigger>
                                <Tooltip.Content>
                                  <p class="max-w-[650px] text-wrap">{suggestion.title}</p>
                                </Tooltip.Content>
                              </Tooltip.Root>
                            </a>
                          </div>
                        {:else if cell.id === "importance"}
                          {@const importance = entities.SuggestionImportanceDetails[major(suggestion.importances)]}
                          <div class="flex max-w-[100px] items-center">
                            <svelte:component
                              this={importance.icon}
                              class="mr-2 h-4 w-4 shrink-0 text-muted-foreground"
                            />
                            <span>{importance.title}</span>
                          </div>
                        {:else if cell.id === "seen"}
                          <span class="flex items-center justify-start gap-1">
                            <Tooltip.Root openDelay={333}>
                              <Tooltip.Trigger class="cursor-text">
                                {simplify(dayjs(suggestion.lastSeenAt).fromNow())}
                              </Tooltip.Trigger>
                              <Tooltip.Content>
                                <p>Last seen</p>
                                <p>{dayjs(suggestion.lastSeenAt).toString()}</p>
                              </Tooltip.Content>
                            </Tooltip.Root>
                            <span class="text-muted-foreground">/</span>
                            <Tooltip.Root openDelay={333}>
                              <Tooltip.Trigger class="cursor-text">
                                {simplify(dayjs(suggestion.firstSeenAt).fromNow(true))} old
                              </Tooltip.Trigger>
                              <Tooltip.Content>
                                <p>First seen</p>
                                <p>{dayjs(suggestion.firstSeenAt).toString()}</p>
                              </Tooltip.Content>
                            </Tooltip.Root>
                          </span>
                        {:else if cell.id === "customers"}
                          <span class="flex justify-end">
                            <Tooltip.Root openDelay={333}>
                              <Tooltip.Trigger class="cursor-text">
                                {unitize(suggestion.customers)}
                              </Tooltip.Trigger>
                              <Tooltip.Content>
                                <p>{suggestion.customers} affected</p>
                              </Tooltip.Content>
                            </Tooltip.Root>
                          </span>
                        {:else if cell.id === "assignee"}
                          {@const assignee = $users.find((user) => user.id === suggestion.assigneeID)}
                          <div class="mr-2 flex justify-end">
                            <Tooltip.Root openDelay={333}>
                              <Tooltip.Trigger class="cursor-default">
                                {#if assignee}
                                  <Picture size="lg" src={assignee.picture} fallback={assignee.name} />
                                {:else}
                                  <CircleUserRound class="h-6 w-6 stroke-[1.5]" />
                                {/if}
                              </Tooltip.Trigger>
                              <Tooltip.Content>
                                <p>{assignee ? `Assigned to ${assignee.name}` : "Unassigned"}</p>
                              </Tooltip.Content>
                            </Tooltip.Root>
                          </div>
                        {:else}
                          <Render of={cell.render()} />
                        {/if}
                      </Table.Cell>
                    </Subscribe>
                  {/each}
                </Table.Row>
              </Subscribe>
            {/each}
          {:else if !load}
            <Table.Row>
              <Table.Cell colspan={columns.length} class="h-24 text-center">
                <div class="flex flex-1 items-center justify-center space-x-4">
                  <Skeleton class="h-12 w-12 rounded-full" />
                  <div class="space-y-2">
                    <Skeleton class="h-4 w-[250px]" />
                    <Skeleton class="h-4 w-[200px]" />
                  </div>
                </div>
              </Table.Cell>
            </Table.Row>
          {:else}
            <Table.Row>
              <Table.Cell colspan={columns.length} class="h-24 text-center">No suggestions found.</Table.Cell>
            </Table.Row>
          {/if}
        {:catch}
          <Table.Row>
            <Table.Cell colspan={columns.length} class="h-24 text-center">Something went wrong!</Table.Cell>
          </Table.Row>
        {/await}
      </Table.Body>
    </Table.Root>
  </div>
  <div class="mt-4 flex items-center justify-between pl-2">
    <div class="flex-1 text-sm text-muted-foreground">
      {selected.length} of {$suggestions.length} suggestions selected.
    </div>
    <div class="flex items-center space-x-4">
      <div class="flex items-center space-x-2">
        <p class="text-sm font-medium">Suggestions per page</p>
        <Select.Root
          selected={{ value: $request.pagination.limit || 100, label: $request.pagination.limit?.toString() || "100" }}
          onSelectedChange={(limit) => ($request.pagination.limit = limit?.value)}
        >
          <Select.Trigger class="h-8 w-[70px]">
            <Select.Value placeholder="Select page size" />
          </Select.Trigger>
          <Select.Content>
            <Select.Item value={10}>10</Select.Item>
            <Select.Item value={25}>25</Select.Item>
            <Select.Item value={50}>50</Select.Item>
            <Select.Item value={100}>100</Select.Item>
          </Select.Content>
        </Select.Root>
      </div>
      <div class="flex w-[100px] items-center justify-center text-sm font-medium">
        Page {currentPage + 1} of {pages.length}
      </div>
      <div class="flex items-center space-x-2">
        <Button
          variant="outline"
          class="h-8 w-8 p-0"
          on:click={() => ($request.pagination.from = pages[(currentPage = 0)])}
          disabled={currentPage <= 0}
        >
          <DoubleArrowLeft size={15} />
        </Button>
        <Button
          variant="outline"
          class="h-8 w-8 p-0"
          on:click={() => ($request.pagination.from = pages[--currentPage])}
          disabled={currentPage <= 0}
        >
          <ChevronLeft size={15} />
        </Button>
        <Button
          variant="outline"
          class="h-8 w-8 p-0"
          on:click={() => ($request.pagination.from = pages[++currentPage])}
          disabled={currentPage >= pages.length - 1}
        >
          <ChevronRight size={15} />
        </Button>
        <Button
          variant="outline"
          class="h-8 w-8 p-0"
          on:click={() => ($request.pagination.from = pages[(currentPage = pages.length - 1)])}
          disabled={currentPage >= pages.length - 1}
        >
          <DoubleArrowRight size={15} />
        </Button>
      </div>
    </div>
  </div>
</div>
<Dialog.Root bind:open={openArchiveSuggestionsDialog}>
  <Dialog.Content>
    <Dialog.Header>
      <Dialog.Title>Are you absolutely sure?</Dialog.Title>
      <p class="pt-2 text-sm text-muted-foreground">
        This will mark the suggestions as completed and hide them. We will notify you of any regression. You can always
        unarchive them again.
      </p>
    </Dialog.Header>
    <Dialog.Footer>
      <Button variant="outline" on:click={() => (openArchiveSuggestionsDialog = false)}>Cancel</Button>
      <Button on:click={() => archiveSuggestions()}>Archive</Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
<Dialog.Root bind:open={openUnarchiveSuggestionsDialog}>
  <Dialog.Content>
    <Dialog.Header>
      <Dialog.Title>Are you absolutely sure?</Dialog.Title>
      <p class="pt-2 text-sm text-muted-foreground">
        This will mark the suggestions as ongoing. You can always archive them again.
      </p>
    </Dialog.Header>
    <Dialog.Footer>
      <Button variant="outline" on:click={() => (openUnarchiveSuggestionsDialog = false)}>Cancel</Button>
      <Button on:click={() => unarchiveSuggestions()}>Unarchive</Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
