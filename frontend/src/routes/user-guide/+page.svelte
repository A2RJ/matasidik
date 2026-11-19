<script lang="ts">
  import { onMount } from 'svelte';
  import {
    AlertTriangle,
    ArrowLeft,
    ArrowRight,
    CheckCircle2,
    ChevronRight,
    LogIn,
    ShieldCheck
  } from '@lucide/svelte';
  import Screenshot from '$lib/components/user-guide/screenshot.svelte';
  import { currentActor } from '$lib/stores/session';
  import { ROLE_LABELS } from '$lib/constants/enums';
  import {
    FLOW_IMAGE,
    GUIDE_TABS,
    ROLE_GUIDES,
    type GuideTabID
  } from '$lib/content/user-guide';

  type RoleID = Exclude<GuideTabID, 'ringkasan'>;

  const roleTabs = GUIDE_TABS.filter((tab) => tab.id !== 'ringkasan') as {
    id: RoleID;
    label: string;
    short: string;
  }[];

  let active = $state<RoleID>('puskesmas');
  let visibleRoleTabs = $derived.by(() => {
    if (!$currentActor) return roleTabs;
    const matching = roleTabs.filter((tab) => ROLE_GUIDES[tab.id].roles.includes($currentActor.role));
    return matching.length > 0 ? matching : roleTabs;
  });
  let activeGuide = $derived(ROLE_GUIDES[active]);
  let visibleAccess = $derived(activeGuide.access.filter((row) => row.access !== 'Tidak'));
  let roleIndex = $derived(visibleRoleTabs.findIndex((tab) => tab.id === active));
  let previousRole = $derived(roleIndex > 0 ? visibleRoleTabs[roleIndex - 1] : null);
  let nextRole = $derived(roleIndex < visibleRoleTabs.length - 1 ? visibleRoleTabs[roleIndex + 1] : null);

  $effect(() => {
    if (!$currentActor) return;
    const match = roleTabs.find((tab) => ROLE_GUIDES[tab.id].roles.includes($currentActor.role));
    if (match) active = match.id;
  });

  function shotsFor(stage: string) {
    return activeGuide.screenshots.filter((shot) => shot.stage === stage);
  }

  function selectRole(id: RoleID, scroll = false) {
    active = id;
    if (typeof window !== 'undefined') {
      window.history.replaceState(null, '', `${window.location.pathname}#${id}`);
      if (scroll) document.getElementById('panduan-role')?.scrollIntoView({ block: 'start' });
    }
  }

  onMount(() => {
    const id = window.location.hash.replace('#', '') as RoleID;
    if (roleTabs.some((tab) => tab.id === id)) active = id;
  });
</script>

