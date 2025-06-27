<script lang="ts">
    import { Button } from "$lib/components/ui/button";
    import * as Popover from "$lib/components/ui/popover";
    import { ScrollArea } from "$lib/components/ui/scroll-area";
    import Pin from "@lucide/svelte/icons/pin";
    import X from "@lucide/svelte/icons/x";

    import { prettyPrintDate } from '$lib/utils/date';

    let {
        pinnedThoughts,
        thoughtSet,
        contentScrollAreaRef
    }: {
        pinnedThoughts: any[],
        thoughtSet: Set<string> | undefined,
        contentScrollAreaRef: HTMLElement | null
    } = $props();

    function gotoPin(id: string) {
        // if thoughtSet already contains this thought, scroll to it
        if (thoughtSet?.has(id)) {

        } else {

        }
    }

    function unpinThought(id: string) {
        pinnedThoughts = pinnedThoughts.filter(thought => thought.id !== id);
        thoughtSet?.delete(id);

        // api call to unpin
    }
</script>

<Popover.Root>
    <Popover.Trigger>
        <Button variant="outline">
            <Pin class="h-4 w-4" />
        </Button>
    </Popover.Trigger>
    <Popover.Content class="w-[300px]">
        <ScrollArea class="h-[400px] w-full">
            <div class="p-2">
                {#if pinnedThoughts.length > 0}
                    {#each pinnedThoughts as thought (thought.id)}
                        <div class="flex items-start w-full gap-2 p-2">
                            <Button
                                variant="ghost"
                                class="flex-1 text-left p-2 min-w-0"
                                onclick={() => gotoPin(thought.id)}
                            >
                                <div class="w-full">
                                    <p class="text-xs leading-tight truncate text-left break-words ellipsis">
                                        {thought.thought}
                                    </p>
                                    <span class="text-[10px] text-muted-foreground mt-1 block">
                                        {prettyPrintDate(thought.created_at)}
                                    </span>
                                </div>
                            </Button>
                            <Button
                                variant="ghost"
                                class="p-1 h-6 w-6 flex-shrink-0 mt-1"
                                onclick={() => unpinThought(thought.id)}
                            >
                                <X />
                            </Button>
                        </div>
                    {/each}
                {:else}                 
                    <div class="p-4 text-center text-muted-foreground">
                        No pinned thoughts
                    </div>
                {/if}
            </div>
        </ScrollArea>
    </Popover.Content>
</Popover.Root>