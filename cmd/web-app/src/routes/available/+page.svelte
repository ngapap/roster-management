<script>

  import ShiftCard from '$lib/components/ShiftCard.svelte';
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';
  import { enhance } from '$app/forms';
  
  let {data, form} = $props();

  // Get available shifts
  const availableShifts = $derived(data.shifts.filter(shift => {
    // If no requests exist, show all shifts
    if (!data.shiftsRequest || data.shiftsRequest.length === 0) return true;
    
    // Check if this shift has any pending requests
    return !data.shiftsRequest.some(request => 
      request.shift_id === shift.id && request.status === 'pending'
    );
  }));
  
  const sortedAvailableShifts = $derived([...availableShifts].sort((a, b) => new Date(a.start_time).getTime() - new Date(b.start_time).getTime()));
  
  // Current and future dates only
  const today = $state(new Date());
  $effect(() => {
    today.setHours(0, 0, 0, 0);
  });
  
  const upcomingShifts = $derived(sortedAvailableShifts.filter(shift => new Date(shift.start_time) >= today));
  
  const thisWeekShifts = $derived(upcomingShifts.filter(shift => {
    const shiftDate = new Date(shift.start_time);
    const dayDiff = Math.floor((shiftDate.getTime() - today.getTime()) / (1000 * 60 * 60 * 24));
    return dayDiff < 7;
  }));
  
  const laterShifts = $derived(upcomingShifts.filter(shift => {
    const shiftDate = new Date(shift.start_time); 
    const dayDiff = Math.floor((shiftDate.getTime() - today.getTime()) / (1000 * 60 * 60 * 24));
    return dayDiff >= 7;
  }));
  
  let activeTab = $state('all');
  let isSubmitting = $state(false);
  
  function setActiveTab(tab) {
    activeTab = tab;
  }

  onMount(() => {
    if (data.status === 401) {
      goto('/login');
    }
  });
</script>

{#if data.error}
  <div class="error-message">
    {data.error}
  </div>
{:else}
  <div>
    <h2 class="text-2xl font-bold mb-6">Available Shifts</h2>
    
    <div class="mb-6">
      <div class="rounded-9px bg-muted border-muted shadow-mini-inset flex w-fit gap-1 p-1 text-sm font-semibold">
        <button
          type="button"
          class="h-8 rounded-[7px] px-4 py-2 transition-colors flex items-center justify-center min-w-[80px] cursor-pointer {activeTab === 'all' ? 'bg-blue-500 text-white shadow-mini' : 'text-gray-600 hover:bg-gray-100'}"
          onclick={() => setActiveTab('all')}
        >
          All
        </button>
        <button
          type="button" 
          class="h-8 rounded-[7px] px-4 py-2 transition-colors flex items-center justify-center min-w-[80px] cursor-pointer {activeTab === 'thisWeek' ? 'bg-blue-500 text-white shadow-mini' : 'text-gray-600 hover:bg-gray-100'}"
          onclick={() => setActiveTab('thisWeek')}
        >
          This Week
        </button>
        <button
          type="button"
          class="h-8 rounded-[7px] px-4 py-2 transition-colors flex items-center justify-center min-w-[80px] cursor-pointer {activeTab === 'later' ? 'bg-blue-500 text-white shadow-mini' : 'text-gray-600 hover:bg-gray-100'}"
          onclick={() => setActiveTab('later')}
        >
          Later
        </button>
      </div>
      
      <div class="mt-6" style="display: {activeTab === 'all' ? 'block' : 'none'}">
        {#if upcomingShifts.length > 0}
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {#each upcomingShifts as shift}
              <form 
                action="?/requestShift" 
                method="POST" 
                use:enhance={() => {
                  isSubmitting = true;
                  return async ({ update,result }) => {
                    await update();
                    isSubmitting = false;
                    console.log("result", result);
                    if (result.type === 'success') {
                      alert('Shift requested successfully');
                    } else {
                      alert(result.data?.error ?? 'Failed to request shift');
                    }
                  };
                }}
              >
                <input type="hidden" name="shiftId" value={shift.id} />
                <ShiftCard shift={shift} isSubmitting={isSubmitting} />
              </form>
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
              <form 
                action="?/requestShift" 
                method="POST" 
                use:enhance={() => {
                  isSubmitting = true;
                  return async ({ update,result }) => {
                    await update();
                    isSubmitting = false;
                    if (result.type === 'success') {
                      alert('Shift requested successfully');
                    } else {
                      alert(result.data?.error ?? 'Failed to request shift');
                    }
                  };
                }}
              >
                <input type="hidden" name="shiftId" value={shift.id} />
                <ShiftCard shift={shift} isSubmitting={isSubmitting} />
              </form>
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
              <form 
                action="?/requestShift" 
                method="POST" 
                use:enhance={() => {
                  isSubmitting = true;
                  return async ({ update,result }) => {
                    await update();
                    isSubmitting = false;
                    if (result.type === 'success') {
                      alert('Shift requested successfully');
                    } else {
                      alert(result.data?.error ?? 'Failed to request shift');
                    }
                  };
                }}
              >
                <input type="hidden" name="shiftId" value={shift.id} />
                <ShiftCard shift={shift} isSubmitting={isSubmitting} />
              </form>
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
{/if} 