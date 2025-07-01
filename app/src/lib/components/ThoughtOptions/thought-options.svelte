<script lang="ts">
    import { Button } from "$lib/components/ui/button/index.js";
    import * as DropdownMenu from "$lib/components/ui/dropdown-menu/index.js";
    import Pin from "@lucide/svelte/icons/pin";
    import Trash from "@lucide/svelte/icons/trash";
    import MoreHorizontal from "@lucide/svelte/icons/more-horizontal";

    let {
        thought,
        onTogglePin,
        onDelete,
    } = $props<{
        thought: any,
        onTogglePin: (id: string) => void;
        onDelete: (id: string) => void;
    }>();
</script>

<DropdownMenu.Root>
    <DropdownMenu.Trigger>
        <Button
            variant="ghost"
            size="sm"
            class="opacity-0 group-hover:opacity-100 h-6 w-6 p-0"
        >
            <MoreHorizontal class="w-3.5 h-3.5" />
        </Button>
    </DropdownMenu.Trigger>
    <DropdownMenu.Content align="end">
        <DropdownMenu.Item onclick={() => onTogglePin(thought.id)}>
            <Pin class="w-3.5 h-3.5 mr-2" />
            {thought.pinned ? "Unpin" : "Pin"}
        </DropdownMenu.Item>
        <DropdownMenu.Item
            onclick={() => onDelete(thought.id)}
            class="text-red-600 focus:text-red-600"
        >
            <Trash class="w-3.5 h-3.5 mr-2" />
            Delete
        </DropdownMenu.Item>
    </DropdownMenu.Content>
</DropdownMenu.Root>
