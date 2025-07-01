<script lang="ts">
    import { Button } from "$lib/components/ui/button/index.js";
    import { Input } from "$lib/components/ui/input/index.js";
    import Send from "@lucide/svelte/icons/send";
    import Paperclip from "@lucide/svelte/icons/paperclip";
    import Search from "@lucide/svelte/icons/search";
    import Brain from "@lucide/svelte/icons/brain";
    import ArrowLeft from "@lucide/svelte/icons/arrow-left";
    import { Textarea } from "$lib/components/ui/textarea/index.js";
    import { toast } from "svelte-sonner";

    import { ContentArea } from '$lib/components/ContentArea';
    import { FilePreview } from '$lib/components/FilePreview';
    import { PinnedThoughts } from '$lib/components/PinnedThoughts';

    import { onMount, tick } from 'svelte';

    let { data } = $props();

    const thoughts = $derived(data.thoughts);

    let hasMoreAbove = $state(data.hasMoreAbove);
    
    let newThought = $state<string>('');
    let fileInput = $state<HTMLInputElement>();
    let pendingFiles = $state<File[]>([]);
    let isUploading = $state<boolean>(false);
    let isLoading = $state<boolean>(false);
    
    let aiMode = $state<boolean>(false);
    let searchQuery = $state<string>('');

    // ai responses
    let aiSummary = $state<string>('');
    let retrievedThoughts = $state<Array<{id: number, content: string, timestamp: string, pinned: boolean}>>([]);

    let contentScrollAreaRef = $state<HTMLElement | null>(null);

    async function sendThought() {
        if (!newThought.trim()) return;
        
        isLoading = true;
        
        try {
            const formData = new FormData();
                
            formData.append('thought', newThought);
            pendingFiles.forEach(file => {
                formData.append('files', file);
            });

            const res = await fetch('/synapse/api/newThought', {
                method: 'POST',
                body: formData
            });

            if (!res.ok) {
                toast.error('Failed to create new thought');
                return;
            }

            const result = await res.json();
            if (result.error) {
                toast.error('Failed to create thought', result.error);
                return;
            }

            thoughts.push(result.thought);

            // scroll to the bottom of the content area
            setTimeout(() => {
                if (contentScrollAreaRef) {
                    const viewport = contentScrollAreaRef.querySelector('[data-slot="scroll-area-viewport"]');
                    console.log('Found viewport:', viewport);
                    if (viewport) {
                        viewport.scrollTop = viewport.scrollHeight;
                    }
                }
            }, 0);
            
            newThought = '';
            pendingFiles = [];
        } catch (error) {
            console.error('Error:', error);
        } finally {
            isLoading = false;
        }
    }

    async function handleSearch() {
        if (!searchQuery.trim()) return;

        isLoading = true;
        aiMode = true; // switch to AI mode on search

        try {
            // this now simulates the full RAG response from your Go backend
            await new Promise(resolve => setTimeout(resolve, 2000));
            
            aiSummary = `Based on your thoughts, here's what I found about "${searchQuery}"...`;
            retrievedThoughts = [
                {
                    id: 99,
                    content: "Why do we even need passwords? Biometric everything should be the default by now. Face ID for literally everything.",
                    timestamp: "2:58 PM",
                    pinned: true
                },
                {
                    id: 98,
                    content: "Just had this random idea about how we could optimize the login flow. What if we pre-filled email based on browser history?",
                    timestamp: "2:15 PM",
                    pinned: false
                }
            ];
        } catch (error) {
            console.error('Error:', error);
            toast.error('Failed to search thoughts. Please try again.');
        } finally {
            isLoading = false;
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

    function exitAIMode() {
        aiMode = false;
        newThought = '';
        searchQuery = '';
        aiSummary = '';
        retrievedThoughts = [];
        setTimeout(() => {
            const textarea = document.querySelector('textarea');
            if (textarea) {
                (textarea as HTMLTextAreaElement).focus();
                (textarea as HTMLTextAreaElement).select();
            }
        });
    }

    function handleSearchKeydown(event: KeyboardEvent) {
        if (event.key === 'Enter') {
            event.preventDefault();
            handleSearch();
        }
    }

    function handleTextareaKeydown(event: KeyboardEvent) {
        if (event.key === 'Enter' && !event.shiftKey) {
            event.preventDefault();
            sendThought();
        }
    }

    // force scroll to bottom when thoughts are set
    onMount(async() => {
        await tick();
        if (thoughts.length > 0 && contentScrollAreaRef) {
            const viewport = contentScrollAreaRef.querySelector('[data-slot="scroll-area-viewport"]');
            if (viewport) {
                viewport.scrollTop = viewport.scrollHeight;
            }
        }
    });
</script>

<svelte:window onkeydown={(e) => {
    if (e.key === 'Escape' && aiMode) {
        exitAIMode();
    }
}} />

<div class="flex flex-col h-full gap-6 w-full">
   <!-- header -->
    <div class="flex items-center justify-between">
        {#if aiMode}
            <!-- ai mode - show back button and AI mode title -->
            <Button variant="ghost" size="sm" onclick={exitAIMode}>
                <ArrowLeft class="w-4 h-4" />
                (ESC)
            </Button>
            <div class="flex items-center gap-2">
                <Brain class="w-5 h-5 text-primary" />
                <span class="font-medium text-primary">AI Search</span>
            </div>
        {:else}
            <!-- thoughts mode - THE SEARCH BAR IS THE ENTRY POINT TO AI MODE -->
            <div class="flex-1 max-w-xs sm:max-w-md">
                <div class="relative">
                    <Search class="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-muted-foreground" />
                    <Input
                        bind:value={searchQuery}
                        placeholder="Search your thoughts..."
                        class="pl-10"
                        onkeydown={handleSearchKeydown}
                    />
                </div>
            </div>
            <PinnedThoughts
                pinnedThoughts={data.pinnedThoughts}
                thoughtSet={data.thoughtSet}
            />
        {/if}
    </div>

    <!-- content area -->
    <ContentArea
        aiMode={aiMode}
        aiSummary={aiSummary}
        retrievedThoughts={retrievedThoughts}
        thoughts={thoughts}
        hasMoreAbove={hasMoreAbove}
        isLoading={isLoading}
        bind:scrollAreaRef={contentScrollAreaRef}
    />

    <!-- input -->
    <div class="flex flex-col items-center gap-2 border rounded-lg p-2">
        <FilePreview bind:pendingFiles={pendingFiles} />
        
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
                placeholder={aiMode ? "Ask AI about your thoughts..." : "What's on your mind?"}
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
                    <div class="animate-spin h-4 w-4 border-2 border-t-transparent border-white rounded-full"></div>
                {:else}
                    <Send class="w-4 h-4" />
                {/if}
            </Button>
        </div>
    </div>

    <!-- hidden file input -->
    <input
        bind:this={fileInput}
        type="file"
        multiple
        class="hidden"
        onchange={onFileSelected}
        accept="image/*,application/pdf,.txt,.doc,.docx"
    />
</div>