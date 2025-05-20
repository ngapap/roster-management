
import { redirect } from '@sveltejs/kit';
import { publicRoutes } from '$lib/publicRoutes';

/** @type {import('@sveltejs/kit').Handle} */
export async function handle({ event, resolve }) {
  const { pathname } = event.url;

  // Check if the route is protected
  const isProtected = !publicRoutes.includes(pathname);

  if (isProtected) {
    const token = event.cookies.get('token');

    if (!token) {
      throw redirect(303, '/login');
    }

    event.locals.user = event.cookies.get('user');

  }

  // Proceed with the request
  return resolve(event);
}
