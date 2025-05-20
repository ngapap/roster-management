<script>
  import { shifts, users } from '$lib/stores/shiftStore.js';
  
  const { request, name, onApprove = () => {}, onReject = () => {}, isAdmin = false } = $props();
  
  // Find the shift data
  const shift = $derived($shifts.find(s => s.id === request.shift_id));
  
  // Find the user data
  const user = $derived($users.find(u => u.id === request.worker_id));
  
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

<div class="bg-background rounded-lg shadow-card p-4 border border-gray-200 cursor-pointer max-w-md">
  <!-- {#if shift && user} -->
    <div class="flex justify-between items-start">
      <div>
        <h3 class="text-lg font-semibold">{request.shift_id}</h3>
        <p class="text-muted-foreground text-sm">Requested by: {name}</p>
        <p class="text-muted-foreground text-sm">Requested on: {formatDate(request.created_at)}</p>
      </div>
      
      <div class="text-sm rounded-full px-2 py-1 font-medium {getStatusClass(request.status)}">
        {request.status.charAt(0).toUpperCase() + request.status.slice(1)}
      </div>
    </div>
    
    {#if isAdmin && request.status === 'pending'}
      <div class="flex gap-2 mt-3">
        <button
          onclick={() => onApprove(request.id)}
          class="bg-green-500 hover:bg-green-700 text-white font-bold py-1 px-3 rounded text-sm"
        >
          Approve
        </button>
        <button
          onclick={() => onReject(request.id)}
          class="bg-red-500 hover:bg-red-700 text-white font-bold py-1 px-3 rounded text-sm"
        >
          Reject
        </button>
      </div>
    {/if}

</div> 