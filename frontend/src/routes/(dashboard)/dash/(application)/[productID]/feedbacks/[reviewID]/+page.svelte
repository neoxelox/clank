<script lang="ts">
  import { page } from "$app/stores";
  import EmotionChart from "$lib/components/cmp/EmotionChart.svelte";
  import { Badge } from "$lib/components/ui/badge";
  import * as Card from "$lib/components/ui/card";
  import * as Carousel from "$lib/components/ui/carousel";
  import { type CarouselAPI } from "$lib/components/ui/carousel/context";
  import { Label } from "$lib/components/ui/label";
  import { Picture } from "$lib/components/ui/picture";
  import { Skeleton } from "$lib/components/ui/skeleton";
  import * as Tooltip from "$lib/components/ui/tooltip";
  import * as entities from "$lib/entities";
  import * as payloads from "$lib/payloads";
  import { api } from "$lib/services/api";
  import { product } from "$lib/stores";
  import dayjs from "$lib/utils/datetime";
  import { imgPlaceholder } from "$lib/utils/image";
  import { capitalize, simplify, titlelize, unitize } from "$lib/utils/string";
  import BadgeCheck from "lucide-svelte/icons/badge-check";
  import Check from "lucide-svelte/icons/check";
  import ExternalLink from "lucide-svelte/icons/external-link";
  import History from "lucide-svelte/icons/history";
  import Languages from "lucide-svelte/icons/languages";
  import Mail from "lucide-svelte/icons/mail";
  import MapPin from "lucide-svelte/icons/map-pin";
  import MessageSquareText from "lucide-svelte/icons/message-square-text";
  import Star from "lucide-svelte/icons/star";
  import ThumbsUp from "lucide-svelte/icons/thumbs-up";
  import { copy } from "svelte-copy";
  import QuestionMark from "svelte-radix/QuestionMark.svelte";
  import QuestionMarkCircled from "svelte-radix/QuestionMarkCircled.svelte";

  let load: Promise<void>;
  let review: entities.Review;

  $: emotions = (review?.emotions || []).reduce((accu, emotion) => ({ ...accu, [emotion]: 1 }), {});

  const loadReview = async () => {
    load = (async () => {
      const response = await api.get<payloads.GetReviewResponse>(
        `/products/${$product.id}/reviews/${$page.params.reviewID}`,
      );
      review = payloads.toReview(response);
    })();
  };
  loadReview();

  let current: number = 1;
  let carousel: CarouselAPI;
  $: if (carousel) carousel.on("select", () => (current = carousel.selectedScrollSnap() + 1));
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
    <div class="flex flex-col items-start justify-center gap-4 border-b pb-4">
      <div class="flex items-start justify-center gap-4">
        <Picture size="3xl" src={review.feedback.customer.picture} fallback={review.feedback.customer.name} />
        <div class="-mt-px flex flex-col items-start justify-center gap-0.5">
          <div class="flex items-center gap-2">
            {#if review.feedback.customer.link}
              <a href={review.feedback.customer.link} target="_blank" rel="noopener noreferrer">
                <h1 class="text-xl font-semibold hover:text-primary">{review.feedback.customer.name}</h1>
              </a>
            {:else}
              <h1 class="text-xl font-semibold">{review.feedback.customer.name}</h1>
            {/if}
            {#if review.feedback.customer.verified}
              <BadgeCheck class="h-4 w-4 text-primary" />
            {/if}
          </div>
          <h3
            class="inline-block cursor-default font-mono text-sm font-normal text-muted-foreground transition-colors hover:text-foreground"
            use:copy={review.id}
          >
            #{review.id}
          </h3>
        </div>
      </div>
      <div class="flex items-center gap-4 text-center text-sm text-muted-foreground">
        <Tooltip.Root openDelay={750}>
          <Tooltip.Trigger class="flex cursor-default items-center gap-1">
            {#each [...Array(Math.floor(review.feedback.metadata.rating || 0))] as _}
              <Star class="h-5 w-5 fill-primary stroke-primary stroke-2" />
            {/each}
            {#each [...Array(5 - Math.floor(review.feedback.metadata.rating || 0))] as _}
              <Star class="h-5 w-5 fill-muted stroke-muted stroke-2" />
            {/each}
          </Tooltip.Trigger>
          <Tooltip.Content>
            {review.feedback.metadata.rating ? `${review.feedback.metadata.rating} stars` : "No rating"}
          </Tooltip.Content>
        </Tooltip.Root>
        {#if review.feedback.metadata.verified}
          <div class="-ml-1 flex items-center gap-1.5">
            <Check class="h-4 w-4" />
            <span>Verified</span>
          </div>
        {/if}
        <div class="flex items-center gap-2">
          <MessageSquareText class="h-4 w-4" />
          <span>
            {(review.feedback.customer.reviews || 0) > 1 ? `${review.feedback.customer.reviews} reviews` : "1 review"}
          </span>
        </div>
        {#if review.feedback.customer.location}
          <div class="flex items-center gap-1.5">
            <MapPin class="h-4 w-4" />
            <span>{review.feedback.customer.location}</span>
          </div>
        {/if}
        {#if review.feedback.customer.email}
          <div class="flex items-center gap-2 hover:text-primary">
            <Mail class="h-4 w-4" />
            <a href={`mailto:${review.feedback.customer.email}`}>{review.feedback.customer.email}</a>
          </div>
        {/if}
        <div class="flex items-center gap-2">
          <ThumbsUp class="mb-px h-4 w-4" />
          <span>{unitize(review.feedback.metadata.votes || 0)}</span>
        </div>
        <div class="flex items-center gap-1.5">
          <History class="mt-px h-4 w-4" />
          <span>{review.feedback.release === entities.NO_RELEASE ? "Unknown" : review.feedback.release}</span>
        </div>
      </div>
    </div>
    <div class="grid h-full grid-cols-3 gap-8">
      <div class="col-span-2 space-y-6 pt-6">
        <div class="space-y-1">
          <Label>Content</Label>
          <p class="whitespace-pre-line text-justify leading-relaxed">
            {#if review.feedback.language === $product.language}
              {review.feedback.content}
            {:else}
              {#if review.feedback.original}
                {review.feedback.content}
              {:else}
                {review.feedback.translation}
              {/if}
              <span class="mt-2 flex items-center gap-2 whitespace-normal">
                {#if review.feedback.original}
                  <button
                    class="inline-flex items-center justify-center text-sm text-primary"
                    on:click={() => (review.feedback.original = false)}
                  >
                    <Languages class="mr-2 h-4 w-4" />
                    View translated
                  </button>
                {:else}
                  <button
                    class="inline-flex items-center justify-center text-sm text-primary"
                    on:click={() => (review.feedback.original = true)}
                  >
                    <Languages class="mr-2 h-4 w-4" />
                    View original ({capitalize(review.feedback.language)})
                  </button>
                {/if}
              </span>
            {/if}
          </p>
        </div>
        {#if (review.feedback.metadata.media?.length || 0) > 0}
          <div class="space-y-3">
            <Label>Media</Label>
            <Carousel.Root opts={{ watchDrag: false, duration: 20 }} bind:api={carousel} class="w-full">
              <div class="absolute -top-2 right-px z-10 flex -translate-y-full items-center justify-end space-x-2">
                <span class="mr-2 text-sm">{current} of {review.feedback.metadata.media.length}</span>
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
                {#each review.feedback.metadata.media as media}
                  <Carousel.Item>
                    <Card.Root class="rounded-lg p-6 shadow-sm">
                      <img
                        src={media}
                        onerror={imgPlaceholder}
                        fetchpriority="auto"
                        loading="lazy"
                        decoding="async"
                        class="aspect-video max-h-[500px] w-full rounded-xl bg-muted object-cover"
                      />
                    </Card.Root>
                  </Carousel.Item>
                {/each}
              </Carousel.Content>
            </Carousel.Root>
          </div>
        {/if}
      </div>
      <div class="space-y-6 border-l pl-6 pt-6">
        <div class="space-y-2">
          <Label>Source</Label>
          {#if review.feedback.metadata.link}
            <a
              href={review.feedback.metadata.link}
              target="_blank"
              rel="noopener noreferrer"
              class="flex items-center gap-1.5 text-primary"
            >
              <svelte:component
                this={entities.FeedbackSourceDetails[review.feedback.source].icon}
                class="mr-1 h-5 w-5 shrink-0"
              />
              <span>View on {entities.FeedbackSourceDetails[review.feedback.source].title}</span>
              <ExternalLink class="h-4 w-4" />
            </a>
          {:else}
            <div class="flex items-center gap-2.5">
              <svelte:component
                this={entities.FeedbackSourceDetails[review.feedback.source].icon}
                class="h-5 w-5 shrink-0"
              />
              <span>{entities.FeedbackSourceDetails[review.feedback.source].title}</span>
            </div>
          {/if}
        </div>
        <div class="space-y-2">
          <Label>Posted at</Label>
          <div>
            {dayjs(review.feedback.postedAt).toString()} ({simplify(dayjs(review.feedback.postedAt).fromNow())})
          </div>
        </div>
        <div class="space-y-2">
          <Tooltip.Root openDelay={333}>
            <Tooltip.Trigger class="mb-3 flex cursor-default items-center">
              <Label>Category</Label>
              <QuestionMarkCircled class="ml-1 h-4 w-4 shrink-0 text-foreground/50" />
            </Tooltip.Trigger>
            <Tooltip.Content>
              <p class="max-w-[300px] text-wrap">
                Feedbacks are smartly classified into categories to enhance organization.
                <a href={`/dash/${$product.id}/settings`} class="text-primary">Modify product categories.</a>
              </p>
            </Tooltip.Content>
          </Tooltip.Root>
          <div>{titlelize(review.category.replaceAll("_", " "))}</div>
        </div>
        <div class="space-y-2">
          <Tooltip.Root openDelay={333}>
            <Tooltip.Trigger class="mb-3 flex cursor-default items-center">
              <Label>Sentiment</Label>
              <QuestionMarkCircled class="ml-1 h-4 w-4 shrink-0 text-foreground/50" />
            </Tooltip.Trigger>
            <Tooltip.Content>
              <p class="max-w-[300px] text-wrap">
                Attitude is smartly inferred from the context of your product and how customers write the feedback.
              </p>
            </Tooltip.Content>
          </Tooltip.Root>
          <div class="flex items-center gap-2">
            <svelte:component
              this={entities.ReviewSentimentDetails[review.sentiment].icon}
              class="h-5 w-5 shrink-0 text-muted-foreground"
            />
            <span>{entities.ReviewSentimentDetails[review.sentiment].title}</span>
          </div>
        </div>
        <div class="space-y-2">
          <Tooltip.Root openDelay={333}>
            <Tooltip.Trigger class="mb-3 flex cursor-default items-center">
              <Label>Intention</Label>
              <QuestionMarkCircled class="ml-1 h-4 w-4 shrink-0 text-foreground/50" />
            </Tooltip.Trigger>
            <Tooltip.Content>
              <p class="max-w-[315px] text-wrap">
                Tendency is smartly approximated from the context of your product and how customers write the feedback.
              </p>
            </Tooltip.Content>
          </Tooltip.Root>
          {#if review.intention !== entities.NO_INTENTION}
            <div class="flex items-center gap-2.5">
              <svelte:component
                this={entities.ReviewIntentionDetails[review.intention].icon}
                class="h-5 w-5 shrink-0 text-muted-foreground"
              />
              <span>{entities.ReviewIntentionDetails[review.intention].title}</span>
            </div>
          {:else}
            <div class="flex items-center gap-2">
              <QuestionMark class="h-5 w-5 shrink-0 text-muted-foreground" />
              <span>Unknown</span>
            </div>
          {/if}
        </div>
        {#if review.keywords.length > 0}
          <div class="space-y-2">
            <Label>Keywords</Label>
            <div
              class="!mt-2.5 flex h-auto max-h-28 w-full flex-wrap gap-2 overflow-y-auto rounded-md border border-input bg-transparent p-3 shadow-sm"
            >
              {#each review.keywords as keyword}
                <Badge
                  variant="secondary"
                  class="flex h-6 shrink-0 select-text items-center justify-center font-medium"
                >
                  {keyword}
                </Badge>
              {/each}
            </div>
          </div>
        {/if}
        {#if review.emotions.length > 0}
          <div class="space-y-2">
            <Tooltip.Root openDelay={333}>
              <Tooltip.Trigger class="mb-4 flex cursor-default items-center">
                <Label>Emotions</Label>
                <QuestionMarkCircled class="ml-1 h-4 w-4 shrink-0 text-foreground/50" />
              </Tooltip.Trigger>
              <Tooltip.Content>
                <p class="max-w-[300px] text-wrap">
                  Feelings are smartly inferred from the context of your product and how customers write the feedback.
                </p>
                <p class="max-w-[300px] text-wrap">
                  Analysis is based on the <a
                    href="https://en.wikipedia.org/wiki/Emotion_classification#Plutchik's_wheel_of_emotions"
                    target="_blank"
                    rel="noopener noreferrer"
                    class="text-primary">Plutchik wheel of emotions</a
                  >.
                </p>
              </Tooltip.Content>
            </Tooltip.Root>
            <div
              class="!mt-3 flex h-[300px] w-full items-center justify-center rounded-lg border border-input bg-transparent px-3 shadow-sm"
            >
              <EmotionChart {emotions} />
            </div>
          </div>
        {/if}
      </div>
    </div>
  </div>
{:catch}
  <div class="flex flex-1 items-center justify-center space-x-4 rounded-lg border border-dashed shadow-sm">
    Something went wrong!
  </div>
{/await}
