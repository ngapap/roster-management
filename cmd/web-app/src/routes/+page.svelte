<script>
  import { shifts, currentUser } from '$lib/stores/shiftStore.js';
  import { authStore } from '$lib/stores/authStore.js';
  import ShiftCard from '$lib/components/ShiftCard.svelte';
  
  // Get user's assigned shifts
  const userShifts = $derived($shifts.filter(shift => shift.assignedTo === $currentUser.id));
  const sortedUserShifts = $derived([...userShifts].sort((a, b) => new Date(a.date).getTime() - new Date(b.date).getTime()));
</script>

{#if $authStore.isAuthenticated}
  <div>
    <h2 class="text-2xl font-bold mb-6">My Schedule</h2>
    
    {#if sortedUserShifts.length > 0}
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {#each sortedUserShifts as shift}
          <ShiftCard {shift} />
        {/each}
      </div>
    {:else}
      <div class="bg-background rounded-card shadow-card p-8 text-center">
        <h3 class="text-xl font-semibold mb-2">No shifts scheduled</h3>
        <p class="text-muted-foreground mb-4">You don't have any shifts assigned to you yet.</p>
        <a href="/available" class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">
          Find Available Shifts
        </a>
      </div>
    {/if}
  </div>
{:else}
  <div class="text-center max-w-3xl mx-auto">
    <h1 class="text-4xl font-bold mb-4">Roster Management System</h1>
    <p class="text-xl text-muted-foreground mb-8">
      The simplest way to manage your team's shifts and scheduling
    </p>
    
    <div class="flex justify-center gap-4 mb-12">
      <a href="/signup" class="bg-blue-500 hover:bg-blue-600 text-white font-bold py-3 px-6 rounded-md text-lg">
        Create Account
      </a>
      <a href="/login" class="bg-muted hover:bg-gray-200 text-foreground font-bold py-3 px-6 rounded-md text-lg">
        Sign In
      </a>
    </div>
    
    <div class="grid grid-cols-1 md:grid-cols-3 gap-8 mt-12">
      <div class="bg-background rounded-card shadow-card p-6">
        <div class="text-blue-500 text-3xl mb-4">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 mx-auto" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
          </svg>
        </div>
        <h2 class="text-xl font-bold mb-2">View Shifts</h2>
        <p class="text-muted-foreground">
          See your scheduled shifts and upcoming work days at a glance
        </p>
      </div>
      
      <div class="bg-background rounded-card shadow-card p-6">
        <div class="text-blue-500 text-3xl mb-4">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 mx-auto" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 15l-2 5L9 9l11 4-5 2zm0 0l5 5M7.188 2.239l.777 2.897M5.136 7.965l-2.898-.777M13.95 4.05l-2.122 2.122m-5.657 5.656l-2.12 2.122" />
          </svg>
        </div>
        <h2 class="text-xl font-bold mb-2">Request Shifts</h2>
        <p class="text-muted-foreground">
          Request available shifts that fit your schedule and preferences
        </p>
      </div>
      
      <div class="bg-background rounded-card shadow-card p-6">
        <div class="text-blue-500 text-3xl mb-4">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 mx-auto" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        <h2 class="text-xl font-bold mb-2">Track Status</h2>
        <p class="text-muted-foreground">
          Track the status of your shift requests and approvals
        </p>
      </div>
    </div>
  </div>
{/if}
