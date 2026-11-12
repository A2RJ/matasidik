<script lang="ts">
  import type { ChartFaskes, ChartSeries } from '$lib/types';
  import { formatPercent } from '$lib/utils';

  let {
    faskes = [],
    series = [],
    height = 340
  }: { faskes?: ChartFaskes[]; series?: ChartSeries[]; height?: number } = $props();

  const colors = ['#FFDB33', '#90A8FF', '#FF90E8', '#B8FF9F', '#FFB86B', '#FF7A7A', '#C4B5FD', '#7DD3FC'];
  const width = 980;
  const margin = { top: 18, right: 16, bottom: 66, left: 46 };
  const plotW = width - margin.left - margin.right;
  const maxValue = 100;
  let plotH = $derived(height - margin.top - margin.bottom);
  let groupW = $derived(faskes.length > 0 ? plotW / faskes.length : plotW);
  let barW = $derived(series.length > 0 ? Math.min(22, (groupW * 0.72) / series.length) : 12);
  const gap = 3;

  let hasData = $derived(series.some((s) => s.values.some((v) => v !== null && v !== undefined)));
</script>

<div class="w-full">
  <div class="mb-2 flex flex-wrap gap-2">
    {#each series as s, i (s.kode)}
      <span class="inline-flex items-center gap-1 border-2 border-black bg-card px-2 py-0.5 font-mono text-[9px] font-bold">
        <span class="inline-block h-2.5 w-2.5 border border-black" style={`background:${colors[i % colors.length]}`}></span>
        {s.kode} · {s.nama}
      </span>
    {/each}
  </div>

  {#if faskes.length === 0 || series.length === 0 || !hasData}
    <div class="border-2 border-dashed border-black/40 bg-muted p-8 text-center font-mono text-xs text-muted-foreground">
      Belum ada nilai IKU per puskesmas pada periode ini.
    </div>
  {:else}
    <svg viewBox={`0 0 ${width} ${height}`} class="w-full border-2 border-black bg-white" role="img" aria-label="Grafik performa IKU antar puskesmas">
      {#each [0, 25, 50, 75, 100] as tick (tick)}
        {@const y = margin.top + plotH - (tick / maxValue) * plotH}
        <line x1={margin.left} y1={y} x2={width - margin.right} y2={y} stroke="#111" stroke-width={tick === 0 ? 2 : 0.6} stroke-dasharray={tick === 0 ? '' : '4 4'} opacity={tick === 0 ? 1 : 0.35} />
        <text x={margin.left - 8} y={y + 4} text-anchor="end" font-family="JetBrains Mono, monospace" font-size="9" fill="#5c5545">{tick}%</text>
      {/each}

      {#each faskes as fk, fi (fk.id)}
        {@const groupX = margin.left + fi * groupW}
        {@const barsW = series.length * barW + (series.length - 1) * gap}
        {@const startX = groupX + (groupW - barsW) / 2}
        {#each series as s, si (s.kode)}
          {@const v = s.values[fi]}
          {@const h = v === null || v === undefined ? 0 : Math.max(0, Math.min(maxValue, v)) / maxValue * plotH}
          {@const x = startX + si * (barW + gap)}
          <g>
            <rect
              x={x}
              y={margin.top + plotH - h}
              width={barW}
              height={h}
              fill={v === null || v === undefined ? '#F1E9D2' : colors[si % colors.length]}
              stroke="#111"
              stroke-width="1.2"
            >
              <title>{fk.nama} · {s.kode} · {v === null || v === undefined ? 'Belum ada data' : formatPercent(v)}</title>
            </rect>
            {#if v !== null && v !== undefined && barW >= 16}
              <text x={x + barW / 2} y={margin.top + plotH - h - 4} text-anchor="middle" font-family="JetBrains Mono, monospace" font-size="7.5" font-weight="700" fill="#111">{Math.round(v)}</text>
            {/if}
          </g>
        {/each}
        <text x={groupX + groupW / 2} y={height - margin.bottom + 18} text-anchor="middle" font-family="JetBrains Mono, monospace" font-size="8.5" fill="#111">{fk.nama}</text>
        <text x={groupX + groupW / 2} y={height - margin.bottom + 31} text-anchor="middle" font-family="JetBrains Mono, monospace" font-size="7.5" fill="#5c5545">{fk.kode}</text>
      {/each}
    </svg>
  {/if}
</div>
