<script lang="ts">
  export let open: boolean;
</script>

<svg
  xmlns="http://www.w3.org/2000/svg"
  aria-hidden="true"
  class="h-3.5 w-3.5 overflow-visible stroke-secondary-foreground"
  fill="none"
  stroke-width={2}
  stroke-linecap="round"
>
  <path
    d="M0 1H14M0 7H14M0 13H14"
    class={"origin-center transition " + (open ? "scale-90 opacity-0 " : "scale-100 opacity-100")}
  />
  <path
    d="M2 2L12 12M12 2L2 12"
    class={"origin-center transition " + (!open ? "scale-90 opacity-0" : "scale-100 opacity-100")}
  />
</svg>
