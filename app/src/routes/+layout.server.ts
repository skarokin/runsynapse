import { redirect } from '@sveltejs/kit'
import type { LayoutServerLoad } from './$types'

export const load: LayoutServerLoad = async ({ locals: { safeGetSession }, url }) => {
    const { session, user } = await safeGetSession()
    
    if (session && user && url.pathname === '/') {
        redirect(303, '/console')
    }
    
    return {
        session,
        user,
        userMetadata: user?.user_metadata || null,
    }
}