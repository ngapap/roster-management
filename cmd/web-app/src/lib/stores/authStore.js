import { writable, derived } from 'svelte/store';
import { goto } from '$app/navigation';
import { currentUser, users } from './shiftStore.js';

// Mock user authentication data (email + password)
const mockUserAuth = [
  { id: 1, email: 'john@example.com', password: 'password123' },
  { id: 2, email: 'jane@example.com', password: 'password123' },
  { id: 3, email: 'admin@example.com', password: 'admin123' }
];

// Authentication state
export const authStore = writable({
  isAuthenticated: false,
  authError: null,
  isLoading: false
});

// Login function
export function login(email, password) {
  authStore.update(state => ({ ...state, isLoading: true, authError: null }));
  
  // Simulate API call delay
  setTimeout(() => {
    const foundUser = mockUserAuth.find(user => 
      user.email.toLowerCase() === email.toLowerCase() && user.password === password
    );
    
    if (foundUser) {
      // Find the corresponding user data
      const userData = mockUserAuth.find(user => user.id === foundUser.id);
      
      // Set the current user first before updating authentication state
      users.subscribe(allUsers => {
        const user = allUsers.find(u => u.id === foundUser.id);
        if (user) {
          currentUser.set(user);
          
          // Update auth state after setting the user
          authStore.update(state => ({ ...state, isAuthenticated: true, isLoading: false }));
          
          // Redirect to home after everything is updated
          setTimeout(() => {
            goto('/');
          }, 10);
        }
      })();
    } else {
      authStore.update(state => ({ 
        ...state, 
        authError: 'Invalid email or password', 
        isLoading: false 
      }));
    }
  }, 800); // Simulate network delay
}

// Signup function
export function signup(name, email, password) {
  authStore.update(state => ({ ...state, isLoading: true, authError: null }));
  
  // Check if email already exists
  setTimeout(() => {
    const emailExists = mockUserAuth.some(user => 
      user.email.toLowerCase() === email.toLowerCase()
    );
    
    if (emailExists) {
      authStore.update(state => ({ 
        ...state, 
        authError: 'Email already in use', 
        isLoading: false 
      }));
      return;
    }
    
    // Create new user
    const newUserId = mockUserAuth.length + 1;
    
    // Add to authentication store
    mockUserAuth.push({
      id: newUserId,
      email,
      password
    });
    
    // Add to user store
    users.update(allUsers => [
      ...allUsers,
      {
        id: newUserId,
        name,
        role: 'employee'
      }
    ]);
    
    // Create the user profile
    const newUser = {
      id: newUserId,
      name,
      role: 'employee'
    };
    
    // Set current user
    currentUser.set(newUser);
    
    // Log in the new user
    authStore.update(state => ({ ...state, isAuthenticated: true, isLoading: false }));
    
    // Redirect to home with a small delay
    setTimeout(() => {
      goto('/');
    }, 10);
    
  }, 800); // Simulate network delay
}

// Logout function
export function logout() {
  authStore.update(state => ({ ...state, isAuthenticated: false }));
  goto('/login');
}

// Check if user is logged in on page load (simulating token validation)
export function initAuth() {
  // In a real app, this would check for a valid token in localStorage
  // and make an API call to validate it
  
  // For demo purposes, we'll default to not authenticated
  authStore.update(state => ({ ...state, isAuthenticated: false }));
} 