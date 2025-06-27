<script lang="ts">
    import { ScrollArea } from '$lib/components/ui/scroll-area';
    import { Badge } from '$lib/components/ui/badge';
    import { Button } from '$lib/components/ui/button';
    import Brain from "@lucide/svelte/icons/brain";
    import Pin from "@lucide/svelte/icons/pin";

    import { prettyPrintDate } from '$lib/utils/date';

    let {
        aiMode,
        aiSummary,
        retrievedThoughts,
        isLoading,
        thoughts,
        scrollAreaRef = $bindable(),
    } = $props();

    function togglePin(thoughtId: string) {
        const thought = thoughts.find((t: any) => t.id === thoughtId);
        if (thought) {
            thought.pinned = !thought.pinned;
        }
    }
</script>

<ScrollArea
    bind:this={scrollAreaRef}
    class="min-h-0 flex-1"
>
    {#if aiMode}
        <div class="space-y-6">
            {#if aiSummary}
                <!-- 1. ai summary -->
                <div class="bg-muted/50 rounded-lg p-4">
                    <div class="flex items-center gap-2 mb-3">
                        <Brain class="w-4 h-4 text-primary" />
                        <span class="text-sm font-medium text-primary"
                            >AI Response</span
                        >
                    </div>
                    <p class="text-sm leading-relaxed whitespace-pre-wrap">
                        {aiSummary}
                    </p>
                </div>

                <!-- retrieved thoughts -->
                {#if retrievedThoughts.length > 0}
                    <div>
                        <h3
                            class="text-xs font-semibold uppercase text-muted-foreground mb-2 px-2"
                        >
                            Sources
                        </h3>
                        <div class="space-y-1">
                            {#each retrievedThoughts as thought}
                                <div
                                    class="group flex items-center gap-2 hover:bg-muted/30 px-2 py-1 rounded transition-colors"
                                >
                                    <div class="flex-1 min-w-0">
                                        <div class="flex items-baseline gap-2">
                                            <span class="text-xs font-medium">
                                                You
                                            </span>
                                            <span class="text-xs text-muted-foreground">
                                                {prettyPrintDate(thought.created_at)}
                                            </span>
                                            {#if thought.pinned}
                                                <Badge
                                                    variant="secondary"
                                                    class="h-4 px-1 text-xs"
                                                >
                                                    <Pin class="w-2 h-2 mr-1" />
                                                    Pinned
                                                </Badge>
                                            {/if}
                                        </div>
                                        <p
                                            class="text-sm leading-relaxed whitespace-pre-wrap"
                                        >
                                            {thought.thought}
                                        </p>
                                    </div>
                                </div>
                            {/each}
                        </div>
                    </div>
                {/if}
            {:else if !isLoading}
                <!-- initial empty state -->
                <div class="text-center text-muted-foreground py-12">
                    <Brain class="w-12 h-12 mx-auto mb-4 opacity-50" />
                    <p class="text-sm">
                        Ask me anything about your thoughts...
                    </p>
                </div>
            {/if}
        </div>
    {:else}
        <!-- thoughts area -->
        <div class="space-y-1">
            {#each thoughts as thought}
                <div
                    class="group flex items-center gap-2 hover:bg-muted/30 px-2 py-1 rounded transition-colors"
                >
                    <div class="flex-1 min-w-0">
                        <div class="flex items-baseline gap-2">
                            <span class="text-xs font-medium">You</span>
                            <span class="text-xs text-muted-foreground">{prettyPrintDate(thought.created_at)}</span>
                            {#if thought.pinned}
                                <Badge
                                    variant="secondary"
                                    class="h-4 px-1 text-xs"
                                >
                                    <Pin class="w-2 h-2 mr-1" />
                                    Pinned
                                </Badge>
                            {/if}
                            <Button
                                variant="ghost"
                                size="sm"
                                class="opacity-0 group-hover:opacity-100 h-5 w-5 p-0 ml-auto"
                                onclick={() => togglePin(thought.id)}
                            >
                                <Pin class="w-3 h-3" />
                            </Button>
                        </div>
                        <p class="text-sm leading-relaxed whitespace-pre-wrap">
                            {thought.thought}
                        </p>
                    </div>
                </div>
            {/each}
        </div>
    {/if}
</ScrollArea>
