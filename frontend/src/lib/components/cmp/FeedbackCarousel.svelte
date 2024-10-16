<script lang="ts">
  import * as Card from "$lib/components/ui/card";
  import * as Carousel from "$lib/components/ui/carousel";
  import { type CarouselAPI } from "$lib/components/ui/carousel/context";
  import { Picture } from "$lib/components/ui/picture";
  import { Skeleton } from "$lib/components/ui/skeleton";
  import * as Tooltip from "$lib/components/ui/tooltip";
  import * as entities from "$lib/entities";
  import * as payloads from "$lib/payloads";
  import { api } from "$lib/services/api";
  import { product } from "$lib/stores";
  import dayjs from "$lib/utils/datetime";
  import { capitalize, simplify, unitize } from "$lib/utils/string";
  import BadgeCheck from "lucide-svelte/icons/badge-check";
  import CalendarRange from "lucide-svelte/icons/calendar-range";
  import Check from "lucide-svelte/icons/check";
  import ExternalLink from "lucide-svelte/icons/external-link";
  import History from "lucide-svelte/icons/history";
  import Languages from "lucide-svelte/icons/languages";
  import Mail from "lucide-svelte/icons/mail";
  import MapPin from "lucide-svelte/icons/map-pin";
  import MessageSquareShare from "lucide-svelte/icons/message-square-share";
  import MessageSquareText from "lucide-svelte/icons/message-square-text";
  import Star from "lucide-svelte/icons/star";
  import ThumbsUp from "lucide-svelte/icons/thumbs-up";

  export let issue: entities.Issue | undefined = undefined;
  export let suggestion: entities.Suggestion | undefined = undefined;

  let load: Promise<void>;
  let feedbacks: entities.Feedback[] = [];
  let next: string | undefined = undefined;

  const loadFeedbacks = async (
    request: payloads.ListIssueFeedbacksRequest | payloads.ListSuggestionFeedbacksRequest,
  ) => {
    load = (async () => {
      const from = request.from ? encodeURIComponent(request.from) : undefined;
      let response: payloads.ListIssueFeedbacksResponse | payloads.ListSuggestionFeedbacksResponse;

      if (issue)
        response = await api.get<payloads.ListIssueFeedbacksResponse>(
          `/products/${$product.id}/issues/${issue.id}/feedbacks?${from ? `from=${from}` : ""}`,
        );
      else if (suggestion)
        response = await api.get<payloads.ListSuggestionFeedbacksResponse>(
          `/products/${$product.id}/suggestions/${suggestion.id}/feedbacks?${from ? `from=${from}` : ""}`,
        );
      else throw new Error("Unknown feedback entity type");

      feedbacks = feedbacks.concat(response.feedbacks.map((feedback) => payloads.toFeedback(feedback)));
      next = response.next ? decodeURIComponent(response.next) : undefined;
    })();
  };
  loadFeedbacks({});

  let current: number = 1;
  let requested: boolean = false;
  let carousel: CarouselAPI;
  $: if (carousel) {
    carousel.on("select", async () => {
      const index = carousel.selectedScrollSnap();
      current = index + 1;

      if (requested) return;
      requested = true;

      if (index === feedbacks.length - 1 && next) {
        loadFeedbacks({ from: next });
        await load;
        setTimeout(() => carousel.scrollTo(index, true), 1);
      }

      requested = false;
    });
  }
</script>

