import { redirect } from '@sveltejs/kit';

export async function POST({ cookies }) {
    cookies.delete('user', { path: '/' });
    cookies.delete('token', { path: '/' });
    throw redirect(303, '/login');
}
