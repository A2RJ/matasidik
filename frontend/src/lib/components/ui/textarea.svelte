<script lang="ts">
  import { cn } from '$lib/utils';
  let {
    value = $bindable(''),
    placeholder = '',
    class: className = '',
    rows = 4,
    name = '',
    id = '',
    required = false
  }: {
    value?: string;
    placeholder?: string;
    class?: string;
    rows?: number;
    name?: string;
    id?: string;
    required?: boolean;
  } = $props();
</script>

<textarea bind:value {placeholder} {rows} {name} {id} {required} class={cn('nb-textarea', className)}></textarea>
