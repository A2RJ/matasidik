<script lang="ts">
  import { goto } from '$app/navigation';
  import { currentActor, bootstrap } from '$lib/stores/session';
  import { actingAs, logout, startActingAs, stopActingAs } from '$lib/stores/auth';
  import { ROLE_LABELS } from '$lib/constants/enums';
  import Button from '$lib/components/ui/button/button.svelte';
  import Select from '$lib/components/ui/select.svelte';
  import Badge from '$lib/components/ui/badge.svelte';
  import { LogOut, ShieldCheck, UserCog } from '@lucide/svelte';

  let target = $state('');
  let busy = $state(false);
  let error = $state('');

  let targets = $derived(
    ($bootstrap?.users ?? []).filter(
      (u) => u.aktif && u.role !== 'admin' && u.id !== $currentActor?.id
    )
  );
  let targetOptions = $derived(targets.map((u) => ({ value: u.id, label: `${u.nama} · ${ROLE_LABELS[u.role]}` })));

  async function act() {
    if (!target) {
      error = 'Pilih user target.';
      return;
    }
    busy = true;
    error = '';
    try {
      await startActingAs(target);
      await goto('/');
    } catch (err) {
      error = err instanceof Error ? err.message : 'ActingAs gagal.';
    } finally {
      busy = false;
    }
  }

  async function stop() {
    busy = true;
    error = '';
    try {
      await stopActingAs();
      await goto('/');
    } catch (err) {
      error = err instanceof Error ? err.message : 'Gagal kembali ke admin.';
    } finally {
      busy = false;
    }
  }

  async function signOut() {
    await logout();
    await goto('/login');
  }
</script>

<div class="flex flex-wrap items-center gap-2">
  {#if $currentActor}
    <div class="border-2 border-black bg-card px-2 py-1">
      <div class="font-mono text-[10px] font-bold uppercase">{ $currentActor.nama }</div>
      <div class="font-mono text-[9px] text-muted-foreground">
        {ROLE_LABELS[$currentActor.role]}{#if $currentActor.position_nama} · {$currentActor.position_nama}{/if}
      </div>
    </div>
    {#if $actingAs}
      <Badge tone="warning">ActingAs · {$actingAs.nama}</Badge>
      <Button size="sm" variant="outline" disabled={busy} onclick={stop}><ShieldCheck size={13} /> Kembali</Button>
    {:else if $currentActor.role === 'admin'}
      <div class="hidden min-w-56 xl:block">
        <Select bind:value={target} options={targetOptions} placeholder="Bertindak sebagai…" />
      </div>
      <Button size="sm" variant="secondary" disabled={busy} onclick={act}><UserCog size={13} /> ActingAs</Button>
    {/if}
    <Button size="sm" variant="ghost" onclick={signOut}><LogOut size={13} /> Keluar</Button>
  {/if}
</div>
{#if error}<div class="w-full font-mono text-[10px] font-bold text-destructive">{error}</div>{/if}
