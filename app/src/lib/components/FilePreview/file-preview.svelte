<script lang="ts">
    import Trash from "@lucide/svelte/icons/trash";
    import { Button } from "$lib/components/ui/button/index.js";

    let { pendingFiles = $bindable() }: { pendingFiles: File[] } = $props();
</script>

{#if pendingFiles.length > 0}
    <div class="w-full bg-muted rounded-lg px-3 py-2 mb-2 flex items-center gap-2 overflow-x-auto min-h-[48px]">
        {#each pendingFiles as file (file.name)}
            <div class="flex items-center gap-1 border rounded px-2 py-1 text-xs shadow-sm h-24">
                {#if file.type.startsWith("image/")}
                    <img
                        src={URL.createObjectURL(file)}
                        alt={file.name}
                        class="w-16 h-16 object-cover rounded"
                        onload={() => URL.revokeObjectURL(file.name)}
                    />
                {/if}
                <span class="truncate max-w-[100px]">{file.name}</span>
                <Button
                    variant="ghost"
                    size="icon"
                    class="h-4 w-4 p-0 text-muted-foreground hover:text-red-500"
                    title="Remove"
                    onclick={() => {
                        pendingFiles = pendingFiles.filter((f: File) => f !== file);
                    }}
                >
                    <Trash class="w-3 h-3" />
                </Button>
            </div>
        {/each}
    </div>
{/if}