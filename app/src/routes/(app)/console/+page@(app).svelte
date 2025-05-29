<script>
    import { Button } from "$lib/components/ui/button/index.js";
    import * as Card from "$lib/components/ui/card/index.js";
    import { Badge } from "$lib/components/ui/badge/index.js";

    let { data } = $props();
    const { projects } = $derived(data);
</script>

<div class="space-y-6 w-full">
    <div>
        <h2 class="text-2xl font-bold">Console</h2>
        <p class="text-sm text-muted-foreground">
            Orchestrate your projects, deployments, and integrations.
        </p>
    </div>
    <div class="space-y-4">
        {#each projects || [] as project}
            <Card.Root>
                <Card.Content class="space-y-4">
                    <div class="flex items-center justify-between">
                        <div class="flex items-center gap-4">
                            <Badge variant="secondary" class="text-xs">{project.state}</Badge>
                            <h3 class="text-lg font-semibold">{project.displayName}</h3>
                        </div>
                    </div>
                </Card.Content> 
                <Card.Footer class="flex items-center justify-between">
                    <Button variant="outline" size="sm" href={`/console/${project.projectID}`}>
                        View Project
                    </Button>
                    <Button variant="destructive" size="sm" href={`/console/${project.projectID}/delete`}>
                        Delete Project
                    </Button>
                </Card.Footer>
            </Card.Root>
        {/each}
    </div>
</div>
