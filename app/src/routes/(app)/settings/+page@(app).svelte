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
    
    let { data } = $props();
    const { userMetadata } = data;

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
            Manage your account preferences and settings.
        </p>
    </div>

    <!-- profile -->
    <div class="space-y-4">
        <p class="text-xl font-semibold">Connected Google Account</p>
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
                        <h3 class="font-medium">Delete your account</h3>
                        <p class="text-sm text-muted-foreground">Permanently delete your account</p>
                    </div>
                </div>
                <Button size="sm" onclick={deleteAccount}>
                    Delete
                </Button>
            </li>
        </ul>
    </div>
</div>