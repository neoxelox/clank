<script lang="ts">
  import { page } from "$app/stores";
  import { Curtain } from "$lib/components/ui/curtain";
  import { Progress } from "$lib/components/ui/progress";
  import { Toaster } from "$lib/components/ui/sonner";
  import * as payloads from "$lib/payloads";
  import { api } from "$lib/services/api";
  import { invitations, organization, user, users } from "$lib/stores";
  import { delay } from "$lib/utils/delay";

  let progress: number = 0;
  let load: Promise<void>;
  (() => {
    if (!$page.url.pathname.startsWith("/dash/signin")) {
      progress = 0;
      setTimeout(() => (progress = 33), 250);
      setTimeout(() => (progress = 66), 500);
      setTimeout(() => (progress = 100), 750);

      load = (async () => {
        let start = new Date();

        const [userResponse, organizationResponse, usersResponse, invitationsResponse] = await Promise.all([
          api.get<payloads.GetMeResponse>("/user"),
          api.get<payloads.GetOrganizationResponse>("/organization"),
          api.get<payloads.ListUsersResponse>("/users"),
          api.get<payloads.ListInvitationsResponse>("/invitations"),
        ]);

        $user = payloads.toUser(userResponse);
        $organization = payloads.toOrganization(organizationResponse);
        $users = usersResponse.users.map((user) => payloads.toUser(user));
        $invitations = invitationsResponse.invitations.map((invitation) => payloads.toInvitation(invitation));

        await delay(1000 - (new Date() - start));
      })();
    }
  })();
</script>

{#await load}
  <main class="no-scrollbar mx-auto flex h-full w-full max-w-7xl flex-col overflow-y-auto px-6 py-16 lg:px-8">
    <div class="mx-auto my-auto flow-root max-w-lg">
      <div class="mx-auto grid w-96 gap-12">
        <h1 class="text-center font-cal text-3xl text-foreground sm:text-4xl md:text-5xl">Welcome back</h1>
        <Progress value={progress} max={100} />
      </div>
    </div>
  </main>
{:then}
  <Curtain />
  <slot />
{:catch}
  <Curtain />
  <slot />
{/await}
<Toaster />
