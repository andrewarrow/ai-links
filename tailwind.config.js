/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['views/*.html','tailwind/*.html',],
  theme: {
    extend: {
      fontFamily: {
        poppins: ['Poppins'],
        montserrat: ['Montserrat'],
      },
      colors: {
        'cream': '#EFDECD',
        'lime': '#8FBC8F'
      },
    },
  },
  plugins: [],
}
