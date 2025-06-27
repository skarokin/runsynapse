import type { PageServerLoad } from './$types';
import { redirect } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ locals: { supabase, user } }) => {
    const userID = user?.id;
    if (!userID) {
        console.error('User ID not found in session');
        throw redirect(302, '/auth/login');
    }

    // idk do something
};