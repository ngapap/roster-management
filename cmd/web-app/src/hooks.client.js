import { isTokenExpired, handleTokenExpiration } from '$lib/utils/auth';
import { authStore } from '$lib/stores/authStore.js';

export async function handle({ event, resolve }) {
    const { expiresAt } = $authStore;
    
    // Check if token is expired
    if (isTokenExpired(expiresAt)) {
        await handleTokenExpiration();
        return new Response(null, {
            status: 303,
            headers: { Location: '/login' }
        });
    }

    return resolve(event);
} 