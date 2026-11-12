<script lang="ts">
  import { api } from '$lib/api/client';
  import { bootstrap, currentActor, periods, selectedPeriod } from '$lib/stores/session';
  import PageHeader from '$lib/components/ui/page-header.svelte';
  import Card from '$lib/components/ui/card/card.svelte';
  import Button from '$lib/components/ui/button/button.svelte';
  import Select from '$lib/components/ui/select.svelte';
  import Textarea from '$lib/components/ui/textarea.svelte';
  import Modal from '$lib/components/ui/modal.svelte';
  import DataTable from '$lib/components/ui/data-table.svelte';
  import Spinner from '$lib/components/ui/spinner.svelte';
  import StatusBadge from '$lib/components/ui/status-badge.svelte';
  import { PAGE_COPY, EMPTY_STATE } from '$lib/constants/labels';
  import { formatDateTime, formatPercent } from '$lib/utils';
  import type { FactorRow, TindakLanjut } from '$lib/types';
  import { MessageSquare, Plus, RefreshCw, CheckSquare } from '@lucide/svelte';

  let tab = $state<'faktor' | 'tl'>('faktor');
  let factors = $state<FactorRow[]>([]);
  let tl = $state<TindakLanjut[]>([]);
  let loading = $state(true);
  let error = $state('');
  let message = $state('');

  let factorOpen = $state(false);
  let selectedFactor = $state<FactorRow | null>(null);
  let pendukung = $state('');
  let penghambat = $state('');

  let tlOpen = $state(false);
  let tlForm = $state({ ikuValueID: '', masalah: '', rekomendasi: '', pj_user_id: '', status: 'open' });

  let periodOptions = $derived($periods.map((p) => ({ value: p.id, label: p.kode })));
  let userOptions = $derived(($bootstrap?.users ?? []).map((u) => ({ value: u.id, label: `${u.nama} · ${u.role}` })));
  let ikuOptions = $derived(factors.map((f) => ({ value: f.iku_value.id, label: `${f.iku_value.label_kode} · ${f.iku_value.label_nama ?? ''}` })));
  let canFactor = $derived($currentActor && ['kabid', 'sekretaris', 'admin'].includes($currentActor.role));
  let canCreateTL = $derived($currentActor && ['pimpinan', 'admin'].includes($currentActor.role));
  let canUpdateTL = $derived($currentActor && ['pimpinan', 'kabid', 'sekretaris', 'admin'].includes($currentActor.role));

  let factorColumns = [
    { key: 'label', label: 'IKU' },
    { key: 'value', label: 'Capaian', class: 'text-right w-28' },
    { key: 'num', label: 'Num/Den', class: 'text-right w-28' },
    { key: 'pendukung', label: 'Faktor pendukung' },
    { key: 'penghambat', label: 'Faktor penghambat' },
    { key: 'action', label: 'Aksi', class: 'w-28' }
  ];
  let tlColumns = [
    { key: 'label', label: 'IKU' },
    { key: 'masalah', label: 'Masalah' },
    { key: 'rekomendasi', label: 'Rekomendasi' },
    { key: 'pj', label: 'PJ', class: 'w-40' },
    { key: 'status', label: 'Status', class: 'w-36' },
    { key: 'updated', label: 'Diperbarui', class: 'w-40' }
  ];

  async function load() {
    if (!$selectedPeriod) return;
    loading = true;
    error = '';
    try {
      [factors, tl] = await Promise.all([
        api.factors($selectedPeriod),
        api.tl($selectedPeriod)
      ]);
    } catch (err) {
      error = err instanceof Error ? err.message : 'Gagal memuat faktor/TL.';
    } finally {
      loading = false;
    }
  }

  $effect(() => {
    if ($selectedPeriod) void load();
  });

  function openFactor(row: FactorRow) {
    selectedFactor = row;
    pendukung = row.faktor?.pendukung ?? '';
    penghambat = row.faktor?.penghambat ?? '';
    factorOpen = true;
  }

  async function saveFactor() {
    if (!selectedFactor) return;
    error = '';
    try {
      await api.upsertFactor(selectedFactor.iku_value.id, {
        periode_id: selectedFactor.iku_value.periode_id,
        scope_type: selectedFactor.iku_value.scope_type,
        scope_id: selectedFactor.iku_value.scope_id,
        pendukung,
        penghambat
      });
      factorOpen = false;
      message = 'Faktor disimpan.';
      await load();
    } catch (err) {
      error = err instanceof Error ? err.message : 'Gagal menyimpan faktor.';
    }
  }

  function openTL() {
    tlForm = { ikuValueID: factors[0]?.iku_value.id ?? '', masalah: '', rekomendasi: '', pj_user_id: '', status: 'open' };
    tlOpen = true;
  }

  async function createTL() {
    error = '';
    try {
      await api.createTL({
        ...tlForm,
        periode_id: $selectedPeriod
      });
      tlOpen = false;
      message = 'Tindak lanjut dibuat.';
      await load();
    } catch (err) {
      error = err instanceof Error ? err.message : 'Gagal membuat TL.';
    }
  }

  async function updateStatus(item: TindakLanjut, status: string) {
    error = '';
    try {
      await api.updateTL(item.id, { status });
      message = `Status TL ${item.label_kode ?? ''} diperbarui.`;
      await load();
    } catch (err) {
      error = err instanceof Error ? err.message : 'Gagal memperbarui TL.';
    }
  }
