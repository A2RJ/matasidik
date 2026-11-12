<script lang="ts">
  import type { Snippet } from 'svelte';
  import { cn } from '$lib/utils';
  import { buttonVariants } from './variants';

  let {
    variant = 'default',
    size = 'default',
    class: className = '',
    type = 'button',
    disabled = false,
    title = '',
    onclick,
    children
  }: {
    variant?: 'default' | 'secondary' | 'outline' | 'ghost' | 'destructive' | 'success' | 'dark';
    size?: 'default' | 'sm' | 'lg' | 'icon';
    class?: string;
    type?: 'button' | 'submit' | 'reset';
    disabled?: boolean;
    title?: string;
    onclick?: (event: MouseEvent) => void;
    children: Snippet;
  } = $props();
</script>

<button
  {type}
  {disabled}
  {title}
  {onclick}
  class={cn(buttonVariants({ variant, size }), className)}
>
  {@render children()}
</button>
