import { redirect } from '@sveltejs/kit';
import crypto from 'crypto';
import type { RequestHandler } from './$types';

export const GET: RequestHandler = async ({ locals: { supabase, safeGetSession } }) => {
    const userID = await safeGetSession().then((session) => session?.user?.id);

    if (!userID) {
        throw redirect(302, '/auth/login');
    }

    // to associate an installation with a user, we generate a state token
    // when github redirects back to us, we can lookup state token in DB and find the userID associated
    const state = crypto.randomBytes(32).toString('hex');

    await supabase.from('github_connection_states').insert({
        state_token: state,
        user_id: userID,
        expires_at: new Date(Date.now() + 10 * 60 * 1000) // 10 min  expiration
    })

    const githubUrl = `https://github.com/apps/runsynapse/installations/new?state=${state}`;
    
    throw redirect(302, githubUrl);
};