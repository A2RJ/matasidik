<script lang="ts">
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';
  import { login } from '$lib/stores/auth';
  import Button from '$lib/components/ui/button/button.svelte';
  import Input from '$lib/components/ui/input.svelte';
  import {
    ShieldCheck,
    LogIn,
    Eye,
    EyeOff,
    Upload,
    LayoutDashboard,
    Network,
    FileSpreadsheet,
    BookOpen,
    ArrowRight
  } from '@lucide/svelte';

  let username = $state('');
  let password = $state('');
  let busy = $state(false);
  let error = $state('');
  let showPassword = $state(false);
  let passwordChanged = $state(false);

  onMount(() => {
    if (sessionStorage.getItem('matasidik_password_changed') === '1') {
      sessionStorage.removeItem('matasidik_password_changed');
      passwordChanged = true;
    }
  });

  const functions = [
    {
      icon: Upload,
      title: 'Upload terpusat',
      text: 'Satu pintu untuk data BNBA dari seluruh unit.'
    },
    {
      icon: LayoutDashboard,
      title: 'Dashboard SPM',
      text: 'Capaian, target, dan status dalam satu layar.'
    },
    {
      icon: Network,
      title: 'Jejak angka',
      text: 'Setiap nilai dapat dilacak sampai file asli.'
    },
    {
      icon: FileSpreadsheet,
      title: 'Laporan & tindak lanjut',
      text: 'Export PDF/Excel, catat faktor, dan pantau TL.'
    }
  ];

  const steps = [
    { no: '1', title: 'Upload', final: false },
    { no: '2', title: 'Verifikasi', final: false },
    { no: '3', title: 'Validasi', final: false },
    { no: '4', title: 'Final', final: true }
  ];

  async function submit(event: SubmitEvent) {
    event.preventDefault();
    if (busy) return;
    busy = true;
    error = '';
    try {
      await login(username.trim(), password);
      await goto('/');
    } catch (err) {
      error = err instanceof Error ? err.message : 'Login gagal. Periksa username dan password.';
    } finally {
      busy = false;
    }
  }
</script>

