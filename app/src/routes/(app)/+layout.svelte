<script lang="ts">
    import { UserDropdown } from "$lib/components/UserDropdown";
    import Logo from "$lib/assets/logo.svelte";
    import { Toaster } from "$lib/components/ui/sonner/index.js";

    import { page } from "$app/state";

    import "../../app.css";

    let { data, children } = $props();
    const { userMetadata } = data;

    const activeTab = $derived(
        page.url.pathname === "/console"
            ? "console"
            : page.url.pathname === "/settings"
            ? "settings"
            : "console"
    );

    const tabs = [
        { id: "console", label: "Console", href: "/console" },
        { id: "settings", label: "Settings", href: "/settings" },
    ];
</script>

<Toaster />

<div class="h-screen w-full flex flex-col font-sans">
    <div class="w-full border-b">
        <header
            class="px-4 py-2 h-full mx-auto flex items-center justify-between"
        >
            <div class="flex gap-8 items-center">
                <a href="/console">
                    <Logo width={20} height={20} />
                </a>
                {#each tabs as tab}
                    <a
                        href={tab.href}
                        class="font-medium text-xs transition-colors duration-200
                            {activeTab === tab.id
                                ? 'text-primary font-semibold'
                                : 'text-muted-foreground hover:text-foreground'}"
                    >
                        {tab.label}
                    </a>
                {/each}
            </div>
            <UserDropdown
                userMetadata={userMetadata}
            />
        </header>
    </div>
    <div class="container mx-auto flex-1">
        <main
            class="px-4 sm:px-0 py-8 flex flex-col items-start justify-start h-full w-full gap-12"
        >
            {@render children()}
        </main>
    </div>
</div>