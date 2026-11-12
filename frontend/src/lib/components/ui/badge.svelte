<script lang="ts">
  import type { Snippet } from 'svelte';
  import { cn } from '$lib/utils';
  import type { Tone } from '$lib/constants/enums';

  let {
    tone = 'neutral',
    class: className = '',
    children
  }: { tone?: Tone; class?: string; children: Snippet } = $props();

  const tones: Record<Tone, string> = {
    neutral: 'bg-muted text-muted-foreground',
    info: 'bg-secondary text-secondary-foreground',
    success: 'bg-success text-success-foreground',
    warning: 'bg-warning text-warning-foreground',
    danger: 'bg-destructive text-destructive-foreground',
    final: 'bg-primary text-primary-foreground',
    dark: 'bg-black text-white',
    slate: 'bg-slate-200 text-slate-700'
  };
</script>

<span
  class={cn(
    'inline-flex items-center border-2 border-black px-2 py-0.5 font-mono text-[10px] font-bold uppercase tracking-[0.06em]',
    tones[tone],
    className
  )}
>
  {@render children()}
</span>
