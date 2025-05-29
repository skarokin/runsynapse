<script lang="ts">
    import * as Avatar from "$lib/components/ui/avatar/index.js";
    import * as DropdownMenu from "$lib/components/ui/dropdown-menu/index.js";
    import LogOutIcon from "@lucide/svelte/icons/log-out";
    import SettingsIcon from "@lucide/svelte/icons/settings";
    import { goto } from "$app/navigation";
    import { onMount } from "svelte";

    let { userMetadata } = $props();

    // idk why but im getting hydration mismatch if I dont do this
    // this was working just fine previously but now I need this??? bruh what is going on
    let mounted = $state(false);
    onMount(() => {
        mounted = true;
    });
</script>

{#if mounted}
    <DropdownMenu.Root>
        <DropdownMenu.Trigger class="hover:cursor-pointer">
            <Avatar.Root class="w-8 h-8">
                <Avatar.Image
                    class="rounded-full"
                    src={userMetadata?.avatar_url}
                    alt="User Avatar"
                />
                <Avatar.Fallback class="bg-primary text-muted-foreground">
                    {userMetadata?.name?.charAt(0)}
                </Avatar.Fallback>
            </Avatar.Root>
        </DropdownMenu.Trigger>
        <DropdownMenu.Content class="w-48 font-sans">
            <DropdownMenu.Group>
                <DropdownMenu.Label>
                    {userMetadata?.name}
                    <div class="text-xs text-muted-foreground font-normal">
                        {userMetadata?.email}
                    </div>
                </DropdownMenu.Label>
                <DropdownMenu.Separator />
                <DropdownMenu.Item onclick={() => goto("/settings")}>
                    <SettingsIcon class="w-4 h-4" />
                    <span>Settings</span>
                </DropdownMenu.Item>
                <DropdownMenu.Separator />
                <DropdownMenu.Item onclick={() => goto("/auth/logout")}>
                    <LogOutIcon class="w-4 h-4" />
                    <span>Log Out</span>
                </DropdownMenu.Item>
            </DropdownMenu.Group>
        </DropdownMenu.Content>
    </DropdownMenu.Root>
{/if}