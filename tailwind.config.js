/** @type {import('tailwindcss').Config} */
export default {
  content: [
    './web/templates/**/*.html',
    './web/assets/**/*.js',
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          50: '#eff6ff',
          500: '#3b82f6',
          600: '#2563eb',
          700: '#1d4ed8',
        }
      }
    },
  },
  plugins: [],
}
