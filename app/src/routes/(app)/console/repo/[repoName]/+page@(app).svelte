<script lang="ts">
    import { Button } from "$lib/components/ui/button/index.js";
    import * as Card from "$lib/components/ui/card/index.js";
    import { Badge } from "$lib/components/ui/badge/index.js";
    import { Separator } from "$lib/components/ui/separator/index.js";
    import GitBranch from "@lucide/svelte/icons/git-branch";
    import Cloud from "@lucide/svelte/icons/cloud";
    import ExternalLink from "@lucide/svelte/icons/external-link";
    import Settings from "@lucide/svelte/icons/settings";
    import ArrowLeft from "@lucide/svelte/icons/arrow-left";
    import Activity from "@lucide/svelte/icons/activity";
    import FileText from "@lucide/svelte/icons/file-text";
    import AlertCircle from "@lucide/svelte/icons/alert-circle";
    import CheckCircle from "@lucide/svelte/icons/check-circle";
    import { toast } from "svelte-sonner";

    import { enhance } from "$app/forms";

    let { data } = $props();
    const { repository, logs, metrics, maxAgeDays, error } = $derived(data);

    let allLogs = $derived([...logs]); // start with initial logs
    let nextPageToken = $state(data.nextPageToken);
    let isLoading = $state(false);

    let loadMoreForm: HTMLFormElement | null = null;

    function getConnectionStatus() {
        return repository.gcp_project_id ? "connected" : "not connected";
    }

    function getConnectionBadgeVariant(): "default" | "outline" {
        return repository.gcp_project_id ? "default" : "outline";
    }

    function getSeverityColor(severity: string) {
        switch (severity?.toLowerCase()) {
            case "error":
                return "text-red-600";
            case "warning":
                return "text-yellow-600";
            case "info":
                return "text-blue-600";
            default:
                return "text-gray-600";
        }
    }

    $effect(() => {
        if (error) {
            toast.error(error);
        }
    });

    function handleLoadMoreResult(result: any) {
        const success = result.result.type === "success";
        const logs = result.result.data?.logs || [];
        const nextPage = result.result.data?.nextPageToken || null;

        if (success && result.result.data) {
            if (logs.length > 0) {
                allLogs = [...allLogs, ...logs];
                nextPageToken = nextPage;
            } else {
                toast.info("No more logs to load");
            }
        } else {
            toast.error("Failed to load more logs");
        }

        isLoading = false;
    }

    function loadMore() {
        if (!nextPageToken || isLoading) return;
        isLoading = true;
        loadMoreForm?.requestSubmit();
    }

    function reset() {
        allLogs = [...logs];
        nextPageToken = data.nextPageToken;
    }
</script>

<!-- hidden form using enhance to easily deserialize data -->
<form
    bind:this={loadMoreForm}
    method="POST"
    action="?/loadMore"
    use:enhance={() => {
        return async ({ result }) => {
            handleLoadMoreResult({ result });
        };
    }}
    class="hidden"
>
    <input type="hidden" name="page" value={nextPageToken || ""} />
    <input type="hidden" name="maxAgeDays" value={maxAgeDays || ""} />
</form>

