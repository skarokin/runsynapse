package events

import (
	"github.com/skarokin/runsynapse/webhooks/supabase"
)


func TriggerCloudBuild(client *supabase.SupabaseClient, payload any) error {
	// 1. extract installation ID and branch from payload
	// 2. get user's refresh token, GCP project ID, and Dockerfile path associated with installation ID
	// 3. using GitHub App credentials, clone the repo (we need credentials for private repos)
	// 4. zip the cloned repo and upload to GCS
	//     - if user does not have GCS bucket, create one (we have their credentials so we can do it)
	// 5. define a Cloud Build config
	// 6. trigger Cloud Build with the config and uploaded zip file
	// 7. return build status to user, along with any logs or artifact
	// 8. cleanup the GCS bucket (delete the zip file we just pushed) and delete cloned repo from disk
	//	  - use defer to ensure cleanup happens no matter what
	return nil
}

func getTrafficAllocation(branch string) string {
    if branch == "main" || branch == "master" {
        return "--traffic=100"
    }
    return "--no-traffic"
}