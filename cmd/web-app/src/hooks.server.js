import { redirect } from '@sveltejs/kit';
import { publicRoutes } from '$lib/publicRoutes';
import { sequence } from '@sveltejs/kit/hooks';

const cors = async ({ event, resolve }) => {
    // Skip CORS for Chrome DevTools requests
    if (event.url.pathname.startsWith('/.well-known/')) {
        return resolve(event);
    }

    const response = await resolve(event);
    response.headers.set('Access-Control-Allow-Origin', event.request.headers.get('origin') || '*');
    response.headers.set('Access-Control-Allow-Methods', 'GET, POST, PUT, DELETE, OPTIONS');
    response.headers.set('Access-Control-Allow-Headers', 'Content-Type, Authorization');
    response.headers.set('Access-Control-Allow-Credentials', 'true');
    return response;
};

const csrf = async ({ event, resolve }) => {
    // Skip CSRF for Chrome DevTools requests
    if (event.url.pathname.startsWith('/.well-known/')) {
        return resolve(event);
    }

    if (event.request.method === 'POST') {
        const origin = event.request.headers.get('origin');
        const referer = event.request.headers.get('referer');
        
        // Allow same-origin requests
        if (!origin || !referer || new URL(referer).origin === origin) {
            return resolve(event);
        }
        
        // For cross-origin requests, check if it's a form submission
        const contentType = event.request.headers.get('content-type');
        if (contentType?.includes('application/x-www-form-urlencoded')) {
            return resolve(event);
        }
    }
    
    return resolve(event);
};

const auth = async ({ event, resolve }) => {
    // Skip auth for Chrome DevTools requests
    if (event.url.pathname.startsWith('/.well-known/')) {
        return resolve(event);
    }

    const { pathname } = event.url;
    const token = event.cookies.get('token');
    const userCookie = event.cookies.get('user');

    // Check if the route is protected
    const isProtected = !publicRoutes.includes(pathname);

    // If no token and not on login/signup page, redirect to login
    if (!token && !pathname.startsWith('/login') && !pathname.startsWith('/signup')) {
        throw redirect(303, '/login');
    }

    // If token exists but on login/signup page, redirect to home
    if (token && (pathname.startsWith('/login') || pathname.startsWith('/signup'))) {
        throw redirect(303, '/');
    }

    if (isProtected && userCookie) {
        try {
            event.locals.user = JSON.parse(userCookie);
            event.locals.token = token;
        } catch (e) {
            console.error('Failed to parse user cookie:', e);
        }
    }

    return resolve(event);
};

export const handle = sequence(cors, auth);
