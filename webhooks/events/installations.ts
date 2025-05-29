import { SupabaseClient } from "@supabase/supabase-js";

export const handleInstallationEvent = async (supabaseClient: SupabaseClient, payload: any) => {
    const action = payload.action;
    const installationID = payload.installation.id;

    if (action == "created") {
        return; // handled on the SvelteKit app side
    } else if (action == "deleted") {
        await deleteInstallation(supabaseClient, installationID);
    } else if (action == "suspend") {
        await suspendInstallation(supabaseClient, installationID);
    } else if (action == "unsuspend") {
        await resumeInstallation(supabaseClient, installationID);
    } else if (action == "new_permissions_accepted") {
        // can justs reuse the repositories event handler
        const repositoriesObject = payload.repositories;
        await addRepositories(supabaseClient, installationID, repositoriesObject);

    }
}

export const handleRepositoryInstallationEvent = async (supabaseClient: SupabaseClient, payload: any) => {
    const action = payload.action;
    const installationID = payload.installation.id;

    if (action == "added") {
        const repositoriesObject = payload.repositories_added;
        await addRepositories(supabaseClient, installationID, repositoriesObject); 
    } else if (action == "removed") {
        const repositoriesObject = payload.repositories_removed;
        await removeRepositories(supabaseClient, installationID, repositoriesObject);
    }
}

const addRepositories = async (
    supabaseClient: SupabaseClient, installationID: number, repositoriesObject: any[]
) => {
    // find the user_id for this installation
    const { data: userData, error: userError } = await supabaseClient
        .from('github_installations')
        .select('user_id')
        .eq('installation_id', installationID)
        .single();

    if (userError) {
        console.error('Error fetching user data for installation:', userError);
        throw userError;
    }

    if (!userData) {
        console.error('No user data found for installation:', installationID);
        throw new Error('User data not found for installation');
    }

    const repositoryRecords = repositoriesObject.map(repo => ({
        installation_id: installationID,
        repo_name: repo.full_name,        
        repo_id: repo.id,                 
        is_private: repo.private,         
        user_id: userData.user_id,
        dockerfile_path: './Dockerfile', // default path, can be updated later
    }));

    const { data, error } = await supabaseClient
        .from('github_repositories')
        .insert(repositoryRecords);

    if (error) {
        console.error('Error adding repositories:', error);
        throw error;
    }

    return data;
}

const removeRepositories = async (
    supabaseClient: SupabaseClient, installationID: number, repositoriesObject: any[]
) => {
    const repositories = repositoriesObject.map(repo => repo.full_name);

    const { data, error } = await supabaseClient
        .from('github_repositories')
        .delete()
        .eq('installation_id', installationID)
        .in('repo_name', repositories);

    if (error) {
        console.error('Error removing repositories:', error);
        throw error;
    }

    return data;
}

const deleteInstallation = async (
    supabaseClient: SupabaseClient, installationID: number
) => {
    const { data, error } = await supabaseClient
        .from('github_installations')
        .delete()
        .eq('installation_id', installationID);

    if (error) {
        console.error('Error deleting installation:', error);
        throw error;
    }

    return data;
}

const suspendInstallation = async (
    supabaseClient: SupabaseClient, installationID: number
) => {
    const { data, error } = await supabaseClient
        .from('github_installations')
        .update({ suspended: true })
        .eq('installation_id', installationID);

    if (error) {
        console.error('Error suspending installation:', error);
        throw error;
    }

    return data;
}

const resumeInstallation = async (
    supabaseClient: SupabaseClient, installationID: number
) => {
    const { data, error } = await supabaseClient
        .from('github_installations')
        .update({ suspended: false })
        .eq('installation_id', installationID);

    if (error) {
        console.error('Error resuming installation:', error);
        throw error;
    }

    return data;
}