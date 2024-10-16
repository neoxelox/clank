<script lang="ts">
  import Persona1Avatar from "$assets/avatars/persona-1.png";
  import Persona2Avatar from "$assets/avatars/persona-2.png";
  import Persona3Avatar from "$assets/avatars/persona-3.png";
  import Persona4Avatar from "$assets/avatars/persona-4.png";
  import { Button } from "$lib/components/ui/button";
  import { Picture } from "$lib/components/ui/picture";
  import * as Tabs from "$lib/components/ui/tabs";
  import { Textarea } from "$lib/components/ui/textarea";
  import * as Tooltip from "$lib/components/ui/tooltip";
  import CornerDownLeft from "lucide-svelte/icons/corner-down-left";
  import Images from "lucide-svelte/icons/images";
  import Paperclip from "lucide-svelte/icons/paperclip";
  import CaretSort from "svelte-radix/CaretSort.svelte";
  import PlusCircled from "svelte-radix/PlusCircled.svelte";

  const personas = [
    {
      id: "cop2f42bklfs60283no0",
      name: "Michael Thompson",
      personality: "Analytical and organized",
      description:
        "Michael uses the app primarily for booking long-distance flights and renting cars at his destinations. He values the app's comprehensive search for over 600 airlines, allowing him to compare and select the best flights that suit his schedule. The easy booking process and access to booking references and flight status updates ensure his business trips run smoothly.",
      image: Persona1Avatar,
    },
    {
      id: "coonm8ldnifcucu7qsu0",
      name: "Rachel Greene",
      personality: "Outgoing and spontaneous",
      description:
        "Rachel loves using the app to find and compare cheap flights for her last-minute weekend getaways. She often books flight + hotel packages to save money and appreciates the free travel guides to scout for Instagram-worthy spots. The real-time flight status notifications keep her updated without any fuss, which fits her busy lifestyle.",
      image: Persona2Avatar,
    },
    {
      id: "cop2h32bklfs60283pb0",
      name: "David Parker",
      personality: "Curious and adventurous",
      description:
        "David uses the app to find multi-city flight deals tailored to his unique itineraries. He often needs flexible travel plans, so he appreciates the app's free cancellation policy for both flights and car rentals. David also uses the car rental feature to secure vehicles for exploring off-the-beaten-path locations and enjoys browsing through the extensive selection of hotels to find budget-friendly and unique accommodations for his travels.",
      image: Persona3Avatar,
    },
    {
      id: "cooo1rddnifc7oc7i680",
      name: "Emma Johnson",
      personality: "Detail-oriented and pragmatic",
      description:
        "Emma relies on the app to manage her flight itineraries and access her boarding passes easily. She appreciates the option to track any flight in real time and check her booking reference and baggage allowance without an internet connection. For her longer trips, she uses the app to book hotels with free cancellation policies, ensuring flexibility with her plans.",
      image: Persona4Avatar,
    },
  ];

  let personaID = personas[1].id;

  const messages = [
    {
      from: "persona",
      text: "I'm Rachel, a social media manager who loves spontaneous weekend getaways. Feel free to ask me anything about using the app for travel planning!",
    },
    {
      from: "user",
      text: "Tell me what you think about the hotel filters on the app.",
    },
    {
      from: "persona",
      text: "Sure! I think the hotel filters are pretty basic. I'd love to see more options, like filtering for boutique hotels or specific amenities like rooftop pools.",
    },
    {
      from: "user",
      text: "How do you feel about our new feature that bundles flights and hotel deals automatically for you?",
    },
    {
      from: "persona",
      text: "Honestly, I'm not sure if I would use that. I like to have some control over my choices, especially for hotels, to ensure they match my aesthetic.",
    },
    {
      from: "user",
      text: "Do you have any suggestions on how we could improve this feature?",
    },
    {
      from: "persona",
      text: "Maybe if you added an option to customize the bundles, like allowing users to select their preferred hotel brands or specific amenities, it would make the feature more appealing.",
    },
  ];
</script>

