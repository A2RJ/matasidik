<script lang="ts">
  import type { Snippet } from 'svelte';
  let { title = 'Belum ada data', description = '', action }: { title?: string; description?: string; action?: Snippet } = $props();
</script>

<div class="border-2 border-dashed border-black/40 bg-card p-10 text-center">
  <div class="mx-auto mb-3 flex h-12 w-12 items-center justify-center border-2 border-black bg-primary text-xl shadow-nbsm">∅</div>
  <h3 class="font-mono text-sm font-bold uppercase tracking-wide">{title}</h3>
  {#if description}<p class="mx-auto mt-2 max-w-md text-sm text-muted-foreground">{description}</p>{/if}
  {#if action}<div class="mt-4 flex justify-center">{@render action()}</div>{/if}
</div>
