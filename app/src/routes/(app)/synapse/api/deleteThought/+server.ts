import type { RequestHandler } from '@sveltejs/kit';
import { GO_API_ENDPOINT, GO_API_KEY } from '$env/static/private';

export const POST: RequestHandler = async ({ request, locals: { user } }) => {
    if (!user) {
        return new Response(JSON.stringify({ error: 'Unauthorized' }), { 
            status: 401,
            headers: { 'Content-Type': 'application/json' }
        });
    }

    const { thought_id } = await request.json();

    try {
        const res = await fetch(`${GO_API_ENDPOINT}/deleteThought`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'X-API-Key': GO_API_KEY
            },
            body: JSON.stringify({
                user_id: user.id,
                thought_id: thought_id
            })
        });

        const result = await res.json();

        if (!result.success) {
            console.error('Failed to delete thought:', result);
            return new Response(JSON.stringify({ error: 'Failed to delete thought' }), {
                status: 400,
                headers: { 'Content-Type': 'application/json' }
            });
        }
        
        return new Response(JSON.stringify(result.success), {
            headers: { 'Content-Type': 'application/json' }
        });

    } catch (error) {
        console.error('Error calling Go backend:', error);
        return new Response(JSON.stringify({ error: 'Failed to create thought' }), {
            status: 500,
            headers: { 'Content-Type': 'application/json' }
        });
    }
};