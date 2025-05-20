<script>
  import { enhance } from '$app/forms';
  import { goto } from '$app/navigation';

  let { form } = $props();

  let isLoading = $state(false);

</script>

<div class="max-w-md mx-auto">
  <div class="text-center mb-8">
    <h1 class="text-2xl font-bold">Sign In</h1>
    <p class="text-muted-foreground mt-2">Enter your credentials to access your account</p>
  </div>
  
  <form action="?/login" method="post" enctype="application/x-www-form-urlencoded" class="bg-background rounded-card shadow-card p-6" use:enhance={() => {
			isLoading = true;
			return async ({ update, result }) => {
				await update();
        isLoading = false;
        console.log("result", result);
        if (result.data?.success) {
          goto(result.data.redirectTo);
        }
			};
		}}>
    {#if form?.error}
      <div class="mb-4 p-3 bg-red-50 border border-red-200 text-red-700 rounded-md text-sm">
        {form.error}
      </div>
    {/if}
    
    <div class="mb-4">
      <label for="email" class="block text-sm font-medium mb-1">Email</label>
      <input 
        type="email" 
        id="email" 
        name="email"
        value={form?.email ?? ''}
        placeholder="you@example.com"
        class="w-full rounded-input border-border-input focus:border-border-input-hover p-2 text-sm"
        required
      />
    </div>
    
    <div class="mb-6">
      <div class="flex justify-between items-center mb-1">
        <label for="password" class="block text-sm font-medium">Password</label>
      </div>
      <input 
        type="password" 
        id="password" 
        name="password"
        placeholder="••••••••"
        class="w-full rounded-input border-border-input focus:border-border-input-hover p-2 text-sm"
        required
        minlength="6"
      />
    </div>
    
    <button 
      type="submit" 
      class="w-full bg-blue-500 hover:bg-blue-600 text-white font-medium py-2 px-4 rounded-md disabled:opacity-50 disabled:cursor-not-allowed"
      disabled={isLoading}
    >
      {isLoading ? 'Signing in...' : 'Sign In'}
    </button>
    
    <div class="mt-4 text-center text-sm">
      <span class="text-muted-foreground">Don't have an account?</span>
      <a href="/signup" class="text-blue-600 hover:underline ml-1">Create one</a>
    </div>
  </form>
  
  <!-- <div class="mt-8 text-center text-sm text-muted-foreground">
    <p class="mb-2">Demo Accounts:</p>
    <div class="grid grid-cols-2 gap-4">
      <div class="bg-muted rounded-card p-3">
        <div class="font-medium mb-1">Employee</div>
        <div class="text-xs">john@example.com</div>
        <div class="text-xs">password123</div>
      </div>
      <div class="bg-muted rounded-card p-3">
        <div class="font-medium mb-1">Admin</div>
        <div class="text-xs">admin@roster.com</div>
        <div class="text-xs">Password.1</div>
      </div>
    </div>
  </div> -->
</div> 