<script>
  import { users } from '$lib/stores/shiftStore.js';
  
  const { onSelect = (id = 0) => {}, excludeUserId = null } = $props();
  
  // Filter to only show employees
  const employees = $derived($users.filter(user => user.role === 'employee' && user.id !== excludeUserId));
</script>

<div class="mt-2">
  <p id="reassign-label" class="block text-sm font-medium mb-1">Reassign to:</p>
  <div class="grid grid-cols-1 gap-2" role="radiogroup" aria-labelledby="reassign-label">
    {#each employees as employee}
      <button 
        on:click={() => onSelect(employee.id)}
        on:keydown={(e) => e.key === 'Enter' && onSelect(employee.id)}
        class="rounded-input border-border-input hover:border-border-input-hover bg-background text-left px-3 py-2 text-sm border flex items-center justify-between"
        role="radio"
        aria-checked="false"
        type="button"
      >
        <span>{employee.name}</span>
        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-muted-foreground" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
        </svg>
      </button>
    {/each}
  </div>
</div> 