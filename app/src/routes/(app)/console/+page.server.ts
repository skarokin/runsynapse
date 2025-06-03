import type { PageServerLoad, Actions } from './$types';
import { redirect, fail } from '@sveltejs/kit';
import { GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET } from '$env/static/private';

import { ProjectsClient } from '@google-cloud/resource-manager';

export const load: PageServerLoad = async ({ locals: { supabase, user } }) => {
    const userID = user?.id;
    if (!userID) {
        console.error('User ID not found in session');
        throw redirect(302, '/auth/login');
    }

    // get google refresh token
    const { data: userData, error: userDataError } = await supabase
        .from('users')
        .select('google_refresh_token')
        .eq('user_id', userID)
        .single();
    
    if (userDataError) {
        console.error('Error fetching user data:', userDataError);
        throw redirect(302, '/settings?error=failed to fetch user data');
    }

    if (!userData?.google_refresh_token) {
        console.error('No Google refresh token found for user:', userID);
        throw redirect(302, '/settings?error=no google refresh token found');
    }

    const [gcpResult, reposResult] = await Promise.allSettled([
        fetchGCPProjects(userData.google_refresh_token),
        fetchConnectedRepos(supabase, userID)
    ]);

    let gcpProjects = [];
    let error = null;

    if (gcpResult.status === 'fulfilled') {
        gcpProjects = gcpResult.value;
    } else {
        console.error('Error fetching GCP projects:', gcpResult.reason);
        const gcpError = gcpResult.reason;
        if (gcpError.code === 400 && gcpError.message?.includes('invalid_grant')) {
            error = "Invalid Google credentials. Please re-authenticate.";
        } else {
            error = 'Failed to fetch GCP projects';
        }
    }

    let connectedRepos = [];
    let installationID = null;

    if (reposResult.status === 'fulfilled') {
        connectedRepos = reposResult.value.repos;
        installationID = reposResult.value.installationID;
    } else {
        console.error('Error fetching GitHub repos:', reposResult.reason);
        // dont fail whole page if repo loading fails 
    }

    return {
        projects: gcpProjects,
        connectedRepos,
        installationID,
        error
    };
};

export const actions: Actions = {
    connectRepo: async ({ request, locals: { supabase, user } }) => {
        const userID = user?.id;
        if (!userID) {
            console.error('User ID not found in session');
            throw redirect(302, '/auth/login');
        }

        const data = await request.formData();
        
        const repoID = data.get('repoID')?.toString();
        const autoDeployMain = data.get('autoDeployMain') === "true";
        const enablePreviewDeployments = data.get('enablePreviewDeployments') === "true";
        const gcpProjectID = data.get('selectedProject')?.toString();
        const dockerfilePath = data.get('dockerfilePath')?.toString() || './Dockerfile';
        const serviceName = data.get('serviceName')?.toString() || 'default-service';

        if (!repoID || !gcpProjectID) {
            console.error('Missing required fields for connecting repo');
            return fail(400, {
                success: false,
                message: 'Missing required fields'
            });
        }

        const { error: updateError } = await supabase
            .from('github_repositories')
            .update({
                gcp_project_id: gcpProjectID,
                auto_deploy_main: autoDeployMain,
                enable_preview_deployments: enablePreviewDeployments,
                dockerfile_path: dockerfilePath,
                service_name: serviceName
            })
            .eq('id', repoID)   // remember - we are using internal repo ID, not the repo ID from GitHub
            .select()
            .single();

        if (updateError) {
            console.error('Error updating repository:', updateError);
            return fail(500, {
                success: false,
                message: 'Failed to connect repository'
            });
        }

        return {
            success: true,
        };

    }
}

// returning GCP projects is lowkey an overfetch but its fast enough to not worry about yet
// if initial load times are slow we can do this on-demand instead of in the load function
async function fetchGCPProjects(refreshToken: string): Promise<any[]> {
    const credentials = {
        type: 'authorized_user',
        client_id: GOOGLE_CLIENT_ID,
        client_secret: GOOGLE_CLIENT_SECRET,
        refresh_token: refreshToken
    };

    const projectsClient = new ProjectsClient({
        credentials: credentials,
    });

    const [projects] = await projectsClient.searchProjects();

    return projects.map((project: any) => ({
        projectID: project.projectId,
        state: project.state,
        displayName: project.displayName,
    }));
}

async function fetchConnectedRepos(supabase: any, userID: string): Promise<{ repos: any[], installationID: number | null }> {
    const { data: githubInstallations, error } = await supabase
        .from('github_installations')
        .select(`
            installation_id,
            installed_at,
            github_repositories (
                id,
                repo_name,
                is_private,
                dockerfile_path,
                created_at,
                gcp_project_id,
                auto_deploy_main,
                enable_preview_deployments,
                service_name
            )
        `)
        .eq('user_id', userID);

    if (error) throw error;

    const connectedRepos = githubInstallations?.flatMap((install: any) => install.github_repositories) || [];
    const installationID = githubInstallations?.[0]?.installation_id || null;

    return {
        repos: connectedRepos.map((repo: any) => ({
            ...repo,
            connected: !!repo.gcp_project_id
        })),
        installationID
    };
}