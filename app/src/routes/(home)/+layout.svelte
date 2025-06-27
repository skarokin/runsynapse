<script lang="ts">
	import Logo from "$lib/assets/logo.svelte";
    import { Button } from "$lib/components/ui/button";
    import LightDark from "$lib/assets/light-dark.svelte";
    import { ModeWatcher, toggleMode } from "mode-watcher";
	import '../../app.css';
	
	let {children } = $props()
</script>

<ModeWatcher />

<div class=" h-screen w-full flex flex-col font-sans">
    <div class="w-full mt-8">
        <header class="px-4 sm:px-0 container mx-auto flex items-center justify-between">
            <div class="flex-1"></div>
            <a class="flex items-center gap-4" href="/">
                <Logo width={48} height={48} />
                <h1 class="text-xl">
                    runsynapse
                </h1>
            </a>
            <!-- div wrap is necessary to ensure button isnt stretched -->
            <div class="flex-1 flex justify-end">
                <Button
                    size="icon"
                    variant="ghost"
                    onclick={toggleMode}
                >
                    <LightDark />
                </Button>
            </div>
        </header>
    </div>
	<div class="container mx-auto flex-1">
		<main class="px-4 sm:px-0 py-4 flex flex-col items-center justify-center h-full gap-12">
			{@render children()}
		</main>
	</div>
</div>
