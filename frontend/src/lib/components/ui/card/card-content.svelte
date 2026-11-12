<script lang="ts">
  import type { Snippet } from 'svelte';
  import { cn } from '$lib/utils';
  let { class: className = '', children }: { class?: string; children: Snippet } = $props();
</script>

<div class={cn('p-4', className)}>
  {@render children()}
</div>
