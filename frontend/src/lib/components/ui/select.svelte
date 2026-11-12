<script lang="ts">
  import { cn } from '$lib/utils';
  let {
    value = $bindable(''),
    options = [],
    placeholder = 'Pilih…',
    class: className = '',
    name = '',
    id = '',
    disabled = false,
    allowEmpty = false
  }: {
    value?: string;
    options?: { value: string; label: string }[];
    placeholder?: string;
    class?: string;
    name?: string;
    id?: string;
    disabled?: boolean;
    allowEmpty?: boolean;
  } = $props();
</script>

<div class="relative">
  <select bind:value {name} {id} {disabled} class={cn('nb-select', className)}>
    {#if placeholder}<option value="" disabled={!allowEmpty}>{placeholder}</option>{/if}
    {#each options as option (option.value)}
      <option value={option.value}>{option.label}</option>
    {/each}
  </select>
  <span class="pointer-events-none absolute right-3 top-1/2 -translate-y-1/2 font-mono text-xs">▾</span>
</div>
