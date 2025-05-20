import { goto } from '$app/navigation';
import { authStore } from '$lib/stores/authStore.js';

export function isTokenExpired(expiresAt) {
    if (!expiresAt) return true;
    return new Date(expiresAt) <= new Date();
}

export async function handleTokenExpiration() {
    // Clear auth store
    authStore.set({
        user: null,
        token: null,
        expiresAt: null,
        authError: 'Session expired. Please login again.'
    });

    // Redirect to login
    await goto('/login');
}

export async function checkAuth(response) {
    if (response.status === 401) {
        await handleTokenExpiration();
        return false;
    }
    return true;
} 