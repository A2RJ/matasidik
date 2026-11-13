<script lang="ts">
  import { onMount } from 'svelte';
  import { Maximize2, Minimize2, RotateCcw, X, ZoomIn, ZoomOut } from '@lucide/svelte';
  import type { HighlightBox } from '$lib/content/user-guide';

  let {
    src,
    alt,
    caption,
    detail = '',
    overlay,
    inline = 'scroll',
    access = ''
  }: {
    src: string;
    alt: string;
    caption: string;
    detail?: string;
    overlay?: HighlightBox;
    inline?: 'scroll' | 'full';
    access?: string;
  } = $props();

  let dialog = $state<HTMLDialogElement | null>(null);
  let zoom = $state(1);
  let isFullscreen = $state(false);
  let expanded = $state(false);
  let zoomLabel = $derived(`${Math.round(zoom * 100)}%`);
  let immersive = $derived(isFullscreen || expanded);

  function open() {
    if (!dialog) return;
    zoom = 1;
    dialog.showModal();
    document.body.style.overflow = 'hidden';
  }

  function close() {
    dialog?.close();
  }

  function reset() {
    zoom = 1;
  }

  function zoomBy(delta: number) {
    zoom = Math.min(3, Math.max(0.5, Math.round((zoom + delta) * 100) / 100));
  }

  async function toggleFullscreen() {
    if (!dialog) return;
    if (document.fullscreenElement) {
      await document.exitFullscreen().catch(() => undefined);
      isFullscreen = false;
      return;
    }
    if (dialog.requestFullscreen) {
      try {
        await dialog.requestFullscreen();
        isFullscreen = true;
        expanded = false;
        return;
      } catch {
        // Browser menolak native fullscreen (mis. iframe/headless): pakai layar penuh CSS.
      }
    }
    expanded = !expanded;
  }

  function onClose() {
    document.body.style.overflow = '';
    zoom = 1;
    isFullscreen = false;
    expanded = false;
  }

  onMount(() => {
    const onChange = () => {
      isFullscreen = document.fullscreenElement === dialog;
    };
    document.addEventListener('fullscreenchange', onChange);
    return () => {
      document.removeEventListener('fullscreenchange', onChange);
      document.body.style.overflow = '';
    };
  });
</script>

<svelte:window
  onkeydown={(event) => {
    if (event.key === 'Escape' && dialog?.open) close();
  }}
/>

