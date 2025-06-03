import type { PageServerLoad, Actions } from './$types';
import { redirect, fail } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ params, locals: { supabase, user } }) => {
    const userID = user?.id;
    if (!userID) {
        console.error('User ID not found in session');
        throw redirect(302, '/auth/login');
    }

    const repoName = decodeURIComponent(params.repoName);

    const { data: repository, error: repoError } = await supabase
        .from('github_repositories')
        .select('*')
        .eq('repo_name', repoName)
        .eq('user_id', userID)
        .single();

    if (repoError || !repository) {
        console.error('Error fetching repository:', repoError);
        throw redirect(302, '/console?error=repository not found');
    }

    return {
        repository
    };
};

export const actions: Actions = {
    updateRepository: async ({ request, params, locals: { supabase, user } }) => {
        const userID = user?.id;
        if (!userID) {
            throw redirect(302, '/auth/login');
        }

        const repoName = decodeURIComponent(params.repoName);
        const data = await request.formData();
        
        const autoDeployMain = data.get('autoDeployMain') === "true";
        const enablePreviewDeployments = data.get('enablePreviewDeployments') === "true";
        const dockerfilePath = data.get('dockerfilePath')?.toString() || './Dockerfile';
        const serviceName = data.get('serviceName')?.toString();

        try {
            const { error: updateError } = await supabase
                .from('github_repositories')
                .update({
                    auto_deploy_main: autoDeployMain,
                    enable_preview_deployments: enablePreviewDeployments,
                    dockerfile_path: dockerfilePath,
                    service_name: serviceName
                })
                .eq('repo_name', repoName)

            if (updateError) {
                throw updateError;
            }

            return {
                success: true,
                message: 'Repository settings updated successfully!'
            };

        } catch (error) {
            console.error('Error updating repository:', error);
            return fail(500, {
                success: false,
                message: 'Failed to update repository settings'
            });
        }
    }
};