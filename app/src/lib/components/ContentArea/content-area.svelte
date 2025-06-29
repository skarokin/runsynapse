<script lang="ts">
    import { ScrollArea } from '$lib/components/ui/scroll-area';
    import { Badge } from '$lib/components/ui/badge';
    import { Button } from '$lib/components/ui/button';
    import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
    import Brain from "@lucide/svelte/icons/brain";
    import Trash from "@lucide/svelte/icons/trash";
    import Pin from "@lucide/svelte/icons/pin";
    import MoreHorizontal from "@lucide/svelte/icons/more-horizontal";
    import { toast } from 'svelte-sonner';

    import { prettyPrintDate } from '$lib/utils/date';

    let {
        aiMode,
        aiSummary,
        retrievedThoughts,
        isLoading,
        thoughts,
        scrollAreaRef = $bindable(),
        hasMoreAbove = $bindable(false),
    } = $props();

    let isLoadingMore = $state(false);  // universal for above or below loading
    let topSentinel = $state<HTMLElement | null>(null);    // tiny element for intersection observer to detect
    let observer: IntersectionObserver | null = null;

    function togglePin(thoughtID: string) {
        const thought = thoughts.find((t: any) => t.id === thoughtID);
        if (thought) {
            thought.pinned = !thought.pinned;
        }
    }

    async function loadMoreThoughts() {
        if (isLoadingMore || aiMode) {
            return;
        }

        isLoadingMore = true;

        try {
            const res = await fetch('/synapse/api/loadMoreThoughts', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    lastThoughtID: thoughts.length > 0 ? thoughts[0].id : null,
                }),
            });

            if (!res.ok) {
                toast.error('Failed to load more thoughts');
                return;
            }

            const result = await res.json();
            if (result.error) {
                toast.error('Failed to load more thoughts', result.error);
                return;
            }

            thoughts = [...result.thoughts, ...thoughts];
            hasMoreAbove = result.hasMoreAbove || false;
        } catch (error) {
            console.error('Error loading more thoughts:', error);
            toast.error('An error occurred while loading more thoughts');
        } finally {
            isLoadingMore = false;
        }

    }

    async function deleteThought(thoughtID: string) {
        const thought = thoughts.find((t: any) => t.id === thoughtID);
        if (!thought) {
            return;
        }

        const res = await fetch('/synapse/api/deleteThought', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                thought_id: thoughtID,
            }),
        });

        if (!res.ok) {
            toast.error('Failed to delete thought');
            return;
        }

        const result = await res.json();
        if (result.error) {
            toast.error('Failed to delete thought', result.error);
            return;
        }

        thoughts = thoughts.filter((t: any) => t.id !== thoughtID);
    }

    $effect(() => {
        if (topSentinel && scrollAreaRef && hasMoreAbove && !aiMode) {
            observer = new IntersectionObserver(
                (entries) => {
                    if (entries[0].isIntersecting && hasMoreAbove && !isLoadingMore) {
                        loadMoreThoughts();
                    }
                },
                {
                    root: scrollAreaRef,
                    rootMargin: '0px',
                    threshold: 0.1
                }
            );
            observer.observe(topSentinel);
        }
        return () => {
            observer?.disconnect();
        };
    });
</script>

<ScrollArea
    bind:ref={scrollAreaRef}
    class="min-h-0 flex-1"
>
    {#if aiMode}
        <div class="space-y-6">
            {#if aiSummary}
                <!-- 1. ai summary -->
                <div class="bg-muted/50 rounded-lg p-4">
                    <div class="flex items-center gap-2 mb-3">
                        <Brain class="w-4 h-4 text-primary" />
                        <span class="text-sm font-medium text-primary">AI Response</span>
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
        <!-- tiny div for intersection observer to detect scroll -->
        <div bind:this={topSentinel} class="h-1"></div>
        {#if isLoadingMore}
            <div class="text-center py-2 text-muted-foreground text-xs">
                Loading more thoughts...
            </div>
        {/if}
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
                            <DropdownMenu.Root>
                                <DropdownMenu.Trigger>
                                    <Button
                                        variant="ghost"
                                        size="sm"
                                        class="opacity-0 group-hover:opacity-100 h-6 w-6 p-0"
                                    >
                                        <MoreHorizontal class="w-3.5 h-3.5" />
                                    </Button>
                                </DropdownMenu.Trigger>
                                <DropdownMenu.Content align="end">
                                    <DropdownMenu.Item onclick={() => togglePin(thought.id)}>
                                        <Pin class="w-3.5 h-3.5 mr-2" />
                                        {thought.pinned ? 'Unpin' : 'Pin'}
                                    </DropdownMenu.Item>
                                    <DropdownMenu.Item 
                                        onclick={() => deleteThought(thought.id)}
                                        class="text-red-600 focus:text-red-600"
                                    >
                                        <Trash class="w-3.5 h-3.5 mr-2" />
                                        Delete
                                    </DropdownMenu.Item>
                                </DropdownMenu.Content>
                            </DropdownMenu.Root>
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