</script>

<PageHeader title={PAGE_COPY.factors.title} description={PAGE_COPY.factors.description}>
  {#snippet actions()}
    <Button variant="outline" onclick={load}><RefreshCw size={14} /> Muat ulang</Button>
    {#if tab === 'tl'}<Button onclick={openTL} disabled={!canCreateTL}><Plus size={14} /> TL baru</Button>{/if}
  {/snippet}
</PageHeader>

<div class="mb-4 border-2 border-black bg-secondary p-3 font-mono text-[11px] font-bold">
  Faktor = analisis kualitatif (pendukung/penghambat), tidak mengubah angka. TL = rencana aksi masalah → rekomendasi → PJ → status.
</div>

{#if !canFactor}
  <div class="mb-4 border-2 border-black bg-warning p-3 font-mono text-xs font-bold">
    Hanya Kabid/Sekretaris/Admin yang boleh mengisi faktor. Ganti aktor di kanan atas untuk mencoba alur ini.
  </div>
{/if}

{#if error}<div class="mb-4 border-2 border-black bg-destructive p-3 font-mono text-xs font-bold">{error}</div>{/if}
{#if message}<div class="mb-4 border-2 border-black bg-success p-3 font-mono text-xs font-bold">{message}</div>{/if}

<div class="mb-4 flex flex-wrap items-end justify-between gap-3">
  <div class="flex flex-wrap gap-2">
    <Button variant={tab === 'faktor' ? 'default' : 'outline'} onclick={() => (tab = 'faktor')}>
      <MessageSquare size={14} /> Faktor ({factors.length})
    </Button>
    <Button variant={tab === 'tl' ? 'default' : 'outline'} onclick={() => (tab = 'tl')}>
      <CheckSquare size={14} /> Tindak Lanjut ({tl.length})
    </Button>
  </div>
  <div class="min-w-52">
    <Select bind:value={$selectedPeriod} options={periodOptions} placeholder="Pilih periode" />
  </div>
</div>

{#if loading}
  <div class="flex items-center justify-center py-20"><Spinner label="Memuat faktor & TL…" /></div>
{:else if tab === 'faktor'}
  <Card>
    {#if factors.length === 0}
      <div class="p-6 text-center text-sm text-muted-foreground">{EMPTY_STATE.factors}</div>
    {:else}
      <DataTable columns={factorColumns} rows={factors}>
        {#snippet cell(item: FactorRow, key)}
          {#if key === 'label'}
            <div class="font-bold">{item.iku_value.label_kode}</div>
            <div class="text-[10px] text-muted-foreground">{item.iku_value.label_nama} · {item.iku_value.status}</div>
          {:else if key === 'value'}
            <span class="font-bold">{formatPercent(item.iku_value.value)}</span>
          {:else if key === 'num'}
            {item.iku_value.numerator}/{item.iku_value.denominator}
          {:else if key === 'pendukung'}
            <span class="text-[11px]">{item.faktor?.pendukung ?? '—'}</span>
          {:else if key === 'penghambat'}
            <span class="text-[11px]">{item.faktor?.penghambat ?? '—'}</span>
          {:else if key === 'action'}
            <Button size="sm" variant="outline" onclick={() => openFactor(item)} disabled={!canFactor}>
              {item.faktor ? 'Ubah' : 'Isi'}
            </Button>
          {/if}
        {/snippet}
      </DataTable>
    {/if}
  </Card>
{:else}
  <Card>
    {#if tl.length === 0}
      <div class="p-6 text-center text-sm text-muted-foreground">Belum ada tindak lanjut.</div>
    {:else}
      <DataTable columns={tlColumns} rows={tl}>
        {#snippet cell(item: TindakLanjut, key)}
          {#if key === 'label'}
            <span class="font-bold">{item.label_kode ?? '—'}</span>
            <div class="text-[10px] text-muted-foreground">{item.label_nama ?? ''}</div>
          {:else if key === 'masalah'}
            <span class="text-[11px]">{item.masalah}</span>
          {:else if key === 'rekomendasi'}
            <span class="text-[11px]">{item.rekomendasi}</span>
          {:else if key === 'pj'}
            {item.pj_nama ?? item.pj_user_id ?? '—'}
          {:else if key === 'status'}
            <StatusBadge kind="tl" value={item.status} />
            {#if canUpdateTL}
              <div class="mt-1">
                <select class="border-2 border-black bg-card px-1 py-0.5 font-mono text-[10px]" value={item.status} onchange={(event) => updateStatus(item, (event.currentTarget as HTMLSelectElement).value)}>
                  <option value="open">Terbuka</option>
                  <option value="in_progress">Dikerjakan</option>
                  <option value="done">Selesai</option>
                  <option value="cancelled">Dibatalkan</option>
                </select>
              </div>
            {/if}
          {:else if key === 'updated'}
            {formatDateTime(item.updated_at)}
          {/if}
        {/snippet}
      </DataTable>
    {/if}
  </Card>
{/if}

<Modal open={factorOpen} title={`Faktor · ${selectedFactor?.iku_value.label_kode ?? ''}`} description="Analisis kualitatif; tidak mengubah kalkulasi." onclose={() => (factorOpen = false)}>
  <div class="space-y-3">
    <div>
      <label class="nb-label" for="pendukung">Faktor pendukung</label>
      <Textarea id="pendukung" bind:value={pendukung} rows={3} placeholder="Hal yang mendukung capaian…" />
    </div>
    <div>
      <label class="nb-label" for="penghambat">Faktor penghambat</label>
      <Textarea id="penghambat" bind:value={penghambat} rows={3} placeholder="Hal yang menghambat capaian…" />
    </div>
    <div class="flex justify-end gap-2">
      <Button variant="outline" onclick={() => (factorOpen = false)}>Batal</Button>
      <Button variant="success" onclick={saveFactor}>Simpan faktor</Button>
    </div>
  </div>
</Modal>

<Modal open={tlOpen} title="Tindak lanjut baru" description="Masalah → rekomendasi → PJ → status." onclose={() => (tlOpen = false)}>
  <div class="grid gap-3">
    <div><label class="nb-label" for="tl-iku">IKU (opsional)</label><Select id="tl-iku" bind:value={tlForm.ikuValueID} options={ikuOptions} placeholder="Tidak terkait IKU tertentu" /></div>
    <div><label class="nb-label" for="tl-masalah">Masalah</label><Textarea id="tl-masalah" bind:value={tlForm.masalah} rows={2} /></div>
    <div><label class="nb-label" for="tl-rekomendasi">Rekomendasi</label><Textarea id="tl-rekomendasi" bind:value={tlForm.rekomendasi} rows={2} /></div>
    <div class="grid gap-3 sm:grid-cols-2">
      <div><label class="nb-label" for="tl-pj">PJ</label><Select id="tl-pj" bind:value={tlForm.pj_user_id} options={userOptions} placeholder="Pilih PJ" /></div>
      <div><label class="nb-label" for="tl-status">Status</label><Select id="tl-status" bind:value={tlForm.status} options={[{value:'open',label:'Terbuka'},{value:'in_progress',label:'Dikerjakan'},{value:'done',label:'Selesai'},{value:'cancelled',label:'Dibatalkan'}]} /></div>
    </div>
    <div class="flex justify-end gap-2">
      <Button variant="outline" onclick={() => (tlOpen = false)}>Batal</Button>
      <Button variant="success" onclick={createTL}>Simpan TL</Button>
    </div>
  </div>
</Modal>
