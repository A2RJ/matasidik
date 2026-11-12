<script lang="ts">
  import { api } from '$lib/api/client';
  import { applyBootstrap, bootstrap, currentActor } from '$lib/stores/session';
  import PageHeader from '$lib/components/ui/page-header.svelte';
  import Card from '$lib/components/ui/card/card.svelte';
  import Button from '$lib/components/ui/button/button.svelte';
  import Badge from '$lib/components/ui/badge.svelte';
  import Input from '$lib/components/ui/input.svelte';
  import Select from '$lib/components/ui/select.svelte';
  import Textarea from '$lib/components/ui/textarea.svelte';
  import Modal from '$lib/components/ui/modal.svelte';
  import DataTable from '$lib/components/ui/data-table.svelte';
  import Spinner from '$lib/components/ui/spinner.svelte';
  import { PAGE_COPY } from '$lib/constants/labels';
  import { formatDateTime, formatPercent } from '$lib/utils';
  import type { Bidang, CapaianTahunan, DashboardCategory, Indicator, Label, Target } from '$lib/types';
  import { Plus, RefreshCw, Database, Tag, Target as TargetIcon, CalendarDays, History } from '@lucide/svelte';

  let activeTab = $state<'indicators' | 'preview' | 'labels' | 'categories' | 'bidang' | 'targets' | 'baselines' | 'periods'>('indicators');
  let loading = $state(false);
  let error = $state('');
  let message = $state('');

  let indicatorOpen = $state(false);
  let labelOpen = $state(false);
  let categoryOpen = $state(false);
  let bidangOpen = $state(false);
  let editingLabel = $state<Label | null>(null);
  let editingCategory = $state<DashboardCategory | null>(null);
  let editingBidang = $state<Bidang | null>(null);
  let targetOpen = $state(false);
  let baselineOpen = $state(false);
  let editingBaseline = $state<CapaianTahunan | null>(null);
  let baselineYear = $state(String(new Date().getFullYear() - 1));
  let baselines = $state<CapaianTahunan[]>([]);
  let baselinesLoading = $state(false);

  let indForm = $state({
    kode: '', nama: '', satuan: '%', definisi: '', numerator_label: '', denominator_label: '',
    labelIDs: [] as string[], primaryLabelID: '', pjBidangID: '', pjTimkerID: ''
  });
  let labelForm = $state({
    kode: '', nama: '', categoryID: '', parentID: '',
    aggregation_method: 'RATIO_OF_SUMS', allow_direct_input: true, is_rollup_only: false
  });
  let categoryForm = $state({ kode: '', nama: '', jenis: 'LABEL', urutan: '0', aktif: true });
  let bidangForm = $state({ kode: '', nama: '', jenis: 'bidang', parentID: '', aktif: true });
  let targetForm = $state({ labelID: '', periodID: '', target_tahunan: '100', target_layanan_tahunan: '' });
  let baselineForm = $state({ target: '', tahun: String(new Date().getFullYear() - 1), nilai: '', catatan: '' });

  let labels = $derived($bootstrap?.labels ?? []);
  let indicators = $derived($bootstrap?.indicators ?? []);
  let targets = $derived($bootstrap?.targets ?? []);
  let periods = $derived($bootstrap?.periods ?? []);
  let categories = $derived($bootstrap?.categories ?? []);
  let bidangList = $derived($bootstrap?.bidang ?? []);
  let isAdmin = $derived($currentActor?.role === 'admin' || $currentActor?.role === 'kabid' || $currentActor?.role === 'sekretaris');
  let canMaster = $derived($currentActor?.role === 'admin');

  let labelOptions = $derived(labels.map((l) => ({
    value: l.id, label: `${l.category_kode ?? l.tipe} · ${l.kode} · ${l.nama}`
  })));
  let categoryOptions = $derived(categories.filter((c) => c.aktif && c.jenis === 'LABEL').map((c) => ({ value: c.id, label: `${c.kode} · ${c.nama}` })));
  let bidangOptions = $derived(bidangList.filter((b) => b.aktif && b.jenis !== 'timker').map((b) => ({ value: b.id, label: `${b.kode} · ${b.nama}` })));
  let timkerOptions = $derived(bidangList.filter((b) => b.aktif && b.jenis === 'timker').map((b) => ({ value: b.id, label: `${b.kode} · ${b.nama}` })));
  let selectedLabelOptions = $derived(labelOptions.filter((o) => indForm.labelIDs.includes(o.value)));
  let periodOptions = $derived(periods.map((p) => ({ value: p.id, label: p.kode })));
  let parentOptions = $derived([{ value: '', label: 'Tanpa parent' }, ...labelOptions]);
  let bidangParentOptions = $derived([
    { value: '', label: 'Tanpa parent' },
    ...bidangList.filter((b) => b.aktif && b.jenis !== 'timker').map((b) => ({ value: b.id, label: `${b.kode} · ${b.nama}` }))
  ]);
  let baselineTargetOptions = $derived([
    ...labels.map((l) => ({ value: `label:${l.id}`, label: `IKU · ${l.kode} · ${l.nama}` })),
    ...indicators.map((i) => ({ value: `indicator:${i.id}`, label: `Indikator · ${i.kode} · ${i.nama}` }))
  ]);
  let kabupatenID = $derived($bootstrap?.wilayah.find((w) => w.level === 'kabupaten')?.id ?? '');

  let indicatorColumns = [
    { key: 'kode', label: 'Kode', class: 'w-32' },
    { key: 'nama', label: 'Indikator' },
    { key: 'satuan', label: 'Satuan', class: 'w-24' },
    { key: 'sifat', label: 'Sifat input', class: 'w-32' },
    { key: 'frekuensi', label: 'Frekuensi', class: 'w-28' },
    { key: 'pj', label: 'PJ' }
  ];
  let previewColumns = [
    { key: 'no', label: 'No.', width: '56px' },
    { key: 'nama', label: 'Nama Indikator', width: '22%' },
    { key: 'satuan', label: 'Satuan', width: '90px' },
    { key: 'definisi', label: 'Definisi Operasional', width: '22%' },
    { key: 'label', label: 'Label Indikator', width: '20%' },
    { key: 'sifat', label: 'Sifat Data', width: '140px' },
    { key: 'bidang', label: 'Bidang/Sekretariat', width: '160px' },
    { key: 'timker', label: 'Timker/Sub Bagian', width: '160px' }
  ];
  let labelColumns = [
    { key: 'kode', label: 'Kode', class: 'w-44' },
    { key: 'nama', label: 'Label' },
    { key: 'category', label: 'Kategori', class: 'w-40' },
    { key: 'aggregation', label: 'Agregasi', class: 'w-36' },
    { key: 'input', label: 'Input', class: 'w-32' },
    { key: 'parent', label: 'Parent', class: 'w-44' },
    { key: 'action', label: 'Aksi', class: 'w-24' }
  ];
  let categoryColumns = [
    { key: 'urutan', label: 'Urutan', class: 'w-20' },
    { key: 'kode', label: 'Kode', class: 'w-40' },
    { key: 'nama', label: 'Kategori' },
    { key: 'jenis', label: 'Jenis', class: 'w-28' },
    { key: 'aktif', label: 'Status', class: 'w-28' },
    { key: 'action', label: 'Aksi', class: 'w-40' }
  ];
  let bidangColumns = [
    { key: 'kode', label: 'Kode', class: 'w-32' },
    { key: 'nama', label: 'Bidang/Timker' },
    { key: 'jenis', label: 'Jenis', class: 'w-28' },
    { key: 'parent', label: 'Parent', class: 'w-48' },
    { key: 'aktif', label: 'Status', class: 'w-28' },
    { key: 'action', label: 'Aksi', class: 'w-24' }
  ];
  let targetColumns = [
    { key: 'label', label: 'Label' },
    { key: 'periode', label: 'Periode', class: 'w-28' },
    { key: 'scope', label: 'Scope', class: 'w-32' },
    { key: 'tahunan', label: 'Target tahunan', class: 'text-right w-32' },
    { key: 'bulanan', label: 'Target bulanan', class: 'text-right w-32' }
  ];
  let baselineColumns = [
    { key: 'target', label: 'Label / Indikator' },
    { key: 'tahun', label: 'Tahun', class: 'w-20' },
    { key: 'scope', label: 'Scope', class: 'w-28' },
    { key: 'nilai', label: 'Capaian', class: 'text-right w-28' },
    { key: 'catatan', label: 'Catatan', class: 'w-52' },
    { key: 'updated', label: 'Diperbarui', class: 'w-44' },
    { key: 'action', label: 'Aksi', class: 'w-32' }
  ];
  let periodColumns = [
    { key: 'kode', label: 'Kode', class: 'w-32' },
    { key: 'tahun', label: 'Tahun', class: 'w-24' },
    { key: 'bulan', label: 'Bulan', class: 'w-24' },
    { key: 'status', label: 'Status', class: 'w-24' }
  ];

  const labelMap = $derived(new Map(labels.map((l) => [l.id, l])));
  const periodMap = $derived(new Map(periods.map((p) => [p.id, p])));
  const bidangMap = $derived(new Map(bidangList.map((b) => [b.id, b])));

  function indicatorLabels(item: Indicator): Label[] {
    return (item.label_ids ?? [])
      .map((id) => labelMap.get(id))
      .filter((label): label is Label => !!label);
  }

  function bidangName(id?: string): string {
    return id ? bidangMap.get(id)?.nama ?? id : '—';
  }

  async function reload() {
    loading = true;
    error = '';
    try {
      applyBootstrap(await api.bootstrap());
    } catch (err) {
      error = err instanceof Error ? err.message : 'Gagal memuat master.';
    } finally {
      loading = false;
    }
  }

  function resetIndicatorForm() {
    indForm = {
      kode: '', nama: '', satuan: '%', definisi: '', numerator_label: '', denominator_label: '',
      labelIDs: [], primaryLabelID: '', pjBidangID: '', pjTimkerID: ''
    };
  }

  async function createIndicator() {
    error = '';
    if (indForm.labelIDs.length === 0) {
      error = 'Pilih minimal satu Label Indikator.';
      return;
    }
    const primaryLabelID = indForm.labelIDs.includes(indForm.primaryLabelID)
      ? indForm.primaryLabelID
      : indForm.labelIDs[0];
    try {
      await api.createIndicator({
        ...indForm,
        label_ids: indForm.labelIDs,
        primary_label_id: primaryLabelID
      });
      indicatorOpen = false;
      message = 'Indikator dibuat.';
      resetIndicatorForm();
      await reload();
    } catch (err) {
      error = err instanceof Error ? err.message : 'Gagal membuat indikator.';
    }
  }

  function openLabel(item?: Label) {
    editingLabel = item ?? null;
    labelForm = item
      ? {
          kode: item.kode,
          nama: item.nama,
          categoryID: item.category_id ?? '',
          parentID: item.parent_id ?? '',
          aggregation_method: item.aggregation_method,
          allow_direct_input: item.allow_direct_input,
          is_rollup_only: item.is_rollup_only
        }
      : {
          kode: '', nama: '', categoryID: categoryOptions[0]?.value ?? '', parentID: '',
          aggregation_method: 'RATIO_OF_SUMS', allow_direct_input: true, is_rollup_only: false
        };
    labelOpen = true;
  }

  async function saveLabel() {
    error = '';
    const releaseID = labels[0]?.master_release_id ?? '';
    try {
      if (editingLabel) {
        await api.updateLabel(editingLabel.id, {
          nama: labelForm.nama,
          category_id: labelForm.categoryID,
          parent_id: labelForm.parentID,
          allow_direct_input: labelForm.allow_direct_input,
          is_rollup_only: labelForm.is_rollup_only
        });
        message = 'Label diperbarui.';
      } else {
        await api.createLabel({
          kode: labelForm.kode,
          nama: labelForm.nama,
          category_id: labelForm.categoryID,
          parent_id: labelForm.parentID,
          master_release_id: releaseID,
          aggregation_method: labelForm.aggregation_method,
          allow_direct_input: labelForm.allow_direct_input,
          is_rollup_only: labelForm.is_rollup_only
        });
        message = 'Label dibuat.';
      }
      labelOpen = false;
      editingLabel = null;
      await reload();
    } catch (err) {
      error = err instanceof Error ? err.message : 'Gagal menyimpan label.';
    }
  }

  function openCategory(item?: DashboardCategory) {
    editingCategory = item ?? null;
    categoryForm = item
      ? { kode: item.kode, nama: item.nama, jenis: item.jenis, urutan: String(item.urutan), aktif: item.aktif }
      : { kode: '', nama: '', jenis: 'LABEL', urutan: '0', aktif: true };
    categoryOpen = true;
  }

  async function saveCategory() {
    error = '';
    const payload = {
      kode: categoryForm.kode,
      nama: categoryForm.nama,
      jenis: categoryForm.jenis,
      urutan: Number(categoryForm.urutan) || 0,
      aktif: categoryForm.aktif
    };
    try {
      if (editingCategory) {
        await api.updateCategory(editingCategory.id, payload);
        message = 'Kategori diperbarui.';
      } else {
        await api.createCategory(payload);
        message = 'Kategori dibuat.';
      }
      categoryOpen = false;
      editingCategory = null;
      await reload();
    } catch (err) {
      error = err instanceof Error ? err.message : 'Gagal menyimpan kategori.';
    }
  }

  async function deactivateCategory(item: DashboardCategory) {
    if (!confirm(`Nonaktifkan kategori ${item.nama}?`)) return;
    error = '';
    try {
      await api.deleteCategory(item.id);
      message = 'Kategori dinonaktifkan.';
      await reload();
    } catch (err) {
      error = err instanceof Error ? err.message : 'Gagal menonaktifkan kategori.';
    }
  }

  function openBidang(item?: Bidang) {
    editingBidang = item ?? null;
    bidangForm = item
      ? { kode: item.kode, nama: item.nama, jenis: item.jenis, parentID: item.parent_id ?? '', aktif: item.aktif }
      : { kode: '', nama: '', jenis: 'bidang', parentID: '', aktif: true };
    bidangOpen = true;
  }

  async function saveBidang() {
    error = '';
    const payload = {
      kode: bidangForm.kode,
      nama: bidangForm.nama,
      jenis: bidangForm.jenis,
      parent_id: bidangForm.parentID,
      aktif: bidangForm.aktif
    };
    try {
      if (editingBidang) {
        await api.updateBidang(editingBidang.id, payload);
        message = 'Bidang/timker diperbarui.';
      } else {
        await api.createBidang(payload);
        message = 'Bidang/timker dibuat.';
      }
      bidangOpen = false;
      editingBidang = null;
      await reload();
    } catch (err) {
      error = err instanceof Error ? err.message : 'Gagal menyimpan bidang/timker.';
    }
  }

  async function saveTarget() {
    error = '';
    try {
      await api.upsertTarget({
        label_id: targetForm.labelID,
        periode_id: targetForm.periodID,
        scope_type: 'KABUPATEN',
        scope_id: 'wil-kab-swb',
        target_tahunan: Number(targetForm.target_tahunan),
        target_layanan_tahunan: targetForm.target_layanan_tahunan
          ? Number(targetForm.target_layanan_tahunan)
          : undefined
      });
      targetOpen = false;
      message = 'Target disimpan.';
      await reload();
    } catch (err) {
      error = err instanceof Error ? err.message : 'Gagal menyimpan target.';
    }
  }

  async function loadBaselines(year = Number(baselineYear)) {
    if (!year || year < 2000) {
      error = 'Tahun capaian tidak valid.';
      return;
    }
    baselinesLoading = true;
    error = '';
    try {
      baselines = await api.capaianTahunan(year);
    } catch (err) {
      error = err instanceof Error ? err.message : 'Gagal memuat capaian tahunan.';
    } finally {
      baselinesLoading = false;
    }
  }

  $effect(() => {
    const year = Number(baselineYear);
    if (activeTab === 'baselines' && year >= 2000) {
      void loadBaselines(year);
    }
  });

  function baselineName(item: CapaianTahunan): string {
    if (item.label_id) {
      return `${item.label_kode ?? item.label_id} · ${item.label_nama ?? ''}`.trim();
    }
    return `${item.indikator_kode ?? item.indikator_id ?? '—'} · ${item.indikator_nama ?? ''}`.trim();
  }

  function openBaseline(item?: CapaianTahunan) {
    editingBaseline = item ?? null;
    baselineForm = item
      ? {
          target: item.label_id ? `label:${item.label_id}` : `indicator:${item.indikator_id}`,
          tahun: String(item.tahun),
          nilai: String(item.nilai),
          catatan: item.catatan ?? ''
        }
      : { target: '', tahun: baselineYear, nilai: '', catatan: '' };
    baselineOpen = true;
  }

  async function saveBaseline() {
    error = '';
    if (!baselineForm.target) {
      error = 'Pilih label/IKU atau indikator.';
      return;
    }
    if (baselineForm.nilai === '' || Number.isNaN(Number(baselineForm.nilai))) {
      error = 'Nilai capaian wajib diisi.';
      return;
    }
    const [kind, targetID] = baselineForm.target.split(':');
    const payload = {
      label_id: kind === 'label' ? targetID : '',
      indikator_id: kind === 'indicator' ? targetID : '',
      tahun: Number(baselineForm.tahun),
      scope_type: 'KABUPATEN',
      scope_id: kabupatenID,
      nilai: Number(baselineForm.nilai),
      catatan: baselineForm.catatan
    };
    try {
      if (editingBaseline) {
        await api.updateCapaianTahunan(editingBaseline.id, payload);
      } else {
        await api.createCapaianTahunan(payload);
      }
      baselineOpen = false;
      baselineYear = String(payload.tahun);
      message = 'Capaian tahunan disimpan.';
      await loadBaselines(payload.tahun);
    } catch (err) {
      error = err instanceof Error ? err.message : 'Gagal menyimpan capaian tahunan.';
    }
  }

  async function deleteBaseline(item: CapaianTahunan) {
    if (!confirm(`Hapus capaian ${item.tahun} untuk ${baselineName(item)}?`)) return;
    error = '';
    try {
      await api.deleteCapaianTahunan(item.id);
      message = 'Capaian tahunan dihapus.';
      await loadBaselines();
    } catch (err) {
      error = err instanceof Error ? err.message : 'Gagal menghapus capaian tahunan.';
    }
  }
