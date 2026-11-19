<script lang="ts">
  import '../app.css';
  import { onMount } from 'svelte';
  import type { Snippet } from 'svelte';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { Menu } from '@lucide/svelte';
  import Sidebar from '$lib/components/layout/sidebar.svelte';
  import UserMenu from '$lib/components/layout/user-menu.svelte';
  import { api } from '$lib/api/client';
  import { applyBootstrap, currentActor } from '$lib/stores/session';
  import { authReady, restoreSession } from '$lib/stores/auth';
  import { APP_NAME } from '$lib/constants/labels';
  import { accessLabelFor, canAccessPath } from '$lib/constants/access';

  let { children }: { children: Snippet } = $props();
  let mobileOpen = $state(false);
  let bootstrapError = $state('');
  let bootstrapped = $state(false);
  let unauthorized = $derived(
    !!$currentActor && !canAccessPath($currentActor.role, $page.url.pathname)
  );
  let deniedLabel = $derived(accessLabelFor($page.url.pathname));

  onMount(async () => {
    await restoreSession();
  });

  async function loadBootstrap() {
    try {
      applyBootstrap(await api.bootstrap());
      bootstrapped = true;
      bootstrapError = '';
    } catch (error) {
      bootstrapError = error instanceof Error ? error.message : 'Gagal memuat master data.';
    }
  }

  $effect(() => {
    if (!$authReady) return;
    const path = $page.url.pathname;
    const publicPaths = ['/login', '/user-guide'];
    if (!$currentActor && !publicPaths.includes(path)) {
      void goto('/login');
      return;
    }
    if ($currentActor && path === '/login') {
      void goto('/');
      return;
    }
    if ($currentActor && !bootstrapped) {
      void loadBootstrap();
    }
  });
</script>

{#if $page.url.pathname === '/login'}
  <div class="min-h-screen">{@render children()}</div>
{:else if !$authReady}
  <div class="flex min-h-screen items-center justify-center">
    <div class="nb-card px-6 py-4 font-mono text-xs font-bold uppercase">Memeriksa sesi…</div>
  </div>
{:else if !$currentActor && $page.url.pathname === '/user-guide'}
  <div class="min-h-screen">{@render children()}</div>
{:else if !$currentActor}
  <div class="flex min-h-screen items-center justify-center">
    <div class="nb-card px-6 py-4 font-mono text-xs font-bold uppercase">Mengalihkan ke login…</div>
  </div>
{:else}
  <div class="min-h-screen">
    <Sidebar open={mobileOpen} onclose={() => (mobileOpen = false)} />

    <div class="flex min-h-screen min-w-0 flex-col md:pl-72">
      <header class="sticky top-0 z-20 flex flex-wrap items-center justify-between gap-3 border-b-2 border-black bg-card px-4 py-3">
        <div class="flex items-center gap-3">
          <button
            class="border-2 border-black bg-primary p-2 shadow-nbsm md:hidden"
            onclick={() => (mobileOpen = true)}
            aria-label="Buka menu"
          >
            <Menu size={18} />
          </button>
          <div>
            <div class="nb-kicker">Ruang kerja</div>
            <div class="font-mono text-sm font-bold uppercase tracking-[0.12em]">{APP_NAME}</div>
          </div>
        </div>
        <UserMenu />
      </header>

      {#if bootstrapError}
        <div class="border-b-2 border-black bg-destructive px-4 py-2 font-mono text-xs font-bold">
          {bootstrapError}
        </div>
      {/if}

      <main class="mx-auto w-full max-w-[1600px] flex-1 p-4 md:p-6">
        {#if unauthorized}
          <div class="mx-auto mt-6 max-w-2xl border-2 border-black bg-card p-6 shadow-nb">
            <div class="nb-kicker">Akses terbatas</div>
            <h1 class="mt-2 text-3xl">Halaman ini tidak tersedia untuk akun Anda</h1>
            <p class="mt-3 text-sm leading-relaxed text-muted-foreground">
              Menu {deniedLabel} hanya dapat dibuka oleh peran yang berwenang. Menu yang tampil di
              samping sudah disesuaikan dengan peran akun Anda.
            </p>
            <a
              href="/"
              class="mt-5 inline-flex h-10 items-center border-2 border-black bg-primary px-4 font-mono text-xs font-bold uppercase shadow-nbsm"
            >
              Kembali ke dashboard
            </a>
          </div>
        {:else}
          {@render children()}
        {/if}
      </main>

      <footer class="border-t-2 border-black bg-card px-4 py-3 font-mono text-[10px] text-muted-foreground md:px-6">
        MATA SIDIK · data periode aktif
      </footer>
    </div>
  </div>
{/if}
