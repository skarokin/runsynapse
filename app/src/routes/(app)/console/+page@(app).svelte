<script lang="ts">
    import { Button } from "$lib/components/ui/button/index.js";
    import * as Card from "$lib/components/ui/card/index.js";
    import { Badge } from "$lib/components/ui/badge/index.js";
    import { Input } from "$lib/components/ui/input/index.js";
    import Search from "@lucide/svelte/icons/search";
    import GitBranch from "@lucide/svelte/icons/git-branch";
    import Cloud from "@lucide/svelte/icons/cloud";
    import ExternalLink from "@lucide/svelte/icons/external-link";

    import { ConnectToGCP } from "$lib/components/ConnectToGCP/index.js";

    let { data } = $props();
    const {
        error,
        projects,
        connectedRepos,
        installationID
    } = data;

    let searchQuery = $state("");
    const filteredRepos = $derived(
        connectedRepos.filter(repo => 
            repo.repo_name.toLowerCase().includes(searchQuery.toLowerCase()) ||
            repo.description?.toLowerCase().includes(searchQuery.toLowerCase())
        )
    );

    function getStatusColor(status: string): string {
        switch(status) {
            case 'deployed': return 'bg-green-500';
            case 'building': return 'bg-yellow-500';
            case 'failed': return 'bg-red-500';
            default: return 'bg-gray-500';
        }
    }

    function getStatusTextColor(status: string): string {
        switch(status) {
            case 'deployed': return 'green-600';
            case 'building': return 'yellow-600';
            case 'failed': return 'red-600';
            default: return 'gray-600';
        }
    }

    function getStatusVariant(status: string): 'outline' | 'default' | 'secondary' | 'destructive' {
        switch(status) {
            case 'deployed': return 'default';
            case 'building': return 'secondary';
            case 'failed': return 'destructive';
            default: return 'outline';
        }
    }
</script>

<div class="space-y-6 w-full">
    <div class="flex items-center justify-between">
        <div>
            <h2 class="text-2xl font-bold">Repositories</h2>
            <p class="text-sm text-muted-foreground">
                Manage your connected repositories and Cloud Run deployments.
            </p>
        </div>
        <Button href="https://github.com/settings/installations/{installationID}">
            Connect Repository
        </Button>
    </div>

    <div class="relative">
        <Search class="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground h-4 w-4" />
        <Input
            placeholder="Search repositories and projects..."
            bind:value={searchQuery}
            class="pl-10"
        />
    </div>

    {#if error}
        <div class="border border-red-200 bg-red-50 rounded-lg p-4">
            <h3 class="font-medium text-red-800">Oops!</h3>
            <p class="text-sm text-red-600">There was a problem loading your repositories: {error}</p>
        </div>
    {/if}

    <div class="grid gap-4">
        {#each filteredRepos as repo (repo.id)}
            <Card.Root class="hover:shadow-md transition-shadow">
                <Card.Content class="p-6">
                    <div class="flex items-start justify-between">
                        <!-- repo info on left side -->
                        <div class="flex items-start gap-4 flex-1">
                            <!-- status ... FUTURE - PULL FROM REPOSITORY_DEPLOYMENTS TABLE FOR STATUS UPDATES, MOST RECENT DEPLOYMENT, ETC -->
                            <div class="flex items-center gap-2 mt-1">
                                <div class="w-2 h-2 rounded-full {repo.gcp_project_id ? getStatusColor('deployed') : 'bg-gray-400'}"></div>
                            </div>
                            
                            <!-- repo details -->
                            <div class="flex-1 space-y-2">
                                <div class="flex items-center gap-3">
                                    <h3 class="text-lg font-semibold">{repo.repo_name}</h3>
                                    <Badge variant="outline" class="text-xs">
                                        <GitBranch class="w-3 h-3 mr-1" />
                                        main
                                    </Badge>
                                    {#if repo.gcp_project_id}
                                        <Badge variant="default" class="text-xs">
                                            connected
                                        </Badge>
                                    {:else}
                                        <Badge variant="outline" class="text-xs">
                                            not connected
                                        </Badge>
                                    {/if}
                                </div>
                                
                                <p class="text-sm text-muted-foreground">
                                    {repo.description || 'No description available'}
                                </p>
                                
                                <div class="flex items-center gap-4 text-xs text-muted-foreground">
                                    <span>{repo.dockerfile_path}</span>
                                    {#if repo.is_private}
                                        <span>•</span>
                                        <Badge variant="outline" class="text-xs">Private</Badge>
                                    {/if}
                                </div>
                            </div>
                        </div>

                        <!-- right side actions -->
                        {#if repo.gcp_project_id}
                            <div class="flex flex-col items-end gap-4">                                      
                                <div class="flex gap-2">
                                    <Button variant="outline" size="sm" href={`/console/repo/${encodeURIComponent(repo.repo_name)}`}>
                                        Configure
                                    </Button>
                                    <Button variant="outline" size="sm" href={`https://github.com/${repo.repo_name}`} target="_blank">
                                        <ExternalLink class="w-3 h-3 mr-1" />
                                        Repo
                                    </Button>
                                </div>
                                <div class="text-right space-y-1">
                                    <div class="flex items-center justify-end gap-2 text-sm">
                                        <Cloud class="w-4 h-4 text-blue-500" />
                                        <span class="font-medium">{repo.gcp_project_id}</span>
                                    </div>
                                    <div class="text-xs text-muted-foreground">
                                        {repo.service_name || 'No service configured'}
                                    </div>
                                </div>
                            </div>
                        {:else}
                            <!-- connect to gcp project -->
                            <div class="flex flex-col items-end gap-4 w-1/6">
                                <!-- uses internal UUID instead of actual repository ID from GitHub -->
                                <ConnectToGCP {projects} {repo} />

                                <div class="text-right space-y-1">
                                    <div class="flex items-center justify-end gap-2 text-sm text-muted-foreground">
                                        <span>No GCP project connected</span>
                                    </div>
                                    <div class="text-xs text-muted-foreground text-pretty">
                                        Connect to enable deployments and configure a Cloud Run service.
                                    </div>
                                </div>
                            </div>
                        {/if}
                    </div>
                </Card.Content>
            </Card.Root>
        {/each}
        
        {#if filteredRepos.length === 0}
            <div class="text-center py-12">
                <div class="text-muted-foreground">
                    {#if searchQuery}
                        No repositories found matching "{searchQuery}"
                    {:else}
                        No repositories connected yet.
                    {/if}
                </div>
                {#if !searchQuery}
                    <Button class="mt-4" href="https://github.com/settings/installations/{installationID}">
                        Connect Your First Repository
                    </Button>
                {/if}
            </div>
        {/if}
    </div>
</div>