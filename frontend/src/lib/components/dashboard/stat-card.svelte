<script lang="ts">
  import type { Snippet } from 'svelte';
  import { cn } from '$lib/utils';

  let {
    label,
    value,
    hint = '',
    tone = 'primary',
    icon
  }: {
    label: string;
    value: string;
    hint?: string;
    tone?: 'primary' | 'secondary' | 'success' | 'warning' | 'danger' | 'muted';
    icon?: Snippet;
  } = $props();

  const tones = {
    primary: 'bg-primary',
    secondary: 'bg-secondary',
    success: 'bg-success',
    warning: 'bg-warning',
    danger: 'bg-destructive',
    muted: 'bg-muted'
  };
</script>

<div class={cn('border-2 border-black p-4 shadow-nb', tones[tone])}>
  <div class="flex items-start justify-between gap-2">
    <div class="font-mono text-[10px] font-bold uppercase tracking-[0.12em]">{label}</div>
    {#if icon}<div class="border-2 border-black bg-card p-1">{@render icon()}</div>{/if}
  </div>
  <div class="mt-3 text-3xl font-bold leading-none">{value}</div>
  {#if hint}<div class="mt-2 font-mono text-[10px] leading-relaxed text-foreground/70">{hint}</div>{/if}
</div>
