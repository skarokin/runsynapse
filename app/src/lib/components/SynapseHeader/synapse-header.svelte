<script lang="ts">
    import { Button } from "$lib/components/ui/button";
    import { Input } from "$lib/components/ui/input";
    import Search from "@lucide/svelte/icons/search";
    import Brain from "@lucide/svelte/icons/brain";
    import ArrowLeft from "@lucide/svelte/icons/arrow-left";

    import { PinnedThoughts } from "$lib/components/PinnedThoughts";

    let {
        aiMode = false,
        exitAIMode = () => {},
        searchQuery = $bindable(""),
        handleSearchKeydown = (event: KeyboardEvent) => {},
        data,
    } = $props();
</script>

<div class="flex items-center justify-between">
    {#if aiMode}
        <!-- ai mode - show back button and AI mode title -->
        <Button variant="ghost" size="sm" onclick={() => exitAIMode()}>
            <ArrowLeft class="w-4 h-4" />
            (ESC)
        </Button>
        <div class="flex items-center gap-2">
            <Brain class="w-5 h-5 text-primary" />
            <span class="font-medium text-primary">AI Search</span>
        </div>
    {:else}
        <!-- thoughts mode - THE SEARCH BAR IS THE ENTRY POINT TO AI MODE -->
        <div class="flex-1 max-w-xs sm:max-w-md">
            <div class="relative">
                <Search
                    class="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-muted-foreground"
                />
                <Input
                    bind:value={searchQuery}
                    placeholder="Search your thoughts..."
                    class="pl-10"
                    onkeydown={(e) => handleSearchKeydown(e)}
                />
            </div>
        </div>
        <PinnedThoughts
            pinnedThoughts={data.pinnedThoughts}
            thoughtSet={data.thoughtSet}
        />
    {/if}
</div>
