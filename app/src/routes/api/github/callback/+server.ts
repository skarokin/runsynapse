import { redirect } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { App } from "octokit";
import { GITHUB_APP_ID, GITHUB_PRIVATE_KEY } from '$env/static/private';

export const GET: RequestHandler = async ({ url, locals: { supabase } }) => {
    const installationId = url.searchParams.get('installation_id');
    const setupAction = url.searchParams.get('setup_action');
    const state = url.searchParams.get('state');

    if (setupAction !== 'install' || !installationId || !state) {
        throw redirect(302, '/settings?github=cancelled');
    }

    // get user ID from state token
    const { data: stateRecord, error: stateError } = await supabase
        .from('github_connection_states')
        .select('user_id')
        .eq('state_token', state)
        .gt('expires_at', new Date().toISOString())
        .single();

    if (stateError || !stateRecord?.user_id) {
        throw redirect(302, '/settings?github=error&reason=invalid_state');
    }

    // fetch repositories accessible to the installation
    const app = new App({
        appId: GITHUB_APP_ID,
        privateKey: GITHUB_PRIVATE_KEY.replace(/\\n/g, '\n'),
    });

    const octokit = await app.getInstallationOctokit(Number(installationId));
    const repos = await octokit.rest.apps.listReposAccessibleToInstallation();

    // upsert installation, delete old repos, insert new repos
    const { data: result, error: rpcError } = await supabase
        .rpc('handle_github_installation', {
            p_user_id: stateRecord.user_id,
            p_installation_id: Number(installationId),
            p_repositories: repos.data.repositories
        });

    if (rpcError || !result?.success) {
        console.error('RPC error:', rpcError || result?.error);
        throw redirect(302, '/settings?github=error&reason=db_error');
    }

    console.log('done!')

    throw redirect(302, '/settings?github=connected');
};