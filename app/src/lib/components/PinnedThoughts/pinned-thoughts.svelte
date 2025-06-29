<script lang="ts">
    import { Button } from "$lib/components/ui/button";
    import * as Popover from "$lib/components/ui/popover";
    import * as HoverCard from "$lib/components/ui/hover-card";
    import { ScrollArea } from "$lib/components/ui/scroll-area";
    import Pin from "@lucide/svelte/icons/pin";
    import X from "@lucide/svelte/icons/x";

    import { prettyPrintDate } from '$lib/utils/date';

    let {
        pinnedThoughts,
        thoughtSet
    }: {
        pinnedThoughts: any[],
        thoughtSet: Set<string> | undefined,
    } = $props();

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
        <h1 class="text-sm">Pinned Thoughts</h1>
        <ScrollArea class="h-[400px] w-full">
            <div class="p-2">
                {#if pinnedThoughts.length > 0}
                    {#each pinnedThoughts as thought (thought.id)}
                        <div class="flex items-start justify-between w-full gap-2 p-2">
                            <HoverCard.Root openDelay={100}>
                                <HoverCard.Trigger>
                                    <div class="w-full">
                                        <p class="text-xs leading-tight truncate text-left break-words ellipsis">
                                            {thought.thought}
                                        </p>
                                        <span class="text-[10px] text-muted-foreground mt-1 block">
                                            {prettyPrintDate(thought.created_at)}
                                        </span>
                                    </div>
                                </HoverCard.Trigger>
                                <HoverCard.Content class="w-full" side="left">
                                    <div class="p-2">
                                        <p class="text-sm leading-relaxed whitespace-pre-wrap">
                                            {thought.thought}
                                        </p>
                                        <span class="text-xs text-muted-foreground">
                                            {prettyPrintDate(thought.created_at)}
                                        </span>
                                    </div>
                                </HoverCard.Content>
                            </HoverCard.Root>
                            <Button
                                variant="ghost"
                                size="sm"
                                class="h-6 w-6 p-0 text-red-500 hover:text-red-600"
                                onclick={() => unpinThought(thought.id)}
                                title="Unpin"
                            >
                                <X class="w-3 h-3" />
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