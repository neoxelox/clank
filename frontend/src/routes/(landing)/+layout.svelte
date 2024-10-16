<script lang="ts">
  import Cookies from "$lib/components/cms/Cookies.svelte";
  import Footer from "$lib/components/cms/Footer.svelte";
  import Header from "$lib/components/cms/Header.svelte";
</script>

<div class="relative flex h-full flex-col overflow-x-clip">
  <Header />
  <main>
    <slot />
  </main>
  <Footer />
</div>
<Cookies />
