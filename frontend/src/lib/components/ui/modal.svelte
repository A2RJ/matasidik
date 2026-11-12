<script lang="ts">
  import type { Snippet } from 'svelte';
  import { X } from '@lucide/svelte';
  let {
    open = false,
    title = '',
    description = '',
    onclose,
    children,
    footer
  }: {
    open?: boolean;
    title?: string;
    description?: string;
    onclose?: () => void;
    children: Snippet;
    footer?: Snippet;
  } = $props();
</script>

{#if open}
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4">
    <button
      type="button"
      class="absolute inset-0 h-full w-full cursor-default bg-transparent"
      aria-label="Tutup modal"
      onclick={onclose}
    ></button>
    <div
      class="nb-card relative max-h-[90vh] w-full max-w-2xl overflow-y-auto nb-scroll"
      role="dialog"
      aria-modal="true"
      aria-label={title}
      tabindex="-1"
    >
      <div class="flex items-start justify-between border-b-2 border-black bg-secondary px-4 py-3">
        <div>
          <h2 class="font-mono text-sm font-bold uppercase tracking-[0.08em]">{title}</h2>
          {#if description}<p class="mt-1 text-xs text-foreground/70">{description}</p>{/if}
        </div>
        <button class="border-2 border-black bg-card p-1 shadow-nbsm" onclick={onclose} aria-label="Tutup">
          <X size={16} />
        </button>
      </div>
      <div class="p-4">{@render children()}</div>
      {#if footer}
        <div class="flex justify-end gap-2 border-t-2 border-black bg-muted px-4 py-3">{@render footer()}</div>
      {/if}
    </div>
  </div>
{/if}
