<script lang="ts">
    import { Button, buttonVariants } from "$lib/components/ui/button";
    import * as Dialog from "$lib/components/ui/dialog";
    import Cloud from "@lucide/svelte/icons/cloud";
    import { Input } from "$lib/components/ui/input";
    import { Label } from "$lib/components/ui/label";
    import * as Select from "$lib/components/ui/select";
    import { Switch } from "$lib/components/ui/switch";
    import { toast } from "svelte-sonner";

    let { projects = [], repo } = $props();

    let open = $state(false);
    let selectedProject = $state("");
    let dockerfilePath = $state(repo.dockerfile_path || "./Dockerfile");
    let autoDeployMain = $state(true);
    let enablePreviewDeployments = $state(true);

    let isSubmitting = $state(false);

    const selectTriggerContent = $derived(
        projects.find((p) => p.projectID === selectedProject)?.displayName ??
            "Select a project",
    );

    const defaultServiceName = repo.repo_name
        .toLowerCase()
        .replace(/[^a-z0-9-]/g, '-')
        .replace(/^-+|-+$/g, '')  
        .replace(/-+/g, '-')
        .slice(0, 49) || "my-service";

    let serviceName = $state(defaultServiceName);

    const isValidServiceName = $derived(
        /^[a-z0-9]([a-z0-9-]*[a-z0-9])?$/.test(serviceName) && 
        serviceName.length >= 1 && 
        serviceName.length <= 49
    );

    async function handleSubmit(e: Event) {
        e.preventDefault();

        if (!selectedProject || !dockerfilePath || !isValidServiceName) {
            return;
        }

        isSubmitting = true;

        try {
            const formData = new FormData(e.target as HTMLFormElement);
            const res = await fetch("?/connectRepo", {
                method: "POST",
                body: formData,
            })
            
            if (!res.ok) {
                const error = await res.json();
                throw new Error(error.message || "Failed to connect repository");
            }

            const data = await res.json();

            console.log(data);

            if (data.type === "success") {
                toast.success("Repository connected successfully!");
                open = false; 
            } else {
                toast.error("Failed to connect repository");
            }
            
        } catch (error) {
            console.error("Error submitting form:", error);
            toast.error("Failed to connect repository");
        } finally {
            isSubmitting = false;
        }
    }

    function handleCancel() {
        open = false;
    }
</script>

<Dialog.Root bind:open>
    <Dialog.Trigger class={buttonVariants({ variant: "outline", size: "sm" })}>
        <Cloud class="w-3 h-3 mr-1 text-blue-500" />
        Connect to GCP
    </Dialog.Trigger>
    <Dialog.Content>
        <Dialog.Header>
            <Dialog.Title>Connect to Google Cloud Platform</Dialog.Title>
            <Dialog.Description>
                Connect this repository and configure basic deployment settings.
            </Dialog.Description>
        </Dialog.Header>
        <form class="space-y-4" onsubmit={handleSubmit}>
            <input type="hidden" name="repoID" value={repo.id} />
            <input type="hidden" name="autoDeployMain" value={autoDeployMain} />
            <input type="hidden" name="enablePreviewDeployments" value={enablePreviewDeployments} />
            <input type="hidden" name="serviceName" value={serviceName} />

            <div class="space-y-1">
                <Label for="project-select">GCP Project <span class="text-red-500">*</span></Label>
                <Select.Root
                    type="single"
                    name="selectedProject"
                    bind:value={selectedProject}
                >
                    <Select.Trigger id="project-select" class="w-full">
                        {selectTriggerContent}
                    </Select.Trigger>
                    <Select.Content>
                        <Select.Group>
                            <Select.Label>Projects</Select.Label>
                            {#each projects as project}
                                <Select.Item value={project.projectID}>
                                    {project.displayName} ({project.projectID})
                                </Select.Item>
                            {/each}
                        </Select.Group>
                    </Select.Content>
                </Select.Root>
                {#if projects.length === 0}
                    <p class="text-xs text-muted-foreground">
                        No projects found.
                        <a href="https://console.cloud.google.com/projectcreate" target="_blank" class="text-blue-600 hover:underline">
                            Create one in Google Cloud Console
                        </a>
                    </p>
                {/if}
            </div>

            <div class="space-y-1">
                <Label for="service-name">Service Name Prefix <span class="text-red-500">*</span></Label>
                <Input
                    id="service-name"
                    name="serviceName"
                    bind:value={serviceName}
                    placeholder={defaultServiceName}
                    class={!isValidServiceName ? "border-red-500" : ""}
                />
                <p class="text-xs text-muted-foreground">
                    Cloud Run service prefix. Branch names will be appended (e.g., {serviceName}-main, {serviceName}-feature-auth).
                </p>
                {#if !isValidServiceName && serviceName}
                    <p class="text-xs text-red-600">
                        Invalid service name. Use lowercase letters, numbers, and hyphens (1-49 chars).
                    </p>
                {/if}
            </div>

            <div class="space-y-1">
                <Label for="dockerfilePath">Dockerfile Path <span class="text-red-500">*</span></Label>
                <Input
                    id="dockerfilePath"
                    name="dockerfilePath"
                    bind:value={dockerfilePath}
                    placeholder="./Dockerfile"
                />
            </div>

            <div class="space-y-2">
                <div class="flex items-center justify-between">
                    <Label for="auto-deploy-main" class="text-sm font-normal">
                        Auto Deploy Main Branch
                    </Label>
                    <Switch id="auto-deploy-main" bind:checked={autoDeployMain} />
                </div>

                <div class="flex items-center justify-between">
                    <Label for="enable-preview-deployments" class="text-sm font-normal">
                        Enable Preview Deployments
                    </Label>
                    <Switch id="enable-preview-deployments" bind:checked={enablePreviewDeployments} />
                </div>
            </div>

            <p class="text-xs text-muted-foreground">
                <strong>Note:</strong> More advanced features can be configured inside the
                Google Cloud Console.
            </p>

            <Dialog.Footer>
                <Button type="button" variant="outline" onclick={handleCancel} disabled={isSubmitting}>
                    Cancel
                </Button>
                <Button type="submit" disabled={!selectedProject || isSubmitting}>
                    {isSubmitting ? 'Connecting...' : 'Connect Repository'}
                </Button>
            </Dialog.Footer>
        </form>
    </Dialog.Content>
</Dialog.Root>