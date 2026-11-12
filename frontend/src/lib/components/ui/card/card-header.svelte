<script lang="ts">
  import type { Snippet } from 'svelte';
  import { cn } from '$lib/utils';
  let { class: className = '', children }: { class?: string; children: Snippet } = $props();
</script>

<div class={cn('border-b-2 border-black bg-muted px-4 py-3', className)}>
  {@render children()}
</div>
