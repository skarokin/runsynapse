<script>
    import * as Card from "$lib/components/ui/card/index.js";
    import { Button } from "$lib/components/ui/button/index.js";
    import { Toaster } from "$lib/components/ui/sonner/index.js";
    import { toast } from "svelte-sonner";

    import Google from "$lib/assets/google.svelte";

    import { MessageSquare, Search, Pin } from "@lucide/svelte";

    import { page } from "$app/state";

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
        Think on the run, connect at rest.
    </h1>
    <h3 class="text-xl text-muted-foreground text-pretty">
        Capture half-formed thoughts as they flow. Scroll through your mental
        stream. Let AI connect the dots between your scattered ideas.
    </h3>
</div>

<!-- features -->
<div class="grid sm:grid-cols-3 gap-8    w-full sm:w-3/4">
    <div class="space-y-3">
        <h3 class="text-lg font-semibold flex items-center gap-2">
            <MessageSquare class="h-5 w-5 text-blue-600 dark:text-blue-400" />
            Stream of consciousness
        </h3>
        <p class="text-muted-foreground text-sm">
            Just like texting yourself. No formatting, no structure. Pure
            thought flow.
        </p>
    </div>

    <div class="space-y-3">
        <h3 class="text-lg font-semibold flex items-center gap-2">
            <Search class="h-5 w-5 text-purple-600 dark:text-purple-400" />
            AI-powered search
        </h3>
        <p class="text-muted-foreground text-sm">
            "Find that thing about productivity" – AI understands context, not
            just keywords.
        </p>
    </div>

    <div class="space-y-3">
        <h3 class="text-lg font-semibold flex items-center gap-2">
            <Pin class="h-5 w-5 text-green-600 dark:text-green-400" />
            Pin the good stuff
        </h3>
        <p class="text-muted-foreground text-sm">
            Important thoughts rise to the top. Everything else flows naturally.
        </p>
    </div>
</div>

<!-- log in -->
<Card.Root class="self-center max-w-md mx-auto">
    <Card.Content>
        <Button class="flex items-center gap-2 w-full" href="/auth/login">
            <p>Sign in with Google</p>
            <Google width={16} height={16} />
        </Button>
    </Card.Content>
    <Card.Footer>
        <p class="text-muted-foreground text-sm text-pretty">
            By continuing, you agree to our
            <a href="/terms" class="hover:text-primary underline">
                Terms of Service
            </a>
            and
            <a href="/privacy" class="hover:text-primary underline">
                Privacy Policy
            </a>.
        </p>
    </Card.Footer>
</Card.Root>

<p class="text-sm text-muted-foreground text-pretty">
    For people who think in thoughts, not essays
</p>