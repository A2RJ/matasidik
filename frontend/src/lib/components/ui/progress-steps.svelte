<script lang="ts">
  import { cn } from '$lib/utils';
  let {
    steps = [],
    current = '',
    percent = 0
  }: { steps?: { key: string; label: string }[]; current?: string; percent?: number } = $props();
</script>

<div class="space-y-3">
  <div class="h-4 w-full border-2 border-black bg-muted">
    <div class="h-full bg-primary transition-all" style={`width:${Math.max(0, Math.min(100, percent))}%`}></div>
  </div>
  <div class="grid gap-2 sm:grid-cols-4">
    {#each steps as step, index (step.key)}
      <div
        class={cn(
          'border-2 border-black px-3 py-2 font-mono text-[10px] font-bold uppercase tracking-wide',
          step.key === current ? 'bg-primary shadow-nbsm' : index < steps.findIndex((s) => s.key === current) ? 'bg-success' : 'bg-card'
        )}
      >
        {index + 1}. {step.label}
      </div>
    {/each}
  </div>
</div>
