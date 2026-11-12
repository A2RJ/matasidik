<script lang="ts">
  import { page } from '$app/stores';
  import { UserRound, X } from '@lucide/svelte';
  import { APP_NAME, APP_REGION, APP_SUBTITLE, NAV_ITEMS } from '$lib/constants/labels';
  import { currentActor } from '$lib/stores/session';
  import { cn } from '$lib/utils';

  let { open = false, onclose }: { open?: boolean; onclose?: () => void } = $props();

  function isActive(href: string, pathname: string): boolean {
    if (href === '/') return pathname === '/';
    return pathname.startsWith(href);
  }
</script>

{#if open}
  <div class="fixed inset-0 z-30 bg-black/30 md:hidden" onclick={onclose} role="presentation"></div>
{/if}

<aside
  class={cn(
    'fixed inset-y-0 left-0 z-40 flex w-72 flex-col border-r-2 border-black bg-card transition-transform md:translate-x-0',
    open ? 'translate-x-0' : '-translate-x-full'
  )}
>
  <div class="border-b-2 border-black bg-primary p-4">
    <div class="flex items-start justify-between">
      <div>
        <div class="font-mono text-[10px] font-bold uppercase tracking-[0.18em]">{APP_REGION}</div>
        <div class="mt-1 text-2xl font-bold leading-none">{APP_NAME}</div>
        <div class="mt-2 text-[11px] font-medium">{APP_SUBTITLE}</div>
      </div>
      <button class="border-2 border-black bg-card p-1 shadow-nbsm md:hidden" onclick={onclose} aria-label="Tutup menu">
        <X size={16} />
      </button>
    </div>
  </div>

  <nav class="flex-1 space-y-1 overflow-y-auto p-3 nb-scroll">
    <div class="nb-kicker mb-2 px-2">Navigasi</div>
    {#each NAV_ITEMS.filter((item) => !item.roles || ($currentActor && item.roles.includes($currentActor.role))) as item (item.href)}
      {@const active = isActive(item.href, $page.url.pathname)}
      {@const Icon = item.icon}
      <a href={item.href} class={cn('nb-sidebar-link', active && 'active')} onclick={onclose}>
        <Icon size={18} strokeWidth={2.4} />
        <span class="flex flex-col">
          <span>{item.label}</span>
          <span class="text-[9px] font-normal normal-case tracking-normal text-muted-foreground">{item.description}</span>
        </span>
      </a>
    {/each}
  </nav>

  <div class="border-t-2 border-black p-3">
    <div class="border-2 border-black bg-secondary p-3 shadow-nbsm">
      <div class="font-mono text-[10px] font-bold uppercase tracking-[0.12em]">Sesi</div>
      <a
        href="/profile"
        class="mt-2 flex items-center gap-2 border-2 border-black bg-card px-2 py-2 font-mono text-[10px] font-bold uppercase shadow-nbsm"
        class:bg-primary={isActive('/profile', $page.url.pathname)}
        onclick={onclose}
      >
        <UserRound size={14} />
        Profil
        <span class="ml-auto text-[9px] font-normal normal-case tracking-normal text-muted-foreground">Akun</span>
      </a>
    </div>
  </div>
</aside>