</script>

<PageHeader title={PAGE_COPY.master.title} description={PAGE_COPY.master.description}>
  {#snippet actions()}
    <Button variant="outline" onclick={reload}><RefreshCw size={14} /> Muat ulang</Button>
    <Button variant={activeTab === 'preview' ? 'default' : 'outline'} onclick={() => (activeTab = 'preview')}><Tag size={14} /> Pratinjau Pemetaan</Button>
    {#if activeTab === 'indicators'}
      <Button onclick={() => (indicatorOpen = true)} disabled={!canMaster}><Plus size={14} /> Indikator</Button>
    {:else if activeTab === 'labels'}
      <Button onclick={() => openLabel()} disabled={!canMaster}><Plus size={14} /> Label</Button>
    {:else if activeTab === 'categories'}
      <Button onclick={() => openCategory()} disabled={!canMaster}><Plus size={14} /> Kategori</Button>
    {:else if activeTab === 'bidang'}
      <Button onclick={() => openBidang()} disabled={!canMaster}><Plus size={14} /> Bidang/Timker</Button>
    {:else if activeTab === 'targets'}
      <Button onclick={() => (targetOpen = true)} disabled={!isAdmin}><Plus size={14} /> Target</Button>
    {:else if activeTab === 'baselines'}
      <Button onclick={() => openBaseline()} disabled={!canMaster}><Plus size={14} /> Capaian</Button>
    {/if}
  {/snippet}
</PageHeader>

{#if error}<div class="mb-4 border-2 border-black bg-destructive p-3 font-mono text-xs font-bold">{error}</div>{/if}
{#if message}<div class="mb-4 border-2 border-black bg-success p-3 font-mono text-xs font-bold">{message}</div>{/if}

<div class="mb-4 flex flex-wrap gap-2">
  <Button variant={activeTab === 'indicators' ? 'default' : 'outline'} onclick={() => (activeTab = 'indicators')}><Database size={14} /> Indikator ({indicators.length})</Button>
  <Button variant={activeTab === 'labels' ? 'default' : 'outline'} onclick={() => (activeTab = 'labels')}><Tag size={14} /> Label & IKU ({labels.length})</Button>
  <Button variant={activeTab === 'categories' ? 'default' : 'outline'} onclick={() => (activeTab = 'categories')}><Tag size={14} /> Kategori ({categories.length})</Button>
  <Button variant={activeTab === 'bidang' ? 'default' : 'outline'} onclick={() => (activeTab = 'bidang')}><Tag size={14} /> Bidang/Timker ({bidangList.length})</Button>
  <Button variant={activeTab === 'targets' ? 'default' : 'outline'} onclick={() => (activeTab = 'targets')}><TargetIcon size={14} /> Target ({targets.length})</Button>
  <Button variant={activeTab === 'baselines' ? 'default' : 'outline'} onclick={() => (activeTab = 'baselines')}><History size={14} /> Capaian tahunan ({baselines.length})</Button>
  <Button variant={activeTab === 'periods' ? 'default' : 'outline'} onclick={() => (activeTab = 'periods')}><CalendarDays size={14} /> Periode ({periods.length})</Button>
</div>

{#if activeTab === 'baselines'}
  <div class="mb-4 flex flex-wrap items-end gap-3 border-2 border-black bg-card p-3 shadow-nbsm">
    <div class="w-32">
      <label class="nb-label" for="baseline-year">Tahun capaian</label>
      <Input id="baseline-year" type="number" bind:value={baselineYear} />
    </div>
    <Button variant="outline" onclick={() => loadBaselines()}>Muat</Button>
    <p class="max-w-xl font-mono text-[10px] text-muted-foreground">
      Capaian tahun X dipakai sebagai pembanding dashboard periode X+1. Contoh: isi 2025 untuk dashboard 2026.
    </p>
  </div>
{/if}

<Card>
  {#if loading}
    <div class="flex items-center justify-center py-16"><Spinner label="Memuat master…" /></div>
  {:else if activeTab === 'indicators'}
    <DataTable columns={indicatorColumns} rows={indicators}>
      {#snippet cell(item: Indicator, key)}
        {#if key === 'kode'}<span class="font-bold">{item.kode}</span>
        {:else if key === 'nama'}<span class="font-semibold">{item.nama}</span><div class="text-[10px] text-muted-foreground">{item.definisi ?? 'Definisi belum diisi'}</div>
        {:else if key === 'satuan'}{item.satuan}
        {:else if key === 'sifat'}{item.sifat_input}
        {:else if key === 'frekuensi'}{item.frekuensi}
        {:else if key === 'pj'}{item.pj_timker ?? item.pj_bidang_id ?? '—'}
        {/if}
      {/snippet}
    </DataTable>
  {:else if activeTab === 'preview'}
    <DataTable columns={previewColumns} rows={indicators} minWidth="1400px" nowrap>
      {#snippet cell(item: Indicator, key)}
        {#if key === 'no'}<span class="font-mono">{indicators.indexOf(item) + 1}</span>
        {:else if key === 'nama'}
          <div class="font-semibold">{item.nama}</div>
          <div class="font-mono text-[10px] text-muted-foreground">{item.kode}</div>
        {:else if key === 'satuan'}{item.satuan}
        {:else if key === 'definisi'}{item.definisi ?? '—'}
        {:else if key === 'label'}
          {@const itemLabels = indicatorLabels(item)}
          {#if itemLabels.length === 0}
            <span class="text-muted-foreground">—</span>
          {:else}
            <div class="flex min-w-0 max-w-full flex-wrap gap-1">
              {#each itemLabels as label (label.id)}
                <Badge class="max-w-full" tone={label.id === item.primary_label_id ? 'success' : 'info'}>
                  {label.category_kode ?? label.tipe}
                </Badge>
                <span class="max-w-full break-all font-mono text-[10px]">{label.kode}</span>
              {/each}
            </div>
          {/if}
        {:else if key === 'sifat'}{item.sifat_input} · {item.frekuensi}
        {:else if key === 'bidang'}{bidangName(item.pj_bidang_id)}
        {:else if key === 'timker'}
          {item.pj_timker_id ? bidangName(item.pj_timker_id) : (item.pj_timker ?? '—')}
        {/if}
      {/snippet}
    </DataTable>
    <div class="border-t-2 border-black bg-card px-4 py-2 font-mono text-[10px] text-muted-foreground">
      Pratinjau pemetaan indikator ke label/kategori dan PJ. Label Indikator berasal dari label/kategori; Bidang dan Timker berasal dari kolom PJ, bukan dari label.
    </div>
  {:else if activeTab === 'labels'}
    <DataTable columns={labelColumns} rows={labels}>
      {#snippet cell(item: Label, key)}
        {#if key === 'kode'}<span class="font-bold">{item.kode}</span>
        {:else if key === 'nama'}{item.nama}
        {:else if key === 'category'}<Badge tone="neutral">{item.category_nama ?? item.category_kode ?? item.tipe}</Badge>
        {:else if key === 'aggregation'}{item.aggregation_method}
        {:else if key === 'input'}
          {#if item.is_rollup_only}<Badge tone="warning">Rollup</Badge>{:else}<Badge tone="success">Direct input</Badge>{/if}
        {:else if key === 'parent'}{item.parent_id ? labelMap.get(item.parent_id)?.kode ?? item.parent_id : '—'}
        {:else if key === 'action'}
          <Button size="sm" variant="outline" onclick={() => openLabel(item)} disabled={!canMaster}>Edit</Button>
        {/if}
      {/snippet}
    </DataTable>
  {:else if activeTab === 'categories'}
    <DataTable columns={categoryColumns} rows={categories}>
      {#snippet cell(item: DashboardCategory, key)}
        {#if key === 'urutan'}<span class="font-mono">{item.urutan}</span>
        {:else if key === 'kode'}<span class="font-bold">{item.kode}</span>
        {:else if key === 'nama'}{item.nama}
        {:else if key === 'jenis'}<Badge tone={item.jenis === 'OWNER' ? 'warning' : 'info'}>{item.jenis}</Badge>
        {:else if key === 'aktif'}<Badge tone={item.aktif ? 'success' : 'neutral'}>{item.aktif ? 'Aktif' : 'Nonaktif'}</Badge>
        {:else if key === 'action'}
          <div class="flex gap-1">
            <Button size="sm" variant="outline" onclick={() => openCategory(item)} disabled={!canMaster}>Edit</Button>
            {#if item.aktif}<Button size="sm" variant="ghost" onclick={() => deactivateCategory(item)} disabled={!canMaster}>Nonaktifkan</Button>{/if}
          </div>
        {/if}
      {/snippet}
    </DataTable>
  {:else if activeTab === 'bidang'}
    <DataTable columns={bidangColumns} rows={bidangList}>
      {#snippet cell(item: Bidang, key)}
        {#if key === 'kode'}<span class="font-bold">{item.kode}</span>
        {:else if key === 'nama'}{item.nama}
        {:else if key === 'jenis'}<Badge tone={item.jenis === 'timker' ? 'warning' : 'info'}>{item.jenis}</Badge>
        {:else if key === 'parent'}{item.parent_id ? bidangMap.get(item.parent_id)?.kode ?? item.parent_id : '—'}
        {:else if key === 'aktif'}<Badge tone={item.aktif ? 'success' : 'neutral'}>{item.aktif ? 'Aktif' : 'Nonaktif'}</Badge>
        {:else if key === 'action'}
          <Button size="sm" variant="outline" onclick={() => openBidang(item)} disabled={!canMaster}>Edit</Button>
        {/if}
      {/snippet}
    </DataTable>
  {:else if activeTab === 'targets'}
    <DataTable columns={targetColumns} rows={targets}>
      {#snippet cell(item: Target, key)}
        {#if key === 'label'}{item.label_id ? labelMap.get(item.label_id)?.kode ?? item.label_id : item.indikator_id ?? '—'}
        {:else if key === 'periode'}{item.periode_id ? periodMap.get(item.periode_id)?.kode ?? item.periode_id : '—'}
        {:else if key === 'scope'}{item.scope_type}/{item.scope_id ?? '—'}
        {:else if key === 'tahunan'}{formatPercent(item.target_tahunan, 0)}
        {:else if key === 'bulanan'}{formatPercent(item.target_bulanan, 2)}
        {/if}
      {/snippet}
    </DataTable>
  {:else if activeTab === 'baselines'}
    {#if baselinesLoading}
      <div class="flex items-center justify-center py-16"><Spinner label="Memuat capaian tahunan…" /></div>
    {:else if baselines.length === 0}
      <div class="p-6 text-center text-sm text-muted-foreground">
        Belum ada capaian tahunan untuk {baselineYear}. Tambahkan capaian per label/IKU atau indikator.
      </div>
    {:else}
      <DataTable columns={baselineColumns} rows={baselines}>
        {#snippet cell(item: CapaianTahunan, key)}
          {#if key === 'target'}
            <div class="font-semibold">{baselineName(item)}</div>
            <div class="text-[10px] text-muted-foreground">{item.label_id ? 'Label / IKU' : 'Indikator'}</div>
          {:else if key === 'tahun'}{item.tahun}
          {:else if key === 'scope'}{item.scope_type}/{item.scope_id ?? '—'}
          {:else if key === 'nilai'}<span class="font-bold">{formatPercent(item.nilai, 2)}</span>
          {:else if key === 'catatan'}{item.catatan ?? '—'}
          {:else if key === 'updated'}
            <div>{formatDateTime(item.updated_at)}</div>
            <div class="text-[10px] text-muted-foreground">{item.updated_by_nama ?? '—'}</div>
          {:else if key === 'action'}
            <div class="flex gap-1">
              <Button size="sm" variant="outline" onclick={() => openBaseline(item)} disabled={!canMaster}>Edit</Button>
              <Button size="sm" variant="ghost" onclick={() => deleteBaseline(item)} disabled={!canMaster}>Hapus</Button>
            </div>
          {/if}
        {/snippet}
      </DataTable>
    {/if}
  {:else}
    <DataTable columns={periodColumns} rows={periods}>
      {#snippet cell(item, key)}
        {#if key === 'kode'}<span class="font-bold">{item.kode}</span>
        {:else if key === 'tahun'}{item.tahun}
        {:else if key === 'bulan'}{item.bulan}
        {:else if key === 'status'}<Badge tone={item.status === 'open' ? 'success' : 'warning'}>{item.status}</Badge>
        {/if}
      {/snippet}
    </DataTable>
  {/if}
</Card>

<Modal open={indicatorOpen} title="Tambah indikator" description="Indikator memilih label/IKU spesifik; kategori mengikuti label." onclose={() => (indicatorOpen = false)}>
  <div class="grid gap-3 sm:grid-cols-2">
    <div><label class="nb-label" for="ik">Kode</label><Input id="ik" bind:value={indForm.kode} placeholder="IND-013" /></div>
    <div><label class="nb-label" for="in">Nama</label><Input id="in" bind:value={indForm.nama} placeholder="Nama indikator" /></div>
    <div><label class="nb-label" for="is">Satuan</label><Input id="is" bind:value={indForm.satuan} placeholder="% / orang / kasus" /></div>
    <div>
      <label class="nb-label" for="il">Label Indikator (bisa lebih dari satu)</label>
      <select id="il" multiple bind:value={indForm.labelIDs} class="nb-select min-h-32 w-full">
        {#each labelOptions as option (option.value)}
          <option value={option.value}>{option.label}</option>
        {/each}
      </select>
      <p class="mt-1 font-mono text-[10px] text-muted-foreground">Gunakan Ctrl/Cmd untuk memilih beberapa label.</p>
    </div>
    <div>
      <label class="nb-label" for="ip">Label utama</label>
      <Select id="ip" bind:value={indForm.primaryLabelID} options={selectedLabelOptions} placeholder="Pilih label utama" />
    </div>
    <div>
      <label class="nb-label" for="pb">PJ Bidang</label>
      <Select id="pb" bind:value={indForm.pjBidangID} options={bidangOptions} placeholder="Pilih bidang" />
    </div>
    <div>
      <label class="nb-label" for="pt">PJ Timker</label>
      <Select id="pt" bind:value={indForm.pjTimkerID} options={timkerOptions} placeholder="Pilih timker" />
    </div>
    <div class="sm:col-span-2"><label class="nb-label" for="id">Definisi</label><Textarea id="id" bind:value={indForm.definisi} rows={3} /></div>
    <div><label class="nb-label" for="num">Label numerator</label><Input id="num" bind:value={indForm.numerator_label} /></div>
    <div><label class="nb-label" for="den">Label denominator</label><Input id="den" bind:value={indForm.denominator_label} /></div>
  </div>
  <div class="mt-4 flex justify-end gap-2">
    <Button variant="outline" onclick={() => (indicatorOpen = false)}>Batal</Button>
    <Button variant="success" onclick={createIndicator}>Simpan indikator</Button>
  </div>
</Modal>

<Modal open={labelOpen} title={editingLabel ? 'Edit label/IKU' : 'Tambah label/IKU'} description="Setiap label/IKU wajib punya satu kategori." onclose={() => (labelOpen = false)}>
  <div class="grid gap-3 sm:grid-cols-2">
    <div><label class="nb-label" for="lk">Kode</label><Input id="lk" bind:value={labelForm.kode} placeholder="IKU-2026-002" disabled={!!editingLabel} /></div>
    <div><label class="nb-label" for="ln">Nama</label><Input id="ln" bind:value={labelForm.nama} /></div>
    <div><label class="nb-label" for="lc">Kategori</label><Select id="lc" bind:value={labelForm.categoryID} options={categoryOptions} placeholder="Pilih kategori" /></div>
    <div><label class="nb-label" for="lp">Parent</label><Select id="lp" bind:value={labelForm.parentID} options={parentOptions} /></div>
    <div><label class="nb-label" for="la">Agregasi</label><Select id="la" bind:value={labelForm.aggregation_method} options={['RATIO_OF_SUMS','SUM','WEIGHTED_AVERAGE','FORMULA','NONE'].map((v)=>({value:v,label:v}))} disabled={!!editingLabel} /></div>
    <div class="flex items-end gap-4 pb-1">
      <label class="flex items-center gap-2 font-mono text-[10px] font-bold uppercase"><input type="checkbox" bind:checked={labelForm.allow_direct_input} /> Direct input</label>
      <label class="flex items-center gap-2 font-mono text-[10px] font-bold uppercase"><input type="checkbox" bind:checked={labelForm.is_rollup_only} /> Rollup</label>
    </div>
  </div>
  <div class="mt-4 flex justify-end gap-2">
    <Button variant="outline" onclick={() => (labelOpen = false)}>Batal</Button>
    <Button variant="success" onclick={saveLabel}>Simpan label</Button>
  </div>
</Modal>

<Modal open={categoryOpen} title={editingCategory ? 'Edit kategori' : 'Tambah kategori'} description="Kategori mengatur filter dashboard dan taksonomi label." onclose={() => (categoryOpen = false)}>
  <div class="grid gap-3 sm:grid-cols-2">
    <div><label class="nb-label" for="ck">Kode</label><Input id="ck" bind:value={categoryForm.kode} placeholder="IKK" disabled={!!editingCategory} /></div>
    <div><label class="nb-label" for="cn">Nama</label><Input id="cn" bind:value={categoryForm.nama} placeholder="IKK" /></div>
    <div><label class="nb-label" for="cj">Jenis</label><Select id="cj" bind:value={categoryForm.jenis} options={[{value:'LABEL',label:'LABEL'},{value:'OWNER',label:'OWNER'}]} disabled={!!editingCategory?.is_system} /></div>
    <div><label class="nb-label" for="cu">Urutan</label><Input id="cu" type="number" bind:value={categoryForm.urutan} /></div>
    <label class="flex items-center gap-2 font-mono text-[10px] font-bold uppercase"><input type="checkbox" bind:checked={categoryForm.aktif} /> Aktif</label>
  </div>
  <div class="mt-4 flex justify-end gap-2">
    <Button variant="outline" onclick={() => (categoryOpen = false)}>Batal</Button>
    <Button variant="success" onclick={saveCategory}>Simpan kategori</Button>
  </div>
</Modal>

<Modal open={bidangOpen} title={editingBidang ? 'Edit bidang/timker' : 'Tambah bidang/timker'} description="Bidang, sekretariat, dan timker memakai satu tabel dengan jenis berbeda." onclose={() => (bidangOpen = false)}>
  <div class="grid gap-3 sm:grid-cols-2">
    <div><label class="nb-label" for="bk">Kode</label><Input id="bk" bind:value={bidangForm.kode} placeholder="KESGA" /></div>
    <div><label class="nb-label" for="bn">Nama</label><Input id="bn" bind:value={bidangForm.nama} placeholder="Timker Kesehatan Keluarga" /></div>
    <div><label class="nb-label" for="bj">Jenis</label><Select id="bj" bind:value={bidangForm.jenis} options={[{value:'bidang',label:'Bidang'},{value:'sekretariat',label:'Sekretariat'},{value:'timker',label:'Timker'}]} /></div>
    <div><label class="nb-label" for="bp">Parent</label><Select id="bp" bind:value={bidangForm.parentID} options={bidangParentOptions} /></div>
    <label class="flex items-center gap-2 font-mono text-[10px] font-bold uppercase"><input type="checkbox" bind:checked={bidangForm.aktif} /> Aktif</label>
  </div>
  <div class="mt-4 flex justify-end gap-2">
    <Button variant="outline" onclick={() => (bidangOpen = false)}>Batal</Button>
    <Button variant="success" onclick={saveBidang}>Simpan bidang/timker</Button>
  </div>
</Modal>

<Modal open={targetOpen} title="Target periode" description="Target bulanan dihitung otomatis = target tahunan / 12." onclose={() => (targetOpen = false)}>
  <div class="grid gap-3 sm:grid-cols-2">
    <div><label class="nb-label" for="tl">Label</label><Select id="tl" bind:value={targetForm.labelID} options={labelOptions} placeholder="Pilih label" /></div>
    <div><label class="nb-label" for="tp">Periode</label><Select id="tp" bind:value={targetForm.periodID} options={periodOptions} placeholder="Pilih periode" /></div>
    <div><label class="nb-label" for="tt">Target tahunan (%)</label><Input id="tt" bind:value={targetForm.target_tahunan} type="number" /></div>
    <div><label class="nb-label" for="tlt">Target layanan tahunan</label><Input id="tlt" bind:value={targetForm.target_layanan_tahunan} type="number" placeholder="Opsional, mis. 2000" /></div>
  </div>
  <div class="mt-4 flex justify-end gap-2">
    <Button variant="outline" onclick={() => (targetOpen = false)}>Batal</Button>
    <Button variant="success" onclick={saveTarget}>Simpan target</Button>
  </div>
</Modal>

<Modal
  open={baselineOpen}
  title={editingBaseline ? 'Edit capaian tahunan' : 'Tambah capaian tahunan'}
  description="Capaian historis per label/IKU atau indikator. Nilai ini menjadi pembanding dashboard tahun berikutnya."
  onclose={() => (baselineOpen = false)}
>
  <div class="grid gap-3 sm:grid-cols-2">
    <div class="sm:col-span-2">
      <label class="nb-label" for="ct">Label / IKU / Indikator</label>
      <Select
        id="ct"
        bind:value={baselineForm.target}
        options={baselineTargetOptions}
        placeholder="Pilih label/IKU atau indikator"
      />
    </div>
    <div>
      <label class="nb-label" for="cy">Tahun capaian</label>
      <Input id="cy" type="number" bind:value={baselineForm.tahun} placeholder="2025" />
    </div>
    <div>
      <label class="nb-label" for="cn">Nilai capaian</label>
      <Input id="cn" type="number" bind:value={baselineForm.nilai} placeholder="90" />
    </div>
    <div class="sm:col-span-2">
      <label class="nb-label" for="cc">Catatan</label>
      <Textarea id="cc" bind:value={baselineForm.catatan} rows={3} placeholder="Opsional" />
    </div>
  </div>
  <div class="mt-4 flex justify-end gap-2">
    <Button variant="outline" onclick={() => (baselineOpen = false)}>Batal</Button>
    <Button variant="success" onclick={saveBaseline}>Simpan capaian</Button>
  </div>
</Modal>
