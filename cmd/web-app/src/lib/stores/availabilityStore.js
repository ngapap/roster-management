import { writable, derived } from 'svelte/store';
import { currentUser } from './shiftStore.js';

// Mock user availability data
// Structure: { userId: { date: { available: boolean, note: string } } }
const mockAvailability = {
  1: {}, // John Doe's availability
  2: {}, // Jane Smith's availability
  3: {}  // Admin's availability
};

// Create availability store
export const userAvailability = writable(mockAvailability);

// Get current user's availability
export const currentUserAvailability = derived(
  userAvailability,
  ($userAvailability) => {
    // Get a safe read of the current user ID
    let userId = 1; // Default to first user
    
    // This will be populated by browser-side code
    try {
      if (typeof window !== 'undefined') {
        currentUser.subscribe(user => {
          if (user && user.id) {
            userId = user.id;
          }
        })();
      }
    } catch (e) {
      console.error('Error accessing currentUser:', e);
    }
    
    return $userAvailability[userId] || {};
  }
);

// Set availability for a specific date
export function setAvailability(userId, date, isAvailable, note = '') {
  userAvailability.update(availability => {
    const userAvail = availability[userId] || {};
    
    // Create or update availability for this date
    userAvail[date] = {
      available: isAvailable,
      note: note
    };
    
    // Update the store
    return {
      ...availability,
      [userId]: userAvail
    };
  });
}

// Check if a user is available on a specific date
export function isUserAvailable(userId, date) {
  // Get current store value
  let isAvailable = true; // Default to available if no preference is set
  
  // Get the current value from the store without updating it
  let currentAvailability;
  userAvailability.subscribe(val => {
    currentAvailability = val;
  })();
  
  const userAvail = currentAvailability[userId] || {};
  const dateAvail = userAvail[date];
  
  // If preference is set, use it, otherwise default to available
  if (dateAvail) {
    isAvailable = dateAvail.available;
  }
  
  return isAvailable;
}

// Get available dates for a user in a date range
export function getAvailableDates(userId, startDate, endDate) {
  const start = new Date(startDate);
  const end = new Date(endDate);
  const availableDates = [];
  
  // Loop through each day in the range
  const current = new Date(start);
  while (current <= end) {
    const dateStr = current.toISOString().split('T')[0];
    
    // Check if available
    if (isUserAvailable(userId, dateStr)) {
      availableDates.push(dateStr);
    }
    
    // Move to next day
    current.setDate(current.getDate() + 1);
  }
  
  return availableDates;
}

// Get unavailable dates (when explicitly marked as unavailable)
export function getUnavailableDates(userId) {
  let unavailableDates = [];
  
  // Get the current value from the store without updating it
  userAvailability.subscribe(availability => {
    const userAvail = availability[userId] || {};
    
    // Filter for dates explicitly marked as unavailable
    unavailableDates = Object.entries(userAvail)
      .filter(([date, pref]) => pref.available === false)
      .map(([date]) => date);
  })();
  
  return unavailableDates;
} 