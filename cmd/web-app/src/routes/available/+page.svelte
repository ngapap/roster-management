<script>
  import { shifts, currentUser, requestShift } from '$lib/stores/shiftStore.js';
  import ShiftCard from '$lib/components/ShiftCard.svelte';
  
  // Get available shifts
  const availableShifts = $derived($shifts.filter(shift => shift.assignedTo === null && shift.status === 'available'));
  const sortedAvailableShifts = $derived([...availableShifts].sort((a, b) => new Date(a.date).getTime() - new Date(b.date).getTime()));
  
  // Current and future dates only
  const today = $state(new Date());
  $effect(() => {
    today.setHours(0, 0, 0, 0);
  });
  
  const upcomingShifts = $derived(sortedAvailableShifts.filter(shift => new Date(shift.date) >= today));
  
  const thisWeekShifts = $derived(upcomingShifts.filter(shift => {
    const shiftDate = new Date(shift.date);
    const dayDiff = Math.floor((shiftDate.getTime() - today.getTime()) / (1000 * 60 * 60 * 24));
    return dayDiff < 7;
  }));
  
  const laterShifts = $derived(upcomingShifts.filter(shift => {
    const shiftDate = new Date(shift.date);
    const dayDiff = Math.floor((shiftDate.getTime() - today.getTime()) / (1000 * 60 * 60 * 24));
    return dayDiff >= 7;
  }));
  
  let activeTab = $state('all');
  
  function setActiveTab(tab) {
    activeTab = tab;
  }
  
  /**
   * Handle shift request for the given shift ID
   * @param {number | undefined} id - The ID of the shift to request
   */
  function handleRequest(id) {
    if (id === undefined) return;
    
    try {
      requestShift(id, $currentUser.id);
      alert('Shift requested successfully.');
    } catch (error) {
      if (error instanceof Error) {
        alert(`Error: ${error.message}`);
      } else {
        alert('An unknown error occurred');
      }
    }
  }
</script>

<div>
  <h2 class="text-2xl font-bold mb-6">Available Shifts</h2>
  
  <div class="mb-6">
    <div class="rounded-9px bg-muted border-muted shadow-mini-inset flex w-fit gap-1 p-1 text-sm font-semibold">
      <button
        type="button"
        class="h-8 rounded-[7px] px-4 py-2 {activeTab === 'all' ? 'shadow-mini bg-white' : ''}"
        onclick={() => setActiveTab('all')}
      >
        All
      </button>
      <button
        type="button" 
        class="h-8 rounded-[7px] px-4 py-2 {activeTab === 'thisWeek' ? 'shadow-mini bg-white' : ''}"
        onclick={() => setActiveTab('thisWeek')}
      >
        This Week
      </button>
      <button
        type="button"
        class="h-8 rounded-[7px] px-4 py-2 {activeTab === 'later' ? 'shadow-mini bg-white' : ''}"
        onclick={() => setActiveTab('later')}
      >
        Later
      </button>
    </div>
    
    <div class="mt-6" style="display: {activeTab === 'all' ? 'block' : 'none'}">
      {#if upcomingShifts.length > 0}
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {#each upcomingShifts as shift}
            <ShiftCard shift={shift} onRequest={handleRequest} />
          {/each}
        </div>
      {:else}
        <div class="bg-background rounded-card shadow-card p-8 text-center">
          <h3 class="text-xl font-semibold">No available shifts</h3>
          <p class="text-muted-foreground">There are no available shifts at the moment.</p>
        </div>
      {/if}
    </div>
    
    <div class="mt-6" style="display: {activeTab === 'thisWeek' ? 'block' : 'none'}">
      {#if thisWeekShifts.length > 0}
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {#each thisWeekShifts as shift}
            <ShiftCard shift={shift} onRequest={handleRequest} />
          {/each}
        </div>
      {:else}
        <div class="bg-background rounded-card shadow-card p-8 text-center">
          <h3 class="text-xl font-semibold">No available shifts</h3>
          <p class="text-muted-foreground">There are no available shifts for this week.</p>
        </div>
      {/if}
    </div>
    
    <div class="mt-6" style="display: {activeTab === 'later' ? 'block' : 'none'}">
      {#if laterShifts.length > 0}
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {#each laterShifts as shift}
            <ShiftCard shift={shift} onRequest={handleRequest} />
          {/each}
        </div>
      {:else}
        <div class="bg-background rounded-card shadow-card p-8 text-center">
          <h3 class="text-xl font-semibold">No available shifts</h3>
          <p class="text-muted-foreground">There are no available shifts for later dates.</p>
        </div>
      {/if}
    </div>
  </div>
</div> 