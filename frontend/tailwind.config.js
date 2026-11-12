/** @type {import('tailwindcss').Config} */
export default {
  content: ['./src/**/*.{html,js,svelte,ts}'],
  theme: {
    extend: {
      colors: {
        background: '#FDF6E3',
        foreground: '#111111',
        card: '#FFFDF5',
        border: '#111111',
        input: '#111111',
        ring: '#111111',
        primary: { DEFAULT: '#FFDB33', foreground: '#111111' },
        secondary: { DEFAULT: '#90A8FF', foreground: '#111111' },
        accent: { DEFAULT: '#FF90E8', foreground: '#111111' },
        destructive: { DEFAULT: '#FF7A7A', foreground: '#111111' },
        muted: { DEFAULT: '#F1E9D2', foreground: '#5C5545' },
        success: { DEFAULT: '#B8FF9F', foreground: '#111111' },
        warning: { DEFAULT: '#FFD166', foreground: '#111111' }
      },
      fontFamily: {
        sans: ['"Space Grotesk Variable"', 'ui-sans-serif', 'system-ui', 'sans-serif'],
        mono: ['"JetBrains Mono"', 'ui-monospace', 'SFMono-Regular', 'monospace']
      },
      boxShadow: {
        nb: '4px 4px 0 0 #111111',
        nbsm: '2px 2px 0 0 #111111',
        nblg: '8px 8px 0 0 #111111'
      },
      borderRadius: {
        DEFAULT: '4px',
        sm: '2px',
        md: '4px',
        lg: '6px',
        xl: '8px'
      }
    }
  },
  plugins: []
};
