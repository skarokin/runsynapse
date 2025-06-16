<script lang="ts">
    import * as Card from "$lib/components/ui/card";
    import { Badge } from "$lib/components/ui/badge";
    import { Separator } from "$lib/components/ui/separator";

    import GitBranch from "@lucide/svelte/icons/git-branch";
    import Cloud from "@lucide/svelte/icons/cloud";
    import Settings from "@lucide/svelte/icons/settings";
    import CheckCircle from "@lucide/svelte/icons/check-circle";
    import AlertCircle from "@lucide/svelte/icons/alert-circle";

    let { repository, metrics } = $props();

    function getConnectionStatus() {
        return repository.gcp_project_id ? "connected" : "not connected";
    }

    function getConnectionBadgeVariant(): "default" | "outline" {
        return repository.gcp_project_id ? "default" : "outline";
    }
</script>

<div class="space-y-4">
    <p class="text-xl font-semibold flex justify-between items-center">
        Service Status
        <Badge variant={getConnectionBadgeVariant()}>
            {getConnectionStatus()}
        </Badge>
    </p>

    <Card.Root>
        <Card.Content class="space-y-4">
            <div class="flex items-center gap-4">
                <div
                    class="w-12 h-12 bg-gray-100 rounded-full flex items-center justify-center"
                >
                    <GitBranch class="w-6 h-6 text-gray-600" />
                </div>
                <div>
                    <h3 class="font-medium">{repository.repo_name}</h3>
                    <p class="text-sm text-muted-foreground">
                        {repository.description || "No description available"}
                    </p>
                    <div class="flex items-center gap-2 mt-1">
                        {#if repository.is_private}
                            <Badge variant="outline" class="text-xs"
                                >Private</Badge
                            >
                        {/if}
                        <Badge variant="outline" class="text-xs">main</Badge>
                    </div>
                </div>
            </div>

            {#if repository.gcp_project_id}
                <Separator />
                <div class="flex items-center gap-4 text-sm">
                    <div class="flex items-center gap-2">
                        <Cloud class="w-4 h-4 text-blue-500" />
                        <span class="font-medium">GCP Project:</span>
                        <span>{repository.gcp_project_id}</span>
                    </div>
                    {#if repository.service_name}
                        <Separator orientation="vertical" class="h-4" />
                        <div class="flex items-center gap-2">
                            <Settings class="w-4 h-4 text-gray-500" />
                            <span class="font-medium">Service:</span>
                            <span>{repository.service_name}</span>
                        </div>
                    {/if}
                    {#if metrics?.uri}
                        <Separator orientation="vertical" class="h-4" />
                        <div class="flex items-center gap-2">
                            {#if metrics.ready}
                                <CheckCircle class="w-4 h-4 text-green-500" />
                                <span class="text-green-600 font-medium"
                                    >Live</span
                                >
                            {:else}
                                <AlertCircle class="w-4 h-4 text-yellow-500" />
                                <span class="text-yellow-600 font-medium"
                                    >Deploying</span
                                >
                            {/if}
                        </div>
                    {/if}
                </div>
            {/if}
        </Card.Content>
    </Card.Root>
</div>
