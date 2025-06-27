import { redirect } from '@sveltejs/kit'
import type { RequestHandler } from './$types'

export const GET: RequestHandler = async ({ locals: { supabase } }) => {
    const { error } = await supabase.auth.signOut()
    
    if (error) {
        redirect(303, '/?error=' + encodeURIComponent("Failed to log out: " + error.message))
    }
    
    redirect(303, '/')
}