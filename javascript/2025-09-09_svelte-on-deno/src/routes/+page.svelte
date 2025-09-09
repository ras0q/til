<script lang="ts">
  import type { PageProps } from "./$types.d.ts";

  let { data }: PageProps = $props();
</script>

<h1>Welcome to SvelteKit</h1>
<p>
  Visit <a href="https://svelte.dev/docs/kit">svelte.dev/docs/kit</a> to read the
  documentation
</p>
<p>Developing in {data.pwd}</p>
