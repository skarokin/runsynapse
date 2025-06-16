import type { PageServerLoad, Actions } from './$types';
import { redirect } from '@sveltejs/kit';
import { GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET } from '$env/static/private';
import { fetchCloudRunLogs, fetchCloudRunMetrics, enableCloudRunAPI } from '$lib/server/cloud_run/cloudRun.server';

// load the first page of logs initially
export const load: PageServerLoad = async ({ params, url, locals: { supabase, user } }) => {
    const userID = user?.id;
    if (!userID) {
        console.error('User ID not found in session');
        throw redirect(302, '/auth/login');
    }

    const repoName = decodeURIComponent(params.repoName);

    let maxAgeDays = url.searchParams.get('maxAgeDays') || null;
    if (!maxAgeDays) {
        maxAgeDays = new Date(Date.now() - 30 * 24 * 60 * 60 * 1000).toISOString();
    }

    // FUTURE - GET BRANCH NAME FROM QUERY PARAMS AND DO A 
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

    // get user's refresh token
    const { data: userData, error: userDataError } = await supabase
        .from('users')
        .select('google_refresh_token')
        .eq('user_id', userID)
        .single();

    if (userDataError || !userData?.google_refresh_token) {
        console.error('No Google refresh token found for user:', userID);
        return {
            repository,
            logs: [],
            metrics: null,
            nextPageToken: null,
            prevPageToken: null,
            error: 'No Google authentication found'
        };
    }

    let logs: any[] = [];
    let encodedNextPageToken = null;
    let metrics = null;
    let error = null;

    if (repository.gcp_project_id && repository.service_name) {
        try {
            const credentials = await getGoogleCredentials(userData.google_refresh_token);

            const [logsResult, metricsResult] = await Promise.allSettled([
                fetchCloudRunLogs(credentials, repository.gcp_project_id, repository.service_name, repository.branch_name || 'main', "", maxAgeDays),
                fetchCloudRunMetrics(credentials, repository.gcp_project_id, repository.service_name)
            ]);

            if (logsResult.status === 'fulfilled') {
                logs = logsResult.value.logs;
                if (logsResult.value.nextPageToken) {
                    encodedNextPageToken = Buffer.from(logsResult.value.nextPageToken, 'binary').toString('base64');
                }
                maxAgeDays = logsResult.value.timeFilter || maxAgeDays; // update maxAgeDays if provided
            } else {
                console.error('Error fetching logs:', logsResult.reason);
            }

            if (metricsResult.status === 'fulfilled') {
                metrics = metricsResult.value;
            } else {
                console.error('Error fetching metrics:', metricsResult.reason);
                // error - cloud run admin api not enabled
                if (metricsResult.reason.code === 7 && metricsResult.reason.errorInfoMetadata.serviceTitle === "Cloud Run Admin API") {
                    console.warn('Cloud Run API not enabled, attempting to enable it...');

                    try {
                        await enableCloudRunAPI(credentials, repository.gcp_project_id);
                    } catch (enableError: any) {
                        if (enableError.reason.code === 9 && enableError.reason.reason === "UREQ_PROJECT_BILLING_NOT_OPEN") {
                            console.warn('Billing not enabled, warn user');
                            error = `Billing not enabled for the project. Please go to ${enableError.reason.errorInfoMetadata.activationURL} to enable billing.`;
                        } else {
                            console.error('Failed to enable Cloud Run API:', enableError);
                            error = 'Failed to enable Cloud Run API. Please ensure you have the necessary permissions and billing is enabled for the project.';
                        }
                    }

                } else if (metricsResult.reason.code === 9 && metricsResult.reason.reason === "UREQ_PROJECT_BILLING_NOT_OPEN") {
                    console.warn('Billing not enabled, warn user');
                    const billingUrl = metricsResult.reason.errorInfoMetadata.activationURL;
                    error = `Billing not enabled for the project. Please go to ${billingUrl} to enable billing.`;
                } else {
                    error = 'Failed to fetch Cloud Run metrics';
                }
            }
        } catch (err: any) {
            console.error('Error fetching Cloud Run data:', err);
            error = 'Failed to fetch Cloud Run data';
        }
    }

    return {
        repository,
        logs,
        metrics,
        nextPageToken: encodedNextPageToken,
        maxAgeDays,
        error
    };
};

async function getGoogleCredentials(refreshToken: string) {
    return {
        type: 'authorized_user',
        client_id: GOOGLE_CLIENT_ID,
        client_secret: GOOGLE_CLIENT_SECRET,
        refresh_token: refreshToken
    };
}

export const actions: Actions = {
    // literally just a copy paste of load function but listens for a next page token in the request body
    loadMore: async ({ params, request, locals: { supabase, user } }) => {
        const userID = user?.id;
        if (!userID) {
            throw redirect(302, '/auth/login');
        }

        const repoName = decodeURIComponent(params.repoName);
        const formData = await request.formData();
        const encodedPage = formData.get('page') as string;
        const maxAgeDays = formData.get('maxAgeDays') as string;

        let page = null;
        if (encodedPage) {
            try {
                page = Buffer.from(encodedPage, 'base64').toString('binary');
            } catch (e) {
                return { error: 'Invalid page token' }; 
            }
        }

        try {
            const { data: repository } = await supabase
                .from('github_repositories')
                .select('*')
                .eq('repo_name', repoName)
                .eq('user_id', userID)
                .single();

            if (!repository) {
                return { error: 'Repository not found' };
            }

            const { data: userData } = await supabase
                .from('users')
                .select('google_refresh_token')
                .eq('user_id', userID)
                .single();

            if (!userData?.google_refresh_token) {
                return { error: 'Google authentication required' };
            }

            const credentials = await getGoogleCredentials(userData.google_refresh_token);

            const logsResult = await fetchCloudRunLogs(
                credentials,
                repository.gcp_project_id,
                repository.service_name,
                repository.branch_name || 'main',
                page,
                maxAgeDays
            );

            let encodedNextPageToken = null;
            if (logsResult.nextPageToken) {
                encodedNextPageToken = Buffer.from(logsResult.nextPageToken, 'binary').toString('base64');
            }

            return {
                logs: logsResult.logs,
                nextPageToken: encodedNextPageToken,
                maxAgeDays: logsResult.timeFilter
            };

        } catch (error) {
            console.error('API Error:', error);
            return { error: 'Failed to fetch logs' }
        }
    }
}