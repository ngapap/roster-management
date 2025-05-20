import { redirect } from '@sveltejs/kit';
import { publicRoutes } from '$lib/publicRoutes';


export function load({ cookies, url }) {

    const data = cookies.get('user');
    let user = {isAuthenticated: false};
    if (data) {
        user = JSON.parse(data);
    }
    
    const requiresAuth = !publicRoutes.includes(url.pathname);

      if (requiresAuth && !user?.isAuthenticated) {
        throw redirect(303, '/login');
      }

      if (user.isAuthenticated && publicRoutes.includes(url.pathname)) {
        if (user.isAdmin) {
          throw redirect(303, '/admin');
        } else {
          throw redirect(303, '/assignment');
        }
      }

      if (user.isAuthenticated && url.pathname === '/' && user.isAdmin) {
        throw redirect(303, '/admin');
      }

      let redirectTo  = '/assignment';

      if (user.isAdmin) {
        redirectTo = '/admin';
      }

    return {
         user: user,
         redirectTo: redirectTo
    };
}

