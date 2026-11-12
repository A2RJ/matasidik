<script lang="ts">
  import type { Snippet } from 'svelte';
  let {
    kicker = 'MATA SIDIK',
    title,
    description = '',
    actions
  }: { kicker?: string; title: string; description?: string; actions?: Snippet } = $props();
</script>

<div class="mb-6 flex flex-col gap-4 border-2 border-black bg-card p-4 shadow-nb md:flex-row md:items-end md:justify-between">
  <div class="max-w-3xl">
    <div class="nb-kicker">{kicker}</div>
    <h1 class="mt-1 text-3xl md:text-4xl">{title}</h1>
    {#if description}<p class="mt-2 max-w-2xl text-sm text-muted-foreground">{description}</p>{/if}
  </div>
  {#if actions}
    <div class="flex flex-wrap gap-2">{@render actions()}</div>
  {/if}
</div>
