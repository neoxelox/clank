<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";
  import { Badge } from "$lib/components/ui/badge";
  import { Button } from "$lib/components/ui/button";
  import { Checkbox } from "$lib/components/ui/checkbox";
  import { DatePicker } from "$lib/components/ui/date-picker";
  import { Filter } from "$lib/components/ui/filter";
  import { Picture } from "$lib/components/ui/picture";
  import * as Select from "$lib/components/ui/select";
  import { Skeleton } from "$lib/components/ui/skeleton";
  import * as Table from "$lib/components/ui/table";
  import * as Tooltip from "$lib/components/ui/tooltip";
  import * as entities from "$lib/entities";
  import * as payloads from "$lib/payloads";
  import { api } from "$lib/services/api";
  import { product } from "$lib/stores";
  import dayjs from "$lib/utils/datetime";
  import { debounce } from "$lib/utils/delay";
  import { clone, compare } from "$lib/utils/object";
  import { simplify, titlelize } from "$lib/utils/string";
  import Star from "lucide-svelte/icons/star";
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
  import QuestionMark from "svelte-radix/QuestionMark.svelte";
  import { writable } from "svelte/store";

  let defaults = (request: payloads.ListReviewsRequest): payloads.ListReviewsRequest => {
    request.filters.seen_start_at = request.filters.seen_start_at || undefined;
    request.filters.seen_end_at = request.filters.seen_end_at || new Date();
    request.orders.recency = request.orders.recency || "DESC";
    request.pagination.limit = request.pagination.limit || 50;

    return request;
  };

  let previousRequest: payloads.ListReviewsRequest;
  let request = writable<payloads.ListReviewsRequest>(defaults(payloads.fromListReviewsQuery($page.url.searchParams)));
  let load: Promise<void>;
  let reviews = writable<entities.Review[]>([]);
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
  const keywordOptions = [];
  const sentimentOptions = Object.values(entities.ReviewSentiment).map((sentiment) => {
    const details = entities.ReviewSentimentDetails[sentiment];
    return {
      value: sentiment,
      label: details.title,
      icon: details.icon,
    };
  });
  const emotionOptions = Object.values(entities.ReviewEmotion).map((emotion) => {
    const details = entities.ReviewEmotionDetails[emotion];
    return {
      value: emotion,
      label: details.title,
      icon: details.icon,
    };
  });
  const intentionOptions = Object.values(entities.ReviewIntention).map((intention) => {
    const details = entities.ReviewIntentionDetails[intention];
    return {
      value: intention,
      label: details.title,
      icon: details.icon,
    };
  });
  intentionOptions.push({
    value: entities.NO_INTENTION,
    label: "Unknown",
    icon: QuestionMark,
  });
  const languageOptions = Object.values(entities.FeedbackLanguage).map((language) => {
    return {
      value: language,
      label: titlelize(language.replaceAll("_", " ")),
    };
  });
  languageOptions.push({
    value: entities.NO_LANGUAGE,
    label: "Unknown",
  });
  const orderOptions = ["ASC", "DESC", undefined];
  let currentPage = 0;

  let loadReviews = debounce((params: string) => {
    load = (async () => {
      const response = await api.get<payloads.ListReviewsResponse>(`/products/${$product.id}/reviews?${params}`);
      $reviews = response.reviews.map((review) => payloads.toReview(review));
      const page = response.next ? decodeURIComponent(response.next) : undefined;
      if (!pages.includes(page)) pages = pages.concat([page]);
    })();
  }, 500);

  $: if ($page.url.searchParams.has("_refresh"))
    $request = defaults(payloads.fromListReviewsQuery($page.url.searchParams));

  let resetQuery = () => {
    previousRequest = undefined;
    $request = defaults(payloads.fromListReviewsQuery(new URLSearchParams()));
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

    const params = payloads.toListReviewsQuery(request).toString();
    await goto(`?${params}`, { noScroll: true, keepFocus: true });
    loadReviews(params);
  });

  const table = createTable(reviews);

  const columns = table.createColumns([
    table.column({
      id: "review",
      accessor: (review) => review.feedback.content,
      header: "",
    }),
    table.column({
      id: "sentiment",
      accessor: "sentiment",
      header: "Sentiment",
    }),
    table.column({
      id: "seen",
      accessor: (review) => review.feedback.postedAt,
      header: "Posted at",
    }),
    table.column({
      id: "intention",
      accessor: "intention",
      header: "Intention",
    }),
    table.column({
      id: "customer",
      accessor: (review) => review.feedback.customer.name,
      header: "Customer",
    }),
  ]);

  const { headerRows, pageRows, tableAttrs, tableBodyAttrs } = table.createViewModel(columns);

  onDestroy(() => {
    unsubscribe();
  });
