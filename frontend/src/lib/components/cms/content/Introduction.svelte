<script lang="ts">
  import { setContext } from "svelte";

  setContext("content-scope", "Introduction");
</script>

<div>
  <slot />
</div>
