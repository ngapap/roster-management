<script>
  import { authStore } from '$lib/stores/authStore.js';
  import { enhance } from '$app/forms';
  import { goto } from '$app/navigation';

  let name = $state('');
  let email = $state('');
  let password = $state('');
  let confirmPassword = $state('');

  let {form} = $props();
  let isLoading = $state(false);
  let showSuccessModal = $state(false);
  
  let passwordsMatch = $derived(password === confirmPassword);
  let isFormValid = $derived(name.length > 0 && email.includes('@') && password.length >= 6 && passwordsMatch);

</script>

<div class="max-w-md mx-auto">
  <div class="text-center mb-8">
    <h1 class="text-2xl font-bold">Create Account</h1>
    <p class="text-muted-foreground mt-2">Sign up to start managing your shifts</p>
  </div>
  
  <form action="?/signup" method="post" class="bg-background/80 backdrop-blur-sm rounded-card shadow-card p-6" use:enhance={() => {
    isLoading = true;
    return async ({ update, result }) => {
      await update();
      isLoading = false;
      console.log(result);
      if (result.data?.success) {
        showSuccessModal = true;
      }
    };
  }}>
  
    {#if form?.error}
      <div class="mb-4 p-3 bg-red-50 border border-red-200 text-red-700 rounded-md text-sm">
        {form.error}
      </div>
    {/if}
    
    <div class="mb-4">
      <label for="name" class="block text-sm font-medium mb-1">Full Name</label>
      <input 
        type="text" 
        id="name"
        name="name"
        bind:value={name}
        placeholder="John Doe"
        class="w-full rounded-input border-border-input focus:border-border-input-hover p-2 text-sm"
        required
      />
    </div>
    
    <div class="mb-4">
      <label for="email" class="block text-sm font-medium mb-1">Email</label>
      <input 
        type="email" 
        id="email"
        name="email"
        bind:value={email}
        placeholder="you@example.com"
        class="w-full rounded-input border-border-input focus:border-border-input-hover p-2 text-sm"
        required
      />
    </div>
    
    <div class="mb-4">
      <label for="password" class="block text-sm font-medium mb-1">Password</label>
      <input 
        type="password" 
        id="password"
        name="password"
        bind:value={password}
        placeholder="••••••••"
        class="w-full rounded-input border-border-input focus:border-border-input-hover p-2 text-sm"
        required
        minlength="6"
      />
      <p class="text-xs text-muted-foreground mt-1">Password must be at least 6 characters</p>
    </div>
    
    <div class="mb-6">
      <label for="confirmPassword" class="block text-sm font-medium mb-1">Confirm Password</label>
      <input 
        type="password" 
        id="confirmPassword" 
        bind:value={confirmPassword}
        placeholder="••••••••"
        class="w-full rounded-input border-border-input focus:border-border-input-hover p-2 text-sm"
        class:border-red-300={confirmPassword && !passwordsMatch}
        required
      />
      {#if confirmPassword && !passwordsMatch}
        <p class="text-xs text-red-600 mt-1">Passwords don't match</p>
      {/if}
    </div>
    
    <button 
      type="submit" 
      class="w-full bg-blue-500 hover:bg-blue-600 text-white font-medium py-2 px-4 rounded-md disabled:opacity-50 disabled:cursor-not-allowed"
      disabled={!isFormValid || isLoading}
    >
      {isLoading ? 'Creating account...' : 'Create Account'}
    </button>
    
    <div class="mt-4 text-center text-sm">
      <span class="text-muted-foreground">Already have an account?</span>
      <a href="/login" class="text-blue-600 hover:underline ml-1">Sign in</a>
    </div>
  </form>

  {#if showSuccessModal}
    <div class="fixed inset-0 flex items-center justify-center bg-opacity-100 backdrop-blur-sm">
      <div class="bg-white p-6 rounded-lg shadow-lg max-w-md w-full">
        <h2 class="text-xl font-bold mb-4">Signup Successful!</h2>
        <p class="mb-4">Your account has been created successfully.</p>
        <button onclick={() => goto('/login')} class="w-full bg-blue-500 hover:bg-blue-600 text-white font-medium py-2 px-4 rounded-md">
          Go to Login
        </button>
      </div>
    </div>
  {/if}
</div> 