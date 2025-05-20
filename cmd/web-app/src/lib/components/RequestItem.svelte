<script>
  import { shifts, users } from '$lib/stores/shiftStore.js';
  
  const { request, onApprove = () => {}, onReject = () => {}, isAdmin = false } = $props();
  
  // Find the shift data
  const shift = $derived($shifts.find(s => s.id === request.shiftId));
  
  // Find the user data
  const user = $derived($users.find(u => u.id === request.userId));
  
  function formatDate(dateStr) {
    return new Date(dateStr).toLocaleDateString('en-US', { 
      weekday: 'short', 
      month: 'short', 
      day: 'numeric' 
    });
  }
  
  function getStatusClass(status) {
    switch(status) {
      case 'pending':
        return 'bg-yellow-100 text-yellow-800';
      case 'approved':
        return 'bg-green-100 text-green-800';
      case 'rejected':
        return 'bg-red-100 text-red-800';
      default:
        return '';
    }
  }
</script>

<div class="rounded-card border-muted bg-background shadow-card p-4 w-full">
  {#if shift && user}
    <div class="flex justify-between items-start">
      <div>
        <h3 class="text-lg font-semibold">{formatDate(shift.date)}</h3>
        <p class="text-muted-foreground text-sm">Requested by: {user.name}</p>
        <p class="text-muted-foreground text-sm">Requested on: {formatDate(request.requestDate)}</p>
      </div>
      
      <div class="text-sm rounded-full px-2 py-1 font-medium {getStatusClass(request.status)}">
        {request.status.charAt(0).toUpperCase() + request.status.slice(1)}
      </div>
    </div>
    
    {#if isAdmin && request.status === 'pending'}
      <div class="flex gap-2 mt-3">
        <button
          on:click={() => onApprove(request.id)}
          class="bg-green-500 hover:bg-green-700 text-white font-bold py-1 px-3 rounded text-sm"
        >
          Approve
        </button>
        <button
          on:click={() => onReject(request.id)}
          class="bg-red-500 hover:bg-red-700 text-white font-bold py-1 px-3 rounded text-sm"
        >
          Reject
        </button>
      </div>
    {/if}
  {:else}
    <p>Invalid request data</p>
  {/if}
</div> 