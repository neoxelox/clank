<script lang="ts">
  import { Transition } from "@rgossiaux/svelte-headlessui";

  setTimeout(() => {
    const component = document.getElementById("curtain");
    component?.parentNode?.removeChild(component);
  }, 1005);
</script>

<Transition
  id="curtain"
  as="div"
  show={true}
  enter="transition-opacity duration-1000"
  enterFrom="opacity-100"
  enterTo="opacity-0"
  class="no-scrollbar absolute z-50 h-full w-full bg-background"
>
</Transition>
