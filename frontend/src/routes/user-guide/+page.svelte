<script lang="ts">
  import { onMount } from 'svelte';
  import {
    AlertTriangle,
    ArrowLeft,
    ArrowRight,
    BookOpen,
    CheckCircle2,
    ChevronRight,
    CircleHelp,
    LogIn,
    MousePointerClick,
    ShieldCheck,
    Workflow
  } from '@lucide/svelte';
  import Screenshot from '$lib/components/user-guide/screenshot.svelte';
  import { currentActor } from '$lib/stores/session';
  import { ROLE_LABELS } from '$lib/constants/enums';
  import {
    DATA_FLOW_LINKS,
    DEMO_ACCOUNTS,
    FLOW_IMAGE,
    FLOW_LEGEND,
    GUIDE_TABS,
    MENU_ACCESS,
    MENU_CATALOG,
    OPEN_DECISIONS,
    ROLE_GUIDES,
    ROLE_MATRIX,
    STATUS_ROWS,
    type GuideTabID,
    type RoleGuide
  } from '$lib/content/user-guide';

  let active = $state<GuideTabID>('ringkasan');

  const roleTabs = GUIDE_TABS.filter((tab) => tab.id !== 'ringkasan');
  const currentIndex = $derived(roleTabs.findIndex((tab) => tab.id === active));
  const previousTab = $derived(currentIndex > 0 ? roleTabs[currentIndex - 1] : null);
  const nextTab = $derived(
    currentIndex >= 0 && currentIndex < roleTabs.length - 1 ? roleTabs[currentIndex + 1] : null
  );
  const activeGuide: RoleGuide | null = $derived(
    active === 'ringkasan' ? null : ROLE_GUIDES[active]
  );

  function selectTab(id: GuideTabID, scroll = false) {
    active = id;
    if (typeof window !== 'undefined') {
      window.history.replaceState(null, '', `${window.location.pathname}#${id}`);
      if (scroll) {
        document.getElementById('guide-top')?.scrollIntoView({ block: 'start' });
      }
    }
  }

  onMount(() => {
    const applyHash = () => {
      const id = window.location.hash.replace('#', '') as GuideTabID;
      if (GUIDE_TABS.some((tab) => tab.id === id)) active = id;
    };
    applyHash();
    window.addEventListener('hashchange', applyHash);
    return () => window.removeEventListener('hashchange', applyHash);
  });
</script>

