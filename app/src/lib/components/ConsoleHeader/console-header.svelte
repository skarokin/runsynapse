<script lang="ts">
import { Button } from '$lib/components/ui/button';
    import ExternalLink from "@lucide/svelte/icons/external-link";
    import Settings from "@lucide/svelte/icons/settings";
    import ArrowLeft from "@lucide/svelte/icons/arrow-left";

    let { repository } = $props();
</script>

<div class="flex items-center justify-between">
    <div class="flex items-center gap-4">
        <Button variant="ghost" size="sm" href="/console">
            <ArrowLeft class="w-4 h-4 mr-2" />
            Back to Repositories
        </Button>
    </div>
    <div class="flex gap-2">
        <Button
            variant="outline"
            size="sm"
            href={`/console/repo/${encodeURIComponent(repository.repo_name)}/settings`}
        >
            <Settings class="w-4 h-4 mr-2" />
            Settings
        </Button>
        <Button
            variant="outline"
            href={`https://github.com/${repository.repo_name}`}
            target="_blank"
            size="sm"
        >
            <ExternalLink class="w-4 h-4 mr-2" />
            View on GitHub
        </Button>
    </div>
</div>
