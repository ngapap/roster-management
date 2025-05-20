/** @type {import('tailwindcss').Config} */
export default {
  content: ['./src/**/*.{html,js,svelte,ts}'],
  theme: {
    extend: {
      colors: {
        'background': '#ffffff',
        'background-alt': '#f9fafb',
        'foreground': '#0f172a',
        'muted': '#f1f5f9',
        'muted-foreground': '#64748b',
        'border-input': '#e2e8f0',
        'border-input-hover': '#94a3b8'
      },
      boxShadow: {
        'card': '0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px 0 rgba(0, 0, 0, 0.06)',
        'mini': '0 1px 2px 0 rgba(0, 0, 0, 0.05)',
        'mini-inset': 'inset 0 1px 2px 0 rgba(0, 0, 0, 0.05)',
        'date-field-focus': '0 0 0 2px rgba(99, 102, 241, 0.3)'
      },
      borderRadius: {
        'card': '0.5rem',
        'input': '0.375rem',
        '5px': '5px',
        '9px': '9px',
        '15px': '15px'
      }
    }
  },
  plugins: []
} 