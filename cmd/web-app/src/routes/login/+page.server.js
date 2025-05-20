import { SERVICE_API_HOST } from '$env/static/private';
import { fail, redirect } from '@sveltejs/kit';
import { goto } from '$app/navigation';

export const actions = {
    login: async ({ request, cookies }) => {
        const formData = await request.formData();
        const email = formData.get('email');
        const password = formData.get('password');
        
        const response = await fetch(`http://${SERVICE_API_HOST}/api/auth/login`, {
            method: 'POST',
            body: JSON.stringify({ email, password }),
        });
        
        if (response.status != 200) {
            return fail(response.status, { 
                email, 
                error: response.statusText ?? 'Login failed. Please try again.'
            });
        }

        const data = await response.json();
        // console.log("data", data);
        // // Set auth token in a secure HTTP-only cookie
        cookies.set('token', data.token, {
            path: '/',
            maxAge: new Date(data.expires_at) - new Date()
        });
        
        // // Set user data in a regular cookie for client access
        cookies.set('user', JSON.stringify({
            id: data.user.id,
            name: data.user.name,
            email: data.user.email,
            isAuthenticated: true,
            isAdmin: data.user.role === 'worker' ? false : true
        }), {
            path: '/',
        });

        // Return success value for client-side handling
        // throw redirect(303, '/');
        return { success: true, redirectTo: '/' };
    }
}