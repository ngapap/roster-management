<script>
    import ShiftCard from '$lib/components/ShiftCard.svelte';
	import { shifts } from '$lib/stores/shiftStore.js';
    
    // Get user's assigned shifts
    let {data} = $props();
    const userShifts = $derived(data.shifts);
    const sortedUserShifts = $derived([...userShifts].sort((a, b) => new Date(a.start_time).getTime() - new Date(b.start_time).getTime()));
  </script>
  
  
    <div>
      <h2 class="text-2xl font-bold mb-6">My Assignment</h2>
      
      {#if data.shifts?.length > 0}
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {#each data.shifts as shift}
            <ShiftCard {shift} />
          {/each}
        </div>
      {:else}
        <div class="bg-background rounded-card shadow-card p-8 text-center">
          <h3 class="text-xl font-semibold mb-2">No shifts assigned</h3>
          <p class="text-muted-foreground mb-4">You don't have any shifts assigned to you yet.</p>
          <a href="/available" class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">
            Find Available Shifts
          </a>
        </div>
      {/if}
    </div>

  