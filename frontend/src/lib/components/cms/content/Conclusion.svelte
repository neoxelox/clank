<script lang="ts">
  import { setContext } from "svelte";

  setContext("content-scope", "Conclusion");
</script>

<div>
  <hr class="mb-4 mt-16 border-t border-border" />
  <slot />
</div>
