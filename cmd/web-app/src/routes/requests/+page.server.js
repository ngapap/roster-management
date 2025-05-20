import { SERVICE_API_HOST } from '$env/static/private';
import { logout } from '$lib/stores/authStore.js';
import { fail } from '@sveltejs/kit';

export async function load({ locals, cookies }) {

    const user = locals.user;
    
    try {
        const [requestsResponse] = await Promise.all([
            fetch(`http://${SERVICE_API_HOST}/api/shift-request/worker/${user.id}`, {
                headers: {
                    'Authorization': `Bearer ${locals.token}`
                }
            })
        ]);

        // If unauthorized, clear the token
        if (requestsResponse.status === 401) {
            cookies.delete('token', { path: '/' });
            cookies.delete('user', { path: '/' });
            return {
                user,
                shiftsRequest: [],
                status: 401,
                error: 'Session expired. Please login again.'
            };
        }

        const [requestsData] = await Promise.all([
            requestsResponse.json()
        ]);

        return {
            user,
            shiftsRequest: requestsData.data,
            status: requestsResponse.status,
            error: requestsData.status !== 200 ? (requestsData.message ?? 'Unknown error.') : null
        };
    } catch (error) {
        console.error('Error fetching data:', error);
        return {
            user,
            shiftsRequest: [],
            status: 500,
            error: 'Failed to fetch data'
        };
    }
}

export const actions = {

};
