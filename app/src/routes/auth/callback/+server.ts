import { redirect } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

export const GET: RequestHandler = async (event) => {
    const {
        url,
        locals: { supabase }
    } = event;
    const code = url.searchParams.get('code') as string;
    const next = url.searchParams.get('next') ?? '/console';

    if (code) {
        const { data, error } = await supabase.auth.exchangeCodeForSession(code)
        if (!error && data.session) {
            const { user } = await event.locals.safeGetSession()
            // IMPORTANT - WE NEED THE PROVIDER TOKEN, NOT SUPABASE'S REFRESH TOKEN (session.refresh_token)
            const googleRefreshToken = data.session.provider_refresh_token;

            if (googleRefreshToken && user) {
                const { error } = await supabase.from('users').upsert({
                    user_id: user.id,
                    google_refresh_token: googleRefreshToken,
                    updated_at: new Date().toISOString()
                }, {
                    onConflict: 'user_id'  // use user_id for conflict resolution
                });
                
                if (error) {
                    console.error('Upsert error:', error);
                } 

            }

            redirect(303, `/${next.slice(1)}`);
        }
    }

    // return the user to an error page with instructions
    throw redirect(303, '/?error=' + encodeURIComponent('Authentication failed'));
};