<script lang="ts">
  import type { Snippet } from 'svelte';
  import { cn } from '$lib/utils';
  let { class: className = '', children }: { class?: string; children: Snippet } = $props();
</script>

<h3 class={cn('font-mono text-xs font-bold uppercase tracking-[0.1em]', className)}>
  {@render children()}
</h3>
