<script>
  import { users } from '$lib/stores/shiftStore.js';
  
  const { isSubmitting, shift } = $props(); 
  console.log("shift", shift);
  function formatTime(time) {
    const date = new Date(time);
    const utcHours = date.getUTCHours();
    const utcMinutes = date.getUTCMinutes().toString().padStart(2, '0');
    const period = utcHours >= 12 ? 'PM' : 'AM';
    const formattedHour = utcHours % 12 || 12;
    return `${formattedHour}:${utcMinutes} ${period} UTC`;
  }
</script>

<div class="bg-background rounded-lg shadow-card p-4 border border-gray-200">
  <div class="flex justify-between items-start mb-4">
    <div>
      <h3 class="text-lg font-semibold text-gray-900">
        {new Date(shift.start_time).toUTCString().split(' ').slice(0,4).join(' ')}
        {new Date(shift.start_time).getUTCDate() !== new Date(shift.end_time).getUTCDate() ? 
          `- ${new Date(shift.end_time).toUTCString().split(' ').slice(0,4).join(' ')}` : 
          ''}
      </h3>
      <p class="text-gray-600 mt-1">{formatTime(shift.start_time)} - {formatTime(shift.end_time)}</p>
      <p class="text-sm text-gray-500 mt-1">ID: {shift.id}</p>
    </div>
  </div>
  
  {#if shift.is_available}
  <div class="flex justify-start">
    <button 
      type="submit"
      class="px-5 py-2.5 bg-blue-500 hover:bg-blue-600 text-white rounded-md text-sm font-medium disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
      disabled={isSubmitting ?? false}
    >
      {isSubmitting ? 'Requesting...' : 'Request Shift'}
    </button>
  </div>
  {/if}
</div>