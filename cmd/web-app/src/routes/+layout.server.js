import { redirect } from '@sveltejs/kit';
import { publicRoutes } from '$lib/publicRoutes';
import { get } from 'svelte/store';

export function load({ cookies, url }) {
	const user = JSON.parse(cookies.get('user'));
  
  const requiresAuth = !publicRoutes.includes(url.pathname);

  if (requiresAuth && !user.isAuthenticated) {
    throw redirect(303, '/login');
  }
  
  if (user.isAuthenticated && publicRoutes.includes(url.pathname)) {
    if (user.isAdmin) {
      
      throw redirect(303, '/admin');
    } else {
      throw redirect(303, '/');
    }
  }

  if (user.isAuthenticated && url.pathname === '/' && user.isAdmin) {
    throw redirect(303, '/admin');
  }

	return {
		 user: user,
     redirectTo: '/'
	};
}
