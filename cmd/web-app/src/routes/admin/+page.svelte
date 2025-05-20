<script>
  import { shifts, requests, updateRequestStatus, reassignShift } from '$lib/stores/shiftStore.js';
  import RequestItem from '$lib/components/RequestItem.svelte';
  import ShiftCard from '$lib/components/ShiftCard.svelte';
  import UserSelector from '$lib/components/UserSelector.svelte';
  
  // Get pending requests
  const pendingRequests = $derived($requests.filter(req => req.status === 'pending'));
  const sortedPendingRequests = $derived([...pendingRequests].sort((a, b) => new Date(a.requestDate).getTime() - new Date(b.requestDate).getTime()));
  
  // Get all assigned shifts
  const assignedShifts = $derived($shifts.filter(shift => shift.assignedTo !== null));
  const sortedAssignedShifts = $derived([...assignedShifts].sort((a, b) => new Date(a.date).getTime() - new Date(b.date).getTime()));
  
  let activeTab = $state('requests');
  let selectedShiftId = null;
  
  function setActiveTab(tab) {
    activeTab = tab;
  }
  
  function handleApprove(requestId) {
    try {
      updateRequestStatus(requestId, 'approved');
      alert('Request approved successfully.');
    } catch (error) {
      alert(`Error: ${error.message}`);
    }
  }
  
  function handleReject(requestId) {
    try {
      updateRequestStatus(requestId, 'rejected');
      alert('Request rejected successfully.');
    } catch (error) {
      alert(`Error: ${error.message}`);
    }
  }
  
  function selectShiftForReassign(shiftId) {
    selectedShiftId = shiftId;
  }
  
  function handleReassign(newUserId) {
    if (!selectedShiftId) return;
    
    try {
      reassignShift(selectedShiftId, newUserId);
      alert('Shift reassigned successfully.');
      selectedShiftId = null;
    } catch (error) {
      alert(`Error: ${error.message}`);
    }
  }
</script>

<div>
  <h2 class="text-2xl font-bold mb-6">Admin Panel</h2>
  
  <div class="mb-6">
    <div class="rounded-9px bg-muted border-muted shadow-mini-inset flex w-fit gap-1 p-1 text-sm font-semibold">
      <button
        type="button"
        class="h-8 rounded-[7px] px-4 py-2 {activeTab === 'requests' ? 'shadow-mini bg-white' : ''}"
        onclick={() => setActiveTab('requests')}
      >
        Pending Requests ({pendingRequests.length})
      </button>
      <button
        type="button"
        class="h-8 rounded-[7px] px-4 py-2 {activeTab === 'shifts' ? 'shadow-mini bg-white' : ''}"
        onclick={() => setActiveTab('shifts')}
      >
        Manage Shifts
      </button>
    </div>
    
    <div class="mt-6" style="display: {activeTab === 'requests' ? 'block' : 'none'}">
      {#if sortedPendingRequests.length > 0}
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          {#each sortedPendingRequests as request}
            <RequestItem 
              request={request} 
              isAdmin={true} 
              onApprove={handleApprove} 
              onReject={handleReject} 
            />
          {/each}
        </div>
      {:else}
        <div class="bg-background rounded-card shadow-card p-8 text-center">
          <h3 class="text-xl font-semibold">No pending requests</h3>
          <p class="text-muted-foreground">There are no pending shift requests at the moment.</p>
        </div>
      {/if}
    </div>
    
    <div class="mt-6" style="display: {activeTab === 'shifts' ? 'block' : 'none'}">
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div class="md:col-span-2">
          <h3 class="text-xl font-semibold mb-4">Assigned Shifts</h3>
          
          {#if sortedAssignedShifts.length > 0}
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              {#each sortedAssignedShifts as shift}
                <button 
                  class="w-full text-left cursor-pointer" 
                  class:ring-2={selectedShiftId === shift.id} 
                  class:ring-blue-500={selectedShiftId === shift.id}
                  class:rounded-card={selectedShiftId === shift.id}
                  onclick={() => selectShiftForReassign(shift.id)}
                  onkeydown={(e) => e.key === 'Enter' && selectShiftForReassign(shift.id)}
                  type="button"
                  aria-pressed={selectedShiftId === shift.id}
                >
                  <ShiftCard shift={shift} />
                </button>
              {/each}
            </div>
          {:else}
            <div class="bg-background rounded-card shadow-card p-8 text-center">
              <h3 class="text-xl font-semibold">No assigned shifts</h3>
              <p class="text-muted-foreground">There are no assigned shifts at the moment.</p>
            </div>
          {/if}
        </div>
        
        <div>
          <h3 class="text-xl font-semibold mb-4">Reassign Shift</h3>
          
          {#if selectedShiftId}
            <div class="bg-background rounded-card shadow-card p-4">
              <p class="text-sm text-muted-foreground mb-2">
                Select an employee to reassign the selected shift:
              </p>
              
              {#if selectedShiftId}
                {@const shift = $shifts.find(s => s.id === selectedShiftId)}
                <UserSelector 
                  excludeUserId={shift?.assignedTo}
                  onSelect={handleReassign}
                />
              {/if}
            </div>
          {:else}
            <div class="bg-background rounded-card shadow-card p-4">
              <p class="text-sm text-muted-foreground">
                Select a shift from the list to reassign it to another employee.
              </p>
            </div>
          {/if}
        </div>
      </div>
    </div>
  </div>
</div> 