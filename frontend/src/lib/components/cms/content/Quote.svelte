<script lang="ts">
  import { setContext } from "svelte";

  setContext("content-scope", "Quote");
</script>

<blockquote>
  <p class="mt-8 border-l-2 border-primary pl-9 font-medium italic">
    <slot />
  </p>
</blockquote>