<div class={"w-full " + ($$props.class || "")}>
  <Carousel.Root opts={{ watchDrag: false, duration: 20 }} bind:api={carousel} class="w-full">
    <div class="absolute -top-2 right-px z-10 flex -translate-y-full items-center justify-end space-x-2">
      <span class="mr-2 text-sm">{current} of {feedbacks.length}</span>
      <Carousel.Previous
        variant="outline"
        class="static flex h-8 w-8 translate-x-0 translate-y-0 rotate-0 items-center justify-center rounded-md p-0"
      />
      <Carousel.Next
        variant="outline"
        class="static flex h-8 w-8 translate-x-0 translate-y-0 rotate-0 items-center justify-center rounded-md p-0"
      />
    </div>
    <Carousel.Content>
      {#await load}
        <Carousel.Item>
          <Card.Root
            class="flex flex-1 items-center justify-center space-x-4 rounded-lg border border-dashed p-6 pr-12 shadow-sm"
          >
            <Skeleton class="h-12 w-12 rounded-full" />
            <div class="space-y-2">
              <Skeleton class="h-4 w-[250px]" />
              <Skeleton class="h-4 w-[200px]" />
            </div>
          </Card.Root>
        </Carousel.Item>
      {:then}
        {#each feedbacks as feedback}
          <Carousel.Item>
            <Card.Root class="rounded-lg p-6 pr-12 shadow-sm">
              <div class="flex items-start gap-4">
                <Picture size="3xl" src={feedback.customer.picture} fallback={feedback.customer.name} class="mt-1" />
                <div class="flex-1">
                  <div class="flex items-center gap-2">
                    {#if feedback.customer.link}
                      <a href={feedback.customer.link} target="_blank" rel="noopener noreferrer">
                        <h3 class="font-semibold hover:text-primary">{feedback.customer.name}</h3>
                      </a>
                    {:else}
                      <h3 class="font-semibold">{feedback.customer.name}</h3>
                    {/if}
                    {#if feedback.customer.verified}
                      <BadgeCheck class="h-4 w-4 text-primary" />
                    {/if}
                  </div>
                  <div class="mt-2 flex items-center gap-4 text-center text-sm text-muted-foreground">
                    <Tooltip.Root openDelay={750}>
                      <Tooltip.Trigger class="flex cursor-default items-center gap-1">
                        {#each [...Array(Math.floor(feedback.metadata.rating || 0))] as _}
                          <Star class="h-5 w-5 fill-primary stroke-primary stroke-2" />
                        {/each}
                        {#each [...Array(5 - Math.floor(feedback.metadata.rating || 0))] as _}
                          <Star class="h-5 w-5 fill-muted stroke-muted stroke-2" />
                        {/each}
                      </Tooltip.Trigger>
                      <Tooltip.Content>
                        {feedback.metadata.rating ? `${feedback.metadata.rating} stars` : "No rating"}
                      </Tooltip.Content>
                    </Tooltip.Root>
                    {#if feedback.metadata.verified}
                      <div class="-ml-1 flex items-center gap-1.5">
                        <Check class="h-4 w-4" />
                        <span>Verified</span>
                      </div>
                    {/if}
                    <div class="flex items-center gap-2">
                      <MessageSquareText class="h-4 w-4" />
                      <span>
                        {(feedback.customer.reviews || 0) > 1 ? `${feedback.customer.reviews} reviews` : "1 review"}
                      </span>
                    </div>
                    {#if feedback.customer.location}
                      <div class="flex items-center gap-1.5">
                        <MapPin class="h-4 w-4" />
                        <span>{feedback.customer.location}</span>
                      </div>
                    {/if}
                    {#if feedback.customer.email}
                      <div class="flex items-center gap-2 hover:text-primary">
                        <Mail class="h-4 w-4" />
                        <a href={`mailto:${feedback.customer.email}`}>{feedback.customer.email}</a>
                      </div>
                    {/if}
                  </div>
                  <p class="mt-4 whitespace-pre-line text-justify text-sm leading-relaxed">
                    {#if feedback.language === $product.language}
                      {feedback.content}
                    {:else}
                      {#if feedback.original}
                        {feedback.content}
                      {:else}
                        {feedback.translation}
                      {/if}
                      <span class="mt-2 flex items-center gap-2 whitespace-normal">
                        {#if feedback.original}
                          <button
                            class="inline-flex items-center justify-center text-primary"
                            on:click={() => (feedback.original = false)}
                          >
                            <Languages class="mr-2 h-4 w-4" />
                            View translated
                          </button>
                        {:else}
                          <button
                            class="inline-flex items-center justify-center text-primary"
                            on:click={() => (feedback.original = true)}
                          >
                            <Languages class="mr-2 h-4 w-4" />
                            View original ({capitalize(feedback.language)})
                          </button>
                        {/if}
                      </span>
                    {/if}
                  </p>
                  <div class="mt-4 flex items-center gap-4 text-center text-sm text-muted-foreground">
                    <div class="flex items-center gap-2">
                      <ThumbsUp class="mb-px h-4 w-4" />
                      <span>{unitize(feedback.metadata.votes || 0)}</span>
                    </div>
                    <div class="flex items-center gap-1.5">
                      <History class="mt-px h-4 w-4" />
                      <span>{feedback.release === entities.NO_RELEASE ? "Unknown" : feedback.release}</span>
                    </div>
                    <Tooltip.Root openDelay={333}>
                      <Tooltip.Trigger class="flex cursor-default items-center gap-2">
                        <CalendarRange class="h-4 w-4" />
                        <span class="cursor-text">{simplify(dayjs(feedback.postedAt).fromNow())}</span>
                      </Tooltip.Trigger>
                      <Tooltip.Content>
                        <p>Posted at</p>
                        <p>{dayjs(feedback.postedAt).toString()}</p>
                      </Tooltip.Content>
                    </Tooltip.Root>
                    {#if feedback.metadata.link}
                      <a
                        href={feedback.metadata.link}
                        target="_blank"
                        rel="noopener noreferrer"
                        class="flex items-center gap-1.5 text-primary"
                      >
                        <ExternalLink class="h-4 w-4" />
                        <span>View on {entities.FeedbackSourceDetails[feedback.source].title}</span>
                      </a>
                    {:else}
                      <div class="flex items-center gap-2">
                        <MessageSquareShare class="mt-px h-4 w-4" />
                        <span>{entities.FeedbackSourceDetails[feedback.source].title}</span>
                      </div>
                    {/if}
                  </div>
                </div>
              </div>
            </Card.Root>
          </Carousel.Item>
        {/each}
      {:catch}
        <Carousel.Item>
          <Card.Root
            class="flex h-[6.1rem] flex-1 items-center justify-center space-x-4 rounded-lg border border-dashed p-6 pr-12 text-center shadow-sm"
          >
            Something went wrong!
          </Card.Root>
        </Carousel.Item>
      {/await}
    </Carousel.Content>
  </Carousel.Root>
</div>
