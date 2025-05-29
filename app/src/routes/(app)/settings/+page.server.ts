import type { PageServerLoad } from './$types';
import { redirect } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ locals: { supabase, safeGetSession } }) => {
    const session = await safeGetSession();

    if (!session?.user) {
        throw redirect(302, '/auth/login');
    }

    // get github installations w/ repositories
    const { data: githubInstallations, error } = await supabase
        .from('github_installations')
        .select(`
            installation_id,
            installed_at,
            github_repositories (
                id,
                repo_id,
                repo_name,
                is_private,
                dockerfile_path,
                created_at
            )
        `)
        .eq('user_id', session.user.id);
    if (error) {
        console.error('Error fetching GitHub installations:', error);
        throw redirect(302, '/settings?github=error&reason=db_error');
    }

    const connectedRepos = githubInstallations?.flatMap(install => install.github_repositories) || [];
    const installationID = githubInstallations?.[0]?.installation_id;

    return {
        isConnected: (githubInstallations?.length || 0) > 0,
        connectedRepos,
        installationID,
        autoDeploy: true,
        onlyMain: false,
    };
};