<div class="min-h-screen bg-background lg:grid lg:grid-cols-2">
  <section
    class="flex flex-col border-b-2 border-black bg-primary px-6 py-8 sm:px-8 lg:min-h-screen lg:border-b-0 lg:border-r-2 lg:px-10 lg:py-10"
  >
    <div class="mx-auto flex w-full max-w-xl flex-1 flex-col">
      <div class="my-auto">
        <div>
          <div class="font-mono text-[10px] font-bold uppercase tracking-[0.18em]">
            Kabupaten Sumbawa Barat
          </div>
          <h1 class="mt-3 text-5xl font-bold leading-none sm:text-6xl">MATA SIDIK</h1>
          <p class="mt-4 max-w-md text-sm font-medium leading-relaxed">
            Manajemen Data dan Sistem Informasi Dinas Kesehatan. Satu tempat untuk mengunggah,
            memeriksa, dan membaca capaian SPM tanpa berpindah aplikasi.
          </p>
        </div>

        <div class="mt-8">
          <div class="border-b-2 border-black pb-2">
            <div class="nb-kicker">Fungsi aplikasi</div>
            <h2 class="mt-1 text-2xl leading-tight">Empat pekerjaan utama, satu alur</h2>
          </div>

          <div class="mt-3 grid gap-3 sm:grid-cols-2">
            {#each functions as item (item.title)}
              {@const Icon = item.icon}
              <div class="border-2 border-black bg-card p-3 shadow-nbsm">
                <div class="flex items-center gap-2">
                  <span class="flex h-8 w-8 shrink-0 items-center justify-center border-2 border-black bg-secondary">
                    <Icon size={16} />
                  </span>
                  <h3 class="font-mono text-[11px] font-bold uppercase leading-tight">{item.title}</h3>
                </div>
                <p class="mt-2 text-xs leading-relaxed">{item.text}</p>
              </div>
            {/each}
          </div>
        </div>

        <div class="mt-6 border-2 border-black bg-card p-4 shadow-nbsm">
          <div class="nb-kicker">Alur aplikasi</div>
          <div class="mt-2 flex flex-nowrap items-center gap-0.5 overflow-x-auto whitespace-nowrap font-mono text-[8.5px] font-bold uppercase sm:gap-1 sm:text-[10px]">
            {#each steps as step, i (step.no)}
              <span
                class="border-2 border-black px-1 py-1 sm:px-1.5"
                class:bg-primary={step.final}
                class:bg-card={!step.final}
              >
                {step.no}. {step.title}
              </span>
              {#if i < steps.length - 1}<span class="text-muted-foreground">→</span>{/if}
            {/each}
          </div>
        </div>
      </div>
    </div>
  </section>

  <section
    class="flex flex-col justify-center bg-background px-6 py-8 sm:px-8 lg:min-h-screen lg:px-10 lg:py-10"
  >
    <div class="mx-auto w-full max-w-xl">
      {#if passwordChanged}
        <div class="mb-3 border-2 border-black bg-success p-3 font-mono text-[10px] font-bold" role="status">
          Password diperbarui. Silakan masuk dengan password baru.
        </div>
      {/if}
      <form class="border-2 border-black bg-card p-4 shadow-nblg sm:p-5" onsubmit={submit}>
        <div class="mb-3 flex items-center gap-2 border-b-2 border-black pb-2">
          <ShieldCheck size={18} />
          <h2 class="font-mono text-sm font-bold uppercase tracking-[0.08em]">Masuk ke ruang kerja</h2>
        </div>
        <div class="space-y-3">
          <div>
            <label class="nb-label" for="username">Username</label>
            <Input
              id="username"
              bind:value={username}
              placeholder="username"
              autocomplete="username"
              required
              oninput={() => (error = '')}
            />
          </div>
          <div>
            <label class="nb-label" for="password">Password</label>
            <div class="relative">
              <Input
                id="password"
                type={showPassword ? 'text' : 'password'}
                bind:value={password}
                placeholder="••••••••"
                autocomplete="current-password"
                required
                class="pr-14"
                oninput={() => (error = '')}
              />
              <button
                type="button"
                class="absolute right-1.5 top-1/2 flex h-8 w-9 -translate-y-1/2 items-center justify-center border-2 border-black bg-card"
                onclick={() => (showPassword = !showPassword)}
                aria-label={showPassword ? 'Sembunyikan password' : 'Tampilkan password'}
                aria-pressed={showPassword}
                title={showPassword ? 'Sembunyikan password' : 'Tampilkan password'}
              >
                {#if showPassword}<EyeOff size={16} />{:else}<Eye size={16} />{/if}
              </button>
            </div>
          </div>
          {#if error}
            <div role="alert" class="border-2 border-black bg-destructive p-2 font-mono text-[10px] font-bold">
              {error}
            </div>
          {/if}
          <Button type="submit" class="w-full" size="lg" disabled={busy}>
            {#if busy}Memeriksa…{:else}<LogIn size={16} /> Masuk{/if}
          </Button>
        </div>
        <details class="mt-4 border-2 border-black bg-muted p-2 font-mono text-[10px]">
          <summary class="cursor-pointer font-bold uppercase">Akun awal</summary>
          <div class="mt-2 space-y-1">
            <div><b>admin</b> · pengelola aplikasi</div>
            <div><b>siti</b> · pengirim data</div>
            <div><b>rina</b> · pemeriksa data</div>
            <div><b>kabid</b> · penyetuju</div>
            <div><b>pimpinan</b> · pembaca dashboard</div>
            <div class="pt-1 text-muted-foreground">Password awal: password123 — ganti sebelum produksi.</div>
          </div>
        </details>
      </form>

      <a
        href="/user-guide"
        class="mt-3 flex items-center justify-between gap-3 border-2 border-black bg-success px-4 py-3 font-mono text-xs font-bold uppercase shadow-nb transition-transform hover:-translate-y-0.5"
      >
        <span class="flex items-center gap-2"><BookOpen size={16} /> Panduan penggunaan</span>
        <span class="flex items-center gap-2">Tata cara per peran <ArrowRight size={14} /></span>
      </a>
      <p class="mt-2 text-center font-mono text-[10px] text-muted-foreground">
        Panduan dapat dibaca tanpa login.
      </p>

      <div class="mt-4 border-2 border-black bg-secondary p-4 shadow-nbsm">
        <div class="font-mono text-[10px] font-bold uppercase tracking-[0.12em]">Ruang kerja</div>
        <p class="mt-1 text-lg font-bold leading-tight">Kerja data yang rapi, keputusan yang tenang.</p>
      </div>
    </div>
  </section>
</div>
