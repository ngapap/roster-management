import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [sveltekit()],
	server: {
		// cors: true,
		// proxy: {
		// 	// Proxy requests starting with /api to your backend server
		// 	'/api': {
		// 		target: process.env.SERVICE_API_HOST ? `http://${process.env.SERVICE_API_HOST}` : 'http://localhost:9001',
		// 		changeOrigin: true,
		// 		rewrite: path => path.replace(/^\/api/, '')
		// 	}
		// }
	}
});
