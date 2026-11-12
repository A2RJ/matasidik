<script lang="ts">
  import type { Snippet } from 'svelte';
  import { cn } from '$lib/utils';
  let { class: className = '', children }: { class?: string; children: Snippet } = $props();
</script>

<div class={cn('flex items-center gap-2 border-t-2 border-black bg-card px-4 py-3', className)}>
  {@render children()}
</div>
