<script>
    import * as Card from "$lib/components/ui/card/index.js";
    import { Button } from "$lib/components/ui/button/index.js";
    import { Toaster } from "$lib/components/ui/sonner/index.js";
    import { toast } from "svelte-sonner";

    import Google from "$lib/assets/google.svelte";

    import { page } from '$app/state';

    const urlParams = page.url.searchParams;
    const error = urlParams.get("error");

    $effect(() => {
        if (error) {
            toast.error(error);
        }
    });
</script>

<Toaster />

<!-- hero -->
<div class="space-y-2 w-full sm:w-3/4">
    <h1 class="text-4xl font-bold leading-tight">
        Intelligent Cloud Run <br class="hidden sm:inline" /> orchestration
    </h1>
    <h3 class="text-xl text-muted-foreground text-pretty">
        Streamline your containerized workloads with better monitoring, seamless
        integrations, and automated workflows.
    </h3>
</div>

<!-- features -->
<div class="grid grid-cols-2 gap-x-16 gap-y-6 w-full sm:w-3/4">
    <div class="space-y-2">
        <h3 class="font-medium">Enhanced monitoring</h3>
        <p class="text-sm text-muted-foreground">
            Better UI for metrics & billing visualization.
        </p>
    </div>
    <div class="space-y-2">
        <h3 class="font-medium">Natively integrated</h3>
        <p class="text-sm text-muted-foreground">
            Direct Cloud Run integration with seamless container view switching.
        </p>
    </div>
    <div class="space-y-2">
        <h3 class="font-medium">Smart alerts</h3>
        <p class="text-sm text-muted-foreground">
            Easy alert setup and monitoring for proactive issue resolution.
        </p>
    </div>
    <div class="space-y-2">
        <h3 class="font-medium">DevOps workflow</h3>
        <p class="text-sm text-muted-foreground">
            GitHub CI/CD integration with preview deployments and secrets
            management.
        </p>
    </div>
</div>

<!-- log in -->
<Card.Root class="self-center">
    <Card.Content>
        <Button
            class="flex items-center gap-2 w-full"
            href="/auth/login"
        >
            <p>Connect a Google Cloud account</p>
            <Google width={16} height={16} />
        </Button>
    </Card.Content>
    <Card.Footer>
        <p class="text-muted-foreground text-sm text-pretty">
            By continuing, you agree to our <a href="/terms" class="hover:text-primary underline">Terms of Service</a>
            and <a href="/privacy" class="hover:text-primary underline">Privacy Policy</a>.
        </p>
    </Card.Footer>
</Card.Root>