</script>

<div class={"w-full " + ($$props.class || "")}>
  <div class="mb-4 flex flex-wrap items-center justify-start gap-4">
    <Filter title="Sources" options={sourceOptions} bind:values={$request.filters.sources} />
    <Filter title="Releases" options={releaseOptions} addable bind:values={$request.filters.releases} />
    <Filter title="Categories" options={categoryOptions} addable bind:values={$request.filters.categories} />
    <Filter title="Keywords" options={keywordOptions} addable bind:values={$request.filters.keywords} />
    {#if $page.url.searchParams.toString().length > 1}
      <Button variant="ghost" on:click={() => resetQuery()}>
        Reset
        <Cross2 class="ml-2 h-4 w-4" />
      </Button>
    {/if}
    <DatePicker
      class="ml-auto"
      placeholder="Posted at"
      maxValue={new Date()}
      bind:startValue={$request.filters.seen_start_at}
      bind:endValue={$request.filters.seen_end_at}
    />
  </div>
  <div class="mb-4 flex flex-wrap items-center justify-start gap-4">
    <Filter title="Sentiments" options={sentimentOptions} bind:values={$request.filters.sentiments} />
    <Filter title="Emotions" options={emotionOptions} bind:values={$request.filters.emotions} />
    <Filter title="Intentions" options={intentionOptions} bind:values={$request.filters.intentions} />
    <Filter title="Languages" options={languageOptions} bind:values={$request.filters.languages} />
    <Button
      variant="outline"
      class="ml-auto"
      on:click={() => ($request.orders.recency = orderOptions[(orderOptions.indexOf($request.orders.recency) + 1) % 3])}
    >
      Recency
      {#if $request.orders.recency === "ASC"}
        <ArrowUp class="ml-2 h-4 w-4" />
      {:else if $request.orders.recency === "DESC"}
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
                    {#if cell.id === "review"}
                      {@const checked =
                        $reviews.length > 0
                          ? selected.length === $reviews.length
                            ? true
                            : selected.length > 0
                              ? "indeterminate"
                              : false
                          : false}
                      <div class="flex items-center justify-start space-x-2">
                        <Checkbox
                          class="ml-2 mr-1"
                          {checked}
                          onCheckedChange={(change) => {
                            if (change && checked !== "indeterminate") selected = $reviews.map((review) => review.id);
                            else selected = [];
                          }}
                        />
                        <span>Feedback</span>
                      </div>
                    {:else if cell.id === "sentiment"}
                      <Tooltip.Root openDelay={333}>
                        <Tooltip.Trigger class="cursor-text">
                          <Render of={cell.render()} />
                        </Tooltip.Trigger>
                        <Tooltip.Content>
                          <p class="max-w-[300px] text-wrap">
                            Attitude is smartly inferred from the context of your product and how customers write the
                            feedback.
                          </p>
                        </Tooltip.Content>
                      </Tooltip.Root>
                    {:else if cell.id === "intention"}
                      <Tooltip.Root openDelay={333}>
                        <Tooltip.Trigger class="cursor-text">
                          <Render of={cell.render()} />
                        </Tooltip.Trigger>
                        <Tooltip.Content>
                          <p class="max-w-[315px] text-wrap">
                            Tendency is smartly approximated from the context of your product and how customers write
                            the feedback.
                          </p>
                        </Tooltip.Content>
                      </Tooltip.Root>
                    {:else if cell.id === "customer"}
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
                    {@const review = cell.row.original}
                    <Subscribe attrs={cell.attrs()} let:attrs>
                      <Table.Cell {...attrs}>
                        {#if cell.id === "review"}
                          {@const included = selected.includes(review.id)}
                          {@const content =
                            review.feedback.language === $product.language
                              ? review.feedback.content
                              : review.feedback.translation}
                          <div class="flex shrink-0 flex-nowrap items-center space-x-2">
                            <Checkbox
                              class="ml-2 mr-1"
                              checked={included}
                              onCheckedChange={(checked) => {
                                if (checked && !included) selected = [...selected, review.id];
                                else if (!checked && included) selected = selected.filter((id) => id !== review.id);
                              }}
                            />
                            <a
                              class="flex max-w-[625px] shrink-0 cursor-pointer flex-nowrap items-center space-x-2"
                              href={`/dash/${$product.id}/feedbacks/${review.id}`}
                            >
                              {#if review.category !== entities.NO_CATEGORY}
                                <Badge variant="outline" class="shrink-0">{review.category.replaceAll("_", " ")}</Badge>
                              {/if}
                              <Tooltip.Root openDelay={750}>
                                <Tooltip.Trigger class="truncate text-nowrap font-medium">
                                  {content}
                                </Tooltip.Trigger>
                                <Tooltip.Content>
                                  <p class="max-w-[625px] text-wrap">{content}</p>
                                </Tooltip.Content>
                              </Tooltip.Root>
                            </a>
                          </div>
                        {:else if cell.id === "sentiment"}
                          {@const sentiment = entities.ReviewSentimentDetails[review.sentiment]}
                          <div class="flex max-w-[100px] items-center">
                            <svelte:component
                              this={sentiment.icon}
                              class="mr-2 h-4 w-4 shrink-0 text-muted-foreground"
                            />
                            <span>{sentiment.title}</span>
                          </div>
                        {:else if cell.id === "seen"}
                          <span class="flex items-center justify-start gap-1">
                            <Tooltip.Root openDelay={333}>
                              <Tooltip.Trigger class="cursor-text">
                                {simplify(dayjs(review.feedback.postedAt).fromNow())}
                              </Tooltip.Trigger>
                              <Tooltip.Content>
                                <p>Posted at</p>
                                <p>{dayjs(review.feedback.postedAt).toString()}</p>
                              </Tooltip.Content>
                            </Tooltip.Root>
                          </span>
                        {:else if cell.id === "intention"}
                          {#if review.intention !== entities.NO_INTENTION}
                            {@const intention = entities.ReviewIntentionDetails[review.intention]}
                            <div class="flex max-w-[180px] items-center">
                              <svelte:component
                                this={intention.icon}
                                class="mr-2 h-4 w-4 shrink-0 text-muted-foreground"
                              />
                              <span>{intention.title}</span>
                            </div>
                          {:else}
                            <div class="flex max-w-[180px] items-center">
                              <QuestionMark class="mr-2 h-4 w-4 shrink-0 text-muted-foreground" />
                              <span>Unknown</span>
                            </div>
                          {/if}
                        {:else if cell.id === "customer"}
                          <div class="mr-2 flex items-center justify-end gap-1.5">
                            <Tooltip.Root openDelay={333}>
                              <Tooltip.Trigger class="cursor-default">
                                <div class="relative">
                                  <Star class="h-6 w-6 fill-muted stroke-muted stroke-2" />
                                  <span
                                    class="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-[47%] text-xs font-semibold text-primary"
                                  >
                                    {Math.floor(review.feedback.metadata.rating || 0)}
                                  </span>
                                </div>
                              </Tooltip.Trigger>
                              <Tooltip.Content>
                                {review.feedback.metadata.rating
                                  ? `${review.feedback.metadata.rating} stars`
                                  : "No rating"}
                              </Tooltip.Content>
                            </Tooltip.Root>
                            <Tooltip.Root openDelay={333}>
                              <Tooltip.Trigger class="cursor-default">
                                <Picture
                                  size="lg"
                                  src={review.feedback.customer.picture}
                                  fallback={review.feedback.customer.name}
                                />
                              </Tooltip.Trigger>
                              <Tooltip.Content>
                                <p>{`Posted by ${review.feedback.customer.name}`}</p>
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
              <Table.Cell colspan={columns.length} class="h-24 text-center">No feedbacks found.</Table.Cell>
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
    <div class="flex-1 text-sm text-muted-foreground">{selected.length} of {$reviews.length} feedbacks selected.</div>
    <div class="flex items-center space-x-4">
      <div class="flex items-center space-x-2">
        <p class="text-sm font-medium">Feedbacks per page</p>
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
