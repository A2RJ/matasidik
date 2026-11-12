<script lang="ts">
  import type { Snippet } from 'svelte';
  import { cn } from '$lib/utils';
  let { for: htmlFor = '', class: className = '', children }: { for?: string; class?: string; children: Snippet } = $props();
</script>

<label for={htmlFor} class={cn('nb-label', className)}>
  {@render children()}
</label>
