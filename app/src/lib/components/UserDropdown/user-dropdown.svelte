<script lang="ts">
    import * as Avatar from "$lib/components/ui/avatar/index.js";
    import * as DropdownMenu from "$lib/components/ui/dropdown-menu/index.js";
    import LogOutIcon from "@lucide/svelte/icons/log-out";
    import SettingsIcon from "@lucide/svelte/icons/settings";
    import { ModeWatcher, toggleMode } from "mode-watcher";
    import LightDark from "$lib/assets/light-dark.svelte";

    import { goto } from "$app/navigation";

    let { userMetadata } = $props();
</script>

<ModeWatcher />

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
            <DropdownMenu.Item onclick={toggleMode}>
                <LightDark />
                <span>Toggle Theme</span>
            </DropdownMenu.Item>
            <DropdownMenu.Separator />
            <DropdownMenu.Item onclick={() => goto("/auth/logout")}>
                <LogOutIcon class="w-4 h-4" />
                <span>Log Out</span>
            </DropdownMenu.Item>
        </DropdownMenu.Group>
    </DropdownMenu.Content>
</DropdownMenu.Root>
