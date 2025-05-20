import { SERVICE_API_HOST } from '$env/static/private';
import {fail} from "@sveltejs/kit";

export const actions ={
    signup: async({request}) => {
        const formData = await request.formData();
        const email = formData.get('email');
        const name = formData.get('name');
        const password = formData.get('password');

        const response = await fetch(`http://${SERVICE_API_HOST}/api/auth/register`, {
            method: 'POST',
            body: JSON.stringify({ email, password, name }),
        });

        const data = await response.json();

        if (data.status != 200) {
            return fail(data.status, { 
                email, 
                error: data.message ?? 'Signup failed. Please try again.'
            });
        }

        return { success: true, redirectTo: '/login' };
    }
}