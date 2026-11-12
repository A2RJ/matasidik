<script lang="ts">
  import type { Snippet } from 'svelte';
  import { cn } from '$lib/utils';
  let { class: className = '', children }: { class?: string; children: Snippet } = $props();
</script>

<div class={cn('border-2 border-black bg-card shadow-nb', className)}>
  {@render children()}
</div>
