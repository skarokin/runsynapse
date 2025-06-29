import type { RequestHandler } from '@sveltejs/kit';
import { GO_API_ENDPOINT, GO_API_KEY } from '$env/static/private';

export const POST: RequestHandler = async ({ request, locals: { user } }) => {
    if (!user) {
        return new Response(JSON.stringify({ error: 'Unauthorized' }), { 
            status: 401,
            headers: { 'Content-Type': 'application/json' }
        });
    }

    const { lastThoughtID } = await request.json();

    const res = await fetch(`${GO_API_ENDPOINT}/loadThoughts`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${GO_API_KEY}`
        },
        body: JSON.stringify({
            user_id: user.id,
            cursor: lastThoughtID
        })
    });

    if (!res.ok) {
        return new Response(JSON.stringify({ error: 'Failed to load more thoughts' }), {
            status: 500,
            headers: { 'Content-Type': 'application/json' }
        });
    }

    const result = await res.json();

    console.log('Load more thoughts result:', result);
    
    return new Response(JSON.stringify(result), {
        status: 200,
        headers: { 'Content-Type': 'application/json' }
    });
}