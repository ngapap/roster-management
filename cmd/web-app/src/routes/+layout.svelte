<script lang="ts">
	import '../app.css';
	import { currentUser, users } from '$lib/stores/shiftStore.js';
	// import { authStore, logout, initAuth } from '$lib/stores/authStore.js';
	import { onMount } from 'svelte';
	
	let { data, children } = $props();	
	
	function switchUser(userId: number) {
		const selectedUser = $users.find(u => u.id === userId);
		if (selectedUser) {
			currentUser.set(selectedUser);
		}
	}
</script>

<div class="min-h-screen bg-background-alt">
	<header class="bg-background shadow-sm">
		<div class="container mx-auto px-4 py-4">
			<div class="flex justify-between items-center">
				<h1 class="text-xl font-bold">Roster Management</h1>
				
				{#if data.user.isAuthenticated}
					<div class="flex items-center gap-4">
						<div class="text-sm">
							Logged in as: <span class="font-medium">{$currentUser.name}</span>
							<span class="ml-2 px-2 py-0.5 bg-blue-100 text-blue-800 rounded-full text-xs font-medium">
								{$currentUser.role}
							</span>
						</div>
						
						<button 

							class="text-sm text-red-600 hover:text-red-800 font-medium"
						>
							Logout
						</button>
					</div>
				{:else}
					<div class="flex items-center gap-4">
						<a href="/login" class="text-sm text-muted-foreground hover:text-foreground font-medium">Login</a>
						<a href="/signup" class="bg-blue-500 hover:bg-blue-600 text-white font-medium py-1 px-3 rounded-md text-sm">
							Sign Up
						</a>
					</div>
				{/if}
			</div>
			
			{#if data.user.isAuthenticated}
				<nav class="mt-4 flex gap-4">
					{#if $currentUser.role === 'admin'}
						<a href="/admin" class="text-muted-foreground hover:text-foreground font-medium">Admin Panel</a>
					{:else}
						<a href="/" class="text-muted-foreground hover:text-foreground font-medium">My Schedule</a>
					{/if}
					<a href="/available" class="text-muted-foreground hover:text-foreground font-medium">Available Shifts</a>
					<a href="/requests" class="text-muted-foreground hover:text-foreground font-medium">My Requests</a>
				</nav>
			{/if}
		</div>
	</header>
	
	<main class="container mx-auto px-4 py-8">
		{@render children?.()}
	</main>
</div>
