<script>
  import RequestItem from '$lib/components/RequestItem.svelte';
  
  let {data} = $props();

  // Filter requests for current user
  const userRequests = $derived(data.shiftsRequest);
  
  // Sort by newest first
  const sortedRequests = $derived([...userRequests].sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime()));
  
  // Filter by status
  const pendingRequests = $derived(sortedRequests.filter(req => req.status === 'pending'));
  const approvedRequests = $derived(sortedRequests.filter(req => req.status === 'approved'));
  const rejectedRequests = $derived(sortedRequests.filter(req => req.status === 'rejected'));
  const notSelectedRequests = $derived(sortedRequests.filter(req => req.status === 'not_selected'));
  
  let activeTab = $state('all');
</script>

<div>
  <h2 class="text-2xl font-bold mb-6">My Shift Requests</h2>
  
  <div class="mb-6">
    <div class="rounded-9px bg-muted border-muted shadow-mini-inset flex w-fit gap-1 p-1 text-sm font-semibold">
      <button
        class="h-8 rounded-[7px] px-4 py-2 transition-colors flex items-center justify-center min-w-[80px] cursor-pointer {activeTab === 'all' ? 'bg-blue-500 text-white shadow-mini' : 'text-gray-600 hover:bg-gray-100'}"
        onclick={() => activeTab = 'all'}
      >
        All ({sortedRequests.length})
      </button>
      <button
        class="h-8 rounded-[7px] px-4 py-2 transition-colors flex items-center justify-center min-w-[80px] cursor-pointer {activeTab === 'pending' ? 'bg-blue-500 text-white shadow-mini' : 'text-gray-600 hover:bg-gray-100'}"
        onclick={() => activeTab = 'pending'}
      >
        Pending ({pendingRequests.length})
      </button>
      <button
        class="h-8 rounded-[7px] px-4 py-2 transition-colors flex items-center justify-center min-w-[80px] cursor-pointer {activeTab === 'approved' ? 'bg-blue-500 text-white shadow-mini' : 'text-gray-600 hover:bg-gray-100'}"
        onclick={() => activeTab = 'approved'}
      >
        Approved ({approvedRequests.length})
      </button>
      <button
        class="h-8 rounded-[7px] px-4 py-2 transition-colors flex items-center justify-center min-w-[80px] cursor-pointer {activeTab === 'rejected' ? 'bg-blue-500 text-white shadow-mini' : 'text-gray-600 hover:bg-gray-100'}"
        onclick={() => activeTab = 'rejected'}
      >
        Rejected ({rejectedRequests.length})
      </button>
      <button
        class="h-8 rounded-[7px] px-4 py-2 transition-colors flex items-center justify-center min-w-[80px] cursor-pointer {activeTab === 'not_selected' ? 'bg-blue-500 text-white shadow-mini' : 'text-gray-600 hover:bg-gray-100'}"
        onclick={() => activeTab = 'not_selected'}
      >
        Not Selected ({notSelectedRequests.length})
      </button>
    </div>
    
    {#if activeTab === 'all'}
      {#if sortedRequests.length > 0}
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          {#each sortedRequests as request}
            <RequestItem {request} name={data.user.name} />
          {/each}
        </div>
      {:else}
        <div class="bg-background rounded-card shadow-card p-8 text-center">
          <h3 class="text-xl font-semibold">No requests</h3>
          <p class="text-muted-foreground mb-4">You haven't requested any shifts yet.</p>
          <a href="/available" class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded cursor-pointer">
            Find Available Shifts
          </a>
        </div>
      {/if}
    {/if}
    
    {#if activeTab === 'pending'}
      {#if pendingRequests.length > 0}
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          {#each pendingRequests as request}
            <RequestItem {request} name={data.user.name} />
          {/each}
        </div>
      {:else}
        <div class="bg-background rounded-card shadow-card p-8 text-center">
          <h3 class="text-xl font-semibold">No pending requests</h3>
          <p class="text-muted-foreground">You don't have any pending shift requests.</p>
        </div>
      {/if}
    {/if}
    
    {#if activeTab === 'approved'}
      {#if approvedRequests.length > 0}
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          {#each approvedRequests as request}
            <RequestItem {request} name={data.user.name} />
          {/each}
        </div>
      {:else}
        <div class="bg-background rounded-card shadow-card p-8 text-center">
          <h3 class="text-xl font-semibold">No approved requests</h3>
          <p class="text-muted-foreground">You don't have any approved shift requests.</p>
        </div>
      {/if}
    {/if}
    
    {#if activeTab === 'rejected'}
      {#if rejectedRequests.length > 0}
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          {#each rejectedRequests as request}
            <RequestItem {request} name={data.user.name} />
          {/each}
        </div>
      {:else}
        <div class="bg-background rounded-card shadow-card p-8 text-center">
          <h3 class="text-xl font-semibold">No rejected requests</h3>
          <p class="text-muted-foreground">You don't have any rejected shift requests.</p>
        </div>
      {/if}
    {/if}

    {#if activeTab === 'not_selected'}
      {#if notSelectedRequests.length > 0}
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          {#each notSelectedRequests as request}
            <RequestItem {request} name={data.user.name} />
          {/each}
        </div>
      {:else}
        <div class="bg-background rounded-card shadow-card p-8 text-center">
          <h3 class="text-xl font-semibold">No not selected requests</h3>
          <p class="text-muted-foreground">You don't have any not selected shift requests.</p>
        </div>
      {/if}
    {/if}
  </div>
</div> 