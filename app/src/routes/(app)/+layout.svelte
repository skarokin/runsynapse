<script lang="ts">
    import Logo from "$lib/assets/logo.svelte";
    import { Toaster } from "$lib/components/ui/sonner/index.js";
    import { UserDropdown } from "$lib/components/UserDropdown";

    import "../../app.css";

    let { data, children } = $props();
    const { userMetadata } = data;
</script>

<Toaster />

<div class="h-screen w-full flex flex-col font-sans">
    <div class="w-full border-b">
        <header
            class="px-4 py-2 h-full mx-auto flex items-center justify-between"
        >
            <div class="flex gap-8 items-center">
                <a href="/synapse">
                    <Logo width={20} height={20} />
                </a>
            </div>
            <UserDropdown
                userMetadata={userMetadata}
            />
        </header>
    </div>
    <div class="container mx-auto flex-1 min-h-0">
        <main
            class="px-4 sm:px-0 h-full w-full py-8"
        >
            {@render children()}
        </main>
    </div>
</div>