<div class="space-y-6 w-full">
    <!-- Header -->
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
            >
                <ExternalLink class="w-4 h-4 mr-2" />
                View on GitHub
            </Button>
        </div>
    </div>

    <!-- Page Title -->
    <div>
        <h2 class="text-2xl font-bold">{repository.repo_name}</h2>
        <p class="text-sm text-muted-foreground">
            Monitor your application's performance, logs, and deployment status.
        </p>
    </div>

    <!-- Repository Overview -->
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
                            {repository.description ||
                                "No description available"}
                        </p>
                        <div class="flex items-center gap-2 mt-1">
                            {#if repository.is_private}
                                <Badge variant="outline" class="text-xs"
                                    >Private</Badge
                                >
                            {/if}
                            <Badge variant="outline" class="text-xs">main</Badge
                            >
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
                                    <CheckCircle
                                        class="w-4 h-4 text-green-500"
                                    />
                                    <span class="text-green-600 font-medium"
                                        >Live</span
                                    >
                                {:else}
                                    <AlertCircle
                                        class="w-4 h-4 text-yellow-500"
                                    />
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

    {#if !repository.gcp_project_id}
        <!-- Not Connected State -->
        <Card.Root class="border-dashed">
            <Card.Content class="py-12">
                <div class="text-center space-y-4">
                    <div
                        class="w-12 h-12 mx-auto bg-gray-100 rounded-full flex items-center justify-center"
                    >
                        <Cloud class="w-6 h-6 text-gray-400" />
                    </div>
                    <div>
                        <h3 class="text-lg font-semibold">
                            Repository Not Connected
                        </h3>
                        <p class="text-sm text-muted-foreground mt-1">
                            Connect this repository to a GCP project to view
                            logs and metrics.
                        </p>
                    </div>
                    <Button href="/console">Go Back to Connect</Button>
                </div>
            </Card.Content>
        </Card.Root>
    {:else}
        <!-- Service Overview -->
        {#if metrics}
            <div class="space-y-4">
                <p class="text-xl font-semibold">Service Overview</p>
                <ul class="flex flex-col gap-4">
                    <li class="flex items-center justify-between">
                        <div class="flex items-center gap-4">
                            <Activity size="20" class="text-muted-foreground" />
                            <div>
                                <h3 class="font-medium">Service URL</h3>
                                <p class="text-sm text-muted-foreground">
                                    Public endpoint for your application
                                </p>
                            </div>
                        </div>
                        {#if metrics.uri}
                            <Button
                                variant="outline"
                                size="sm"
                                href={metrics.uri}
                                target="_blank"
                            >
                                <ExternalLink class="w-3 h-3 mr-1" />
                                Visit
                            </Button>
                        {:else}
                            <Badge variant="outline">No URI</Badge>
                        {/if}
                    </li>
                    <Separator />
                    <li class="flex items-center justify-between">
                        <div class="flex items-center gap-4">
                            <Settings size="20" class="text-muted-foreground" />
                            <div>
                                <h3 class="font-medium">Generation</h3>
                                <p class="text-sm text-muted-foreground">
                                    Current deployment version
                                </p>
                            </div>
                        </div>
                        <Badge variant="outline"
                            >#{metrics.generation || "Unknown"}</Badge
                        >
                    </li>
                </ul>
            </div>
        {/if}

        <!-- Recent Logs -->
        <div class="space-y-4">
            <div class="flex items-center justify-between">
                <div class="flex items-center gap-3">
                    <p class="text-xl font-semibold">Recent Logs</p>
                    <Badge variant="outline" class="text-xs">
                        {allLogs.length} entries
                    </Badge>
                </div>
                
                <div class="flex items-center gap-2">
                    <Button
                        variant="ghost"
                        size="sm"
                        onclick={reset}
                        disabled={isLoading}
                    >
                        Reset
                    </Button>
                    {#if nextPageToken}
                        <Button
                            variant="outline"
                            size="sm"
                            onclick={loadMore}
                            disabled={isLoading}
                        >
                            {isLoading ? 'Loading...' : 'Load More'}
                        </Button>
                    {/if}
                </div>
            </div>

            {#if error}
                <Card.Root class="border-yellow-200 bg-yellow-50">
                    <Card.Content class="p-4">
                        <div class="flex items-center gap-2">
                            <AlertCircle class="w-4 h-4 text-yellow-600" />
                            <span class="text-yellow-800 font-medium">Warning</span>
                        </div>
                        <p class="text-sm text-yellow-700 mt-1">{error}</p>
                    </Card.Content>
                </Card.Root>
            {/if}

            {#if allLogs.length > 0}
               <Card.Root>
                    <Card.Content class="p-0">
                        <div class="divide-y divide-gray-100 px-3 py-2">
                            {#each allLogs as log, i}
                                <div class="flex items-center py-1 gap-2 text-xs">
                                    <!-- tiny severity dot -->
                                    <div class="w-1.5 h-1.5 rounded-full flex-shrink-0 {
                                        log.severity?.toLowerCase() === 'error' ? 'bg-red-500' :
                                        log.severity?.toLowerCase() === 'warning' ? 'bg-yellow-500' :
                                        log.severity?.toLowerCase() === 'info' ? 'bg-blue-500' :
                                        'bg-gray-400'
                                    }"></div>
                                    
                                    <!-- timestamp -->
                                    <span class="text-gray-500 font-mono flex-shrink-0 w-28 mr-2">
                                        {new Date(log.timestamp).toLocaleString('en-US', {
                                            month: 'numeric',
                                            day: 'numeric', 
                                            hour: '2-digit',
                                            minute: '2-digit',
                                            second: '2-digit',
                                            hour12: false
                                        })}
                                    </span>
                                    
                                    <!-- severity -->
                                    <span class="flex-shrink-0 font-medium mr-2 w-12 {getSeverityColor(log.severity)}">
                                        {(log.severity || "INFO")}
                                    </span>
                                    
                                    <!-- payload type badge -->
                                    {#if log.payloadType && log.payloadType !== 'textPayload'}
                                        <span class="bg-gray-200 text-gray-600 px-1 py-0 rounded text-xs flex-shrink-0">
                                            {log.payloadType.substring(0, 5)}
                                        </span>
                                    {/if}
                                    
                                    <!-- log message - single line with ellipsis -->
                                    <span class="font-mono text-gray-800 flex-1 truncate">
                                        {log.data}
                                    </span>
                                </div>
                            {/each}
                        </div>
                    </Card.Content>
                </Card.Root>
                
                <!-- load more also at bottom -->
                {#if nextPageToken}
                    <div class="flex justify-center pt-4">
                        <Button
                            variant="outline"
                            onclick={loadMore}
                            disabled={isLoading}
                            class="w-full max-w-xs"
                        >
                            {isLoading ? 'Loading...' : 'Load More Logs'}
                        </Button>
                    </div>
                {/if}
            {:else if !error}
                <Card.Root class="border-dashed">
                    <Card.Content class="py-8">
                        <div class="text-center">
                            <FileText class="w-8 h-8 mx-auto text-gray-400 mb-2" />
                            <p class="text-sm text-muted-foreground">
                                No recent logs found
                            </p>
                        </div>
                    </Card.Content>
                </Card.Root>
            {/if}
        </div>
    {/if}
</div>