<div class="flex flex-1 flex-col">
  <!-- TODO: Uncomment this when ready -->
  <!-- {#if $organization.plan === entities.OrganizationPlan.STARTER || $organization.plan === entities.OrganizationPlan.TRIAL || $organization.plan === entities.OrganizationPlan.DEMO}
    <div
      class="mb-4 flex w-full flex-row items-center justify-center gap-x-4 gap-y-2 rounded-lg border border-border bg-muted/40 px-3.5 py-2.5 shadow-sm"
    >
      {#if $organization.plan === entities.OrganizationPlan.STARTER}
        <h3 class="text-wrap font-semibold leading-none tracking-tight text-muted-foreground">
          Unlock interactive customer personas
        </h3>
      {:else if $organization.plan === entities.OrganizationPlan.TRIAL}
        <h3 class="text-wrap font-semibold leading-none tracking-tight text-muted-foreground">
          Unlock interactive customer personas
        </h3>
      {:else if $organization.plan === entities.OrganizationPlan.DEMO}
        <h3 class="text-wrap font-semibold leading-none tracking-tight text-muted-foreground">
          While in demo mode you can't interact with customer personas
        </h3>
      {/if}
      <Button href="/#pricing" variant="outline" size="sm" class="rounded-full px-4 py-2.5 text-sm font-semibold">
        Upgrade now &rarr;
      </Button>
    </div>
  {/if} -->
  <div class="relative flex flex-1">
    <!-- TODO: Delete this when ready -->
    <div
      class="absolute -left-6 -top-6 z-30 flex h-[calc(100%+3rem)] w-[calc(100%+3rem)] cursor-not-allowed select-none items-center justify-center backdrop-blur-sm"
    >
      <h1 class="text-2xl font-semibold text-muted-foreground/90">Customer Personas will be available soon...</h1>
    </div>
    <div class="grid w-full grid-cols-2 gap-4 lg:grid-cols-3">
      <div
        class="flex flex-1 flex-col items-center justify-center rounded-lg border bg-muted/40 p-6 shadow-sm lg:col-span-2"
      >
        <div class="flex w-full flex-1 flex-col">
          <div class="flex-shrink-0 cursor-pointer rounded-lg border bg-background px-4 py-2">
            <div class="flex w-full items-center justify-start gap-3">
              <Picture size="xl" src={personas[1].image} fallback={personas[1].name} />
              <div class="-mt-px flex items-center justify-center gap-1 overflow-hidden">
                <div class="flex shrink-0 items-center gap-2">
                  <h1 class="text-lg font-semibold">{personas[1].name}</h1>
                </div>
                <h3
                  class="inline-block w-full overflow-hidden text-ellipsis text-nowrap text-left text-sm font-normal text-muted-foreground"
                >
                  / {personas[1].description}
                </h3>
              </div>
              <CaretSort class="ml-auto h-5 w-5 shrink-0 opacity-50" />
            </div>
          </div>
          <div class="relative flex w-full flex-1">
            <div
              class="pointer-events-none absolute left-0 top-0 z-10 h-full w-full bg-muted"
              style="mask-image: linear-gradient(to bottom, rgba(241, 245, 249, 0.4) 0%, rgba(0, 0, 0, 0) 10%, rgba(0, 0, 0, 0) 90%, rgba(241, 245, 249, 0.4) 100%);"
            ></div>
            <div class="absolute left-0 top-0 flex h-full w-full flex-col gap-6 overflow-y-auto pb-6 pt-4">
              {#each messages as message}
                <div class={`flex ${message.from === "user" ? "justify-end" : "justify-start"}`}>
                  <div
                    class={`relative max-w-[80%] rounded-lg p-3 ${
                      message.from === "user" ? "bg-primary text-primary-foreground" : "bg-muted"
                    }`}
                  >
                    {message.text}
                    <div
                      class={`absolute h-0 w-0 border-8 ${
                        message.from === "user"
                          ? "-bottom-2.5 right-0 border-primary border-b-transparent border-l-transparent"
                          : "-top-2.5 left-0 border-muted border-r-transparent border-t-transparent"
                      }`}
                    ></div>
                  </div>
                </div>
              {/each}
            </div>
          </div>
        </div>
        <div
          class="relative w-full shrink-0 self-end overflow-hidden rounded-lg border bg-background focus-within:ring-1 focus-within:ring-ring"
        >
          <Textarea
            placeholder={`Ask anything to ${personas[1].name.split(" ")[0]}...`}
            class="mb-3 min-h-12 resize-none border-0 p-3 shadow-none focus-visible:ring-0"
          />
          <div class="flex items-center gap-2 p-3 pt-0">
            <Tooltip.Root>
              <Tooltip.Trigger asChild let:builder>
                <Button builders={[builder]} variant="outline" size="icon">
                  <Paperclip class="size-4" />
                </Button>
              </Tooltip.Trigger>
              <Tooltip.Content side="top">Attach File</Tooltip.Content>
            </Tooltip.Root>
            <Tooltip.Root>
              <Tooltip.Trigger asChild let:builder>
                <Button builders={[builder]} variant="outline" size="icon">
                  <Images class="size-4" />
                </Button>
              </Tooltip.Trigger>
              <Tooltip.Content side="top">Attach Image</Tooltip.Content>
            </Tooltip.Root>
            <Button class="ml-auto gap-1.5">
              Send
              <CornerDownLeft class="size-3.5" />
            </Button>
          </div>
        </div>
      </div>
      <div class="relative flex flex-1 flex-col items-center justify-start rounded-lg border bg-muted/40 shadow-sm">
        <Tabs.Root
          class="absolute left-0 top-0 h-full w-full overflow-y-auto"
          bind:value={personaID}
          orientation="vertical"
        >
          <Tabs.List class="h-auto w-full flex-col bg-transparent p-6">
            {#each personas as persona}
              <Tabs.Trigger
                value={persona.id}
                class={"z-10 h-24 w-full shrink-0 rounded-lg " +
                  (personaID === persona.id ? "border opacity-100" : "opacity-50")}
              >
                <div class="flex w-full items-center justify-start gap-4 px-4">
                  <Picture size="3xl" src={persona.image} fallback={persona.name} />
                  <div class="-mt-px flex flex-col items-start justify-center gap-0.5">
                    <div class="flex items-center gap-2">
                      <h1 class="text-xl font-semibold">{persona.name}</h1>
                    </div>
                    <h3 class="inline-block text-sm text-muted-foreground">{persona.personality}</h3>
                  </div>
                </div>
              </Tabs.Trigger>
            {/each}
            <div class="absolute bottom-0 left-0 z-0 flex h-full w-full p-6">
              <Tabs.Trigger value="none" class="h-24 w-full shrink-0 self-end rounded-lg border-2 border-dashed">
                <div class="flex w-full items-center justify-start gap-4 px-4">
                  <PlusCircled class="h-12 w-12 shrink-0 text-muted-foreground" />
                  <div class="-mt-px flex flex-col items-start justify-center gap-0.5">
                    <div class="flex items-center gap-2">
                      <h1 class="text-xl font-semibold text-muted-foreground">Customer Persona</h1>
                    </div>
                    <h3 class="inline-block text-sm text-muted-foreground">Synthesize another Persona</h3>
                  </div>
                </div>
              </Tabs.Trigger>
            </div>
          </Tabs.List>
        </Tabs.Root>
      </div>
    </div>
  </div>
</div>
