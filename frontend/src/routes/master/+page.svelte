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
  import type { CapaianTahunan, Indicator, Label, Target } from '$lib/types';
  import { Plus, RefreshCw, Database, Tag, Target as TargetIcon, CalendarDays, History } from '@lucide/svelte';

  let activeTab = $state<'indicators' | 'labels' | 'targets' | 'baselines' | 'periods'>('indicators');
  let loading = $state(false);
  let error = $state('');
  let message = $state('');

  let indicatorOpen = $state(false);
  let labelOpen = $state(false);
  let targetOpen = $state(false);
  let baselineOpen = $state(false);
  let editingBaseline = $state<CapaianTahunan | null>(null);
  let baselineYear = $state(String(new Date().getFullYear() - 1));
  let baselines = $state<CapaianTahunan[]>([]);
  let baselinesLoading = $state(false);

  let indForm = $state({ kode: '', nama: '', satuan: '%', definisi: '', numerator_label: '', denominator_label: '', labelID: '' });
  let labelForm = $state({ kode: '', nama: '', tipe: 'IKU', parentID: '', aggregation_method: 'RATIO_OF_SUMS', allow_direct_input: true, is_rollup_only: false });
  let targetForm = $state({ labelID: '', periodID: '', target_tahunan: '100' });
  let baselineForm = $state({ target: '', tahun: String(new Date().getFullYear() - 1), nilai: '', catatan: '' });

  let labels = $derived($bootstrap?.labels ?? []);
  let indicators = $derived($bootstrap?.indicators ?? []);
  let targets = $derived($bootstrap?.targets ?? []);
  let periods = $derived($bootstrap?.periods ?? []);
  let isAdmin = $derived($currentActor?.role === 'admin' || $currentActor?.role === 'kabid' || $currentActor?.role === 'sekretaris');
  let canMaster = $derived($currentActor?.role === 'admin');

  let labelOptions = $derived(labels.map((l) => ({ value: l.id, label: `${l.kode} · ${l.nama}` })));
  let periodOptions = $derived(periods.map((p) => ({ value: p.id, label: p.kode })));
  let parentOptions = $derived([{ value: '', label: 'Tanpa parent' }, ...labelOptions]);
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
  let labelColumns = [
    { key: 'kode', label: 'Kode', class: 'w-44' },
    { key: 'nama', label: 'Label' },
    { key: 'tipe', label: 'Tipe', class: 'w-24' },
    { key: 'aggregation', label: 'Agregasi', class: 'w-36' },
    { key: 'input', label: 'Input', class: 'w-32' },
    { key: 'parent', label: 'Parent', class: 'w-44' }
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

  async function createIndicator() {
    error = '';
    try {
      await api.createIndicator({
        ...indForm,
        label_ids: indForm.labelID ? [indForm.labelID] : []
      });
      indicatorOpen = false;
      message = 'Indikator dibuat.';
      indForm = { kode: '', nama: '', satuan: '%', definisi: '', numerator_label: '', denominator_label: '', labelID: '' };
      await reload();
    } catch (err) {
      error = err instanceof Error ? err.message : 'Gagal membuat indikator.';
    }
  }

  async function createLabel() {
    error = '';
    const releaseID = labels[0]?.master_release_id ?? '';
    try {
      await api.createLabel({
        kode: labelForm.kode,
        nama: labelForm.nama,
        tipe: labelForm.tipe,
        parent_id: labelForm.parentID,
        master_release_id: releaseID,
        aggregation_method: labelForm.aggregation_method,
        allow_direct_input: labelForm.allow_direct_input,
        is_rollup_only: labelForm.is_rollup_only
      });
      labelOpen = false;
      message = 'Label dibuat.';
      labelForm = { kode: '', nama: '', tipe: 'IKU', parentID: '', aggregation_method: 'RATIO_OF_SUMS', allow_direct_input: true, is_rollup_only: false };
      await reload();
    } catch (err) {
      error = err instanceof Error ? err.message : 'Gagal membuat label.';
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
        target_tahunan: Number(targetForm.target_tahunan)
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
    {#if activeTab === 'indicators'}
      <Button onclick={() => (indicatorOpen = true)} disabled={!canMaster}><Plus size={14} /> Indikator</Button>
    {:else if activeTab === 'labels'}
      <Button onclick={() => (labelOpen = true)} disabled={!canMaster}><Plus size={14} /> Label</Button>
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
  {:else if activeTab === 'labels'}
    <DataTable columns={labelColumns} rows={labels}>
      {#snippet cell(item: Label, key)}
        {#if key === 'kode'}<span class="font-bold">{item.kode}</span>
        {:else if key === 'nama'}{item.nama}
        {:else if key === 'tipe'}<Badge tone="neutral">{item.tipe}</Badge>
        {:else if key === 'aggregation'}{item.aggregation_method}
        {:else if key === 'input'}
          {#if item.is_rollup_only}<Badge tone="warning">Rollup</Badge>{:else}<Badge tone="success">Direct input</Badge>{/if}
        {:else if key === 'parent'}{item.parent_id ? labelMap.get(item.parent_id)?.kode ?? item.parent_id : '—'}
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

<Modal open={indicatorOpen} title="Tambah indikator" description="Indikator baru harus punya kode, satuan, dan label." onclose={() => (indicatorOpen = false)}>
  <div class="grid gap-3 sm:grid-cols-2">
    <div><label class="nb-label" for="ik">Kode</label><Input id="ik" bind:value={indForm.kode} placeholder="IND-013" /></div>
    <div><label class="nb-label" for="in">Nama</label><Input id="in" bind:value={indForm.nama} placeholder="Nama indikator" /></div>
    <div><label class="nb-label" for="is">Satuan</label><Input id="is" bind:value={indForm.satuan} placeholder="% / orang / kasus" /></div>
    <div><label class="nb-label" for="il">Label / IKU</label><Select id="il" bind:value={indForm.labelID} options={labelOptions} placeholder="Pilih label" /></div>
    <div class="sm:col-span-2"><label class="nb-label" for="id">Definisi</label><Textarea id="id" bind:value={indForm.definisi} rows={3} /></div>
    <div><label class="nb-label" for="num">Label numerator</label><Input id="num" bind:value={indForm.numerator_label} /></div>
    <div><label class="nb-label" for="den">Label denominator</label><Input id="den" bind:value={indForm.denominator_label} /></div>
  </div>
  <div class="mt-4 flex justify-end gap-2">
    <Button variant="outline" onclick={() => (indicatorOpen = false)}>Batal</Button>
    <Button variant="success" onclick={createIndicator}>Simpan indikator</Button>
  </div>
</Modal>

<Modal open={labelOpen} title="Tambah label" description="Label bisa berupa IKU, IKK, SPM, atau taksonomi lain." onclose={() => (labelOpen = false)}>
  <div class="grid gap-3 sm:grid-cols-2">
    <div><label class="nb-label" for="lk">Kode</label><Input id="lk" bind:value={labelForm.kode} placeholder="IKU-2026-002" /></div>
    <div><label class="nb-label" for="ln">Nama</label><Input id="ln" bind:value={labelForm.nama} /></div>
    <div><label class="nb-label" for="lt">Tipe</label><Select id="lt" bind:value={labelForm.tipe} options={['IKU','IKK','SPM','RIBK','PROGRAM'].map((v)=>({value:v,label:v}))} /></div>
    <div><label class="nb-label" for="lp">Parent</label><Select id="lp" bind:value={labelForm.parentID} options={parentOptions} /></div>
    <div><label class="nb-label" for="la">Agregasi</label><Select id="la" bind:value={labelForm.aggregation_method} options={['RATIO_OF_SUMS','SUM','WEIGHTED_AVERAGE','FORMULA','NONE'].map((v)=>({value:v,label:v}))} /></div>
    <div class="flex items-end gap-4 pb-1">
      <label class="flex items-center gap-2 font-mono text-[10px] font-bold uppercase"><input type="checkbox" bind:checked={labelForm.allow_direct_input} /> Direct input</label>
      <label class="flex items-center gap-2 font-mono text-[10px] font-bold uppercase"><input type="checkbox" bind:checked={labelForm.is_rollup_only} /> Rollup</label>
    </div>
  </div>
  <div class="mt-4 flex justify-end gap-2">
    <Button variant="outline" onclick={() => (labelOpen = false)}>Batal</Button>
    <Button variant="success" onclick={createLabel}>Simpan label</Button>
  </div>
</Modal>

<Modal open={targetOpen} title="Target periode" description="Target bulanan dihitung otomatis = target tahunan / 12." onclose={() => (targetOpen = false)}>
  <div class="grid gap-3 sm:grid-cols-2">
    <div><label class="nb-label" for="tl">Label</label><Select id="tl" bind:value={targetForm.labelID} options={labelOptions} placeholder="Pilih label" /></div>
    <div><label class="nb-label" for="tp">Periode</label><Select id="tp" bind:value={targetForm.periodID} options={periodOptions} placeholder="Pilih periode" /></div>
    <div><label class="nb-label" for="tt">Target tahunan</label><Input id="tt" bind:value={targetForm.target_tahunan} type="number" /></div>
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