<figure class="overflow-hidden border-2 border-black bg-card shadow-nb">
  <button
    type="button"
    class="group relative block w-full cursor-zoom-in text-left"
    onclick={open}
    aria-label={`Buka ${alt} dalam ukuran penuh`}
  >
    <div
      class="bg-white"
      class:max-h-[620px]={inline === 'scroll'}
      class:overflow-auto={inline === 'scroll'}
      class:nb-scroll={inline === 'scroll'}
    >
      <div class="relative min-w-full">
        <img {src} {alt} class="block w-full" loading="lazy" />
        {#if overlay}
          <div
            class="pointer-events-none absolute border-4 border-green-600 shadow-[0_0_0_2px_rgba(255,255,255,0.9)]"
            style={`left:${overlay.left}%;top:${overlay.top}%;width:${overlay.width}%;height:${overlay.height}%;`}
          >
            <span class="absolute left-1 top-1 bg-green-600 px-1.5 py-0.5 font-mono text-[9px] font-bold uppercase tracking-[0.08em] text-white">
              {overlay.label}
            </span>
          </div>
        {/if}
      </div>
    </div>
    <span class="pointer-events-none absolute right-2 top-2 flex items-center gap-1 border-2 border-black bg-primary px-2 py-1 font-mono text-[9px] font-bold uppercase shadow-nbsm">
      <Maximize2 size={12} /> Zoom / layar penuh
    </span>
  </button>

  <figcaption class="border-t-2 border-black bg-card px-3 py-3">
    <div class="flex flex-wrap items-center justify-between gap-2">
      <div class="font-mono text-[10px] font-bold uppercase tracking-[0.08em]">{caption}</div>
      <span class="border border-black/30 bg-muted px-1.5 py-0.5 font-mono text-[9px] uppercase text-muted-foreground">
        klik untuk zoom
      </span>
    </div>
    {#if detail}<p class="mt-1 text-xs leading-relaxed text-muted-foreground">{detail}</p>{/if}
    {#if access}
      <div class="mt-2 border-2 border-black bg-secondary px-2 py-1 font-mono text-[10px] leading-relaxed">
        <b>Cara akses:</b> {access}
      </div>
    {/if}
  </figcaption>
</figure>

<dialog
  bind:this={dialog}
  onclose={onClose}
  class="m-auto max-h-[96vh] w-[96vw] max-w-[1500px] overflow-hidden border-2 border-black bg-card p-0 shadow-nb backdrop:bg-black/70"
  class:w-[100vw]={expanded}
  class:h-[100vh]={expanded}
  class:max-h-none={expanded}
  class:max-w-none={expanded}
  class:border-0={expanded}
  aria-label={`Pratinjau ${caption}`}
>
  <div class="flex flex-wrap items-center justify-between gap-2 border-b-2 border-black bg-secondary px-3 py-2">
    <div class="min-w-0">
      <div class="nb-kicker">Pratinjau screenshot</div>
      <div class="truncate font-mono text-[10px] font-bold">{caption}</div>
    </div>
    <div class="flex flex-wrap items-center gap-1">
      <button
        type="button"
        class="inline-flex h-8 items-center gap-1 border-2 border-black bg-card px-2 font-mono text-[10px] font-bold uppercase shadow-nbsm disabled:opacity-40"
        onclick={() => zoomBy(-0.25)}
        disabled={zoom <= 0.5}
        aria-label="Perkecil"
      >
        <ZoomOut size={14} />
      </button>
      <span class="min-w-12 border-2 border-black bg-card px-2 py-1 text-center font-mono text-[10px] font-bold">
        {zoomLabel}
      </span>
      <button
        type="button"
        class="inline-flex h-8 items-center gap-1 border-2 border-black bg-card px-2 font-mono text-[10px] font-bold uppercase shadow-nbsm disabled:opacity-40"
        onclick={() => zoomBy(0.25)}
        disabled={zoom >= 3}
        aria-label="Perbesar"
      >
        <ZoomIn size={14} />
      </button>
      <button
        type="button"
        class="inline-flex h-8 items-center gap-1 border-2 border-black bg-card px-2 font-mono text-[10px] font-bold uppercase shadow-nbsm"
        onclick={reset}
      >
        <RotateCcw size={13} /> Fit
      </button>
      <button
        type="button"
        class="inline-flex h-8 items-center gap-1 border-2 border-black bg-card px-2 font-mono text-[10px] font-bold uppercase shadow-nbsm"
        onclick={toggleFullscreen}
      >
        {#if immersive}<Minimize2 size={13} /> Keluar{:else}<Maximize2 size={13} /> Layar penuh{/if}
      </button>
      <button
        type="button"
        class="inline-flex h-8 items-center gap-1 border-2 border-black bg-destructive px-2 font-mono text-[10px] font-bold uppercase shadow-nbsm"
        onclick={close}
      >
        <X size={14} /> Tutup
      </button>
    </div>
  </div>

  <div class="overflow-auto bg-neutral-900 p-3 nb-scroll" style="max-height: calc(100vh - 150px);">
    <div class="relative mx-auto" style={`width:${zoom * 100}%;`}>
      <img {src} {alt} class="block w-full" />
      {#if overlay}
        <div
          class="pointer-events-none absolute border-4 border-green-500 shadow-[0_0_0_2px_rgba(255,255,255,0.9)]"
          style={`left:${overlay.left}%;top:${overlay.top}%;width:${overlay.width}%;height:${overlay.height}%;`}
        >
          <span class="absolute left-1 top-1 bg-green-600 px-1.5 py-0.5 font-mono text-[9px] font-bold uppercase tracking-[0.08em] text-white">
            {overlay.label}
          </span>
        </div>
      {/if}
    </div>
  </div>

  <div class="border-t-2 border-black bg-card px-3 py-2 text-xs text-muted-foreground">
    {detail || caption}{#if access} · <b>Cara akses:</b> {access}{/if} · Gunakan tombol zoom atau Layar penuh untuk membaca detail tabel.
  </div>
</dialog>
