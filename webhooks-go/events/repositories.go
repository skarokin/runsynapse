package events

import (
    "context"
    "log"
    "fmt"
    "encoding/json"

    "github.com/skarokin/runsynapse/webhooks/supabase"
)

func HandleInstallationRepositoriesEvent(client *supabase.SupabaseClient, payload any) error {
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
    case "added":
        repositoriesAdded, ok := payloadMap["repositories_added"].([]any)
        if !ok {
            return fmt.Errorf("missing repositories_added field")
        }
        return addRepositories(client, int(installationID), repositoriesAdded)

    case "removed":
        repositoriesRemoved, ok := payloadMap["repositories_removed"].([]any)
        if !ok {
            return fmt.Errorf("missing repositories_removed field")
        }
        return removeRepositories(client, int(installationID), repositoriesRemoved)

    default:
        log.Printf("Unhandled repository installation action: %s", action)
        return nil
    }
}

func addRepositories(client *supabase.SupabaseClient, installationID int, repositories []any) error {
    ctx := context.Background()

    reposJSON, err := json.Marshal(repositories)
    if err != nil {
        return fmt.Errorf("failed to marshal repositories: %w", err)
    }
    
    _, err = client.Pool.Exec(ctx, 
        "SELECT add_repositories($1, $2)", 
        installationID, string(reposJSON))
    
    if err != nil {
        return fmt.Errorf("failed to add repositories: %w", err)
    }
    
    log.Printf("Added %d repositories for installation %d", len(repositories), installationID)
    return nil
}

func removeRepositories(client *supabase.SupabaseClient, installationID int, repositories []any) error {
    ctx := context.Background()

    var repoNames []string
    for _, repo := range repositories {
        if repoMap, ok := repo.(map[string]any); ok {
            if fullName, ok := repoMap["full_name"].(string); ok {
                repoNames = append(repoNames, fullName)
            }
        }
    }

    if len(repoNames) == 0 {
        return nil
    }

    _, err := client.Pool.Exec(ctx, 
        "SELECT remove_repositories($1, $2)", 
        installationID, repoNames)
    
    if err != nil {
        return fmt.Errorf("failed to remove repositories: %w", err)
    }
    
    log.Printf("Removed %d repositories for installation %d", len(repoNames), installationID)
    return nil
}