<script>
  import { users } from '$lib/stores/shiftStore.js';
  
  const { shift, onRequest = (id = 0) => {} } = $props();
  
  let assignedUser = $state(null);
  
  $effect(() => {
    if (shift.assignedTo) {
      assignedUser = $users.find(user => user.id === shift.assignedTo);
    } else {
      assignedUser = null;
    }
  });
  
  function formatTime(time) {
    const [hours, minutes] = time.split(':');
    const hour = parseInt(hours);
    const period = hour >= 12 ? 'PM' : 'AM';
    const formattedHour = hour % 12 || 12;
    return `${formattedHour}:${minutes} ${period}`;
  }
</script>

<div class="rounded-card border-muted bg-background shadow-card p-4 w-full flex flex-col">
  <div class="flex justify-between items-start mb-2">
    <div>
      <h3 class="text-lg font-semibold">{new Date(shift.date).toLocaleDateString('en-US', { weekday: 'short', month: 'short', day: 'numeric' })}</h3>
      <p class="text-muted-foreground text-sm">{formatTime(shift.startTime)} - {formatTime(shift.endTime)}</p>
    </div>
    
    <div class="text-sm rounded-full px-2 py-1 font-medium" class:bg-green-100={shift.status === 'approved'} class:text-green-800={shift.status === 'approved'} class:bg-blue-100={shift.status === 'available'} class:text-blue-800={shift.status === 'available'}>
      {shift.status === 'approved' ? 'Assigned' : 'Available'}
    </div>
  </div>
  
  {#if assignedUser}
    <p class="text-sm text-muted-foreground mt-1">Assigned to: {assignedUser.name}</p>
  {/if}
  
  {#if shift.status === 'available'}
    <button onclick={() => onRequest(shift.id)} class="mt-2 bg-blue-500 hover:bg-blue-700 text-white font-bold py-1 px-3 rounded text-sm">
      Request Shift
    </button>
  {/if}
</div> 