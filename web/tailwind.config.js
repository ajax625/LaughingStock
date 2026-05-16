/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{ts,tsx}'],
  theme: {
    extend: {
      colors: {
        long: '#16a34a',
        short: '#dc2626',
      },
    },
  },
  plugins: [],
}
