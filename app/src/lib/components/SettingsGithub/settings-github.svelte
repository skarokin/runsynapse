<script lang="ts">
    import GithubIcon from "@lucide/svelte/icons/github";
    import ExternalLinkIcon from "@lucide/svelte/icons/external-link";
    import ChevronsUpDownIcon from "@lucide/svelte/icons/chevrons-up-down";

    import { Button } from "$lib/components/ui/button/index.js";
    import * as Card from "$lib/components/ui/card/index.js";
    import { Separator } from "$lib/components/ui/separator/index.js";
    import * as Collapsible from "$lib/components/ui/collapsible/index.js";

    import { SettingsGithubRepo } from "$lib/components/SettingsGithubRepo";

    let { isConnected, connectedRepos, installationID } = $props();

    let isExpanded = $state(false);

    const firstRepo = connectedRepos?.[0];
    const restRepos = connectedRepos?.slice(1) || [];
</script>

<Card.Root>
    <Card.Content class="space-y-4">
        {#if isConnected && connectedRepos?.length > 0}
            <Collapsible.Root>
                <div class="mb-4 hover:pointer relative cursor-pointer">
                    <SettingsGithubRepo repo={firstRepo} />
                </div>
                {#if !isExpanded && restRepos.length > 0}
                    <Collapsible.Trigger
                        onclick={() => isExpanded = !isExpanded}
                        class="w-full flex items-center justify-between hover:cursor-pointer"
                    >
                        <p class="text-sm font-medium text-muted-foreground">
                            Show {restRepos.length} more repos
                        </p>
                        <ChevronsUpDownIcon class="w-4 h-4 text-muted-foreground" />
                    </Collapsible.Trigger>
                {/if}
                <Collapsible.Content>
                    <div class="flex flex-col gap-4">
                        {#each restRepos as repo (repo.id)}
                            <SettingsGithubRepo
                                repo={repo} 
                            />
                        {/each}
                    </div>
                </Collapsible.Content>
                {#if isExpanded}
                    <Collapsible.Trigger
                        onclick={() => isExpanded = !isExpanded}
                        class="mt-4 w-full flex items-center justify-between hover:cursor-pointer"
                    >
                        <p class="text-sm font-medium text-muted-foreground">
                            Show less
                        </p>
                        <ChevronsUpDownIcon class="w-4 h-4 text-muted-foreground" />
                    </Collapsible.Trigger>
                {/if}
            </Collapsible.Root>

            <Separator />

            <div class="flex items-center justify-between">
                <div class="flex items-center space-x-2">
                    <p class="text-sm font-medium">
                        {connectedRepos.length} repositor{connectedRepos.length !== 1 ? "ies": "y"} connected
                    </p>
                </div>
                <Button
                    variant="outline"
                    size="sm"
                    href="https://github.com/settings/installations/{installationID}"
                    target="_blank"
                >
                    <ExternalLinkIcon class="w-3 h-3 mr-1" />
                    Manage Access
                </Button>
            </div>
        {:else if isConnected}
            <div class="text-center py-8">
                <GithubIcon
                    class="w-12 h-12 text-muted-foreground mx-auto mb-4"
                />
                <p class="text-sm text-muted-foreground mb-4">
                    No repositories connected yet
                </p>
                <Button
                    variant="outline"
                    href="https://github.com/settings/installations"
                    target="_blank"
                >
                    <ExternalLinkIcon class="w-3 h-3 mr-1" />
                    Add Repositories
                </Button>
            </div>
        {:else}
            <div class="text-center py-8">
                <GithubIcon
                    class="w-12 h-12 text-muted-foreground mx-auto mb-4"
                />
                <p class="text-sm text-muted-foreground mb-4">
                    Connect your GitHub account to start deploying repositories
                </p>
                <Button href="/api/github/setup">
                    <GithubIcon class="w-4 h-4 mr-2" />
                    Connect GitHub
                </Button>
            </div>
        {/if}
    </Card.Content>
</Card.Root>
