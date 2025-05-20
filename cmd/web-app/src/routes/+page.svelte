<script>
  import { goto } from '$app/navigation';
  
  let {data} = $props();
  
  function handleGetStarted() {
    if (data.user.isAdmin) {
      goto('/admin');
    } else {
      goto('/assignment');
    }
  }
  console.log(data);
</script>

<div class="min-h-screen flex flex-col items-center justify-center p-4">
  <div class="max-w-2xl text-center">
    <h1 class="text-4xl font-bold mb-4">Welcome to Roster Management</h1>
    <p class="text-xl text-muted-foreground mb-8">
      Streamline your shift management and scheduling process
    </p>
    
    <div class="flex flex-col sm:flex-row gap-4 justify-center">
      <button
        class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded cursor-pointer"
        onclick={handleGetStarted}
      >
        {data.user.isAdmin ? 'Go to Dashboard' : 'Get Started'}
      </button>
      
      {#if !data.user.isAuthenticated}
        <a
          href="/login"
          class="px-6 py-3 border border-blue-500 text-blue-500 hover:bg-blue-50 font-medium rounded-md transition-colors"
        >
          Sign In
        </a>
      {/if}
    </div>
  </div>
  
  <div class="mt-16 grid grid-cols-1 md:grid-cols-3 gap-8 max-w-4xl">
    <div class="p-6 bg-background/80 backdrop-blur-sm rounded-card shadow-card">
      <h3 class="text-xl font-semibold mb-2">Easy Scheduling</h3>
      <p class="text-muted-foreground">Create and manage shifts with an intuitive interface</p>
    </div>
    
    <div class="p-6 bg-background/80 backdrop-blur-sm rounded-card shadow-card">
      <h3 class="text-xl font-semibold mb-2">Team Collaboration</h3>
      <p class="text-muted-foreground">Coordinate with your team members seamlessly</p>
    </div>
    
    <div class="p-6 bg-background/80 backdrop-blur-sm rounded-card shadow-card">
      <h3 class="text-xl font-semibold mb-2">Real-time Updates</h3>
      <p class="text-muted-foreground">Stay informed with instant notifications</p>
    </div>
  </div>
</div> 