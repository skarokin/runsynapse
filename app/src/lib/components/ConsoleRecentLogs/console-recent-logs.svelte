<script lang="ts">
    import { Badge } from "$lib/components/ui/badge";
    import { Button } from "$lib/components/ui/button";
    import * as Card from "$lib/components/ui/card";

    import FileText from "@lucide/svelte/icons/file-text";
    import AlertCircle from "@lucide/svelte/icons/alert-circle";

    let { allLogs, nextPageToken, isLoading, error, loadMore, reset } = $props();

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
</script>

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
                    {isLoading ? "Loading..." : "Load More"}
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
                            <div
                                class="w-1.5 h-1.5 rounded-full flex-shrink-0 {log.severity?.toLowerCase() ===
                                'error'
                                    ? 'bg-red-500'
                                    : log.severity?.toLowerCase() === 'warning'
                                      ? 'bg-yellow-500'
                                      : log.severity?.toLowerCase() === 'info'
                                        ? 'bg-blue-500'
                                        : 'bg-gray-400'}"
                            ></div>

                            <!-- timestamp -->
                            <span
                                class="text-gray-500 font-mono flex-shrink-0 w-28 mr-2"
                            >
                                {new Date(log.timestamp).toLocaleString(
                                    "en-US",
                                    {
                                        month: "numeric",
                                        day: "numeric",
                                        hour: "2-digit",
                                        minute: "2-digit",
                                        second: "2-digit",
                                        hour12: false,
                                    },
                                )}
                            </span>

                            <!-- severity -->
                            <span
                                class="flex-shrink-0 font-medium mr-2 w-12 {getSeverityColor(
                                    log.severity,
                                )}"
                            >
                                {log.severity || "INFO"}
                            </span>

                            <!-- payload type badge -->
                            {#if log.payloadType && log.payloadType !== "textPayload"}
                                <span
                                    class="bg-gray-200 text-gray-600 px-1 py-0 rounded text-xs flex-shrink-0"
                                >
                                    {log.payloadType.substring(0, 5)}
                                </span>
                            {/if}

                            <!-- log message - single line with ellipsis -->
                            <span
                                class="font-mono text-gray-800 flex-1 truncate"
                            >
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
                    {isLoading ? "Loading..." : "Load More Logs"}
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
