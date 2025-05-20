<script>
  import { requests, currentUser, shifts } from '$lib/stores/shiftStore.js';
  import RequestItem from '$lib/components/RequestItem.svelte';
  import { Tabs } from 'bits-ui';
  
  // Filter requests for current user
  $: userRequests = $requests.filter(req => req.userId === $currentUser.id);
  
  // Sort by newest first
  $: sortedRequests = [...userRequests].sort((a, b) => new Date(b.requestDate).getTime() - new Date(a.requestDate).getTime());
  
  // Filter by status
  $: pendingRequests = sortedRequests.filter(req => req.status === 'pending');
  $: approvedRequests = sortedRequests.filter(req => req.status === 'approved');
  $: rejectedRequests = sortedRequests.filter(req => req.status === 'rejected');
  
  let activeTab = 'all';
</script>

<div>
  <h2 class="text-2xl font-bold mb-6">My Shift Requests</h2>
  
  <Tabs.Root bind:value={activeTab} class="mb-6">
    <Tabs.List class="rounded-9px bg-muted border-muted shadow-mini-inset flex w-fit gap-1 p-1 text-sm font-semibold">
      <Tabs.Trigger
        value="all"
        class="data-[state=active]:shadow-mini data-[state=active]:bg-white h-8 rounded-[7px] px-4 py-2"
      >
        All ({sortedRequests.length})
      </Tabs.Trigger>
      <Tabs.Trigger
        value="pending"
        class="data-[state=active]:shadow-mini data-[state=active]:bg-white h-8 rounded-[7px] px-4 py-2"
      >
        Pending ({pendingRequests.length})
      </Tabs.Trigger>
      <Tabs.Trigger
        value="approved"
        class="data-[state=active]:shadow-mini data-[state=active]:bg-white h-8 rounded-[7px] px-4 py-2"
      >
        Approved ({approvedRequests.length})
      </Tabs.Trigger>
      <Tabs.Trigger
        value="rejected"
        class="data-[state=active]:shadow-mini data-[state=active]:bg-white h-8 rounded-[7px] px-4 py-2"
      >
        Rejected ({rejectedRequests.length})
      </Tabs.Trigger>
    </Tabs.List>
    
    <Tabs.Content value="all">
      {#if sortedRequests.length > 0}
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          {#each sortedRequests as request}
            <RequestItem {request} />
          {/each}
        </div>
      {:else}
        <div class="bg-background rounded-card shadow-card p-8 text-center">
          <h3 class="text-xl font-semibold">No requests</h3>
          <p class="text-muted-foreground mb-4">You haven't requested any shifts yet.</p>
          <a href="/available" class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">
            Find Available Shifts
          </a>
        </div>
      {/if}
    </Tabs.Content>
    
    <Tabs.Content value="pending">
      {#if pendingRequests.length > 0}
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          {#each pendingRequests as request}
            <RequestItem {request} />
          {/each}
        </div>
      {:else}
        <div class="bg-background rounded-card shadow-card p-8 text-center">
          <h3 class="text-xl font-semibold">No pending requests</h3>
          <p class="text-muted-foreground">You don't have any pending shift requests.</p>
        </div>
      {/if}
    </Tabs.Content>
    
    <Tabs.Content value="approved">
      {#if approvedRequests.length > 0}
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          {#each approvedRequests as request}
            <RequestItem {request} />
          {/each}
        </div>
      {:else}
        <div class="bg-background rounded-card shadow-card p-8 text-center">
          <h3 class="text-xl font-semibold">No approved requests</h3>
          <p class="text-muted-foreground">You don't have any approved shift requests.</p>
        </div>
      {/if}
    </Tabs.Content>
    
    <Tabs.Content value="rejected">
      {#if rejectedRequests.length > 0}
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          {#each rejectedRequests as request}
            <RequestItem {request} />
          {/each}
        </div>
      {:else}
        <div class="bg-background rounded-card shadow-card p-8 text-center">
          <h3 class="text-xl font-semibold">No rejected requests</h3>
          <p class="text-muted-foreground">You don't have any rejected shift requests.</p>
        </div>
      {/if}
    </Tabs.Content>
  </Tabs.Root>
</div> 