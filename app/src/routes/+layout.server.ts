import { redirect } from '@sveltejs/kit'
import type { LayoutServerLoad } from './$types'

export const load: LayoutServerLoad = async ({ locals: { user }, url }) => {
    if (user && url.pathname === '/') {
        redirect(303, '/synapse')
    }
    
    return {
        user,
        userMetadata: user?.user_metadata || null,
    }
}