<script lang="ts">
    import { Button } from "$lib/components/ui/button";
    import { Textarea } from "$lib/components/ui/textarea/index.js";
    import Send from "@lucide/svelte/icons/send";
    import Paperclip from "@lucide/svelte/icons/paperclip";
    
    import { toast } from "svelte-sonner";

    let {
        aiMode,
        newThought,
        onNewThought = (thought: any) => {},
        onAISearch = (query: string) => {},
    } = $props();

    function handleTextareaKeydown(event: KeyboardEvent) {
        if (event.key === "Enter" && !event.shiftKey) {
            event.preventDefault();
            handleSubmit();
        }
    }

    async function handleSubmit() {
        if (!inputText.trim()) return;

        if (aiMode) {
            // --- In AI Mode, call the AI handler ---
            onAIFollowUp(inputText);
            inputText = ''; // Clear the input
        } else {
            // --- In Normal Mode, call the new thought handler ---
            isLoading = true;
            try {
                const formData = new FormData();
                formData.append("thought", inputText);
                pendingFiles.forEach((file) => formData.append("files", file));
                
                // The onNewThought function now handles the API call
                await onNewThought(formData);

                inputText = "";
                pendingFiles = [];
            } catch (error) {
                toast.error("Failed to create thought.");
            } finally {
                isLoading = false;
            }
        }
    }

    function handleFileUpload() {
        fileInput?.click();
    }

    function onFileSelected(event: Event) {
        const target = event.target as HTMLInputElement;
        const files = target.files;
        if (files && files.length > 0) {
            pendingFiles = [...pendingFiles, ...Array.from(files)];
            target.value = '';
        }
    }
    
    let fileInput = $state<HTMLInputElement>();
    let pendingFiles = $state<File[]>([]);
    let isUploading = $state<boolean>(false);
    let isLoading = $state<boolean>(false);

</script>

<div class="w-full flex flex-row items-center gap-2">
    {#if !aiMode}
        <!-- if not AI mode, allow file uploads -->
        <Button
            variant="ghost"
            size="sm"
            class="h-8 w-8 p-0 shrink-0"
            onclick={handleFileUpload}
            disabled={isUploading}
        >
            <Paperclip class="w-4 h-4" />
        </Button>
    {/if}
    <Textarea
        bind:value={newThought}
        placeholder={aiMode
            ? "Ask AI about your thoughts..."
            : "What's on your mind?"}
        class="border-0 focus-visible:ring-0 focus-visible:ring-offset-0 min-h-[40px] resize-none shadow-none overflow-y-hidden flex-1"
        onkeydown={handleTextareaKeydown}
        rows={1}
        disabled={isLoading}
    />
    <Button
        size="sm"
        class="h-8"
        onclick={sendThought}
        disabled={!newThought.trim() || isLoading}
    >
        {#if isLoading}
            <div
                class="animate-spin h-4 w-4 border-2 border-t-transparent border-white rounded-full"
            ></div>
        {:else}
            <Send class="w-4 h-4" />
        {/if}
    </Button>
</div>

<!-- hidden file input -->
<input
    bind:this={fileInput}
    type="file"
    multiple
    onchange={onFileSelected}
    accept="image/*,application/pdf,.txt,.doc,.docx"
    class="hidden"
/>