import type { PageServerLoad, Actions } from './$types';
import { redirect } from '@sveltejs/kit';
import { GO_API_ENDPOINT, GO_API_KEY } from '$env/static/private';

import type { Database } from '$lib/database.types';
type Thought = Database['public']['Tables']['user_thoughts']['Row'];

/*
 * remember - using API Gateway for Go lambda so dont forget the API key in any backend requests
 */

export const load: PageServerLoad = async ({ locals: { supabase, user } }) => {
    if (!user) {
        console.error('User not found in session');
        throw redirect(302, '/auth/login');
    }

    const res = await fetch(`${GO_API_ENDPOINT}/loadFunction`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            // 'X-API-Key': GO_API_KEY
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

    console.log(data)

    return {
        thoughts: data.thoughts,
        thoughtSet: new Set<string>(data.thoughts.map((thought: Thought) => thought.id)),
        pinnedThoughts: data.pinned_thoughts || [],
        moreAbove: data.more_above || false,
        moreBelow: data.more_below || false,
    }
};

// gotoPin, jumptoBottom, and moreThoughts are handled in $lib/api
export const actions: Actions = {
    newThought: async ({ request, locals: { supabase, user } }) => {
        if (!user) {
            console.error('User not found in session');
            throw redirect(302, '/auth/login');
        }

        const formData = await request.formData();
        // if any file is greater than 10MB, and there are more than 5 files, error
        const files = formData.getAll('files') as File[];
        if (files.length > 5) {
            console.error('Too many files uploaded');
            return { error: 'Too many files uploaded' };
        }
        for (const file of files) {
            if (file.size > 10 * 1024 * 1024) { //
                console.error('File too large:', file.name);
                return { error: 'File too large: ' + file.name };
            }
        }

        formData.append('user_id', user.id);

        try {
            const res = await fetch(`${GO_API_ENDPOINT}/newThought`, {
                method: 'POST',
                headers: {
                    // 'X-API-Key': GO_API_KEY
                },
                body: formData
            });

            if (!res.ok) {
                console.error('Failed to create new thought:', res.status, res.statusText);
                return { error: 'Failed to create new thought' };
            }

            const data = await res.json();
            console.log('data:', data);

            return { success: true, thought: data.thought };
        } catch (error) {
            console.error('Error creating new thought:', error);
            return { error: 'Error creating new thought' };
        }

    },
    deleteThought: async ({ request, locals: { supabase, user } }) => {

    },
    pinThought: async ({ request, locals: { supabase, user } }) => {

    },
    unpinThought: async ({ request, locals: { supabase, user } }) => {

    },
};