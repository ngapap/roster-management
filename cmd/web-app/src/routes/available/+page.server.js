import { SERVICE_API_HOST } from '$env/static/private';
import { fail } from '@sveltejs/kit';

export async function load({ locals, cookies }) {
    const user = locals.user;
    
    try {
        const [shiftsResponse, requestsResponse] = await Promise.all([
            fetch(`http://${SERVICE_API_HOST}/api/shift/available`, {
                headers: {
                    'Authorization': `Bearer ${locals.token}`
                }
            }),
            fetch(`http://${SERVICE_API_HOST}/api/shift-request/worker/${user.id}`, {
                headers: {
                    'Authorization': `Bearer ${locals.token}`
                }
            })
        ]);

        // If unauthorized, clear the token
        if (shiftsResponse.status === 401 || requestsResponse.status === 401) {
            cookies.delete('token', { path: '/' });
            cookies.delete('user', { path: '/' });
            return {
                user,
                shifts: [],
                shiftsRequest: [],
                status: 401,
                error: 'Session expired. Please login again.'
            };
        }

        const [shiftsData, requestsData] = await Promise.all([
            shiftsResponse.json(),
            requestsResponse.json()
        ]);

        return {
            user,
            shifts: shiftsData.data,
            shiftsRequest: requestsData.data,
            status: shiftsResponse.status,
            error: shiftsData.status !== 200 ? (shiftsData.message ?? 'Unknown error.') : null
        };
    } catch (error) {
        console.error('Error fetching data:', error);
        return {
            user,
            shifts: [],
            shiftsRequest: [],
            status: 500,
            error: 'Failed to fetch data'
        };
    }
}

export const actions = {
    requestShift: async ({ request, locals, cookies }) => {
        const formData = await request.formData();
        const shiftId = formData.get('shiftId');
        console.log("shiftId", shiftId);
        try {
            const response = await fetch(`http://${SERVICE_API_HOST}/api/shift-request/`, {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${locals.token}`,
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    shift_id: shiftId,
                    worker_id: locals.user.id
                })
            });
            console.log("data", response);
            const data = await response.json();
            
            if (response.status === 401) {
                cookies.delete('token', { path: '/' });
                cookies.delete('user', { path: '/' });
                console.log("Session expired. Please login again.");
                return fail(401, { error: 'Session expired. Please login again.' });
            }

            if (response.status !== 200) {
                return fail(response.status, { error: data.message ?? 'Failed to request shift' });
            }

            return { success: true, data: data.data };
        } catch (error) {
            console.error('Error requesting shift:', error);
            return fail(500, { error: 'Failed to request shift' });
        }
    }
};
