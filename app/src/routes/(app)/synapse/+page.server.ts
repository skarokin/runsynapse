import type { PageServerLoad } from './$types';
import { redirect } from '@sveltejs/kit';
import { GO_API_ENDPOINT, GO_API_KEY } from '$env/static/private';

import type { Database } from '$lib/database.types';
type Thought = Database['public']['Tables']['user_thoughts']['Row'];

/*
 * remember - using API Gateway for Go lambda so dont forget the API key in any backend requests
 */

export const load: PageServerLoad = async ({ locals: { user } }) => {
    if (!user) {
        console.error('User not found in session');
        throw redirect(302, '/auth/login');
    }

    const res = await fetch(`${GO_API_ENDPOINT}/loadFunction`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'X-API-Key': GO_API_KEY
        },
        body: JSON.stringify({
            user_id: user.id,
            limit: 25,
            cursor: '',
            order: 'before',
        })
    });

    if (!res.ok) {
        console.error('Failed to load thoughts:', res.status, res.statusText);
        return { error: 'Failed to load thoughts' };
    }

    const data = await res.json();
    if (!data || !Array.isArray(data.thoughts)) {
        console.error('Invalid response format:', data);
        return { error: 'Invalid response format' };
    }
    
    return {
        thoughts: data.thoughts,
        hasMoreAbove: data.more_above || false,
        hasMoreBelow: data.more_below || false,
        thoughtSet: new Set<string>(data.thoughts.map((thought: Thought) => thought.id)),
        pinnedThoughts: data.pinned_thoughts || [],
    }
};