package events

import (
    "context"
    "log"
	"fmt"

    "github.com/skarokin/runsynapse/webhooks/supabase"
)

func HandleInstallationEvent(client *supabase.SupabaseClient, payload any) error {
    payloadMap, ok := payload.(map[string]any)
    if !ok {
        return fmt.Errorf("invalid payload format")
    }

    action, ok := payloadMap["action"].(string)
    if !ok {
        return fmt.Errorf("missing action field")
    }

    installation, ok := payloadMap["installation"].(map[string]any)
    if !ok {
        return fmt.Errorf("missing installation field")
    }

    installationID, ok := installation["id"].(float64)
    if !ok {
        return fmt.Errorf("missing installation id")
    }

    switch action {
    case "created":
        log.Printf("Installation created: %d (handled on SvelteKit side)", int(installationID))
        return nil // this is handled on the SvelteKit app side

    case "deleted":
        return deleteInstallation(client, int(installationID))

    case "suspend":
        return suspendInstallation(client, int(installationID))

    case "unsuspend":
        return resumeInstallation(client, int(installationID))

    case "new_permissions_accepted":
        repositories, ok := payloadMap["repositories"].([]any)
        if !ok {
            return fmt.Errorf("missing repositories field")
        }
        return addRepositories(client, int(installationID), repositories)

    default:
        log.Printf("Unhandled installation action: %s", action)
        return nil
    }
}

func deleteInstallation(client *supabase.SupabaseClient, installationID int) error {
    ctx := context.Background()
    
    _, err := client.Pool.Exec(ctx, "SELECT delete_installation($1)", installationID)
    if err != nil {
        return fmt.Errorf("failed to delete installation: %w", err)
    }

    log.Printf("Deleted installation: %d", installationID)
    return nil
}

func suspendInstallation(client *supabase.SupabaseClient, installationID int) error {
    ctx := context.Background()
    
    _, err := client.Pool.Exec(ctx, "SELECT suspend_installation($1)", installationID)
    if err != nil {
        return fmt.Errorf("failed to suspend installation: %w", err)
    }

    log.Printf("Suspended installation: %d", installationID)
    return nil
}

func resumeInstallation(client *supabase.SupabaseClient, installationID int) error {
    ctx := context.Background()
    
    _, err := client.Pool.Exec(ctx, "SELECT resume_installation($1)", installationID)
    if err != nil {
        return fmt.Errorf("failed to resume installation: %w", err)
    }

    log.Printf("Resumed installation: %d", installationID)
    return nil
}