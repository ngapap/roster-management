import { writable } from 'svelte/store';
import { isUserAvailable } from './availabilityStore.js';

// Get current date and some future dates
const today = new Date();
const futureDate1 = new Date();
futureDate1.setDate(today.getDate() + 1);
const futureDate2 = new Date();
futureDate2.setDate(today.getDate() + 3);
const futureDate3 = new Date();
futureDate3.setDate(today.getDate() + 7);
const futureDate4 = new Date();
futureDate4.setDate(today.getDate() + 10);

// Format date as YYYY-MM-DD
function formatDate(date) {
  return date.toISOString().split('T')[0];
}

// Mock user data
const mockUsers = [
  { id: 1, name: 'John Doe', role: 'employee' },
  { id: 2, name: 'Jane Smith', role: 'employee' },
  { id: 3, name: 'Admin User', role: 'admin' }
];

// Mock shift data
const mockShifts = [
  { 
    id: 1, 
    date: formatDate(today), 
    startTime: '09:00', 
    endTime: '17:00',
    assignedTo: 1,
    status: 'approved'
  },
  { 
    id: 2, 
    date: formatDate(futureDate1), 
    startTime: '09:00', 
    endTime: '17:00',
    assignedTo: null,
    status: 'available'
  },
  { 
    id: 3, 
    date: formatDate(futureDate2), 
    startTime: '09:00', 
    endTime: '17:00',
    assignedTo: null,
    status: 'available'
  },
  { 
    id: 4, 
    date: formatDate(futureDate3), 
    startTime: '09:00', 
    endTime: '17:00',
    assignedTo: 2,
    status: 'approved'
  },
  { 
    id: 5, 
    date: formatDate(futureDate4), 
    startTime: '09:00', 
    endTime: '17:00',
    assignedTo: null,
    status: 'available'
  }
];

// Mock shift requests
const mockRequests = [
  {
    id: 1,
    shiftId: 3,
    userId: 1,
    status: 'pending',
    requestDate: '2023-07-31'
  }
];

// Create stores
export const users = writable(mockUsers);
export const shifts = writable(mockShifts);
export const requests = writable(mockRequests);
export const currentUser = writable(mockUsers[0]); // Default to first employee

// Helper function to check if a shift request would create conflicts
export function checkForConflicts(userId, date, existingShifts) {
  // Check if user already has a shift on this date
  const hasShiftOnSameDay = existingShifts.some(shift => 
    shift.assignedTo === userId && shift.date === date
  );
  
  if (hasShiftOnSameDay) {
    return { valid: false, reason: 'You already have a shift scheduled for this day' };
  }
  
  // Check for pending requests on the same day
  const hasPendingRequestOnSameDay = requests.update(allRequests => {
    return allRequests.some(req => 
      req.userId === userId && 
      req.status === 'pending' && 
      existingShifts.find(s => s.id === req.shiftId)?.date === date
    );
  });
  
  if (hasPendingRequestOnSameDay) {
    return { valid: false, reason: 'You already have a pending request for this day' };
  }
  
  // Check if user has marked themselves as unavailable for this date
  const userAvailable = isUserAvailable(userId, date);
  if (userAvailable === false) {
    return { valid: false, reason: 'You have marked yourself as unavailable for this date' };
  }
  
  // Check weekly limit (max 5 shifts per week)
  const weekStart = new Date(date);
  weekStart.setDate(weekStart.getDate() - weekStart.getDay()); // Get to Sunday
  
  const weekEnd = new Date(weekStart);
  weekEnd.setDate(weekEnd.getDate() + 6); // Get to Saturday
  
  const shiftsThisWeek = existingShifts.filter(shift => {
    const shiftDate = new Date(shift.date);
    return shift.assignedTo === userId && 
           shiftDate >= weekStart && 
           shiftDate <= weekEnd;
  });
  
  if (shiftsThisWeek.length >= 5) {
    return { valid: false, reason: 'You have reached the maximum of 5 shifts per week' };
  }
  
  return { valid: true };
}

// Request a shift
export function requestShift(shiftId, userId) {
  shifts.update(allShifts => {
    const shift = allShifts.find(s => s.id === shiftId);
    
    // Check if shift is available
    if (!shift || shift.assignedTo !== null) {
      throw new Error('Shift is not available');
    }
    
    // Check for conflicts
    const conflict = checkForConflicts(userId, shift.date, allShifts);
    if (!conflict.valid) {
      throw new Error(conflict.reason);
    }
    
    // Add request
    requests.update(allRequests => {
      const newRequest = {
        id: allRequests.length + 1,
        shiftId,
        userId,
        status: 'pending',
        requestDate: new Date().toISOString().split('T')[0]
      };
      return [...allRequests, newRequest];
    });
    
    return allShifts;
  });
}

// Approve or reject a shift request
export function updateRequestStatus(requestId, newStatus) {
  requests.update(allRequests => {
    const requestIndex = allRequests.findIndex(r => r.id === requestId);
    if (requestIndex === -1) return allRequests;
    
    const updatedRequest = { ...allRequests[requestIndex], status: newStatus };
    const updatedRequests = [...allRequests];
    updatedRequests[requestIndex] = updatedRequest;
    
    // If approved, update the shift assignment
    if (newStatus === 'approved') {
      const request = updatedRequest;
      shifts.update(allShifts => {
        const shiftIndex = allShifts.findIndex(s => s.id === request.shiftId);
        if (shiftIndex === -1) return allShifts;
        
        const updatedShifts = [...allShifts];
        updatedShifts[shiftIndex] = {
          ...updatedShifts[shiftIndex],
          assignedTo: request.userId,
          status: 'approved'
        };
        
        return updatedShifts;
      });
      
      // Reject any other pending requests for this shift
      requests.update(reqs => {
        return reqs.map(r => {
          if (r.id !== requestId && r.shiftId === updatedRequest.shiftId && r.status === 'pending') {
            return { ...r, status: 'rejected' };
          }
          return r;
        });
      });
    }
    
    return updatedRequests;
  });
}

// Reassign a shift
export function reassignShift(shiftId, newUserId) {
  shifts.update(allShifts => {
    const shiftIndex = allShifts.findIndex(s => s.id === shiftId);
    if (shiftIndex === -1) return allShifts;
    
    // Check for conflicts for the new user
    const shift = allShifts[shiftIndex];
    const conflict = checkForConflicts(newUserId, shift.date, allShifts);
    if (!conflict.valid) {
      throw new Error(conflict.reason);
    }
    
    const updatedShifts = [...allShifts];
    updatedShifts[shiftIndex] = {
      ...updatedShifts[shiftIndex],
      assignedTo: newUserId,
      status: 'approved'
    };
    
    return updatedShifts;
  });
}

// Convert local date/time to UTC
export function toUTC(date, time) {
  const localDateTime = new Date(`${date}T${time}`);
  return new Date(localDateTime.getTime() - localDateTime.getTimezoneOffset() * 60000);
}

// Convert UTC to local date/time
export function fromUTC(utcDateTime) {
  const localDateTime = new Date(utcDateTime);
  const date = localDateTime.toISOString().split('T')[0];
  const time = localDateTime.toTimeString().split(' ')[0].substring(0, 5);
  return { date, time };
} 