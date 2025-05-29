import type { PageServerLoad } from './$types';
import { redirect } from '@sveltejs/kit';

import { GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET } from '$env/static/private';

import { ProjectsClient } from '@google-cloud/resource-manager';

export const load: PageServerLoad = async ({ locals: { supabase, safeGetSession } }) => {
    const session = await safeGetSession();

    if (!session?.user) {
        throw redirect(302, '/auth/login');
    }

    // get google refresh token
    const { data: userData, error } = await supabase
        .from('users')
        .select('google_refresh_token')
        .eq('user_id', session.user.id)
        .single();
    
    if (error) {
        console.error('Error fetching user data:', error);
        throw redirect(302, '/settings?error=failed to fetch user data');
    }

    if (!userData?.google_refresh_token) {
        console.error('No Google refresh token found for user:', session.user.id);
        throw redirect(302, '/settings?error=no google refresh token found');
    }

    let gcpProjects = [];

    // get gcp projects
    try {
        const credentials = {
            type: 'authorized_user',
            client_id: GOOGLE_CLIENT_ID,
            client_secret: GOOGLE_CLIENT_SECRET,
            refresh_token: userData.google_refresh_token
        };

        const projectsClient = new ProjectsClient({
            credentials: credentials,
        });

        const [projects] = await projectsClient.searchProjects();

        gcpProjects = projects.map((project: any) => {
            return {
                projectID: project.projectId,
                state: project.state,
                displayName: project.displayName,
            }
        });
    } catch (error: any) {
        console.error('Error fetching GCP projects:', error);

        return {
            error: 'Failed to fetch GCP projects'
        };
    }

    return {
        projects: gcpProjects,
    }
};