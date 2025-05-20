import { SERVICE_API_HOST } from '$env/static/private';
import { redirect } from '@sveltejs/kit';

export async function load({ locals, cookies }) {
    const user = locals.user;
    console.log("user", locals.token);
    const response = await fetch(`http://${SERVICE_API_HOST}/api/shift/worker/${user.id}`, {
        headers: {
            'Authorization': `Bearer ${locals.token}`
        }
    });
    
    // If unauthorized, clear the token
    if (response.status === 401) {
        cookies.delete('token', { path: '/' });
        cookies.delete('user', { path: '/' });
    }

    const data = await response.json();
    console.log("data", data);

    return {
        user,
        shifts: data.data,
        error: data.status !== 200 ? (data.message ?? 'Unknown error.') : null
    };

}