<div id="guide-top" class="mx-auto w-full max-w-[1400px] px-4 py-6 md:px-6 md:py-8">
  <section class="border-2 border-black bg-primary p-5 shadow-nb md:p-7">
    <div class="flex flex-wrap items-start justify-between gap-4">
      <div class="max-w-4xl">
        <div class="font-mono text-[10px] font-bold uppercase tracking-[0.18em]">
          MATA SIDIK · Panduan penggunaan
        </div>
        <h1 class="mt-2 text-4xl leading-none md:text-6xl">Tata cara aplikasi dari awal sampai akhir</h1>
        <p class="mt-4 max-w-3xl text-sm font-medium leading-relaxed md:text-base">
          Panduan ini dibagi per peran/level, dilengkapi katalog menu, hak akses, dan keterkaitan
          fitur/menu/flow/data. Setiap modul menampilkan alur lengkap aplikasi yang sama, dengan
          <b>kotak hijau</b> sebagai penanda titik mulai alur kerja role tersebut. Screenshot menyertakan
          tabel/form utuh dan dapat di-zoom atau dibuka layar penuh.
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
        <a
          href="#ringkasan"
          onclick={() => selectTab('ringkasan')}
          class="inline-flex h-10 items-center gap-2 border-2 border-black bg-success px-4 font-mono text-xs font-bold uppercase shadow-nb"
        >
          <BookOpen size={14} /> Mulai dari ringkasan
        </a>
      </div>
    </div>

    <div class="mt-6 grid gap-2 sm:grid-cols-3">
      <div class="border-2 border-black bg-card px-3 py-2">
        <div class="font-mono text-[9px] font-bold uppercase text-muted-foreground">Modul peran</div>
        <div class="mt-1 text-xl font-bold">4 modul · 6 role aktif</div>
      </div>
      <div class="border-2 border-black bg-card px-3 py-2">
        <div class="font-mono text-[9px] font-bold uppercase text-muted-foreground">Alur</div>
        <div class="mt-1 text-xl font-bold">1 WebP · stack-diagram</div>
      </div>
      <div class="border-2 border-black bg-card px-3 py-2">
        <div class="font-mono text-[9px] font-bold uppercase text-muted-foreground">Screenshot</div>
        <div class="mt-1 text-xl font-bold">Zoom + layar penuh</div>
      </div>
    </div>
  </section>

  <nav class="sticky top-0 z-20 -mx-4 mt-4 border-y-2 border-black bg-background/95 px-4 py-2 backdrop-blur md:-mx-6 md:px-6" aria-label="Modul panduan">
    <div class="flex flex-wrap gap-2">
      {#each GUIDE_TABS as tab (tab.id)}
        <button
          type="button"
          class="inline-flex items-center gap-1.5 border-2 border-black px-3 py-1.5 font-mono text-[10px] font-bold uppercase shadow-nbsm"
          class:bg-primary={active === tab.id}
          class:bg-card={active !== tab.id}
          aria-pressed={active === tab.id}
          onclick={() => selectTab(tab.id, true)}
        >
          {#if tab.id === 'ringkasan'}<CircleHelp size={12} />{/if}
          {tab.short}
        </button>
      {/each}
    </div>
  </nav>

  {#if active === 'ringkasan'}
    <div class="mt-5 space-y-5">
      <section class="grid gap-3 md:grid-cols-3">
        <div class="border-2 border-black bg-card p-4 shadow-nbsm">
          <div class="flex items-center gap-2 font-mono text-[10px] font-bold uppercase">
            <MousePointerClick size={16} /> 1. Pilih modul peran
          </div>
          <p class="mt-2 text-sm leading-relaxed text-muted-foreground">
            Klik kartu role di bawah atau tab di atas. Modul role menampilkan langkah kerja, keterkaitan
            dengan role lain, dan screenshot yang relevan.
          </p>
        </div>
        <div class="border-2 border-black bg-card p-4 shadow-nbsm">
          <div class="flex items-center gap-2 font-mono text-[10px] font-bold uppercase">
            <Workflow size={16} /> 2. Baca alur lengkap
          </div>
          <p class="mt-2 text-sm leading-relaxed text-muted-foreground">
            Setiap modul memakai satu WebP alur keseluruhan. Kotak hijau menunjukkan titik mulai role;
            garis putus-putus menunjukkan pengembalian data dan dukungan lintas peran.
          </p>
        </div>
        <div class="border-2 border-black bg-card p-4 shadow-nbsm">
          <div class="flex items-center gap-2 font-mono text-[10px] font-bold uppercase">
            <MousePointerClick size={16} /> 3. Zoom screenshot
          </div>
          <p class="mt-2 text-sm leading-relaxed text-muted-foreground">
            Klik screenshot untuk memperbesar. Gunakan tombol zoom, Fit, atau Layar penuh. Tabel tetap
            ditampilkan penuh dan aksi role ditandai kotak hijau.
          </p>
        </div>
      </section>

      <section class="border-2 border-black bg-card shadow-nb">
        <div class="border-b-2 border-black bg-muted px-4 py-3">
          <div class="nb-kicker">Peta modul</div>
          <h2 class="mt-1 text-2xl">Pilih role/level yang ingin dipelajari</h2>
        </div>
        <div class="grid gap-3 p-4 md:grid-cols-2 xl:grid-cols-3">
          {#each roleTabs as tab (tab.id)}
            {@const guide = ROLE_GUIDES[tab.id as Exclude<GuideTabID, 'ringkasan'>]}
            <button
              type="button"
              class="flex h-full flex-col border-2 border-black bg-card p-4 text-left shadow-nbsm transition-transform hover:-translate-y-0.5 hover:bg-primary/20"
              onclick={() => selectTab(tab.id, true)}
            >
              <div class="flex items-center justify-between gap-2">
                <span class="font-mono text-[9px] font-bold uppercase text-muted-foreground">{guide.level}</span>
                <span class="border border-black px-1.5 py-0.5 font-mono text-[9px] font-bold uppercase">
                  {guide.roles.map((role) => ROLE_LABELS[role]).join(' / ')}
                </span>
              </div>
              <h3 class="mt-2 text-xl leading-tight">{guide.title}</h3>
              <p class="mt-2 flex-1 text-sm leading-relaxed text-muted-foreground">{guide.summary}</p>
              <span class="mt-3 inline-flex items-center gap-1 font-mono text-[10px] font-bold uppercase">
                Buka modul <ChevronRight size={13} />
              </span>
            </button>
          {/each}
          <div class="border-2 border-dashed border-black/50 bg-muted p-4">
            <div class="font-mono text-[9px] font-bold uppercase text-muted-foreground">Catatan role</div>
            <h3 class="mt-2 text-lg leading-tight">PJ Program</h3>
            <p class="mt-2 text-sm leading-relaxed text-muted-foreground">
              PJ Program mengikuti modul Puskesmas; scope faskes mengikuti akun. Role
              <b>iku_owner</b> ada di rancangan, tetapi belum memiliki akun/layar khusus pada pilot.
            </p>
          </div>
        </div>
      </section>

      <section class="border-2 border-black bg-card shadow-nb">
        <div class="flex flex-wrap items-center justify-between gap-2 border-b-2 border-black bg-secondary px-4 py-3">
          <div>
            <div class="nb-kicker">Alur lengkap aplikasi</div>
            <h2 class="mt-1 text-2xl">Dari login sampai keputusan dan tindak lanjut</h2>
          </div>
          <span class="border-2 border-black bg-card px-2 py-1 font-mono text-[9px] font-bold uppercase">
            stack-diagram
          </span>
        </div>
        <div class="p-4">
          <Screenshot
            src={FLOW_IMAGE.src}
            alt={FLOW_IMAGE.alt}
            caption="Alur lengkap MATA SIDIK (tanpa penanda role)"
            detail="Alur utama berada di baris role. Pada modul role, gambar yang sama diberi kotak hijau di titik mulai role tersebut."
          />
          <div class="mt-3 flex flex-wrap gap-3 border-2 border-black bg-muted p-3 font-mono text-[10px] font-bold uppercase">
            {#each FLOW_LEGEND as item (item.text)}
              <span class="inline-flex items-center gap-2">
                <span
                  class="inline-block h-3 w-5 border-2"
                  class:border-green-600={item.swatch === 'green'}
                  class:bg-green-200={item.swatch === 'green'}
                  class:border-black={item.swatch !== 'green'}
                  class:bg-primary={item.swatch === 'blue'}
                  class:border-dashed={item.swatch === 'dashed'}
                ></span>
                {item.text}
              </span>
            {/each}
          </div>
        </div>
      </section>

      <section class="border-2 border-black bg-card shadow-nb">
        <div class="border-b-2 border-black bg-muted px-4 py-3">
          <div class="nb-kicker">Katalog menu & fitur</div>
          <h2 class="mt-1 text-2xl">Menu yang membentuk alur aplikasi</h2>
          <p class="mt-1 text-xs text-muted-foreground">
            Setiap menu membaca/menulis data yang berbeda. Tabel ini menjelaskan fungsi, data, dan bagian lain
            yang terkait agar tidak ada menu yang dipahami terpisah dari alurnya.
          </p>
        </div>
        <div class="overflow-x-auto nb-scroll">
          <table class="nb-table" style="min-width: 1180px;">
            <thead>
              <tr>
                <th>Menu / fitur</th>
                <th>Fungsi utama</th>
                <th>Data yang dibaca / ditulis</th>
                <th>Terkait dengan</th>
              </tr>
            </thead>
            <tbody>
              {#each MENU_CATALOG as menu (menu.menu)}
                <tr>
                  <td>
                    <div class="font-bold">{menu.menu}</div>
                    <div class="font-mono text-[10px] text-muted-foreground">{menu.route}</div>
                  </td>
                  <td>{menu.function}</td>
                  <td>{menu.data}</td>
                  <td>{menu.related}</td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      </section>

      <section class="border-2 border-black bg-card shadow-nb">
        <div class="border-b-2 border-black bg-muted px-4 py-3">
          <div class="nb-kicker">Hak akses per menu</div>
          <h2 class="mt-1 text-2xl">Siapa membuka menu apa</h2>
          <p class="mt-1 text-xs text-muted-foreground">
            Hak akses menu berbeda dari hak akses aksi. Satu role dapat melihat menu tetapi hanya boleh
            melakukan aksi tertentu di dalamnya.
          </p>
        </div>
        <div class="overflow-x-auto nb-scroll">
          <table class="nb-table" style="min-width: 980px;">
            <thead>
              <tr>
                <th>Menu / fitur</th>
                <th>Puskesmas/PJ</th>
                <th>Verifikator</th>
                <th>Kabid/Sekretaris</th>
                <th>Pimpinan</th>
              </tr>
            </thead>
            <tbody>
              {#each MENU_ACCESS as row (row.menu)}
                <tr>
                  <td class="font-semibold">{row.menu}</td>
                  <td>{row.pkm}</td>
                  <td>{row.verifier}</td>
                  <td>{row.kabid}</td>
                  <td>{row.pimpinan}</td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
        <div class="border-t-2 border-black bg-card px-4 py-2 font-mono text-[10px] text-muted-foreground">
          * Kebijakan produksi; pada pilot pembatasan dashboard per bidang belum sepenuhnya ditegakkan.
          Menu Panduan bersifat publik; menu lain memerlukan sesi.
        </div>
      </section>

      <section class="border-2 border-black bg-card shadow-nb">
        <div class="border-b-2 border-black bg-muted px-4 py-3">
          <div class="nb-kicker">Keterkaitan fitur, menu, flow, dan data</div>
          <h2 class="mt-1 text-2xl">Satu data, banyak tampilan</h2>
          <p class="mt-1 text-xs text-muted-foreground">
            Perubahan di satu fitur memengaruhi fitur lain. Tabel ini menunjukkan perpindahan data dan
            bagian yang harus diperiksa saat ada perubahan atau masalah.
          </p>
        </div>
        <div class="overflow-x-auto nb-scroll">
          <table class="nb-table" style="min-width: 980px;">
            <thead>
              <tr>
                <th>Dari</th>
                <th>Ke</th>
                <th>Perantara / ID</th>
                <th>Keterkaitan</th>
              </tr>
            </thead>
            <tbody>
              {#each DATA_FLOW_LINKS as row (row.from + row.to)}
                <tr>
                  <td class="font-semibold">{row.from}</td>
                  <td class="font-semibold">{row.to}</td>
                  <td class="font-mono text-[10px]">{row.via}</td>
                  <td>{row.note}</td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
        <div class="border-t-2 border-black bg-secondary px-4 py-3 font-mono text-[10px] leading-relaxed">
          Rantai trace: DASHBOARD_ROW → IKU_VALUE → IKU_VALUE_DETAIL → INDICATOR_VALUE →
          SUBMISSION_ROW → SUBMISSION_REVISION → FILE_OBJECT. Jika satu mata rantai tidak ada, status
          menjadi NO_DATA/PARTIAL, bukan nol.
        </div>
      </section>

      <section class="border-2 border-black bg-card shadow-nb">
        <div class="border-b-2 border-black bg-muted px-4 py-3">
          <div class="nb-kicker">Matriks peran</div>
          <h2 class="mt-1 text-2xl">Siapa boleh melakukan apa</h2>
          <p class="mt-1 text-xs text-muted-foreground">
            Ringkasan RBAC dari aturan proyek untuk role yang dibahas di panduan ini.
          </p>
        </div>
        <div class="overflow-x-auto nb-scroll">
          <table class="nb-table">
            <thead>
              <tr>
                <th>Aksi</th>
                <th>Puskesmas/PJ</th>
                <th>Verifikator</th>
                <th>Kabid/Sekretaris</th>
                <th>Pimpinan</th>
              </tr>
            </thead>
            <tbody>
              {#each ROLE_MATRIX as row (row.action)}
                <tr>
                  <td class="font-semibold">{row.action}</td>
                  <td>{row.pkm}</td>
                  <td>{row.verifier}</td>
                  <td>{row.kabid}</td>
                  <td>{row.pimpinan}</td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      </section>

      <section class="border-2 border-black bg-warning p-4 shadow-nbsm">
        <div class="flex items-center gap-2 font-mono text-[10px] font-bold uppercase">
          <AlertTriangle size={16} /> Catatan implementasi pilot
        </div>
        <p class="mt-2 text-sm leading-relaxed">
          Pembatasan tampilan dashboard per bidang untuk Verifikator/Kabid/Sekretaris belum diterapkan
          sepenuhnya di backend pilot; tampilan dapat mencakup seluruh faskes. Matriks di atas adalah
          kebijakan RBAC yang mengikat untuk produksi (<code>docs/rules.md</code> §10). Jangan memakai
          tampilan pilot sebagai pengganti kontrol akses produksi. Layar SPECIAL_CHANGE/dual approval dan
          notifikasi SLA juga belum tersedia.
        </p>
      </section>

      <section class="grid gap-4 lg:grid-cols-2">
        <div class="border-2 border-black bg-card shadow-nb">
          <div class="border-b-2 border-black bg-muted px-4 py-3">
            <div class="nb-kicker">Jalur data</div>
            <h2 class="mt-1 text-xl">Dari file upload sampai dashboard</h2>
          </div>
          <div class="space-y-2 p-4 text-sm leading-relaxed">
            <p><b>1.</b> Puskesmas upload BNBA → API menyimpan file + job → worker mem-parse, memvalidasi, dan mengklasifikasi.</p>
            <p><b>2.</b> Preview dan koreksi IKU leaf → kirim → verifikator memeriksa → Kabid/Sekretaris menetapkan final.</p>
            <p><b>3.</b> Final memicu kalkulasi. IKU induk memakai <code>RATIO_OF_SUMS</code>, bukan rata-rata persentase.</p>
            <p><b>4.</b> Dashboard membaca nilai materialized. Setiap baris SPM membawa ID kontrak untuk trace sampai <code>FILE_OBJECT</code>.</p>
            <p class="border-2 border-black bg-secondary p-2 font-mono text-[10px] leading-relaxed">
              DASHBOARD_ROW → IKU_VALUE → IKU_VALUE_DETAIL → INDICATOR_VALUE → SUBMISSION_ROW → SUBMISSION_REVISION → FILE_OBJECT
            </p>
          </div>
        </div>

        <div class="border-2 border-black bg-card shadow-nb">
          <div class="border-b-2 border-black bg-muted px-4 py-3">
            <div class="nb-kicker">Akun demo pilot</div>
            <h2 class="mt-1 text-xl">Untuk mencoba alur</h2>
            <p class="mt-1 text-xs text-muted-foreground">Password awal semua akun: <code>password123</code> — ganti sebelum produksi.</p>
          </div>
          <div class="overflow-x-auto nb-scroll">
            <table class="nb-table">
              <thead><tr><th>Username</th><th>Role</th><th>Dipakai untuk</th></tr></thead>
              <tbody>
                {#each DEMO_ACCOUNTS as account (account.username)}
                  <tr><td class="font-bold">{account.username}</td><td>{account.role}</td><td>{account.use}</td></tr>
                {/each}
              </tbody>
            </table>
          </div>
        </div>
      </section>
    </div>
  {:else if activeGuide}
    <div class="mt-5 space-y-5">
      <section class="border-2 border-black bg-card p-4 shadow-nb md:p-5">
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
            <span class="inline-flex items-center gap-1 border-2 border-green-700 bg-green-100 px-2 py-1 font-mono text-[10px] font-bold uppercase text-green-900">
              <CheckCircle2 size={12} /> {activeGuide.startLabel}
            </span>
          </div>
        </div>
        <div class="mt-4 border-2 border-black bg-muted p-3 font-mono text-[10px] leading-relaxed">
          <b>Akun demo:</b> {activeGuide.accounts.join(' · ')}
        </div>
      </section>

      <section class="border-2 border-black bg-card shadow-nb">
        <div class="flex flex-wrap items-center justify-between gap-2 border-b-2 border-black bg-secondary px-4 py-3">
          <div>
            <div class="nb-kicker">Alur lengkap aplikasi · {activeGuide.title}</div>
            <h2 class="mt-1 text-2xl">Kotak hijau menandai titik mulai role ini</h2>
          </div>
          <span class="border-2 border-green-700 bg-green-100 px-2 py-1 font-mono text-[9px] font-bold uppercase text-green-900">
            mulai: {activeGuide.startLabel}
          </span>
        </div>
        <div class="p-4">
          <Screenshot
            src={FLOW_IMAGE.src}
            alt={FLOW_IMAGE.alt}
            caption={`Alur lengkap MATA SIDIK — penanda ${activeGuide.title}`}
            detail="Gambar yang sama dipakai di semua modul. Kotak hijau menunjukkan awal alur kerja role; garis putus-putus menunjukkan return atau dukungan lintas peran."
            overlay={activeGuide.flowOverlay}
          />
          <div class="mt-3 grid gap-2 border-2 border-black bg-muted p-3 text-xs leading-relaxed sm:grid-cols-2">
            <p><b>Keterkaitan:</b> {activeGuide.relations[0]?.actor ?? 'Role lain'} dan role lain dijelaskan pada tabel di bawah.</p>
            <p><b>Status akhir role ini:</b> lihat langkah dan catatan. Final selalu mengunci data dan memicu kalkulasi.</p>
          </div>
        </div>
      </section>

      <section class="border-2 border-black bg-card shadow-nb">
        <div class="border-b-2 border-black bg-muted px-4 py-3">
          <div class="nb-kicker">Langkah kerja</div>
          <h2 class="mt-1 text-2xl">Dari login sampai selesai</h2>
        </div>
        <ol class="divide-y-2 divide-black/10">
          {#each activeGuide.steps as step, index (step.title)}
            <li class="grid gap-2 p-4 md:grid-cols-[56px_1fr]">
              <div class="flex h-10 w-10 items-center justify-center border-2 border-black bg-primary font-mono text-sm font-bold shadow-nbsm">
                {String(index + 1).padStart(2, '0')}
              </div>
              <div>
                <div class="flex flex-wrap items-center gap-2">
                  <h3 class="text-lg font-bold">{step.title}</h3>
                  {#if step.where}<code class="border border-black/30 bg-muted px-1.5 py-0.5 font-mono text-[10px]">{step.where}</code>{/if}
                </div>
                <p class="mt-1 text-sm leading-relaxed text-muted-foreground">{step.body}</p>
              </div>
            </li>
          {/each}
        </ol>
      </section>

      <section class="border-2 border-black bg-card shadow-nb">
        <div class="border-b-2 border-black bg-muted px-4 py-3">
          <div class="nb-kicker">Hak akses & menu role ini</div>
          <h2 class="mt-1 text-2xl">Menu, akses, dan aksi/data</h2>
          <p class="mt-1 text-xs text-muted-foreground">
            Hak akses menu dan hak akses aksi dipisahkan. Role dapat membuka menu untuk membaca, tetapi
            aksi tulis/approve tetap mengikuti status dan kewenangan.
          </p>
        </div>
        <div class="overflow-x-auto nb-scroll">
          <table class="nb-table" style="min-width: 780px;">
            <thead><tr><th>Menu / fitur</th><th>Akses</th><th>Aksi & data</th></tr></thead>
            <tbody>
              {#each activeGuide.access as row (row.menu)}
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
          <div class="nb-kicker">Keterkaitan fitur, flow, dan data</div>
          <h2 class="mt-1 text-2xl">Bagian lain yang terhubung</h2>
          <p class="mt-1 text-xs text-muted-foreground">
            Menu ini tidak berdiri sendiri. Perubahan data, status, atau master di satu tempat akan terlihat
            di tempat lain sesuai rantai di bawah.
          </p>
        </div>
        <div class="overflow-x-auto nb-scroll">
          <table class="nb-table">
            <thead><tr><th>Fitur / bagian</th><th>Keterkaitan</th></tr></thead>
            <tbody>
              {#each activeGuide.featureLinks as link (link.feature)}
                <tr><td class="font-semibold">{link.feature}</td><td>{link.relation}</td></tr>
              {/each}
            </tbody>
          </table>
        </div>
      </section>

      <section class="border-2 border-black bg-card shadow-nb">
        <div class="border-b-2 border-black bg-muted px-4 py-3">
          <div class="nb-kicker">Keterkaitan antar role</div>
          <h2 class="mt-1 text-2xl">Siapa menerima, mengembalikan, dan menindaklanjuti</h2>
        </div>
        <div class="overflow-x-auto nb-scroll">
          <table class="nb-table">
            <thead><tr><th>Role / pihak</th><th>Hubungan dengan {activeGuide.title}</th></tr></thead>
            <tbody>
              {#each activeGuide.relations as relation (relation.actor)}
                <tr><td class="font-semibold">{relation.actor}</td><td>{relation.relation}</td></tr>
              {/each}
            </tbody>
          </table>
        </div>
      </section>

      <section class="space-y-4">
        <div class="border-2 border-black bg-card px-4 py-3 shadow-nbsm">
          <div class="nb-kicker">Screenshot elemen terkait</div>
          <h2 class="mt-1 text-2xl">Tabel/form tetap utuh, aksi role ditandai hijau</h2>
          <p class="mt-1 text-xs text-muted-foreground">
            Klik screenshot untuk zoom dan layar penuh. Kotak hijau sudah ditempel pada gambar saat capture,
            sehingga konteks tabel tetap terlihat.
          </p>
        </div>
        {#each activeGuide.screenshots as shot (shot.src)}
          <div class="space-y-2">
            <Screenshot
              src={shot.src}
              alt={shot.alt}
              caption={shot.caption}
              detail={shot.detail}
              inline={shot.inline ?? 'scroll'}
              access={shot.access}
            />
            <div class="border-2 border-green-700 bg-green-50 px-3 py-2 text-xs leading-relaxed text-green-950">
              <b>Penanda hijau:</b> {shot.highlight}
            </div>
          </div>
        {/each}
      </section>

      <section class="border-2 border-black bg-warning p-4 shadow-nbsm">
        <div class="flex items-center gap-2 font-mono text-[10px] font-bold uppercase">
          <AlertTriangle size={16} /> Batasan & keputusan terbuka
        </div>
        <ul class="mt-2 space-y-1.5 text-sm leading-relaxed">
          {#each activeGuide.notes as note (note)}
            <li class="flex gap-2"><span class="mt-1 h-1.5 w-1.5 shrink-0 bg-black"></span><span>{note}</span></li>
          {/each}
        </ul>
      </section>

      <div class="flex flex-wrap items-center justify-between gap-2 border-2 border-black bg-card p-3">
        {#if previousTab}
          <button
            type="button"
            class="inline-flex items-center gap-1 border-2 border-black bg-card px-3 py-2 font-mono text-[10px] font-bold uppercase shadow-nbsm"
            onclick={() => selectTab(previousTab.id, true)}
          >
            <ArrowLeft size={13} /> {previousTab.short}
          </button>
        {:else}
          <span></span>
        {/if}
        {#if nextTab}
          <button
            type="button"
            class="inline-flex items-center gap-1 border-2 border-black bg-primary px-3 py-2 font-mono text-[10px] font-bold uppercase shadow-nbsm"
            onclick={() => selectTab(nextTab.id, true)}
          >
            {nextTab.short} <ArrowRight size={13} />
          </button>
        {:else}
          <button
            type="button"
            class="inline-flex items-center gap-1 border-2 border-black bg-primary px-3 py-2 font-mono text-[10px] font-bold uppercase shadow-nbsm"
            onclick={() => selectTab('ringkasan', true)}
          >
            Kembali ke ringkasan <ArrowRight size={13} />
          </button>
        {/if}
      </div>
    </div>
  {/if}

  <section class="mt-6 grid gap-4 lg:grid-cols-2">
    <div class="border-2 border-black bg-card shadow-nb">
      <div class="border-b-2 border-black bg-muted px-4 py-3">
        <div class="nb-kicker">Status workflow</div>
        <h2 class="mt-1 text-xl">Tujuh status formal</h2>
      </div>
      <div class="overflow-x-auto nb-scroll">
        <table class="nb-table">
          <thead><tr><th>Status</th><th>Arti</th><th>Aktor</th><th>Langkah berikut</th></tr></thead>
          <tbody>
            {#each STATUS_ROWS as row (row.status)}
              <tr><td class="font-mono font-bold">{row.status}</td><td>{row.meaning}</td><td>{row.actor}</td><td>{row.next}</td></tr>
            {/each}
          </tbody>
        </table>
      </div>
    </div>

    <div class="border-2 border-black bg-card shadow-nb">
      <div class="border-b-2 border-black bg-muted px-4 py-3">
        <div class="nb-kicker">Keputusan terbuka</div>
        <h2 class="mt-1 text-xl">Jangan ditebak</h2>
        <p class="mt-1 text-xs text-muted-foreground">Daftar ini mengikuti docs/rules.md §14. Tanyakan Dinas sebelum mengubah perilaku aplikasi.</p>
      </div>
      <div class="max-h-[420px] overflow-auto nb-scroll">
        <table class="nb-table">
          <thead><tr><th>Topik</th><th>Status</th><th>Tindakan</th></tr></thead>
          <tbody>
            {#each OPEN_DECISIONS as item (item.topic)}
              <tr><td class="font-semibold">{item.topic}</td><td>{item.status}</td><td>{item.action}</td></tr>
            {/each}
          </tbody>
        </table>
      </div>
    </div>
  </section>

  <section class="mt-6 border-2 border-black bg-card shadow-nb">
    <div class="border-b-2 border-black bg-muted px-4 py-3">
      <div class="nb-kicker">Pemeliharaan & batasan</div>
      <h2 class="mt-1 text-xl">Panduan mengikuti aplikasi</h2>
      <p class="mt-1 text-xs text-muted-foreground">
        Panduan ini mengikuti implementasi pilot per 13 November 2026. Jika menu, flow, atau screenshot
        berubah, perbarui halaman ini bersamaan. Diagram memakai stack-diagram; screenshot diambil dari
        aplikasi pilot.
      </p>
    </div>
    <div class="border-t-2 border-black bg-secondary px-4 py-3 font-mono text-[10px] leading-relaxed">
      Panduan ini memakai istilah dan menu aplikasi. Hak akses di sini menjelaskan perilaku role, bukan
      pengganti kontrol akses backend. Untuk perubahan data final, gunakan mekanisme khusus yang berlaku;
      jangan melakukan edit langsung.
    </div>
  </section>

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


