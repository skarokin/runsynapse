<script lang="ts">
	import Logo from "$lib/assets/logo.svelte";
	import '../../app.css';
	
	let {children } = $props()
</script>

<div class="h-screen w-full flex flex-col font-sans">
    <div class="w-full mt-8">
        <header class="px-4 sm:px-0 container mx-auto flex items-center justify-center">
			<a class="flex items-center gap-4" href="/">
            	<Logo width={48} height={48} />
				<h1 class="text-xl">
					runsynapse
				</h1>
			</a>
        </header>
    </div>
	<div class="container mx-auto flex-1">
		<main class="px-4 sm:px-0 py-4 flex flex-col items-center justify-center h-full gap-12">
			{@render children()}
		</main>
	</div>
</div>
