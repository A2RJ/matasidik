<script lang="ts">
  import Button from '$lib/components/ui/button/button.svelte';
  import Badge from '$lib/components/ui/badge.svelte';
  import Select from '$lib/components/ui/select.svelte';
  import StatusBadge from '$lib/components/ui/status-badge.svelte';
  import Spinner from '$lib/components/ui/spinner.svelte';
  import { formatNumber, formatPercent } from '$lib/utils';
  import type { Label, SubmissionRow } from '$lib/types';
  import { Save } from '@lucide/svelte';

  let {
    row,
    labels = [],
    value = $bindable(''),
    saving = false,
    onsave
  }: {
    row: SubmissionRow;
    labels?: Label[];
    value?: string;
    saving?: boolean;
    onsave?: () => void;
  } = $props();

  let selectedName = $derived(labels.find((label) => label.kode === value)?.nama ?? '');
</script>

<div class="border-b-2 border-black p-4 last:border-b-0">
  <div class="flex flex-wrap items-start justify-between gap-3">
    <div class="flex items-start gap-3">
      <span class="border-2 border-black bg-primary px-2 py-0.5 font-mono text-xs font-bold shadow-nbsm">R{row.row_no}</span>
      <div>
        <div class="font-semibold">{row.indikator_nama ?? row.kode_indikator_raw ?? '—'}</div>
        <div class="mt-1 flex flex-wrap gap-1">
          <Badge tone="neutral">{row.kode_iku_raw ?? '—'}</Badge>
          <Badge tone="neutral">{row.scope_type}/{row.scope_code ?? '—'}</Badge>
          {#if row.matched_indikator_id}<Badge tone="info">{row.indikator_kode ?? row.matched_indikator_id}</Badge>{/if}
        </div>
      </div>
    </div>
    <div class="text-right">
      <div class="font-mono text-lg font-bold">{formatPercent(row.nilai)}</div>
      <div class="font-mono text-[10px] text-muted-foreground">{formatNumber(row.numerator)}/{formatNumber(row.denominator)}</div>
      <div class="mt-2"><StatusBadge kind="classification" value={row.classification_status} /></div>
    </div>
  </div>

  {#if row.error_message}
    <div class="mt-2 border-2 border-black px-2 py-1 font-mono text-[10px] font-bold" class:bg-warning={row.classification_status === 'MATCHED'} class:bg-destructive={row.classification_status !== 'MATCHED'}>
      {row.error_message}
    </div>
  {/if}

  <div class="mt-3 grid items-end gap-3 md:grid-cols-[1fr_280px_110px]">
    <div>
      <div class="nb-label">IKU leaf terpilih</div>
      <Select bind:value options={labels.map((label) => ({ value: label.kode, label: label.kode }))} placeholder="Pilih IKU leaf" />
      <div class="mt-1 min-h-4 font-mono text-[10px] text-muted-foreground">
        {selectedName || 'Pilih IKU leaf yang valid untuk indikator ini'}
      </div>
    </div>
    <div class="hidden md:block">
      <div class="nb-label">Klasifikasi</div>
      <div class="border-2 border-black bg-muted px-2 py-2 font-mono text-[10px]">
        {row.matched_label_id ? row.iku_kode ?? row.matched_label_id : 'Belum ada IKU'}
      </div>
    </div>
    <Button variant="outline" disabled={saving} onclick={onsave}>
      {#if saving}<Spinner label="…" size={14} />{:else}<Save size={14} /> Simpan{/if}
    </Button>
  </div>
</div>
