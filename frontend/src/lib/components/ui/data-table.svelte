<script lang="ts">
  import type { Snippet } from 'svelte';
  import { cn } from '$lib/utils';

  let {
    columns,
    rows,
    cell,
    rowKey = (row: any, index: number) => String(row?.id ?? index),
    empty = 'Tidak ada data.',
    rowClass = () => '',
    onrowclick,
    fixed = false,
    nowrap = false,
    minWidth = ''
  }: {
    columns: { key: string; label: string; class?: string; width?: string }[];
    rows: any[];
    cell: Snippet<[any, string]>;
    rowKey?: (row: any, index: number) => string;
    empty?: string;
    rowClass?: (row: any) => string;
    onrowclick?: (row: any) => void;
    fixed?: boolean;
    nowrap?: boolean;
    minWidth?: string;
  } = $props();
</script>

<div class="overflow-x-auto nb-scroll">
  <table
    class={cn('nb-table', fixed && 'table-fixed', nowrap && 'table-nowrap')}
    style={minWidth ? `min-width:${minWidth}` : undefined}
  >
    <colgroup>
      {#each columns as column (column.key)}
        <col style={column.width ? `width:${column.width}` : undefined} />
      {/each}
    </colgroup>
    <thead>
      <tr>
        {#each columns as column (column.key)}
          <th class={column.class}>{column.label}</th>
        {/each}
      </tr>
    </thead>
    <tbody>
      {#if rows.length === 0}
        <tr>
          <td colspan={columns.length} class="py-10 text-center text-muted-foreground">{empty}</td>
        </tr>
      {:else}
        {#each rows as row, index (rowKey(row, index))}
          <tr
            class={cn(onrowclick ? 'cursor-pointer' : '', rowClass(row))}
            onclick={() => onrowclick?.(row)}
            onkeydown={(event) => {
              if (event.key === 'Enter') onrowclick?.(row);
            }}
            role={onrowclick ? 'button' : undefined}
            tabindex={onrowclick ? 0 : undefined}
          >
            {#each columns as column (column.key)}
              <td class={column.class}>{@render cell(row, column.key)}</td>
            {/each}
          </tr>
        {/each}
      {/if}
    </tbody>
  </table>
</div>
