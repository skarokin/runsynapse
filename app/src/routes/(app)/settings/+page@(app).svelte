<script lang="ts">
    import { Badge } from "$lib/components/ui/badge/index.js";
    import { Button } from "$lib/components/ui/button/index.js";
    import * as Card from "$lib/components/ui/card/index.js";
    import * as Avatar from "$lib/components/ui/avatar/index.js";
    import { Separator } from "$lib/components/ui/separator/index.js";
    import { Switch } from "$lib/components/ui/switch/index.js";
    import { toast } from "svelte-sonner";

    import LogOut from "@lucide/svelte/icons/log-out";
    import Trash2 from "@lucide/svelte/icons/trash-2";

    import { SettingsGithub } from "$lib/components/SettingsGithub";

    import { goto } from "$app/navigation";
    import { page } from "$app/state";
    
    let { data } = $props();
    const { userMetadata, connectedRepos, isConnected, installationID } = data;
    
    let autoDeployEnabled = $state(data.autoDeploy ?? false);
    let onlyMainEnabled = $state(data.onlyMain ?? false);

    $effect(() => {
        const urlParams = page.url.searchParams;
        const githubStatus = urlParams.get("github");
        if (githubStatus === "connected") {
            toast.success(`Github account connected successfully!`);
            goto('/settings', { replaceState: true });
        } else if (githubStatus === "error" || githubStatus === "cancelled") {
            toast.error(`Failed to connect Github account: ${urlParams.get("reason")}`);
            goto('/settings', { replaceState: true });
        }
    })

    async function updatePreference(setting: string, value: boolean) {
        // update the preference in the backend
    }

    async function deleteAccount() {
        // delete this account :( noo pls dont leave
    }
</script>

<div class="space-y-6 w-full">
    <div>
        <h2 class="text-2xl font-bold">Account Settings</h2>
        <p class="text-sm text-muted-foreground">
            Manage your account preferences, GitHub integration, and more.
        </p>
    </div>

    <!-- profile -->
    <div class="space-y-4">
        <p class="text-xl font-semibold">Profile</p>
        <Card.Root>      
            <Card.Content class="space-y-4">
                <div class="flex items-center gap-4">
                    <Avatar.Root>
                        <Avatar.Image src={userMetadata?.avatar_url} alt="Avatar" />
                        <Avatar.Fallback>{userMetadata?.name?.[0] || "U"}</Avatar.Fallback>
                    </Avatar.Root>
                    <div>
                        <h3 class="font-medium">{userMetadata?.name}</h3>
                        <p class="text-sm text-muted-foreground">{userMetadata?.email}</p>
                    </div>
                </div>
            </Card.Content>
        </Card.Root>
    </div>  

    <!-- preferences -->
    <div class="space-y-4">
        <p class="text-xl font-semibold">Preferences</p>
        <ul class="flex flex-col gap-4">
            <li class="flex items-center justify-between">
                <div>
                    <h3 class="font-medium">Auto-deploy</h3>
                    <p class="text-sm text-muted-foreground">Deploy automatically on git push.</p>
                </div>
                <Switch class="ml-4 hover:cursor-pointer" bind:checked={autoDeployEnabled} />
            </li>
            <Separator />
            <li class="flex items-center justify-between">
                <div>
                    <h3 class="font-medium">Only deploy main branch</h3>
                    <p class="text-sm text-muted-foreground">Disables staging/preview deployments from other branches.</p>
                </div>
                <Switch class="ml-4 hover:cursor-pointer" bind:checked={onlyMainEnabled} />
            </li>
        </ul>
    </div>
    
    <!-- github integration -->
    <div class="space-y-4">
        <p class="text-xl font-semibold flex justify-between items-center">
            GitHub Integration
            <Badge variant={isConnected ? "default" : "secondary"}>
                {isConnected ? "Connected" : "Disconnected"}
            </Badge>
        </p>

        <SettingsGithub
            connectedRepos={connectedRepos}
            isConnected={isConnected}
            installationID={installationID}
        />
    </div>

    <!-- account actions -->
    <div class="space-y-4">
        <p class="text-xl font-semibold">Account Actions</p>
        <ul class="flex flex-col gap-4">
            <li class="flex items-center justify-between">
                <div class="flex items-center gap-4">
                    <LogOut size="20" class="text-muted-foreground" />
                    <div>
                        <h3 class="font-medium">Log out</h3>
                        <p class="text-sm text-muted-foreground">Log out of your account</p>
                    </div>
                </div>
                <Button variant="outline" size="sm" href="/auth/logout">
                    Log out
                </Button>
            </li>
            <Separator />
            <li class="flex items-center justify-between">
                <div class="flex items-center gap-4">
                    <Trash2 size="20" class="text-destructive" />
                    <div>
                        <h3 class="font-medium text-destructive">Delete your account</h3>
                        <p class="text-sm text-muted-foreground">Permanently delete your account</p>
                    </div>
                </div>
                <Button variant="destructive" size="sm" onclick={deleteAccount}>
                    Delete
                </Button>
            </li>
        </ul>
    </div>
</div>