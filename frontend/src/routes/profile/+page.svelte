<script lang="ts">
  import { goto } from '$app/navigation';
  import { api } from '$lib/api/client';
  import { bootstrap, currentActor } from '$lib/stores/session';
  import { actingAs, clearAuth } from '$lib/stores/auth';
  import PageHeader from '$lib/components/ui/page-header.svelte';
  import Card from '$lib/components/ui/card/card.svelte';
  import Button from '$lib/components/ui/button/button.svelte';
  import Badge from '$lib/components/ui/badge.svelte';
  import Input from '$lib/components/ui/input.svelte';
  import PasswordInput from '$lib/components/ui/password-input.svelte';
  import { ROLE_LABELS } from '$lib/constants/enums';
  import { formatDateTime } from '$lib/utils';
  import { ArrowLeft, KeyRound, Save, ShieldCheck, UserRound } from '@lucide/svelte';

  let nama = $state('');
  let profileBusy = $state(false);
  let profileError = $state('');
  let profileMessage = $state('');

  let currentPassword = $state('');
  let newPassword = $state('');
  let confirmPassword = $state('');
  let passwordBusy = $state(false);
  let passwordError = $state('');
  let passwordMessage = $state('');

  $effect(() => {
    if ($currentActor) nama = $currentActor.nama;
  });

  let isActing = $derived(!!$actingAs);
  let initial = $derived(($currentActor?.nama ?? '?').trim().slice(0, 2).toUpperCase());
  let faskesNama = $derived(
    $currentActor?.faskes_id
      ? ($bootstrap?.faskes.find((f) => f.id === $currentActor?.faskes_id)?.nama ??
        $currentActor.faskes_id)
      : 'Dinas / tanpa faskes'
  );
  let bidangNama = $derived(
    $currentActor?.bidang_id
      ? ($bootstrap?.bidang.find((b) => b.id === $currentActor?.bidang_id)?.nama ??
        $currentActor.bidang_id)
      : '—'
  );

  async function saveProfile(event: SubmitEvent) {
    event.preventDefault();
    if (isActing) return;
    const value = nama.trim();
    if (value.length < 2) {
      profileError = 'Nama minimal 2 karakter.';
      return;
    }
    profileBusy = true;
    profileError = '';
    profileMessage = '';
    try {
      const user = await api.updateProfile(value);
      currentActor.set(user);
      profileMessage = 'Profil berhasil disimpan.';
    } catch (err) {
      profileError = err instanceof Error ? err.message : 'Gagal menyimpan profil.';
    } finally {
      profileBusy = false;
    }
  }

  async function changePassword(event: SubmitEvent) {
    event.preventDefault();
    if (isActing) return;
    passwordError = '';
    passwordMessage = '';
    if (newPassword.length < 8) {
      passwordError = 'Password baru minimal 8 karakter.';
      return;
    }
    if (newPassword !== confirmPassword) {
      passwordError = 'Konfirmasi password tidak sama.';
      return;
    }
    passwordBusy = true;
    try {
      await api.changePassword(currentPassword, newPassword);
      sessionStorage.setItem('matasidik_password_changed', '1');
      clearAuth();
      await goto('/login');
    } catch (err) {
      passwordError = err instanceof Error ? err.message : 'Gagal mengubah password.';
    } finally {
      passwordBusy = false;
    }
  }
</script>

<PageHeader
  kicker="Akun / profil"
  title="Profil akun"
  description="Kelola nama tampilan dan password akun Anda. Peran, jabatan, dan faskes tetap diatur oleh admin."
