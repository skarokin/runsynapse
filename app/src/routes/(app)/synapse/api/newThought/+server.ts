import type { RequestHandler } from '@sveltejs/kit';
import { GO_API_ENDPOINT } from '$env/static/private';

export const POST: RequestHandler = async ({ request, locals: { user } }) => {
    if (!user) {
        return new Response(JSON.stringify({ error: 'Unauthorized' }), { 
            status: 401,
            headers: { 'Content-Type': 'application/json' }
        });
    }

    const formData = await request.formData();
    
    formData.append('user_id', user.id);

    try {
        const res = await fetch(`${GO_API_ENDPOINT}/newThought`, {
            method: 'POST',
            body: formData
        });

        const result = await res.json();
        
        return new Response(JSON.stringify(result), {
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