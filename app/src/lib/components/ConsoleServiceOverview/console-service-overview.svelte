<script lang="ts">
    import { Badge } from "$lib/components/ui/badge";
    import { Button } from "$lib/components/ui/button";
    import { Separator } from "$lib/components/ui/separator";
    import Activity from "@lucide/svelte/icons/activity";
    import Settings from "@lucide/svelte/icons/settings";
    import ExternalLink from "@lucide/svelte/icons/external-link";

    let { metrics } = $props();
</script>

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
            <Badge variant="outline">#{metrics.generation || "Unknown"}</Badge>
        </li>
    </ul>
</div>
