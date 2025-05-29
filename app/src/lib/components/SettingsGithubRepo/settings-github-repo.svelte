<script lang="ts">
    import LockIcon from "@lucide/svelte/icons/lock";
    import UnlockIcon from "@lucide/svelte/icons/unlock";
    import FileCodeIcon from "@lucide/svelte/icons/file-code";
    import CalendarIcon from "@lucide/svelte/icons/calendar";
    import ExternalLinkIcon from "@lucide/svelte/icons/external-link";

    import { Button } from "$lib/components/ui/button/index.js";
    import { Badge } from "$lib/components/ui/badge/index.js";

    let {
        repo,
    } = $props();

    function formatDate(dateString: string): string {
        const date = new Date(dateString);
        return date.toLocaleDateString("en-US", {
            year: "numeric",
            month: "2-digit",
            day: "2-digit",
        });
    }
</script>

<div
    class="flex items-center justify-between p-3 border rounded-lg hover:bg-muted/50 transition-colors"
>
    <div class="flex items-center space-x-3 flex-1 min-w-0">
        <div class="flex-shrink-0">
            {#if repo.is_private}
                <LockIcon class="w-4 h-4 text-muted-foreground" />
            {:else}
                <UnlockIcon class="w-4 h-4 text-green-600" />
            {/if}
        </div>

        <div class="flex-1">
            <div class="flex items-center space-x-2">
                <p class="font-medium text-sm truncate">
                    {repo.repo_name}
                </p>
                <Badge
                    variant={repo.is_private ? "secondary" : "outline"}
                    class="text-xs"
                >
                    {repo.is_private ? "Private" : "Public"}
                </Badge>
            </div>

            <div class="flex items-center space-x-4 mt-1">
                <div
                    class="flex items-center space-x-1 text-xs text-muted-foreground"
                >
                    <FileCodeIcon class="w-3 h-3" />
                    <span class="font-mono">{repo.dockerfile_path}</span>
                </div>

                <div
                    class="flex items-center space-x-1 text-xs text-muted-foreground"
                >
                    <CalendarIcon class="w-3 h-3" />
                    <span>{formatDate(repo.created_at)}</span>
                </div>
            </div>
        </div>
    </div>

    <button class="flex items-center gap-2 flex-shrink-0" onclick={(e) => e.stopPropagation()}>
        <Button
            variant="ghost"
            size="sm"
            href="https://github.com/{repo.repo_name}"
            target="_blank"
            class="h-8 w-8 p-0"
        >
            <ExternalLinkIcon class="w-3 h-3" />
        </Button>
    </button>
</div>