<div id="guide-top" class="mx-auto w-full max-w-[1400px] px-4 py-6 md:px-6 md:py-8">
  <header class="border-2 border-black bg-primary p-5 shadow-nb md:p-7">
    <div class="flex flex-wrap items-start justify-between gap-4">
      <div class="max-w-4xl">
        <div class="font-mono text-[10px] font-bold uppercase tracking-[0.18em]">
          MATA SIDIK · Panduan penggunaan
        </div>
        <h1 class="mt-2 text-4xl leading-none md:text-6xl">Panduan alur kerja per peran</h1>
        <p class="mt-4 max-w-3xl text-sm font-medium leading-relaxed md:text-base">
          Dokumentasi alur bisnis MATA SIDIK dari input sampai tindak lanjut. Pilih peran di bawah
          untuk melihat posisi peran pada alur, menu dan hak akses yang dipakai, serta tutorial langkah
          kerja dengan data, status, dan screenshot komponen.
        </p>
      </div>
      <div class="flex flex-wrap gap-2">
        {#if $currentActor}
          <a
            href="/"
            class="inline-flex h-10 items-center gap-2 border-2 border-black bg-card px-4 font-mono text-xs font-bold uppercase shadow-nb"
          >
            <ArrowLeft size={14} /> Dashboard
          </a>
        {:else}
          <a
            href="/login"
            class="inline-flex h-10 items-center gap-2 border-2 border-black bg-card px-4 font-mono text-xs font-bold uppercase shadow-nb"
          >
            <LogIn size={14} /> Masuk
          </a>
        {/if}
      </div>
    </div>
  </header>

  <section id="pilih-peran" class="mt-5 border-2 border-black bg-card shadow-nb">
    <div class="border-b-2 border-black bg-muted px-4 py-3">
      <div class="nb-kicker">Pilihan peran</div>
      <h2 class="mt-1 text-2xl">Pilih peran yang ingin dipelajari</h2>
      <p class="mt-1 text-xs text-muted-foreground">
        Panduan di bawah akan mengikuti peran yang dipilih. Setiap peran memiliki alur, menu, dan
        langkah kerja yang berbeda.
      </p>
    </div>
    <div class="grid gap-3 p-4 md:grid-cols-2 xl:grid-cols-4">
      {#each visibleRoleTabs as tab (tab.id)}
        {@const guide = ROLE_GUIDES[tab.id]}
        <button
          type="button"
          class="flex h-full flex-col border-2 border-black bg-card p-4 text-left shadow-nbsm transition-transform hover:-translate-y-0.5 hover:bg-primary/20"
          class:border-green-600={active === tab.id}
          class:border-4={active === tab.id}
          class:bg-green-50={active === tab.id}
          onclick={() => selectRole(tab.id, true)}
          aria-pressed={active === tab.id}
        >
          <div class="flex items-center justify-between gap-2">
            <span class="font-mono text-[9px] font-bold uppercase text-muted-foreground">{guide.level}</span>
            {#if active === tab.id}
              <span class="border border-green-700 bg-green-100 px-1.5 py-0.5 font-mono text-[9px] font-bold uppercase text-green-900">
                Dibuka
              </span>
            {/if}
          </div>
          <h3 class="mt-2 text-xl leading-tight">{guide.title}</h3>
          <p class="mt-2 flex-1 text-xs leading-relaxed text-muted-foreground">{guide.summary}</p>
          <span class="mt-3 inline-flex items-center gap-1 font-mono text-[10px] font-bold uppercase">
            Buka panduan <ChevronRight size={13} />
          </span>
        </button>
      {/each}
    </div>
    <div class="border-t-2 border-black bg-muted px-4 py-2 text-xs text-muted-foreground">
      PJ Program mengikuti modul Puskesmas untuk input BNBA. Pengajuan IKU hanya dapat dilakukan Operator Puskesmas.
      Menu dan tindakan yang tampil pada akun Anda mengikuti peran masing-masing.
    </div>
  </section>

  <div id="panduan-role" class="mt-6 space-y-6">
    <section class="border-2 border-black bg-card shadow-nb">
      <div class="p-4 md:p-5">
        <div class="flex flex-wrap items-start justify-between gap-3">
          <div>
            <div class="nb-kicker">{activeGuide.level}</div>
            <h2 class="mt-1 text-3xl md:text-4xl">{activeGuide.title}</h2>
            <p class="mt-2 max-w-3xl text-sm leading-relaxed text-muted-foreground">{activeGuide.summary}</p>
          </div>
          <div class="flex flex-wrap gap-2">
            {#each activeGuide.roles as role (role)}
              <span class="border-2 border-black bg-secondary px-2 py-1 font-mono text-[10px] font-bold uppercase">
                {ROLE_LABELS[role]}
              </span>
            {/each}
          </div>
        </div>
      </div>

      <div class="border-t-2 border-black p-4 md:p-5">
        <div class="flex flex-wrap items-start justify-between gap-2">
          <div>
            <div class="nb-kicker">Posisi pada alur</div>
            <h3 class="mt-1 text-2xl">Posisi {activeGuide.title} dalam alur aplikasi</h3>
            <p class="mt-1 max-w-3xl text-xs leading-relaxed text-muted-foreground">
              Lane peran ini ditandai hijau. Garis putus-putus menunjukkan pengembalian data atau
              dukungan lintas peran.
            </p>
          </div>
          <span class="inline-flex items-center gap-1 border-2 border-green-700 bg-green-100 px-2 py-1 font-mono text-[10px] font-bold uppercase text-green-900">
            <CheckCircle2 size={12} /> Titik mulai: {activeGuide.startLabel}
          </span>
        </div>
        <div class="mt-3">
          <Screenshot
            src={FLOW_IMAGE.src}
            alt={FLOW_IMAGE.alt}
            caption={`Alur bisnis MATA SIDIK — ${activeGuide.title}`}
            detail="Lane peran yang sedang dibuka ditandai hijau. Garis putus-putus menunjukkan pengembalian data atau dukungan lintas peran."
            overlay={activeGuide.flowOverlay}
          />
        </div>
      </div>
    </section>

    <section class="border-2 border-black bg-card shadow-nb">
      <div class="border-b-2 border-black bg-muted px-4 py-3">
        <div class="nb-kicker">Menu & hak akses</div>
        <h2 class="mt-1 text-2xl">Menu yang dipakai {activeGuide.title}</h2>
        <p class="mt-1 text-xs text-muted-foreground">
          Hanya menu yang relevan untuk peran ini. Hak akses baca dan aksi dibedakan sesuai status
          pengajuan.
        </p>
      </div>
      <div class="overflow-x-auto nb-scroll">
        <table class="nb-table" style="min-width: 720px;">
          <thead><tr><th>Menu / fitur</th><th>Akses</th><th>Aksi & data</th></tr></thead>
          <tbody>
            {#each visibleAccess as row (row.menu)}
              <tr>
                <td class="font-semibold">{row.menu}</td>
                <td>{row.access}</td>
                <td>{row.detail}</td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    </section>

    <section class="border-2 border-black bg-card shadow-nb">
      <div class="border-b-2 border-black bg-muted px-4 py-3">
        <div class="nb-kicker">Tutorial alur kerja</div>
        <h2 class="mt-1 text-2xl">Langkah menggunakan aplikasi</h2>
        <p class="mt-1 text-xs text-muted-foreground">
          Urutan langkah, data/status yang berubah, tujuan berikutnya, dan komponen yang perlu dilihat.
        </p>
      </div>
      <div class="divide-y-2 divide-black/10">
        {#each activeGuide.workflowDetail as step, index (step.stage)}
          <div class="grid gap-3 p-4 md:grid-cols-[44px_1fr]">
            <div class="flex h-9 w-9 items-center justify-center border-2 border-black bg-primary font-mono text-sm font-bold shadow-nbsm">
              {String(index + 1).padStart(2, '0')}
            </div>
            <div class="min-w-0">
              <h3 class="text-lg font-bold">{step.stage}</h3>
              <p class="mt-1 text-sm leading-relaxed text-muted-foreground">{step.actions}</p>
              <div class="mt-2 grid gap-2 sm:grid-cols-2">
                <div class="border-2 border-black bg-muted p-2">
                  <div class="nb-kicker">Data/status</div>
                  <div class="mt-0.5 font-mono text-[10px]">{step.data}</div>
                </div>
                <div class="border-2 border-black bg-muted p-2">
                  <div class="nb-kicker">Lanjut ke</div>
                  <div class="mt-0.5 text-sm font-semibold">{step.next}</div>
                </div>
              </div>
              {#each shotsFor(step.stage) as shot (shot.src)}
                <div class="mt-3">
                  <Screenshot
                    src={shot.src}
                    alt={shot.alt}
                    caption={shot.caption}
                    detail={shot.detail}
                    access={shot.access}
                    inline={shot.inline ?? 'scroll'}
                    size={shot.size ?? 'wide'}
                  />
                  <p class="mt-1 border-2 border-green-700 bg-green-50 px-2 py-1 text-xs leading-relaxed text-green-950">
                    <b>Penanda hijau:</b> {shot.highlight}
                  </p>
                </div>
              {/each}
            </div>
          </div>
        {/each}
      </div>
    </section>

    <section class="border-2 border-black bg-warning p-4 shadow-nbsm">
      <div class="flex items-center gap-2 font-mono text-[10px] font-bold uppercase">
        <AlertTriangle size={16} /> Catatan penting {activeGuide.title}
      </div>
      <ul class="mt-2 space-y-1.5 text-sm leading-relaxed">
        {#each activeGuide.notes as note (note)}
          <li class="flex gap-2"><span class="mt-1 h-1.5 w-1.5 shrink-0 bg-black"></span><span>{note}</span></li>
        {/each}
      </ul>
    </section>

    <div class="flex flex-wrap items-center justify-between gap-2 border-2 border-black bg-card p-3">
      {#if previousRole}
        <button
          type="button"
          class="inline-flex items-center gap-1 border-2 border-black bg-card px-3 py-2 font-mono text-[10px] font-bold uppercase shadow-nbsm"
          onclick={() => selectRole(previousRole.id, true)}
        >
          <ArrowLeft size={13} /> {previousRole.short}
        </button>
      {:else}
        <span></span>
      {/if}
      {#if nextRole}
        <button
          type="button"
          class="inline-flex items-center gap-1 border-2 border-black bg-primary px-3 py-2 font-mono text-[10px] font-bold uppercase shadow-nbsm"
          onclick={() => selectRole(nextRole.id, true)}
        >
          {nextRole.short} <ArrowRight size={13} />
        </button>
      {:else}
        <button
          type="button"
          class="inline-flex items-center gap-1 border-2 border-black bg-primary px-3 py-2 font-mono text-[10px] font-bold uppercase shadow-nbsm"
          onclick={() => selectRole(visibleRoleTabs[0]?.id ?? 'puskesmas', true)}
        >
          Kembali ke {visibleRoleTabs[0]?.short ?? 'Puskesmas'} <ArrowRight size={13} />
        </button>
      {/if}
    </div>
  </div>

  <footer class="mt-6 flex flex-wrap items-center justify-between gap-2 border-t-2 border-black px-1 py-3 font-mono text-[10px] text-muted-foreground">
    <span class="inline-flex items-center gap-1"><ShieldCheck size={12} /> MATA SIDIK · panduan pilot</span>
    <span class="inline-flex items-center gap-2">
      {#if $currentActor}
        <a href="/" class="nb-link">Kembali ke dashboard</a>
      {:else}
        <a href="/login" class="nb-link">Masuk ke aplikasi</a>
      {/if}
    </span>
  </footer>
</div>
