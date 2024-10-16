<script lang="ts">
  import { setContext } from "svelte";

  setContext("content-scope", "Body");
</script>

<div class="mt-4">
  <slot />
</div>
