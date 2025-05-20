import { SERVICE_API_HOST } from '$env/static/private';
import { fail } from '@sveltejs/kit';

export const actions = {
    login: async ({ request, cookies, fetch }) => {
        const formData = await request.formData();
        const email = formData.get('email');
        const password = formData.get('password');
        
        try {
            const response = await fetch(`http://${SERVICE_API_HOST}/api/auth/login`, {
                method: 'POST',
                body: JSON.stringify({ email, password }),
                headers: {
                    'Content-Type': 'application/json'
                }
            });

            const data = await response.json();
            if (data.status !== 200) {
                return fail(data.status, { 
                    email, 
                    error: data.message ?? 'Login failed. Please try again.'
                });
            }

            // Calculate maxAge in seconds
            const expiresAt = new Date(data.data.expires_at);
            const maxAge = Math.floor((expiresAt - new Date()) / 1000);
            
            if (maxAge <= 0) {
                console.error('Invalid expiration time:', data.data.expires_at);
                return fail(500, {
                    email,
                    error: 'Invalid session expiration time'
                });
            }

            // Set auth token in a secure HTTP-only cookie
            cookies.set('token', data.data.token, {
                path: '/',
                // maxAge: maxAge,
                httpOnly: true,
                secure: process.env.NODE_ENV === 'production',
                sameSite: 'lax'
            });
            
            var isAdmin = data.data.user["is_admin"];
            
            // Set user data in a regular cookie for client access
            cookies.set('user', JSON.stringify({
                id: data.data.user.id,
                name: data.data.user.name,
                email: data.data.user.email,
                isAuthenticated: true,
                isAdmin: isAdmin,
            }), {
                path: '/',
                maxAge: maxAge,
                httpOnly: true,
                secure: process.env.NODE_ENV === 'production',
                sameSite: 'lax'
            });

            let redirectTo = '/assignment';
            if (isAdmin) {
               redirectTo = '/admin';
            } 

            return { success: true, redirectTo: redirectTo};
        } catch (error) {
            console.error('Login error:', error);
            return fail(500, {
                email,
                error: 'An unexpected error occurred. Please try again.'
            });
        }
    }
}