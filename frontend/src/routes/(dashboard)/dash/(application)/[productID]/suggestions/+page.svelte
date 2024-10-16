<script lang="ts">
  import SuggestionTable from "$lib/components/cmp/SuggestionTable.svelte";
  import { Button } from "$lib/components/ui/button";
  import * as entities from "$lib/entities";
  import { organization } from "$lib/stores";
</script>

<div class="h-full w-full">
  {#if $organization.plan === entities.OrganizationPlan.TRIAL || $organization.plan === entities.OrganizationPlan.DEMO}
    <div
      class="mb-4 flex w-full flex-row items-center justify-center gap-x-4 gap-y-2 rounded-lg border border-border bg-muted/40 px-3.5 py-2.5 shadow-sm"
    >
      {#if $organization.plan === entities.OrganizationPlan.TRIAL}
        <h3 class="text-wrap font-semibold leading-none tracking-tight text-muted-foreground">
          Analyze thousands of feedbacks a month per organization
        </h3>
      {:else if $organization.plan === entities.OrganizationPlan.DEMO}
        <h3 class="text-wrap font-semibold leading-none tracking-tight text-muted-foreground">
          While in demo mode new feedbacks won't be analyzed
        </h3>
      {/if}
      <Button href="/#pricing" variant="outline" size="sm" class="rounded-full px-4 py-2.5 text-sm font-semibold">
        Upgrade now &rarr;
      </Button>
    </div>
  {/if}
  <SuggestionTable />
</div>
