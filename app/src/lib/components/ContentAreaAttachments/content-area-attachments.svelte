<script lang="ts">
    import { Button } from "$lib/components/ui/button";
    import { ExternalLink, Download, FileText } from "@lucide/svelte";

    let { attachments } = $props();

    function getFileTypeFromUrl(
        url: string,
    ): "image" | "pdf" | "video" | "audio" | "other" {
        const extension = url.split(".").pop()?.toLowerCase();

        if (["jpg", "jpeg", "png", "gif", "webp"].includes(extension || "")) {
            return "image";
        } else if (extension === "pdf") {
            return "pdf";
        } else if (["mp4", "webm", "ogg"].includes(extension || "")) {
            return "video";
        } else if (["mp3", "wav", "ogg"].includes(extension || "")) {
            return "audio";
        }
        return "other";
    }

    function getFileNameFromUrl(url: string): string {
        const parts = url.split("/");
        const filename = parts[parts.length - 1];
        // Remove the hash part but keep the extension
        const withoutHash = filename.split("_");
        if (withoutHash.length > 1) {
            return withoutHash[0] + "." + filename.split(".").pop();
        }
        return filename;
    }

    function openImageModal(url: string) {
        // You can implement a modal or just open in new tab
        window.open(url, "_blank");
    }
</script>

{#if attachments && attachments.length > 0}
    <div class="mt-3 space-y-2">
        {#each attachments as attachment}
            {@const fileType = getFileTypeFromUrl(attachment)}
            {@const fileName = getFileNameFromUrl(attachment)}

            {#if fileType === "image"}
                <!-- Image Preview -->
                <div class="relative max-w-sm">
                    <button
                        type="button"
                        class="block w-full text-left p-0 border-0 bg-transparent cursor-pointer hover:opacity-90 transition-opacity focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 rounded-lg"
                        onclick={() => openImageModal(attachment)}
                        aria-label="Open image {fileName} in new tab"
                    >
                        <img
                            src={attachment}
                            alt={fileName}
                            class="rounded-lg border max-w-full h-auto"
                            loading="lazy"
                        />
                    </button>
                    <Button
                        variant="ghost"
                        size="sm"
                        class="absolute top-2 right-2 h-8 w-8 p-0 bg-black/50 hover:bg-black/70 text-white"
                        onclick={() => openImageModal(attachment)}
                    >
                        <ExternalLink class="h-4 w-4" />
                    </Button>
                </div>
            {:else if fileType === "video"}
                <!-- Video Player -->
                <div class="max-w-md">
                    <video
                        src={attachment}
                        controls
                        class="rounded-lg border w-full"
                        preload="metadata"
                    >
                        <track kind="captions" src="" label="No captions available" />
                        Your browser does not support the video tag.
                    </video>
                </div>
            {:else if fileType === "audio"}
                <!-- Audio Player -->
                <div class="max-w-md">
                    <audio
                        src={attachment}
                        controls
                        class="w-full"
                        preload="metadata"
                    >
                        Your browser does not support the audio tag.
                    </audio>
                </div>
            {:else}
                <!-- Other File Types (PDF, documents, etc.) -->
                <div
                    class="flex items-center gap-2 p-2 bg-muted/30 rounded border max-w-sm"
                >
                    <FileText
                        class="h-4 w-4 text-muted-foreground flex-shrink-0"
                    />
                    <div class="flex-1 min-w-0">
                        <p class="text-sm font-medium truncate">{fileName}</p>
                        <p class="text-xs text-muted-foreground">
                            {fileType === "pdf" ? "PDF Document" : "File"}
                        </p>
                    </div>
                    <Button
                        variant="ghost"
                        size="sm"
                        class="h-8 w-8 p-0 flex-shrink-0"
                        onclick={() => window.open(attachment, "_blank")}
                    >
                        <Download class="h-4 w-4" />
                    </Button>
                </div>
            {/if}
        {/each}
    </div>
{/if}
