<script lang="ts">
  let { label = 'Memuat…', size = 20 }: { label?: string; size?: number } = $props();
</script>

<span class="inline-flex items-center gap-2 font-mono text-xs text-muted-foreground">
  <span
    class="inline-block animate-spin border-2 border-black border-t-primary"
    style={`width:${size}px;height:${size}px`}
  ></span>
  {label}
</span>