>
  {#snippet actions()}
    <a
      href="/"
      class="inline-flex h-10 items-center gap-2 border-2 border-black bg-card px-4 font-mono text-xs font-bold uppercase shadow-nbsm"
    >
      <ArrowLeft size={14} /> Dashboard
    </a>
  {/snippet}
</PageHeader>

{#if isActing}
  <div class="mb-4 border-2 border-black bg-warning p-3 font-mono text-xs font-bold">
    Anda sedang actingAs sebagai {$currentActor?.nama}. Kembali ke akun admin untuk mengubah profil.
  </div>
{/if}

<div class="grid gap-4 xl:grid-cols-[340px_1fr]">
  <Card>
    <div class="border-b-2 border-black bg-primary p-4">
      <div class="flex items-center gap-3">
        <div
          class="flex h-16 w-16 shrink-0 items-center justify-center border-2 border-black bg-card font-mono text-2xl font-bold shadow-nbsm"
        >
          {initial}
        </div>
        <div class="min-w-0">
          <div class="font-mono text-[10px] font-bold uppercase tracking-[0.12em]">Akun saya</div>
          <h2 class="mt-1 truncate text-2xl">{ $currentActor?.nama ?? '—' }</h2>
          <div class="font-mono text-[10px] text-muted-foreground">@{ $currentActor?.username ?? '—' }</div>
        </div>
      </div>
    </div>

    <div class="space-y-3 p-4 text-sm">
      <div class="flex flex-wrap items-center gap-2">
        <Badge tone="info">{ROLE_LABELS[$currentActor?.role ?? 'pkm_operator']}</Badge>
        {#if $currentActor?.position_nama}<Badge tone="neutral">{$currentActor.position_nama}</Badge>{/if}
      </div>
      <dl class="space-y-2">
        <div class="flex items-start justify-between gap-3 border-b border-black/10 pb-2">
          <dt class="font-mono text-[10px] font-bold uppercase text-muted-foreground">Faskes</dt>
          <dd class="text-right font-semibold">{faskesNama}</dd>
        </div>
        <div class="flex items-start justify-between gap-3 border-b border-black/10 pb-2">
          <dt class="font-mono text-[10px] font-bold uppercase text-muted-foreground">Bidang</dt>
          <dd class="text-right font-semibold">{bidangNama}</dd>
        </div>
        <div class="flex items-start justify-between gap-3">
          <dt class="font-mono text-[10px] font-bold uppercase text-muted-foreground">Login terakhir</dt>
          <dd class="text-right font-semibold">{formatDateTime($currentActor?.last_login_at)}</dd>
        </div>
      </dl>
    </div>

    <div class="flex items-start gap-2 border-t-2 border-black bg-muted p-3 text-xs">
      <UserRound size={14} class="mt-0.5 shrink-0" />
      <p>Username, peran, jabatan, dan faskes hanya dapat diubah oleh admin.</p>
    </div>
  </Card>

  <div class="space-y-4">
    <Card>
      <form onsubmit={saveProfile}>
        <div class="border-b-2 border-black bg-muted px-4 py-3">
          <div class="nb-kicker">Identitas</div>
          <h2 class="mt-1 text-xl">Data akun</h2>
        </div>
        <div class="grid gap-3 p-4 sm:grid-cols-2">
          <div>
            <label class="nb-label" for="nama">Nama lengkap</label>
            <Input
              id="nama"
              bind:value={nama}
              disabled={isActing || profileBusy}
              required
              oninput={() => (profileError = '')}
            />
          </div>
          <div>
            <label class="nb-label" for="username">Username</label>
            <Input id="username" value={$currentActor?.username ?? ''} disabled />
          </div>
        </div>
        <div class="flex flex-wrap items-center justify-between gap-3 border-t-2 border-black px-4 py-3">
          <div class="font-mono text-[10px] font-bold">
            {#if profileError}<span class="text-destructive">{profileError}</span>{:else if profileMessage}<span class="text-success-foreground">{profileMessage}</span>{/if}
          </div>
          <Button type="submit" disabled={isActing || profileBusy}>
            <Save size={14} /> {profileBusy ? 'Menyimpan…' : 'Simpan profil'}
          </Button>
        </div>
      </form>
    </Card>

    <Card>
      <form onsubmit={changePassword}>
        <div class="flex items-center gap-2 border-b-2 border-black bg-muted px-4 py-3">
          <KeyRound size={18} />
          <div>
            <div class="nb-kicker">Keamanan</div>
            <h2 class="mt-1 text-xl">Ubah password</h2>
          </div>
        </div>
        <div class="grid gap-3 p-4 lg:grid-cols-3">
          <div>
            <label class="nb-label" for="current-password">Password saat ini</label>
            <PasswordInput
              id="current-password"
              bind:value={currentPassword}
              autocomplete="current-password"
              disabled={isActing || passwordBusy}
              required
              oninput={() => (passwordError = '')}
            />
          </div>
          <div>
            <label class="nb-label" for="new-password">Password baru</label>
            <PasswordInput
              id="new-password"
              bind:value={newPassword}
              autocomplete="new-password"
              disabled={isActing || passwordBusy}
              required
              oninput={() => (passwordError = '')}
            />
          </div>
          <div>
            <label class="nb-label" for="confirm-password">Ulangi password baru</label>
            <PasswordInput
              id="confirm-password"
              bind:value={confirmPassword}
              autocomplete="new-password"
              disabled={isActing || passwordBusy}
              required
              oninput={() => (passwordError = '')}
            />
          </div>
        </div>
        <div class="border-t-2 border-black bg-muted px-4 py-3 text-xs">
          <div class="flex items-start gap-2">
            <ShieldCheck size={14} class="mt-0.5 shrink-0" />
            <p>
              Minimal 8 karakter. Setelah password berubah, semua sesi lain dicabut dan Anda
              diminta masuk kembali.
            </p>
          </div>
        </div>
        <div class="flex flex-wrap items-center justify-between gap-3 border-t-2 border-black px-4 py-3">
          <div class="font-mono text-[10px] font-bold">
            {#if passwordError}<span class="text-destructive">{passwordError}</span>{:else if passwordMessage}<span class="text-success-foreground">{passwordMessage}</span>{/if}
          </div>
          <Button type="submit" variant="secondary" disabled={isActing || passwordBusy}>
            <KeyRound size={14} /> {passwordBusy ? 'Memproses…' : 'Ubah password'}
          </Button>
        </div>
      </form>
    </Card>
  </div>
</div>
