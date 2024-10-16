<script lang="ts">
  import { setContext } from "svelte";

  setContext("content-scope", "Description");
</script>

<p class="mt-6 text-xl leading-8">
  <slot />
</p>
