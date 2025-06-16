<script lang="ts">
    import { Button } from "$lib/components/ui/button/index.js";
    import * as Card from "$lib/components/ui/card/index.js";
    import Cloud from "@lucide/svelte/icons/cloud";
    import { toast } from "svelte-sonner";

    import { ConsoleHeader } from "$lib/components/ConsoleHeader";
    import { ConsoleServiceOverview } from "$lib/components/ConsoleServiceOverview";
    import { ConsoleRepositoryOverview } from "$lib/components/ConsoleRepositoryOverview";
    import { ConsoleRecentLogs } from "$lib/components/ConsoleRecentLogs";

    import { enhance } from "$app/forms";

    let { data } = $props();
    const { repository, logs, metrics, maxAgeDays, error } = $derived(data);

    let allLogs = $derived([...logs]); // start with initial logs
    let nextPageToken = $state(data.nextPageToken);
    let isLoading = $state(false);

    let loadMoreForm: HTMLFormElement | null = null;

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
    <ConsoleHeader {repository} />

    <!-- Page Title -->
    <div>
        <h2 class="text-2xl font-bold">{repository.repo_name}</h2>
        <p class="text-sm text-muted-foreground">
            Monitor your application's performance, logs, and deployment status.
        </p>
    </div>

    <!-- Repository Overview -->
    <ConsoleRepositoryOverview
        repository={repository}
        metrics={metrics}
    />

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
            <ConsoleServiceOverview {metrics} />
        {/if}

        <!-- Recent Logs -->
        <ConsoleRecentLogs
            allLogs={allLogs}
            nextPageToken={nextPageToken}
            isLoading={isLoading}
            error={error}
            loadMore={loadMore}
            reset={reset}
        />

    {/if}
</div>
