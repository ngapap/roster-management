<script>
  import { authStore, signup } from '$lib/stores/authStore.js';
  
  let name = '';
  let email = '';
  let password = '';
  let confirmPassword = '';
  
  $: passwordsMatch = password === confirmPassword;
  $: isFormValid = name.length > 0 && email.includes('@') && password.length >= 6 && passwordsMatch;
  
  function handleSubmit() {
    if (isFormValid) {
      signup(name, email, password);
    }
  }
</script>

<div class="max-w-md mx-auto">
  <div class="text-center mb-8">
    <h1 class="text-2xl font-bold">Create Account</h1>
    <p class="text-muted-foreground mt-2">Sign up to start managing your shifts</p>
  </div>
  
  <form on:submit|preventDefault={handleSubmit} class="bg-background rounded-card shadow-card p-6">
    {#if $authStore.authError}
      <div class="mb-4 p-3 bg-red-50 border border-red-200 text-red-700 rounded-md text-sm">
        {$authStore.authError}
      </div>
    {/if}
    
    <div class="mb-4">
      <label for="name" class="block text-sm font-medium mb-1">Full Name</label>
      <input 
        type="text" 
        id="name" 
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
      disabled={!isFormValid || $authStore.isLoading}
    >
      {$authStore.isLoading ? 'Creating account...' : 'Create Account'}
    </button>
    
    <div class="mt-4 text-center text-sm">
      <span class="text-muted-foreground">Already have an account?</span>
      <a href="/login" class="text-blue-600 hover:underline ml-1">Sign in</a>
    </div>
  </form>
</div> 