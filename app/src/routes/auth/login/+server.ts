import { redirect } from '@sveltejs/kit'

import type { RequestHandler } from './$types';

export const GET: RequestHandler = async ({ url, locals: { supabase } }) => {
    const { data, error } = await supabase.auth.signInWithOAuth({
        provider: 'google',
        options: {
            redirectTo: `${url.origin}/auth/callback`,
            scopes: 'openid email profile https://www.googleapis.com/auth/cloudplatformprojects',
            queryParams: {
                access_type: 'offline', // required for refresh token
                prompt: 'consent',      // forces consent screen to get refresh token
            }
        }
    })

    if (error) {
        redirect(303, '/?error=' + encodeURIComponent("Failed to authenticate: " + error.message))
    }

    if (data.url) {
        redirect(303, data.url)
    }

    redirect(303, '/?error=' + encodeURIComponent('Authentication failed'))
}